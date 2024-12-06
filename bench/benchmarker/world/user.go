package world

import (
	"errors"
	"fmt"
	"log/slog"
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
	// Ride 進行中のライド
	Ride *Ride
	// RegisteredData サーバーに登録されているユーザー情報
	RegisteredData RegisteredUserData
	// PaymentToken 支払いトークン
	PaymentToken string
	// RideHistory ライド履歴
	RideHistory []*Ride
	// TotalEvaluation 完了したライドの平均評価
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
	if u.Ride != nil {
		return fmt.Sprintf("User{id=%d,totalReqs=%d,reqId=%d}", u.ID, len(u.RideHistory), u.Ride.ID)
	}
	return fmt.Sprintf("User{id=%d,totalReqs=%d}", u.ID, len(u.RideHistory))
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

	// 進行中のライドが存在
	case u.Ride != nil:
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

		switch u.Ride.Statuses.User {
		case RideStatusMatching:
			// マッチングされるまで待機する
			// 一向にマッチングされない場合は、このユーザーの行動はハングする
			break

		case RideStatusEnRoute:
			// 椅子が到着するまで待つ
			// 一向に到着しない場合は、このユーザーの行動はハングする
			break

		case RideStatusPickup:
			// 椅子が出発するのを待つ
			// 一向に到着しない場合は、このユーザーの行動はハングする
			break

		case RideStatusCarrying:
			// 椅子が到着するのを待つ
			// 一向に到着しない場合は、このユーザーの行動はハングする
			break

		case RideStatusArrived:
			// 送迎の評価及び支払いがまだの場合は行う
			if !u.Ride.Evaluated.Load() {
				score := u.Ride.CalculateEvaluation().Score()

				err := u.Client.BrowserAccess(ctx, benchrun.FRONTEND_PATH_SCENARIO_CLIENT_EVALUATION)
				if err != nil {
					return WrapCodeError(ErrorCodeFailedToEvaluate, err)
				}

				u.Ride.Statuses.Lock()
				res, err := u.Client.SendEvaluation(ctx, u.Ride, score)
				if err != nil {
					u.Ride.Statuses.Unlock()
					return WrapCodeError(ErrorCodeFailedToEvaluate, err)
				}

				// サーバーが評価を受理したので完了状態になるのを待機する
				u.Ride.CompletedAt = ctx.CurrentTime()
				u.Ride.ServerCompletedAt = res.CompletedAt
				u.Ride.Statuses.Desired = RideStatusCompleted
				u.Ride.Evaluated.Store(true)
				if requests := len(u.RideHistory); requests == 1 {
					u.Region.TotalEvaluation.Add(int32(score))
				} else {
					u.Region.TotalEvaluation.Add(int32((u.TotalEvaluation+score)/requests - u.TotalEvaluation/(requests-1)))
				}
				u.TotalEvaluation += score
				u.Ride.Chair.Owner.CompletedRequest.Append(u.Ride)
				u.Ride.Chair.Owner.TotalSales.Add(int64(u.Ride.Sales()))
				u.Ride.Chair.Owner.SubScore.Add(int64(u.Ride.Score()))
				u.World.PublishEvent(&EventRequestCompleted{Request: u.Ride})

				u.Ride.Statuses.Unlock()
			}

		case RideStatusCompleted:
			// 進行中のライドが無い状態にする
			u.Ride = nil

			// 通知コネクションを切る
			if u.notificationConn != nil {
				u.notificationConn.Close()
				u.notificationConn = nil
			}
		}

	// 進行中のライドは存在しないが、ユーザーがアクティブ状態
	case u.Ride == nil && u.State == UserStateActive:
		if count := len(u.RideHistory); (count == 1 && u.TotalEvaluation <= 1) || float64(u.TotalEvaluation)/float64(count) <= 2 {
			// 初回利用で評価1なら離脱
			// 2回以上利用して平均評価が2以下の場合は離脱
			if u.Region.UserLeave(u) {
				break
			}
			// Region内の最低ユーザー数を下回るならそのまま残る
		}

		// 過去のライドを確認する
		err := u.CheckRequestHistory(ctx)
		if err != nil {
			return WrapCodeError(ErrorCodeFailedToCheckRequestHistory, err)
		}

		// ライドを作成する
		// TODO 作成する条件・頻度
		err = u.CreateRide(ctx)
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
	if len(res.Requests) != len(u.RideHistory) {
		return fmt.Errorf("ライドの数が想定数と一致していません: expected=%d, got=%d", len(u.RideHistory), len(res.Requests))
	}

	historyMap := lo.KeyBy(u.RideHistory, func(r *Ride) string { return r.ServerID })
	for _, req := range res.Requests {
		expected, ok := historyMap[req.ID]
		if !ok {
			return fmt.Errorf("想定されないライドが含まれています: id=%s", req.ID)
		}
		if !req.DestinationCoordinate.Equals(expected.DestinationPoint) || !req.PickupCoordinate.Equals(expected.PickupPoint) {
			return fmt.Errorf("ライドの座標情報が正しくありません: id=%s", req.ID)
		}
		if req.Fare != expected.Fare() {
			return fmt.Errorf("ライドの運賃が正しくありません: id=%s", req.ID)
		}
		if req.Evaluation != expected.CalculateEvaluation().Score() {
			return fmt.Errorf("ライドの評価が正しくありません: id=%s", req.ID)
		}
		if req.Chair.ID != expected.Chair.ServerID || req.Chair.Name != expected.Chair.RegisteredData.Name || req.Chair.Model != expected.Chair.Model.Name || req.Chair.Owner != expected.Chair.Owner.RegisteredData.Name {
			return fmt.Errorf("ライドの椅子の情報が正しくありません: id=%s", req.ID)
		}
		if !req.CompletedAt.Equal(expected.ServerCompletedAt) {
			return fmt.Errorf("ライドの完了日時が正しくありません: id=%s", req.ID)
		}
	}

	return nil
}

