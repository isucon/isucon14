package world

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"slices"
	"sync"
	"time"

	"github.com/isucon/isucon14/bench/benchrun"
	"github.com/isucon/isucon14/bench/internal/concurrent"
	"github.com/samber/lo"
)

type UserState int

const (
	UserStateInactive UserState = iota
	UserStatePaymentMethodsNotRegister
	UserStateActive
)

type UserID int

type User struct {
	// ID ベンチマーカー内部ユーザーID
	ID UserID
	// ServerID サーバー上でのユーザーID
	ServerID string
	// World Worldへの逆参照
	World *World
	// Region ユーザーが居る地域
	Region *Region
	// State ユーザーの状態
	State UserState
	// Request 進行中の配椅子・送迎リクエスト
	Request *Request
	// RegisteredData サーバーに登録されているユーザー情報
	RegisteredData RegisteredUserData
	// PaymentToken 支払いトークン
	PaymentToken string
	// RequestHistory リクエスト履歴
	RequestHistory []*Request
	// TotalEvaluation 完了したリクエストの平均評価
	TotalEvaluation int
	// Client webappへのクライアント
	Client UserClient
	// Rand 専用の乱数
	Rand *rand.Rand
	// Invited 招待されたユーザーか
	Invited bool
	// InvitingLock 招待ロック
	InvitingLock sync.Mutex
	// InvCodeUsedCount 招待コードの使用回数
	InvCodeUsedCount int
	// UnusedInvCoupons 未使用の招待クーポンの数
	UnusedInvCoupons int
	// tickDone 行動が完了しているかどうか
	tickDone tickDone
	// notificationConn 通知ストリームコネクション
	notificationConn NotificationStream
	// notificationQueue 通知キュー。毎Tickで最初に処理される
	notificationQueue chan NotificationEvent
	// validatedRideNotificationEvent 最新のバリデーション済みの通知イベント情報
	validatedRideNotificationEvent *UserNotificationEvent
}

type RegisteredUserData struct {
	UserName       string
	FirstName      string
	LastName       string
	DateOfBirth    string
	InvitationCode string
}

func (u *User) String() string {
	if u.Request != nil {
		return fmt.Sprintf("User{id=%d,totalReqs=%d,reqId=%d}", u.ID, len(u.RequestHistory), u.Request.ID)
	}
	return fmt.Sprintf("User{id=%d,totalReqs=%d}", u.ID, len(u.RequestHistory))
}

func (u *User) SetID(id UserID) {
	u.ID = id
}

func (u *User) GetServerID() string {
	return u.ServerID
}

