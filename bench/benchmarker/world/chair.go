package world

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"slices"

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
	// Provider 椅子を所有している事業者
	Provider *Provider
	// Model 椅子のモデル
	Model *ChairModel
	// State 椅子の状態
	State ChairState
	// Location 椅子の位置情報
	Location ChairLocation
	// RegisteredData サーバーに登録されている椅子情報
	RegisteredData RegisteredChairData
	// ServerRequestID 進行中のリクエストのサーバー上でのID
	ServerRequestID null.String
	// Request 進行中のリクエスト
	Request *Request
	// RequestHistory 引き受けたリクエストの履歴
	RequestHistory []*Request
	// Client webappへのクライアント
	Client ChairClient
	// Rand 専用の乱数
	Rand *rand.Rand
	// tickDone 行動が完了しているかどうか
	tickDone tickDone
	// notificationConn 通知ストリームコネクション
	notificationConn NotificationStream
	// notificationQueue 通知キュー。毎Tickで最初に処理される
	notificationQueue chan NotificationEvent
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
			// 配椅子要求を受理するか、拒否する
			// TODO: 拒否ロジック
			if c.State == ChairStateActive {
				c.Request.Statuses.Lock()

				err := c.Client.SendAcceptRequest(ctx, c, c.Request)
				if err != nil {
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

				c.RequestHistory = append(c.RequestHistory, c.Request)
				if !c.Request.User.Region.Contains(c.Location.Current()) {
					ctx.ContestantLogger().Warn("Userが居るRegionの外部に存在するChairがマッチングされました", slog.Int("distance", c.Request.PickupPoint.DistanceTo(c.Location.Current())))
				}
			} else {
				err := c.Client.SendDenyRequest(ctx, c, c.Request.ServerID)
				if err != nil {
					return WrapCodeError(ErrorCodeFailedToDenyRequest, err)
				}

				// サーバーに要求を拒否の通知が通ったので状態をリセット
				c.Request = nil
				c.ServerRequestID = null.String{}
			}

		case RequestStatusDispatching:
			// 配椅子位置に向かう
			time := ctx.CurrentTime()
			c.Location.MoveTo(&LocationEntry{
				Coord: c.moveToward(c.Request.PickupPoint),
				Time:  time,
			})
			if c.Location.Current().Equals(c.Request.PickupPoint) {
				// 配椅子位置に到着
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
			c.Location.MoveTo(&LocationEntry{
				Coord: c.moveToward(c.Request.DestinationPoint),
				Time:  time,
			})
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
			// 進行中のリクエストが無い状態にする
			c.Request = nil
			c.ServerRequestID = null.String{}
		}

	// オファーされたリクエストが存在するが、詳細を未取得
	case c.Request == nil && c.ServerRequestID.Valid:
		req := c.World.RequestDB.GetByServerID(c.ServerRequestID.String)
		if req == nil {
			// ベンチマーク外で作成されたリクエストがアサインされた場合は処理できないので一律で拒否る
			err := c.Client.SendDenyRequest(ctx, c, c.ServerRequestID.String)
			if err != nil {
				return WrapCodeError(ErrorCodeFailedToDenyRequest, err)
			}

			c.ServerRequestID = null.String{}
		} else {
			// TODO detailレスポンス検証
			_, err := c.Client.GetRequestByChair(ctx, c, c.ServerRequestID.String)
			if err != nil {
				return WrapCodeError(ErrorCodeFailedToGetRequestDetail, err)
			}

			// 椅子がリクエストを正常に認識する
			c.Request = req
		}

	// 進行中のリクエストが存在せず、稼働中
	case c.State == ChairStateActive:
		// TODO: deactivateタイミング
		//err := c.Client.SendDeactivate(ctx, c)
		//if err != nil {
		//	return WrapCodeError(ErrorCodeFailedToDeactivate, err)
		//}
		//
		//// 退勤
		//c.State = ChairStateInactive
		//// 通知コネクションを切断
		//c.notificationConn.Close()
		//c.notificationConn = nil

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

		// 出勤
		c.Location.PlaceTo(&LocationEntry{
			Coord: c.Location.Initial,
			Time:  ctx.CurrentTime(),
		})
		c.State = ChairStateActive
	}

	if c.Location.Dirty() {
		// 動いた場合に自身の座標をサーバーに送信
		res, err := c.Client.SendChairCoordinate(ctx, c)
		if err != nil {
			return WrapCodeError(ErrorCodeFailedToSendChairCoordinate, err)
		}
		c.Location.SetServerTime(res.RecordedAt) // FIXME: ここの反映(ロック)が遅れて、総移動距離の計算が１つずれる場合がある
		c.Location.ResetDirtyFlag()
	}
	return nil
}