func (u *User) CreateRide(ctx *Context) error {
	if u.Ride != nil {
		panic("ユーザーに進行中のライドがあるのにも関わらず、ライドを新規作成しようとしている")
	}

	u.InvitingLock.Lock()
	defer u.InvitingLock.Unlock()

	pickup, dest := RandomTwoCoordinateWithRand(u.Region, u.Rand.IntN(100)+5, u.Rand)

	ride := &Ride{
		User:             u,
		PickupPoint:      pickup,
		DestinationPoint: dest,
		RequestedAt:      ctx.CurrentTime(),
		Statuses: RideStatuses{
			Desired: RideStatusMatching,
			Chair:   RideStatusMatching,
			User:    RideStatusMatching,
		},
	}

	useInvCoupon := false
	switch {
	// 初回利用の割引を適用
	case len(u.RideHistory) == 0:
		ride.Discount = 3000

	// 招待された側のクーポンを適用
	case len(u.RideHistory) == 1 && u.Invited:
		ride.Discount = 1500

	// 招待した側のクーポンを適用
	case u.UnusedInvCoupons > 0:
		ride.Discount = 1000
		useInvCoupon = true
	}

	checkDistance := 50
	now := time.Now()
	nearby, err := u.Client.GetNearbyChairs(ctx, pickup, checkDistance)
	if err != nil {
		return WrapCodeError(ErrorCodeIncorrectNearbyChairs, err)
	}
	if err := u.World.checkNearbyChairsResponse(now, pickup, checkDistance, nearby); err != nil {
		return WrapCodeError(ErrorCodeIncorrectNearbyChairs, err)
	}
	if len(nearby.Chairs) == 0 {
		// 近くに椅子が無いので配車をやめる
		return nil
	}

	estimation, err := u.Client.GetEstimatedFare(ctx, pickup, dest)
	if err != nil {
		return WrapCodeError(ErrorCodeFailedToCreateRequest, err)
	}
	if ride.ActualDiscount() != estimation.Discount || ride.Fare() != estimation.Fare {
		return WrapCodeError(ErrorCodeFailedToCreateRequest, errors.New("ライド料金の見積もり金額が誤っています"))
	}

	res, err := u.Client.SendCreateRequest(ctx, ride)
	if err != nil {
		return WrapCodeError(ErrorCodeFailedToCreateRequest, err)
	}
	ride.ServerID = res.ServerRequestID
	u.Ride = ride
	u.RideHistory = append(u.RideHistory, ride)
	u.World.RideDB.Create(ride)
	if useInvCoupon {
		u.UnusedInvCoupons--
	}
	return nil
}