func (u *User) Tick(ctx *Context) error {
	if u.tickDone.DoOrSkip() {
		return nil
	}
	defer u.tickDone.Done()

	// 通知キューを順番に処理する
	for event := range concurrent.TryIter(u.notificationQueue) {
		err := u.HandleNotification(event)
		if err != nil {
			return err
		}
	}

	switch {
	// 支払いトークンが未登録
	case u.State == UserStatePaymentMethodsNotRegister:
		err := u.Client.BrowserAccess(ctx, benchrun.FRONTEND_PATH_SCENARIO_CLIENT_REGISTER_3)
		if err != nil {
			return WrapCodeError(ErrorCodeFailedToRegisterPaymentMethods, err)
		}

		// トークン登録を試みる
		err = u.Client.RegisterPaymentMethods(ctx, u)
		if err != nil {
			return WrapCodeError(ErrorCodeFailedToRegisterPaymentMethods, err)
		}

		// 成功したのでアクティブ状態にする
		u.State = UserStateActive

	// 進行中のリクエストが存在
	case u.Request != nil:
		if u.notificationConn == nil {
			// 通知コネクションが無い場合は繋いでおく
			conn, err := u.Client.ConnectUserNotificationStream(ctx, u, func(event NotificationEvent) {
				if !concurrent.TrySend(u.notificationQueue, event) {
					slog.Error("通知受け取りチャンネルが詰まってる", slog.String("user_server_id", u.ServerID))
					u.notificationQueue <- event
				}
			})
			if err != nil {
				return WrapCodeError(ErrorCodeFailedToConnectNotificationStream, err)
			}
			u.notificationConn = conn
		}

		switch u.Request.Statuses.User {
		case RequestStatusMatching:
			// マッチングされるまで待機する
			// 30秒待ってもマッチされない場合は、サービスとして重大な問題があるのでクリティカルエラーとして落とす
			if time.Now().Sub(u.Request.BenchRequestedAt) >= 30*time.Second {
				return CodeError(ErrorCodeMatchingTimeout)
			}

		case RequestStatusDispatching:
			// 椅子が到着するまで待つ
			// 一向に到着しない場合は、このユーザーの行動はハングする
			break

		case RequestStatusDispatched:
			// 椅子が出発するのを待つ
			// 一向に到着しない場合は、このユーザーの行動はハングする
			break

		case RequestStatusCarrying:
			// 椅子が到着するのを待つ
			// 一向に到着しない場合は、このユーザーの行動はハングする
			break

		case RequestStatusArrived:
			// 送迎の評価及び支払いがまだの場合は行う
			if !u.Request.Evaluated.Load() {
				score := u.Request.CalculateEvaluation().Score()

				err := u.Client.BrowserAccess(ctx, benchrun.FRONTEND_PATH_SCENARIO_CLIENT_EVALUATION)
				if err != nil {
					return WrapCodeError(ErrorCodeFailedToEvaluate, err)
				}

				u.Request.Statuses.Lock()
				res, err := u.Client.SendEvaluation(ctx, u.Request, score)
				if err != nil {
					u.Request.Statuses.Unlock()
					if errors.Is(err, context.DeadlineExceeded) {
						return WrapCodeError(ErrorCodeEvaluateTimeout, err)
					}
					return WrapCodeError(ErrorCodeFailedToEvaluate, err)
				}

				// サーバーが評価を受理したので完了状態になるのを待機する
				u.Request.CompletedAt = ctx.CurrentTime()
				u.Request.ServerCompletedAt = res.CompletedAt
				u.Request.Statuses.Desired = RequestStatusCompleted
				u.Request.Evaluated.Store(true)
				if requests := len(u.RequestHistory); requests == 1 {
					u.Region.TotalEvaluation.Add(int32(score))
				} else {
					u.Region.TotalEvaluation.Add(int32((u.TotalEvaluation+score)/requests - u.TotalEvaluation/(requests-1)))
				}
				u.TotalEvaluation += score
				u.Request.Chair.Owner.CompletedRequest.Append(u.Request)
				u.Request.Chair.Owner.TotalSales.Add(int64(u.Request.Sales()))
				u.Request.Chair.Owner.SubScore.Add(int64(u.Request.Score()))
				u.World.PublishEvent(&EventRequestCompleted{Request: u.Request})

				u.Request.Statuses.Unlock()
			}

		case RequestStatusCompleted:
			// 進行中のリクエストが無い状態にする
			u.Request = nil

			// 通知コネクションを切る
			if u.notificationConn != nil {
				u.notificationConn.Close()
				u.notificationConn = nil
			}
		}

	// 進行中のリクエストは存在しないが、ユーザーがアクティブ状態
	case u.Request == nil && u.State == UserStateActive:
		if count := len(u.RequestHistory); (count == 1 && u.TotalEvaluation <= 1) || float64(u.TotalEvaluation)/float64(count) <= 2 {
			// 初回利用で評価1なら離脱
			// 2回以上利用して平均評価が2以下の場合は離脱
			if u.Region.UserLeave(u) {
				break
			}
			// Region内の最低ユーザー数を下回るならそのまま残る
		}

		// 過去のリクエストを確認する
		err := u.CheckRequestHistory(ctx)
		if err != nil {
			return WrapCodeError(ErrorCodeFailedToCheckRequestHistory, err)
		}

		// リクエストを作成する
		// TODO 作成する条件・頻度
		err = u.CreateRequest(ctx)
		if err != nil {
			return err
		}

	// 離脱ユーザーは何もしない
	case u.State == UserStateInactive:
		break
	}
	return nil
}

func (u *User) Deactivate() {
	u.State = UserStateInactive
	if u.notificationConn != nil {
		u.notificationConn.Close()
		u.notificationConn = nil
	}
	u.World.PublishEvent(&EventUserLeave{User: u})
}

