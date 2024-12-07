package world

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/guregu/null/v5"
	"github.com/isucon/isucon14/bench/internal/concurrent"
)

type ChairState int

const (
	ChairStateInactive ChairState = iota
	ChairStateActive
)

type ChairID int

type Chair struct {
	// ID ベンチマーカー内部椅子ID
	ID ChairID
	// ServerID サーバー上での椅子ID
	ServerID string
	// World Worldへの逆参照
	World *World
	// Region 椅子がいる地域
	Region *Region
	// Owner 椅子を所有しているオーナー
	Owner *Owner
	// Model 椅子のモデル
	Model *ChairModel
	// State 椅子の状態
	State ChairState
	// Location 椅子の位置情報
	Location ChairLocation
	// RegisteredData サーバーに登録されている椅子情報
	RegisteredData RegisteredChairData
	// matchingData マッチング通知情報
	matchingData *ChairNotificationEventMatched
	// Request 進行中のリクエスト
	Request *Request
	// RequestHistory 引き受けたリクエストの履歴
	RequestHistory *concurrent.SimpleSlice[*Request]
	// Client webappへのクライアント
	Client ChairClient
	// Rand 専用の乱数
	Rand *rand.Rand
	// ActivatedAt Active化レスポンスが返ってきた日時
	ActivatedAt time.Time
	// tickDone 行動が完了しているかどうか
	tickDone tickDone
	// notificationConn 通知ストリームコネクション
	notificationConn NotificationStream
	// notificationQueue 通知キュー。毎Tickで最初に処理される
	notificationQueue chan NotificationEvent
	// forceStopped 強制停止されているかどうか
	forceStopped bool

	// detour 今回のリクエストで迂回するかどうか
	detour bool
	// detoured 今回のリクエストで迂回したかどうか
	detoured bool
	// detourPoint 迂回ポイント
	detourPoint Coordinate
	// detourIn Dispaching or Carryingのどっちで迂回するか
	detourIn RequestStatus
}

type RegisteredChairData struct {
	Name string
}

func (c *Chair) String() string {
	return fmt.Sprintf("Chair{id=%d,c=%s}", c.ID, c.Location.Current())
}

func (c *Chair) SetID(id ChairID) {
	c.ID = id
}

func (c *Chair) GetServerID() string {
	return c.ServerID
}