func (u *User) ChangeRideStatus(status RideStatus, serverRequestID string, validator func() error) error {
	ride := u.Ride
	if ride == nil {
		if status == RideStatusCompleted {
			// 履歴を見て、過去扱っていたRideに向けてのCOMPLETED通知であれば無視する
			for _, r := range slices.Backward(u.RideHistory) {
				if r.ServerID == serverRequestID && r.Statuses.Desired == RideStatusCompleted {
					r.Statuses.User = RideStatusCompleted
					return nil
				}
			}
		}
		return WrapCodeError(ErrorCodeUserNotRequestingButStatusChanged, fmt.Errorf("user_id: %s, got: %v", u.ServerID, status))
	}

	u.Ride.Statuses.Lock()
	defer u.Ride.Statuses.Unlock()
	if u.Ride.Statuses.User != status && u.Ride.Statuses.Desired != status {
		// 現在認識しているユーザーの状態で無いかつ、想定状態ではない状態に遷移しようとしている場合
		if ride.Statuses.User == RideStatusMatching && u.Ride.Statuses.Desired == RideStatusPickup && status == RideStatusEnRoute {
			// ユーザーにENROUTEが送られる前に、椅子が配車位置に到着している場合があるが、その時にPICKUPを受け取ることを許容する
		} else if u.Ride.Statuses.User == RideStatusPickup && u.Ride.Statuses.Desired == RideStatusArrived && status == RideStatusCarrying {
			// もう到着しているが、ユーザー側の通知が遅延していて、PICKUP状態からまだCARRYINGに遷移してないときは、CARRYINGを許容する
		} else if u.Ride.Statuses.Desired == RideStatusPickup && u.Ride.Statuses.User == RideStatusPickup && status == RideStatusCarrying {
			// ユーザーがPICKUPを受け取った状態で、椅子が出発リクエストを送った後、ベンチマーカーのDesiredステータスの変更を行う前にユーザー側にCarrying通知が届いてしまうことがあるがこれは許容する
		} else if status == RideStatusCompleted {
			// 履歴を見て、過去扱っていたRequestに向けてのCOMPLETED通知であれば無視する
			for _, r := range slices.Backward(u.RideHistory) {
				if r.ServerID == serverRequestID && r.Statuses.Desired == RideStatusCompleted {
					r.Statuses.User = RideStatusCompleted
					return nil
				}
			}
			return WrapCodeError(ErrorCodeUnexpectedUserRequestStatusTransitionOccurred, fmt.Errorf("ride_id: %v, expect: %v, got: %v (current: %v)", ride.ServerID, ride.Statuses.Desired, status, ride.Statuses.User))
		} else {
			return WrapCodeError(ErrorCodeUnexpectedUserRequestStatusTransitionOccurred, fmt.Errorf("ride_id: %v, expect: %v, got: %v (current: %v)", ride.ServerID, ride.Statuses.Desired, status, ride.Statuses.User))
		}
	}

	if err := validator(); err != nil {
		return WrapCodeError(ErrorCodeIncorrectUserNotificationData, err)
	}

	u.Ride.Statuses.User = status

	return nil
}