func (u *User) CheckRequestHistory(ctx *Context) error {
	err := u.Client.BrowserAccess(ctx, benchrun.FRONTEND_PATH_SCENARIO_CLIENT_CHECK_HISTORY_1)
	if err != nil {
		return WrapCodeError(ErrorCodeFailedToCheckRequestHistory, err)
	}
	err = u.Client.BrowserAccess(ctx, benchrun.FRONTEND_PATH_SCENARIO_CLIENT_CHECK_HISTORY_2)
	if err != nil {
		return WrapCodeError(ErrorCodeFailedToCheckRequestHistory, err)
	}

	res, err := u.Client.GetRequests(ctx)
	if err != nil {
		return err
	}
	if len(res.Requests) != len(u.RequestHistory) {
		return fmt.Errorf("ライドの数が想定数と一致していません: expected=%d, got=%d", len(u.RequestHistory), len(res.Requests))
	}

	historyMap := lo.KeyBy(u.RequestHistory, func(r *Request) string { return r.ServerID })
	for _, req := range res.Requests {
		expected, ok := historyMap[req.ID]
		if !ok {
			return fmt.Errorf("想定されないライドが含まれています: id=%s", req.ID)
		}
		if !req.DestinationCoordinate.Equals(expected.DestinationPoint) || !req.PickupCoordinate.Equals(expected.PickupPoint) {
			return fmt.Errorf("ライドの座標情報が期待したものと異なります: id=%s", req.ID)
		}
		if req.Fare != expected.Fare() {
			return fmt.Errorf("ライドの運賃が期待したものと異なります: id=%s", req.ID)
		}
		if req.Evaluation != expected.CalculateEvaluation().Score() {
			return fmt.Errorf("ライドの評価が期待したものと異なります: id=%s", req.ID)
		}
		if req.Chair.ID != expected.Chair.ServerID || req.Chair.Name != expected.Chair.RegisteredData.Name || req.Chair.Model != expected.Chair.Model.Name || req.Chair.Owner != expected.Chair.Owner.RegisteredData.Name {
			return fmt.Errorf("ライドの椅子の情報が期待したものと異なります: id=%s", req.ID)
		}
		if !req.CompletedAt.Equal(expected.ServerCompletedAt) {
			return fmt.Errorf("ライドの完了日時が期待したものと異なります: id=%s", req.ID)
		}
	}

	return nil
}

func (u *User) CreateRequest(ctx *Context) error {
	if u.Request != nil {
		panic("ユーザーに進行中のリクエストがあるのにも関わらず、リクエストを新規作成しようとしている")
	}

	u.InvitingLock.Lock()
	defer u.InvitingLock.Unlock()

	pickup, dest := RandomTwoCoordinateWithRand(u.Region, u.Rand.IntN(100)+5, u.Rand)

	req := &Request{
		User:             u,
		PickupPoint:      pickup,
		DestinationPoint: dest,
		RequestedAt:      ctx.CurrentTime(),
		Statuses: RequestStatuses{
			Desired: RequestStatusMatching,
			Chair:   RequestStatusMatching,
			User:    RequestStatusMatching,
		},
	}

	useInvCoupon := false
	switch {
	// 初回利用の割引を適用
	case len(u.RequestHistory) == 0:
		req.Discount = 3000

	// 招待された側のクーポンを適用
	case len(u.RequestHistory) == 1 && u.Invited:
		req.Discount = 1500

	// 招待した側のクーポンを適用
	case u.UnusedInvCoupons > 0:
		req.Discount = 1000
		useInvCoupon = true
	}

	checkDistance := 50
	now := time.Now()
	nearby, err := u.Client.GetNearbyChairs(ctx, pickup, checkDistance)
	if err != nil {
		return WrapCodeError(ErrorCodeWrongNearbyChairs, err)
	}
	if err := u.World.checkNearbyChairsResponse(now, pickup, checkDistance, nearby); err != nil {
		return WrapCodeError(ErrorCodeWrongNearbyChairs, err)
	}
	if len(nearby.Chairs) == 0 {
		// 近くに椅子が無いので配車をやめる
		return nil
	}

	estimation, err := u.Client.GetEstimatedFare(ctx, pickup, dest)
	if err != nil {
		return WrapCodeError(ErrorCodeFailedToCreateRequest, err)
	}
	if req.ActualDiscount() != estimation.Discount || req.Fare() != estimation.Fare {
		return WrapCodeError(ErrorCodeFailedToCreateRequest, errors.New("ライド料金の見積もり金額が誤っています"))
	}

	res, err := u.Client.SendCreateRequest(ctx, req)
	if err != nil {
		return WrapCodeError(ErrorCodeFailedToCreateRequest, err)
	}
	req.ServerID = res.ServerRequestID
	req.BenchRequestedAt = time.Now()
	u.Request = req
	u.RequestHistory = append(u.RequestHistory, req)
	u.World.RequestDB.Create(req)
	if useInvCoupon {
		u.UnusedInvCoupons--
	}
	return nil
}