func (c *Chair) Tick(ctx *Context) error {
	if c.tickDone.DoOrSkip() {
		return nil
	}
	defer c.tickDone.Done()

	if c.forceStopped {
		if c.notificationConn != nil {
			c.notificationConn.Close()
			c.notificationConn = nil
		}
		return nil
	}

	// 通知キューを順番に処理する
	for event := range concurrent.TryIter(c.notificationQueue) {
		err := c.HandleNotification(event)
		if err != nil {
			return err
		}
	}

	switch {
	// 進行中のリクエストが存在
	case c.Request != nil:
		switch c.Request.Statuses.Chair {
		case RequestStatusMatching:
			// Active状態なら配車要求をACKする
			// そうでないなら、応答せずにハングさせる
			if c.State == ChairStateActive {
				c.Request.Statuses.Lock()

				c.Request.BenchRequestAcceptTime = time.Now()
				err := c.Client.SendAcceptRequest(ctx, c, c.Request)
				if err != nil {
					c.Request.BenchRequestAcceptTime = time.Time{}
					c.Request.Statuses.Unlock()

					return WrapCodeError(ErrorCodeFailedToAcceptRequest, err)
				}

				// サーバーに要求を受理の通知が通ったので配椅子地に向かう
				c.Request.Chair = c
				c.Request.Statuses.Desired = RequestStatusDispatching
				c.Request.Statuses.Chair = RequestStatusDispatching
				c.Request.StartPoint = null.ValueFrom(c.Location.Current())
				c.Request.MatchedAt = ctx.CurrentTime()

				c.Request.Statuses.Unlock()

				c.RequestHistory.Append(c.Request)
			}

		case RequestStatusDispatching:
			// 配車位置に向かう
			time := ctx.CurrentTime()
			if c.detour && c.detourIn == RequestStatusDispatching && !c.detoured {
				// 迂回する予定でまだ迂回してない場合
				if c.Location.Current().Equals(c.detourPoint) {
					// 迂回ポイントに着いた次の移動は配車位置から離れる方向に行う
					c.Location.MoveTo(&LocationEntry{
						Coord: c.moveOppositeTo(c.Request.PickupPoint),
						Time:  time,
					})
					c.detoured = true
				} else {
					// 迂回ポイントに向かう
					c.Location.MoveTo(&LocationEntry{
						Coord: c.moveToward(c.detourPoint),
						Time:  time,
					})
				}
			} else {
				// 配車位置に向かう
				c.Location.MoveTo(&LocationEntry{
					Coord: c.moveToward(c.Request.PickupPoint),
					Time:  time,
				})
			}

			if c.Location.Current().Equals(c.Request.PickupPoint) {
				// 配車位置に到着
				c.Request.Statuses.Desired = RequestStatusDispatched
				c.Request.Statuses.Chair = RequestStatusDispatched
				c.Request.DispatchedAt = time
			}

		case RequestStatusDispatched:
			// 乗客を乗せて出発しようとする
			if c.Request.Statuses.User != RequestStatusDispatched {
				// ただし、ユーザーに到着通知が行っていないとユーザーは乗らない振る舞いをするので
				// ユーザー側の状態が変わるまで待機する
				// 一向にユーザーの状態が変わらない場合は、この椅子の行動はハングする
				break
			}

			err := c.Client.SendDepart(ctx, c.Request)
			if err != nil {
				return WrapCodeError(ErrorCodeFailedToDepart, err)
			}

			// サーバーがdepartを受理したので出発する
			c.Request.Statuses.Desired = RequestStatusCarrying
			c.Request.Statuses.Chair = RequestStatusCarrying
			c.Request.PickedUpAt = ctx.CurrentTime()

		case RequestStatusCarrying:
			// 目的地に向かう
			time := ctx.CurrentTime()
			if c.detour && c.detourIn == RequestStatusCarrying && !c.detoured {
				// 迂回する予定でまだ迂回してない場合
				if c.Location.Current().Equals(c.detourPoint) {
					// 迂回ポイントに着いた次の移動は目的地から離れる方向に行う
					c.Location.MoveTo(&LocationEntry{
						Coord: c.moveOppositeTo(c.Request.DestinationPoint),
						Time:  time,
					})
					c.detoured = true
				} else {
					// 迂回ポイントに向かう
					c.Location.MoveTo(&LocationEntry{
						Coord: c.moveToward(c.detourPoint),
						Time:  time,
					})
				}
			} else {
				// 目的地に向かう
				c.Location.MoveTo(&LocationEntry{
					Coord: c.moveToward(c.Request.DestinationPoint),
					Time:  time,
				})
			}

			if c.Location.Current().Equals(c.Request.DestinationPoint) {
				// 目的地に到着
				c.Request.Statuses.Desired = RequestStatusArrived
				c.Request.Statuses.Chair = RequestStatusArrived
				c.Request.ArrivedAt = time
				break
			}

		case RequestStatusArrived:
			// 客が評価するまで待機する
			// 一向に評価されない場合は、この椅子の行動はハングする
			break

		case RequestStatusCompleted:
			slog.Warn("unexpected state")
			break
		}

	// アサインされたリクエストが存在するが、詳細を未取得
	case c.Request == nil && c.matchingData != nil:
		req := c.World.RequestDB.GetByServerID(c.matchingData.ServerRequestID)
		if req == nil {
			// ロックの関係でまだRequestDBに入ってないreqのため、次のtickで処理する
			// ベンチマーク外で作成されたリクエストがアサインされた場合はどうしようもできないのでハングする
			return nil
		}

		if !c.matchingData.Destination.Equals(req.DestinationPoint) ||
			!c.matchingData.Pickup.Equals(req.PickupPoint) ||
			c.matchingData.User.ID != req.User.ServerID ||
			c.matchingData.User.Name != req.User.RegisteredData.FirstName+" "+req.User.RegisteredData.LastName {
			c.forceStopped = true
			return CodeError(ErrorCodeChairReceivedDataIsWrong)
		}

		// 椅子がリクエストを正常に認識する
		c.Request = req
		// 10%の確率で迂回させる(最短距離より1単位速度分だけ遠回しさせる)
		c.detour = c.Rand.Float64() < 0.1
		c.detoured = false
		if c.detour {
			if c.Rand.IntN(2) == 0 {
				c.detourIn = RequestStatusDispatching
				c.detourPoint = CalculateRandomDetourPoint(c.Location.Current(), c.Request.PickupPoint, c.Model.Speed, c.Rand)
			} else {
				c.detourIn = RequestStatusCarrying
				c.detourPoint = CalculateRandomDetourPoint(c.Request.PickupPoint, c.Request.DestinationPoint, c.Model.Speed, c.Rand)
			}
		}

	// 進行中のリクエストが存在せず、稼働中
	case c.State == ChairStateActive:
		c.World.EmptyChairs.Add(c)
		break

	// 未稼働
	case c.State == ChairStateInactive:
		if c.notificationConn == nil {
			// 先に通知コネクションを繋いでおく
			conn, err := c.Client.ConnectChairNotificationStream(ctx, c, func(event NotificationEvent) {
				if !concurrent.TrySend(c.notificationQueue, event) {
					slog.Error("通知受け取りチャンネルが詰まってる", slog.String("chair_server_id", c.ServerID))
					c.notificationQueue <- event
				}
			})
			if err != nil {
				return WrapCodeError(ErrorCodeFailedToConnectNotificationStream, err)
			}
			c.notificationConn = conn
		}

		err := c.Client.SendActivate(ctx, c)
		if err != nil {
			return WrapCodeError(ErrorCodeFailedToActivate, err)
		}
		defer func() { c.ActivatedAt = time.Now() }()

		// 出勤
		c.Location.PlaceTo(&LocationEntry{
			Coord: c.Location.Initial,
			Time:  ctx.CurrentTime(),
		})
		c.State = ChairStateActive
	}

	if c.Location.Dirty() {
		// 動いた場合に自身の座標をサーバーに送信。成功するまでリトライし続ける
		err := backoff.Retry(func() error {
			res, err := c.Client.SendChairCoordinate(ctx, c)
			if err != nil {
				err = WrapCodeError(ErrorCodeFailedToSendChairCoordinate, err)
				go c.World.PublishEvent(&EventSoftError{Error: err})
				return err
			}
			c.Location.SetServerTime(res.RecordedAt)
			c.Location.ResetDirtyFlag()
			return nil
		}, backoff.NewExponentialBackOff())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Chair) moveToward(target Coordinate) Coordinate {
	return c.Location.Current().MoveToward(target, c.Model.Speed, c.Rand)
}

func (c *Chair) moveOppositeTo(target Coordinate) (to Coordinate) {
	current := c.Location.Current()
	to = current

	moveX := 0
	moveY := 0
	switch {
	case target.X == current.X:
		moveY = c.Model.Speed
	case target.Y == current.Y:
		moveX = c.Model.Speed
	default:
		if c.Rand.IntN(2) == 0 {
			moveX = c.Model.Speed
		} else {
			moveY = c.Model.Speed
		}
	}

	switch {
	case current.X < target.X:
		to.X -= moveX
	case current.X > target.X:
		to.X += moveX
	}

	switch {
	case current.Y < target.Y:
		to.Y -= moveY
	case current.Y > target.Y:
		to.Y += moveY
	}

	return to
}

func (c *Chair) moveRandom() (to Coordinate) {
	prev := c.Location.Current()

	// 移動量の決定
	x := c.Rand.IntN(c.Model.Speed + 1)
	y := c.Model.Speed - x

	// 移動方向の決定
	left, right := c.Region.RangeX()
	bottom, top := c.Region.RangeY()

	switch c.Rand.IntN(4) {
	case 0:
		x *= -1
		if prev.X+x < left {
			x *= -1 // 逆側に戻す
		}
		if top < prev.Y+y {
			y *= -1 // 逆側に戻す
		}
	case 1:
		y *= -1
		if right < prev.X+x {
			x *= -1 // 逆側に戻す
		}
		if prev.Y+y < bottom {
			y *= -1 // 逆側に戻す
		}
	case 2:
		x *= -1
		y *= -1
		if prev.X+x < left {
			x *= -1 // 逆側に戻す
		}
		if prev.Y+y < bottom {
			y *= -1 // 逆側に戻す
		}
	case 3:
		if right < prev.X+x {
			x *= -1 // 逆側に戻す
		}
		if top < prev.Y+y {
			y *= -1 // 逆側に戻す
		}
		break
	}

	return C(prev.X+x, prev.Y+y)
}

func (c *Chair) HandleNotification(event NotificationEvent) error {
	switch data := event.(type) {
	case *ChairNotificationEventMatched:
		if c.matchingData != nil && c.matchingData.ServerRequestID != data.ServerRequestID {
			// 椅子が別のリクエストを保持している
			return WrapCodeError(ErrorCodeChairAlreadyHasRequest, fmt.Errorf("chair_id: %s, current_ride_id: %s, got: %s", c.ServerID, c.matchingData.ServerRequestID, data.ServerRequestID))
		}
		c.World.EmptyChairs.Delete(c)
		c.matchingData = data

	case *ChairNotificationEventDispatching:
		if err := c.ValidateChairNotificationEvent(data.ServerRequestID, data.ChairNotificationEvent); err != nil {
			return WrapCodeError(ErrorCodeChairReceivedDataIsWrong, err)
		}

	case *ChairNotificationEventDispatched:
		if err := c.ValidateChairNotificationEvent(data.ServerRequestID, data.ChairNotificationEvent); err != nil {
			return WrapCodeError(ErrorCodeChairReceivedDataIsWrong, err)
		}

	case *ChairNotificationEventCarrying:
		if err := c.ValidateChairNotificationEvent(data.ServerRequestID, data.ChairNotificationEvent); err != nil {
			return WrapCodeError(ErrorCodeChairReceivedDataIsWrong, err)
		}

	case *ChairNotificationEventArrived:
		if err := c.ValidateChairNotificationEvent(data.ServerRequestID, data.ChairNotificationEvent); err != nil {
			return WrapCodeError(ErrorCodeChairReceivedDataIsWrong, err)
		}

	case *ChairNotificationEventCompleted:
		request := c.Request
		if request == nil {
			// 履歴を見て、過去扱っていたRequestに向けてのCOMPLETED通知であれば無視する
			for _, r := range c.RequestHistory.BackwardIter() {
				if r.ServerID == data.ServerRequestID && r.Statuses.Desired == RequestStatusCompleted {
					r.Statuses.Chair = RequestStatusCompleted
					return nil
				}
			}
			return WrapCodeError(ErrorCodeChairNotAssignedButStatusChanged, fmt.Errorf("ride_id: %s, got: %v", data.ServerRequestID, RequestStatusCompleted))
		}

		c.Request.Statuses.RLock()
		defer c.Request.Statuses.RUnlock()
		if request.Statuses.Desired != RequestStatusCompleted {
			return WrapCodeError(ErrorCodeUnexpectedChairRequestStatusTransitionOccurred, fmt.Errorf("ride_id: %s, expect: %v, got: %v", request.ServerID, request.Statuses.Desired, RequestStatusCompleted))
		}
		if err := c.ValidateChairNotificationEvent(data.ServerRequestID, data.ChairNotificationEvent); err != nil {
			return WrapCodeError(ErrorCodeChairReceivedDataIsWrong, err)
		}

		request.Statuses.Chair = RequestStatusCompleted

		// 進行中のリクエストが無い状態にする
		c.Request = nil
		c.matchingData = nil
		c.World.EmptyChairs.Add(c)
	}
	return nil
}

func (c *Chair) ValidateChairNotificationEvent(rideID string, event ChairNotificationEvent) error {
	if c.matchingData == nil {
		return fmt.Errorf("進行中のライドがないときに進行中状態の通知が届きました (ride_id: %s)", rideID)
	}
	if event.User.ID != c.matchingData.User.ID {
		return fmt.Errorf("ユーザーのIDが一致しません。(ride_id: %s, got: %s, want: %s", rideID, event.User.ID, c.matchingData.User.ID)
	}
	if event.User.Name != c.matchingData.User.Name {
		return fmt.Errorf("ユーザーの名前が一致しません。(ride_id: %s, got: %s, want: %s)", rideID, event.User.Name, c.matchingData.User.Name)
	}
	if !event.Pickup.Equals(c.matchingData.Pickup) {
		return fmt.Errorf("配車位置が一致しません。(ride_id: %s, got: %s, want: %s)", rideID, event.Pickup, c.matchingData.Pickup)
	}
	if !event.Destination.Equals(c.matchingData.Destination) {
		return fmt.Errorf("目的地が一致しません。(ride_id: %s, got: %s, want: %s)", rideID, event.Destination, c.matchingData.Destination)
	}
	return nil
}
