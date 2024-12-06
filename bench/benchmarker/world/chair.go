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
	matchingData *ChairNotificationEventMatching
	// Ride 進行中のライド
	Ride *Ride
	// RideHistory 引き受けたライドの履歴
	RideHistory *concurrent.SimpleSlice[*Ride]
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

	// detour 今回のライドで迂回するかどうか
	detour bool
	// detoured 今回のライドで迂回したかどうか
	detoured bool
	// detourPoint 迂回ポイント
	detourPoint Coordinate
	// detourIn Dispatching or Carryingのどっちで迂回するか
	detourIn RideStatus
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
	case c.Ride != nil:
		switch c.Ride.Statuses.Chair {
		case RideStatusMatching:
			// Active状態なら配車要求をACKする
			// そうでないなら、応答せずにハングさせる
			if c.State == ChairStateActive {
				c.Ride.Statuses.Lock()

				err := c.Client.SendAcceptRequest(ctx, c, c.Ride)
				if err != nil {
					c.Ride.Statuses.Unlock()

					return WrapCodeError(ErrorCodeFailedToAcceptRide, err)
				}

				// サーバーに要求を受理の通知が通ったので配椅子地に向かう
				c.Ride.Chair = c
				c.Ride.Statuses.Desired = RideStatusEnRoute
				c.Ride.Statuses.Chair = RideStatusEnRoute
				c.Ride.StartPoint = null.ValueFrom(c.Location.Current())
				c.Ride.MatchedAt = ctx.CurrentTime()
				c.Ride.BenchMatchedAt = time.Now()

				c.Ride.Statuses.Unlock()

				c.RideHistory.Append(c.Ride)
				if !c.Ride.User.Region.Contains(c.Location.Current()) {
					ctx.ContestantLogger().Warn("ユーザーのリージョン外部に存在する椅子がマッチングされました", slog.Int("distance", c.Ride.PickupPoint.DistanceTo(c.Location.Current())))
				}
			}

		case RideStatusEnRoute:
			// 配車位置に向かう
			time := ctx.CurrentTime()
			if c.detour && c.detourIn == RideStatusEnRoute && !c.detoured {
				// 迂回する予定でまだ迂回してない場合
				if c.Location.Current().Equals(c.detourPoint) {
					// 迂回ポイントに着いた次の移動は配車位置から離れる方向に行う
					c.Location.MoveTo(&LocationEntry{
						Coord: c.moveOppositeTo(c.Ride.PickupPoint),
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
					Coord: c.moveToward(c.Ride.PickupPoint),
					Time:  time,
				})
			}

			if c.Location.Current().Equals(c.Ride.PickupPoint) {
				// 配車位置に到着
				c.Ride.Statuses.Desired = RideStatusPickup
				c.Ride.Statuses.Chair = RideStatusPickup
				c.Ride.DispatchedAt = time
			}

		case RideStatusPickup:
			// 乗客を乗せて出発しようとする
			if c.Ride.Statuses.User != RideStatusPickup {
				// ただし、ユーザーに到着通知が行っていないとユーザーは乗らない振る舞いをするので
				// ユーザー側の状態が変わるまで待機する
				// 一向にユーザーの状態が変わらない場合は、この椅子の行動はハングする
				break
			}

			err := c.Client.SendDepart(ctx, c.Ride)
			if err != nil {
				return WrapCodeError(ErrorCodeFailedToDepart, err)
			}

			// サーバーがdepartを受理したので出発する
			c.Ride.Statuses.Desired = RideStatusCarrying
			c.Ride.Statuses.Chair = RideStatusCarrying
			c.Ride.PickedUpAt = ctx.CurrentTime()

		case RideStatusCarrying:
			// 目的地に向かう
			time := ctx.CurrentTime()
			if c.detour && c.detourIn == RideStatusCarrying && !c.detoured {
				// 迂回する予定でまだ迂回してない場合
				if c.Location.Current().Equals(c.detourPoint) {
					// 迂回ポイントに着いた次の移動は目的地から離れる方向に行う
					c.Location.MoveTo(&LocationEntry{
						Coord: c.moveOppositeTo(c.Ride.DestinationPoint),
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
					Coord: c.moveToward(c.Ride.DestinationPoint),
					Time:  time,
				})
			}

			if c.Location.Current().Equals(c.Ride.DestinationPoint) {
				// 目的地に到着
				c.Ride.Statuses.Desired = RideStatusArrived
				c.Ride.Statuses.Chair = RideStatusArrived
				c.Ride.ArrivedAt = time
				break
			}

		case RideStatusArrived:
			// 客が評価するまで待機する
			// 一向に評価されない場合は、この椅子の行動はハングする
			break

		case RideStatusCompleted:
			// c.HandleNotificationでCompletedを受け取った際に c.Ride = nil にしている
			slog.Warn("unexpected state")
			break
		}

	// アサインされたリクエストが存在するが、詳細を未取得
	case c.Ride == nil && c.matchingData != nil:
		req := c.World.RideDB.GetByServerID(c.matchingData.ServerRideID)
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
		c.Ride = req
		// 10%の確率で迂回させる(最短距離より1単位速度分だけ遠回しさせる)
		c.detour = c.Rand.Float64() < 0.1
		c.detoured = false
		if c.detour {
			if c.Rand.IntN(2) == 0 {
				c.detourIn = RideStatusEnRoute
				c.detourPoint = CalculateRandomDetourPoint(c.Location.Current(), c.Ride.PickupPoint, c.Model.Speed, c.Rand)
			} else {
				c.detourIn = RideStatusCarrying
				c.detourPoint = CalculateRandomDetourPoint(c.Ride.PickupPoint, c.Ride.DestinationPoint, c.Model.Speed, c.Rand)
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
			c.Location.SetServerTime(res.RecordedAt) // FIXME: ここの反映(ロック)が遅れて、総移動距離の計算が１つずれる場合がある
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
	case *ChairNotificationEventMatching:
		if c.matchingData != nil && c.matchingData.ServerRideID != data.ServerRideID {
			// 椅子が別のリクエストを保持している
			slog.Debug(fmt.Sprintf("code:%d", ErrorCodeChairAlreadyHasRequest), slog.Any("ride", c.Ride))
			return WrapCodeError(ErrorCodeChairAlreadyHasRequest, fmt.Errorf("chair_id: %s, current_ride_id: %s, got: %s", c.ServerID, c.matchingData.ServerRideID, data.ServerRideID))
		}

		c.World.EmptyChairs.Delete(c)
		c.matchingData = data

	case *ChairNotificationEventCompleted:
		if err := c.ValidateChairNotificationEvent(data.ServerRideID, data.ChairNotificationEvent); err != nil {
			return WrapCodeError(ErrorCodeIncorrectChairNotificationData, err)
		}

		if c.Ride == nil {
			// 履歴を見て、過去扱っていたRequestに向けてのCOMPLETED通知であれば無視する
			for _, r := range c.RideHistory.BackwardIter() {
				if r.ServerID == data.ServerRideID && r.Statuses.Desired == RideStatusCompleted {
					r.Statuses.Chair = RideStatusCompleted
					return nil
				}
			}
			return WrapCodeError(ErrorCodeChairNotAssignedButStatusChanged, fmt.Errorf("ride_id: %s, got: %v", data.ServerRideID, RideStatusCompleted))
		}

		c.Ride.Statuses.RLock()
		defer c.Ride.Statuses.RUnlock()
		if c.Ride.Statuses.Desired != RideStatusCompleted {
			return WrapCodeError(ErrorCodeUnexpectedChairRequestStatusTransitionOccurred, fmt.Errorf("ride_id: %s, expect: %v, got: %v", c.Ride.ServerID, c.Ride.Statuses.Desired, RideStatusCompleted))
		}
		c.Ride.Statuses.Chair = RideStatusCompleted

		// 進行中のリクエストが無い状態にする
		c.Ride = nil
		c.matchingData = nil
		c.World.EmptyChairs.Add(c)
	}

	return nil
}

func (c *Chair) ValidateChairNotificationEvent(rideID string, event ChairNotificationEvent) error {
	if event.User.ID == c.matchingData.User.ID {
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