func (u *User) ChangeRequestStatus(status RequestStatus, serverRequestID string, validator func() error) error {
	request := u.Request
	if request == nil {
		if status == RequestStatusCompleted {
			// 履歴を見て、過去扱っていたRequestに向けてのCOMPLETED通知であれば無視する
			for _, r := range slices.Backward(u.RequestHistory) {
				if r.ServerID == serverRequestID && r.Statuses.Desired == RequestStatusCompleted {
					r.Statuses.User = RequestStatusCompleted
					return nil
				}
			}
		}
		return WrapCodeError(ErrorCodeUserNotRequestingButStatusChanged, fmt.Errorf("user_id: %s, got: %v", u.ServerID, status))
	}
	request.Statuses.RLock()
	defer request.Statuses.RUnlock()
	if request.Statuses.User != status && request.Statuses.Desired != status {
		// 現在認識しているユーザーの状態で無いかつ、想定状態ではない状態に遷移しようとしている場合
		if request.Statuses.User == RequestStatusMatching && request.Statuses.Desired == RequestStatusDispatched && status == RequestStatusDispatching {
			// ユーザーにDispatchingが送られる前に、椅子が到着している場合があるが、その時にDispatchingを受け取ることを許容する
		} else if request.Statuses.User == RequestStatusDispatched && request.Statuses.Desired == RequestStatusArrived && status == RequestStatusCarrying {
			// もう到着しているが、ユーザー側の通知が遅延していて、DISPATCHED状態からまだCARRYINGに遷移してないときは、CARRYINGを許容する
		} else if request.Statuses.Desired == RequestStatusDispatched && request.Statuses.User == RequestStatusDispatched && status == RequestStatusCarrying {
			// ユーザーがDispatchedを受け取った状態で、椅子が出発リクエストを送った後、ベンチマーカーのDesiredステータスの変更を行う前にユーザー側にCarrying通知が届いてしまうことがあるがこれは許容する
		} else if status == RequestStatusCompleted {
			// 履歴を見て、過去扱っていたRequestに向けてのCOMPLETED通知であれば無視する
			for _, r := range slices.Backward(u.RequestHistory) {
				if r.ServerID == serverRequestID && r.Statuses.Desired == RequestStatusCompleted {
					r.Statuses.User = RequestStatusCompleted
					return nil
				}
			}
			return WrapCodeError(ErrorCodeUnexpectedUserRequestStatusTransitionOccurred, fmt.Errorf("ride_id: %v, expect: %v, got: %v (current: %v)", request.ServerID, request.Statuses.Desired, status, request.Statuses.User))
		} else {
			return WrapCodeError(ErrorCodeUnexpectedUserRequestStatusTransitionOccurred, fmt.Errorf("ride_id: %v, expect: %v, got: %v (current: %v)", request.ServerID, request.Statuses.Desired, status, request.Statuses.User))
		}
	}

	if validator != nil {
		if err := validator(); err != nil {
			return WrapCodeError(ErrorCodeUserReceivedDataIsWrong, err)
		}
	}

	request.Statuses.User = status
	return nil
}