// HandleNotification 通知イベントを処理して、自身の状態を推移させる
func (u *User) HandleNotification(event NotificationEvent) error {
	switch data := event.(type) {
	case *UserNotificationEventMatching:
		err := u.ChangeRideStatus(RideStatusMatching, data.ServerRideID, func() error {
			if u.Ride.Statuses.User == RideStatusMatching && u.validatedRideNotificationEvent != nil {
				// バリデーション済みの結果と比較する
				return fmt.Errorf("ride_status (before: %s, after: %s) %w", u.Ride.Statuses.User, RideStatusMatching, compareUserNotificationEvent(data.ServerRideID, *u.validatedRideNotificationEvent, data.UserNotificationEvent))
			} else {
				if err := u.ValidateNotificationEvent(data.ServerRideID, data.UserNotificationEvent); err != nil {
					return fmt.Errorf("ride_status (before: %s, after: %s) %w", u.Ride.Statuses.User, RideStatusMatching, err)
				}
				// 新規登録 or 前回の結果を上書き
				u.validatedRideNotificationEvent = &data.UserNotificationEvent
				return nil
			}
		})
		if err != nil {
			return err
		}
	case *UserNotificationEventEnRoute:
		if u.validatedRideNotificationEvent == nil {
			slog.Error("validatedRideNotificationEventがnilです", slog.String("ride.Statuses.User", u.Ride.Statuses.User.String()), slog.String("data.ServerRideID", data.ServerRideID), slog.Int("len(RideHistory)", len(u.RideHistory)))

		}
		err := u.ChangeRideStatus(RideStatusEnRoute, data.ServerRideID, func() error {
			// MATCHING時に受け取った値と変わっていないはずなので比較のみを行う
			return fmt.Errorf("ride_status (before: %s, after: %s) %w", u.Ride.Statuses.User, RideStatusEnRoute, compareUserNotificationEvent(data.ServerRideID, *u.validatedRideNotificationEvent, data.UserNotificationEvent))
		})
		if err != nil {
			return err
		}

	case *UserNotificationEventPickup:
		err := u.ChangeRideStatus(RideStatusPickup, data.ServerRideID, func() error {
			// MATCHING時に受け取った値と変わっていないはずなので比較のみを行う
			return fmt.Errorf("ride_status (before: %s, after: %s) %w", u.Ride.Statuses.User, RideStatusPickup, compareUserNotificationEvent(data.ServerRideID, *u.validatedRideNotificationEvent, data.UserNotificationEvent))
		})
		if err != nil {
			return err
		}

		u.Ride.Statuses.User = RideStatusPickup

	case *UserNotificationEventCarrying:
		err := u.ChangeRideStatus(RideStatusCarrying, data.ServerRideID, func() error {
			// MATCHING時に受け取った値と変わっていないはずなので比較のみを行う
			return fmt.Errorf("ride_status (before: %s, after: %s) %w", u.Ride.Statuses.User, RideStatusCarrying, compareUserNotificationEvent(data.ServerRideID, *u.validatedRideNotificationEvent, data.UserNotificationEvent))
		})
		if err != nil {
			return err
		}

	case *UserNotificationEventArrived:
		err := u.ChangeRideStatus(RideStatusArrived, data.ServerRideID, func() error {
			// MATCHING時に受け取った値と変わっていないはずなので比較のみを行う
			return fmt.Errorf("ride_status (before: %s, after: %s) %w", u.Ride.Statuses.User, RideStatusArrived, compareUserNotificationEvent(data.ServerRideID, *u.validatedRideNotificationEvent, data.UserNotificationEvent))
		})
		if err != nil {
			return err
		}

	case *UserNotificationEventCompleted:
		err := u.ChangeRideStatus(RideStatusCompleted, data.ServerRideID, func() error {
			if u.Ride.Statuses.User != RideStatusCompleted {
				// ARRIVEDからCOMPLETEDに遷移した際には、chair.statsが変化しているはずのなので再度バリデーションを行う
				if err := u.ValidateNotificationEvent(data.ServerRideID, data.UserNotificationEvent); err != nil {
					return fmt.Errorf("ride_status (before: %s, after: %s) %w", u.Ride.Statuses.User, RideStatusCompleted, err)
				}
				u.validatedRideNotificationEvent = &data.UserNotificationEvent
				return nil
			} else {
				// バリデーション済みの結果と比較する
				return fmt.Errorf("ride_status (before: %s, after: %s) %w", u.Ride.Statuses.User, RideStatusCompleted, compareUserNotificationEvent(data.ServerRideID, *u.validatedRideNotificationEvent, data.UserNotificationEvent))
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) ValidateNotificationEvent(rideID string, serverSide UserNotificationEvent) error {
	if !serverSide.Pickup.Equals(u.Ride.PickupPoint) {
		return fmt.Errorf("配車位置が一致しません。(ride_id: %s, got: %s, want: %s)", rideID, serverSide.Pickup, u.Ride.PickupPoint)
	}
	if !serverSide.Destination.Equals(u.Ride.DestinationPoint) {
		return fmt.Errorf("目的地が一致しません。(ride_id: %s, got: %s, want: %s)", rideID, serverSide.Destination, u.Ride.DestinationPoint)
	}

	if serverSide.Fare != u.Ride.Fare() {
		return fmt.Errorf("運賃が一致しません。(ride_id: %s, got: %d, want: %d)", rideID, serverSide.Fare, u.Ride.Fare())
	}

	if serverSide.Chair == nil {
		return fmt.Errorf("椅子情報がありません。(ride_id: %s)", rideID)
	}

	serverSideChair := *serverSide.Chair

	var chair *Chair
	if u.Ride.Chair == nil {
		//	MATCHING時には椅子情報が存在しない
		chair = u.World.ChairDB.GetByServerID(serverSideChair.ID)
		if chair == nil {
			return fmt.Errorf("存在しない椅子が返却されました。(ride_id: %s, chair_id: %s)", rideID, serverSide.Chair.ID)
		}
	} else {
		//	COMPLETED時には椅子情報が存在する
		chair = u.Ride.Chair
	}

	if serverSideChair.Name != chair.RegisteredData.Name {
		return fmt.Errorf("椅子の名前が一致しません。(ride_id: %s, chair_id: %s, got: %s, want: %s)", rideID, serverSide.Chair.ID, serverSide.Chair.Name, u.Ride.Chair.RegisteredData.Name)
	}
	if serverSideChair.Model != chair.Model.Name {
		return fmt.Errorf("椅子のモデルが一致しません。(ride_id: %s, chair_id: %s, got: %s, want: %s)", rideID, serverSide.Chair.ID, serverSide.Chair.Model, u.Ride.Chair.Model)
	}

	totalRideCount := 0
	totalEvaluation := 0
	for _, r := range chair.RideHistory.Iter() {
		if r.Evaluated.Load() {
			totalRideCount++
			totalEvaluation += r.CalculateEvaluation().Score()
		}
	}

	if serverSideChair.Stats.TotalRidesCount != totalRideCount {
		return fmt.Errorf("椅子の総乗車回数が一致しません。(ride_id: %s, chair_id: %s, got: %d, want: %d)", rideID, serverSide.Chair.ID, serverSide.Chair.Stats.TotalRidesCount, totalRideCount)
	}
	if totalRideCount != 0 && serverSideChair.Stats.TotalEvaluationAvg != float64(totalEvaluation)/float64(totalRideCount) {
		return fmt.Errorf("椅子の評価の平均が一致しません。(ride_id: %s, chair_id: %s, got: %f, want: %f)", rideID, serverSide.Chair.ID, serverSide.Chair.Stats.TotalEvaluationAvg, float64(totalEvaluation)/float64(totalRideCount))
	} else if serverSideChair.Stats.TotalEvaluationAvg != 0 {
		return fmt.Errorf("椅子の評価の平均が一致しません。(ride_id: %s, chair_id: %s, got: %f, want: %f)", rideID, serverSide.Chair.ID, serverSide.Chair.Stats.TotalEvaluationAvg, 0.0)
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
	if new.Chair.Stats.TotalEvaluationAvg != old.Chair.Stats.TotalEvaluationAvg {
		return fmt.Errorf("椅子の評価の平均が一致しません。(ride_id: %s, chair_id: %s, got: %f, want: %f)", rideID, new.Chair.ID, new.Chair.Stats.TotalEvaluationAvg, old.Chair.Stats.TotalEvaluationAvg)
	}

	return nil
}