func (c *Chair) AssignRequest(serverRequestID string) error {
	if c.ServerRequestID.Valid && c.ServerRequestID.String != serverRequestID {
		// 椅子が別のリクエストを保持している
		return WrapCodeError(ErrorCodeChairAlreadyHasRequest, fmt.Errorf("server chair id: %s, current request: %s (%v), got: %s", c.ServerID, c.ServerRequestID.String, c.Request, serverRequestID))
	}
	c.ServerRequestID = null.StringFrom(serverRequestID)
	return nil
}

func (c *Chair) moveToward(target Coordinate) (to Coordinate) {
	prev := c.Location.Current()
	to = c.Location.Current()

	// ランダムにx, y方向で近づける
	x := c.Rand.IntN(c.Model.Speed + 1)
	y := c.Model.Speed - x
	remain := 0

	switch {
	case prev.X < target.X:
		xDiff := target.X - (prev.X + x)
		if xDiff < 0 {
			// X座標で追い越すので、追い越す分をyの移動に加える
			to.X = target.X
			y += -xDiff
		} else {
			to.X += x
		}
	case prev.X > target.X:
		xDiff := (prev.X - x) - target.X
		if xDiff < 0 {
			// X座標で追い越すので、追い越す分をyの移動に加える
			to.X = target.X
			y += -xDiff
		} else {
			to.X -= x
		}
	default:
		y = c.Model.Speed
	}

	switch {
	case prev.Y < target.Y:
		yDiff := target.Y - (prev.Y + y)
		if yDiff < 0 {
			to.Y = target.Y
			remain += -yDiff
		} else {
			to.Y += y
		}
	case prev.Y > target.Y:
		yDiff := (prev.Y - y) - target.Y
		if yDiff < 0 {
			to.Y = target.Y
			remain += -yDiff
		} else {
			to.Y -= y
		}
	default:
		remain = y
	}

	if remain > 0 {
		x = remain
		switch {
		case to.X < target.X:
			xDiff := target.X - (to.X + x)
			if xDiff < 0 {
				to.X = target.X
			} else {
				to.X += x
			}
		case to.X > target.X:
			xDiff := (to.X - x) - target.X
			if xDiff < 0 {
				to.X = target.X
			} else {
				to.X -= x
			}
		}
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
		err := c.AssignRequest(data.ServerRequestID)
		if err != nil {
			return err
		}

	case *ChairNotificationEventCompleted:
		request := c.Request
		if request == nil {
			// 履歴を見て、過去扱っていたRequestに向けてのCOMPLETED通知であれば無視する
			for _, r := range slices.Backward(c.RequestHistory) {
				if r.ServerID == data.ServerRequestID && r.Statuses.Desired == RequestStatusCompleted {
					r.Statuses.Chair = RequestStatusCompleted
					return nil
				}
			}
			return WrapCodeError(ErrorCodeChairNotAssignedButStatusChanged, fmt.Errorf("request server id: %v, got: %v", data.ServerRequestID, RequestStatusCompleted))
		}
		if request.Statuses.Desired != RequestStatusCompleted {
			return WrapCodeError(ErrorCodeUnexpectedChairRequestStatusTransitionOccurred, fmt.Errorf("request server id: %v, expect: %v, got: %v", request.ServerID, request.Statuses.Desired, RequestStatusCompleted))
		}
		request.Statuses.Chair = RequestStatusCompleted
	}
	return nil
}