func (u *User) HandleNotification(event NotificationEvent) error {
	switch data := event.(type) {
	case *UserNotificationEventMatching:
		err := u.ChangeRequestStatus(RequestStatusMatching, data.ServerRequestID, func() error {
			return u.ValidateNotificationEvent(data.ServerRequestID, data.UserNotificationEvent, true)
		})
		if err != nil {
			return err
		}
	case *UserNotificationEventDispatching:
		err := u.ChangeRequestStatus(RequestStatusDispatching, data.ServerRequestID, func() error {
			if err := u.ValidateNotificationEvent(data.ServerRequestID, data.UserNotificationEvent, false); err != nil {
				return err
			}
			u.validatedRideNotificationEvent = &data.UserNotificationEvent
			return nil
		})
		if err != nil {
			return err
		}
	case *UserNotificationEventDispatched:
		err := u.ChangeRequestStatus(RequestStatusDispatched, data.ServerRequestID, func() error {
			if u.validatedRideNotificationEvent != nil {
				return compareUserNotificationEvent(data.ServerRequestID, *u.validatedRideNotificationEvent, data.UserNotificationEvent)
			}
			if err := u.ValidateNotificationEvent(data.ServerRequestID, data.UserNotificationEvent, false); err != nil {
				return err
			}
			u.validatedRideNotificationEvent = &data.UserNotificationEvent
			return nil
		})
		if err != nil {
			return err
		}
	case *UserNotificationEventCarrying:
		err := u.ChangeRequestStatus(RequestStatusCarrying, data.ServerRequestID, func() error {
			if u.validatedRideNotificationEvent != nil {
				return compareUserNotificationEvent(data.ServerRequestID, *u.validatedRideNotificationEvent, data.UserNotificationEvent)
			}
			if err := u.ValidateNotificationEvent(data.ServerRequestID, data.UserNotificationEvent, false); err != nil {
				return err
			}
			u.validatedRideNotificationEvent = &data.UserNotificationEvent
			return nil
		})
		if err != nil {
			return err
		}
	case *UserNotificationEventArrived:
		err := u.ChangeRequestStatus(RequestStatusArrived, data.ServerRequestID, func() error {
			if u.validatedRideNotificationEvent != nil {
				return compareUserNotificationEvent(data.ServerRequestID, *u.validatedRideNotificationEvent, data.UserNotificationEvent)
			}
			if err := u.ValidateNotificationEvent(data.ServerRequestID, data.UserNotificationEvent, false); err != nil {
				return err
			}
			u.validatedRideNotificationEvent = &data.UserNotificationEvent
			return nil
		})
		if err != nil {
			return err
		}
	case *UserNotificationEventCompleted:
		err := u.ChangeRequestStatus(RequestStatusCompleted, data.ServerRequestID, func() error {
			return u.ValidateNotificationEvent(data.ServerRequestID, data.UserNotificationEvent, false)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) ValidateNotificationEvent(rideID string, serverSide UserNotificationEvent, ignoreChair bool) error {
	if !serverSide.Pickup.Equals(u.Request.PickupPoint) {
		return fmt.Errorf("配車位置が一致しません。(ride_id: %s, got: %s, want: %s)", rideID, serverSide.Pickup, u.Request.PickupPoint)
	}
	if !serverSide.Destination.Equals(u.Request.DestinationPoint) {
		return fmt.Errorf("目的地が一致しません。(ride_id: %s, got: %s, want: %s)", rideID, serverSide.Destination, u.Request.DestinationPoint)
	}

	if serverSide.Fare != u.Request.Fare() {
		return fmt.Errorf("運賃が一致しません。(ride_id: %s, got: %d, want: %d)", rideID, serverSide.Fare, u.Request.Fare())
	}

	if ignoreChair {
		return nil
	}

	if serverSide.Chair == nil {
		return fmt.Errorf("椅子情報がありません。(ride_id: %s)", rideID)
	}

	serverSideChair := serverSide.Chair
	chair := u.World.ChairDB.GetByServerID(serverSideChair.ID)
	if chair == nil {
		return fmt.Errorf("想定していない椅子が返却されました。(ride_id: %s, chair_id: %s)", rideID, serverSide.Chair.ID)
	}

	if serverSideChair.Name != chair.RegisteredData.Name {
		return fmt.Errorf("椅子の名前が一致しません。(ride_id: %s, chair_id: %s, got: %s, want: %s)", rideID, serverSide.Chair.ID, serverSide.Chair.Name, u.Request.Chair.RegisteredData.Name)
	}
	if serverSideChair.Model != chair.Model.Name {
		return fmt.Errorf("椅子のモデルが一致しません。(ride_id: %s, chair_id: %s, got: %s, want: %s)", rideID, serverSide.Chair.ID, serverSide.Chair.Model, u.Request.Chair.Model.Name)
	}

	totalRideCount := 0
	totalEvaluation := 0
	for _, r := range chair.RequestHistory.Iter() {
		if r.Evaluated.Load() {
			totalRideCount++
			totalEvaluation += r.CalculateEvaluation().Score()
		}
	}

	if serverSideChair.Stats.TotalRidesCount != totalRideCount {
		return fmt.Errorf("椅子の総乗車回数が一致しません。(ride_id: %s, chair_id: %s, got: %d, want: %d)", rideID, serverSide.Chair.ID, serverSide.Chair.Stats.TotalRidesCount, totalRideCount)
	}
	if totalRideCount > 0 {
		if !almostEqual(serverSideChair.Stats.TotalEvaluationAvg, float64(totalEvaluation)/float64(totalRideCount), 0.01) {
			return fmt.Errorf("椅子の評価の平均が一致しません。(ride_id: %s, chair_id: %s, got: %f, want: %f)", rideID, serverSide.Chair.ID, serverSide.Chair.Stats.TotalEvaluationAvg, float64(totalEvaluation)/float64(totalRideCount))
		}
	} else {
		if serverSideChair.Stats.TotalEvaluationAvg != 0 {
			return fmt.Errorf("椅子の評価の平均が一致しません。(ride_id: %s, chair_id: %s, got: %f, want: %f)", rideID, serverSide.Chair.ID, serverSide.Chair.Stats.TotalEvaluationAvg, 0.0)
		}
	}

	return nil
}

// compareUserNotificationEvent validation済みのUserNotificationEventと比較して、一致しない場合はエラーを返す
func compareUserNotificationEvent(rideID string, old, new UserNotificationEvent) error {
	if !new.Pickup.Equals(old.Pickup) {
		return fmt.Errorf("配車位置が一致しません。(ride_id: %s, got: %s, want: %s)", rideID, new.Pickup, old.Pickup)
	}
	if !new.Destination.Equals(old.Destination) {
		return fmt.Errorf("目的地が一致しません。(ride_id: %s, got: %s, want: %s)", rideID, new.Destination, old.Destination)
	}

	if new.Fare != old.Fare {
		return fmt.Errorf("運賃が一致しません。(ride_id: %s, got: %d, want: %d)", rideID, new.Fare, old.Fare)
	}

	if new.Chair == nil {
		return fmt.Errorf("椅子情報がありません。(ride_id: %s)", rideID)
	}

	if new.Chair.ID != old.Chair.ID {
		return fmt.Errorf("椅子のIDが一致しません。(ride_id: %s, got: %s, want: %s)", rideID, new.Chair.ID, old.Chair.ID)
	}
	if new.Chair.Name != old.Chair.Name {
		return fmt.Errorf("椅子の名前が一致しません。(ride_id: %s, chair_id: %s, got: %s, want: %s)", rideID, new.Chair.ID, new.Chair.Name, old.Chair.Name)
	}
	if new.Chair.Model != old.Chair.Model {
		return fmt.Errorf("椅子のモデルが一致しません。(ride_id: %s, chair_id: %s, got: %s, want: %s)", rideID, new.Chair.ID, new.Chair.Model, old.Chair.Model)
	}

	if new.Chair.Stats.TotalRidesCount != old.Chair.Stats.TotalRidesCount {
		return fmt.Errorf("椅子の総乗車回数が一致しません。(ride_id: %s, chair_id: %s, got: %d, want: %d)", rideID, new.Chair.ID, new.Chair.Stats.TotalRidesCount, old.Chair.Stats.TotalRidesCount)
	}
	if !almostEqual(new.Chair.Stats.TotalEvaluationAvg, old.Chair.Stats.TotalEvaluationAvg, 0.01) {
		return fmt.Errorf("椅子の評価の平均が一致しません。(ride_id: %s, chair_id: %s, got: %f, want: %f)", rideID, new.Chair.ID, new.Chair.Stats.TotalEvaluationAvg, old.Chair.Stats.TotalEvaluationAvg)
	}

	return nil
}

func almostEqual(a, b, epsilon float64) bool {
	// 絶対誤差と相対誤差を組み合わせて比較
	absDiff := math.Abs(a - b)
	if absDiff <= epsilon {
		return true
	}
	// 相対誤差の比較
	maxAbs := math.Max(math.Abs(a), math.Abs(b))
	return absDiff <= epsilon*maxAbs
}
