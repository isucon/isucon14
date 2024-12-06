package world

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/guregu/null/v5"
)

const (
	// InitialFare 初乗り運賃
	InitialFare = 500
	// FarePerDistance １距離あたりの運賃
	FarePerDistance = 100
)

type RideStatus int

func (r RideStatus) String() string {
	switch r {
	case RideStatusMatching:
		return "MATCHING"
	case RideStatusEnRoute:
		return "ENROUTE"
	case RideStatusPickup:
		return "PICKUP"
	case RideStatusCarrying:
		return "CARRYING"
	case RideStatusArrived:
		return "ARRIVED"
	case RideStatusCompleted:
		return "COMPLETED"
	default:
		return "UNKNOWN"
	}
}

const (
	RideStatusMatching RideStatus = iota
	RideStatusEnRoute
	RideStatusPickup
	RideStatusCarrying
	RideStatusArrived
	RideStatusCompleted
)

type RideID int

type Ride struct {
	// ID ベンチマーカー内部ライドID
	ID RideID
	// ServerID サーバー上でのライドID
	ServerID string
	// User リクエストしたユーザー
	User *User
	// PickupPoint 配椅子位置
	PickupPoint Coordinate
	// DestinationPoint 目的地
	DestinationPoint Coordinate
	// Discount 最大割引額
	Discount int

	// Chair 割り当てられた椅子。割り当てられるまでnil
	Chair *Chair
	// StartPoint 椅子の初期位置。割り当てられるまでnil
	StartPoint null.Value[Coordinate]

	// RequestedAt 配車リクエストを行った時間
	RequestedAt int64
	// MatchedAt マッチングが完了した時間。割り当てられるまで0
	MatchedAt int64
	// DispatchedAt 配車位置についた時間。割り当てられるまで0
	DispatchedAt int64
	// PickedUpAt ピックアップされ出発された時間。割り当てられるまで0
	PickedUpAt int64
	// ArrivedAt 目的地に到着した時間。割り当てられるまで0
	ArrivedAt int64
	// CompletedAt ライドが正常に完了した時間。割り当てられるまで0
	CompletedAt int64
	// ServerCompletedAt サーバー側で記録されている完了時間
	ServerCompletedAt time.Time
	// BenchMatchedAt ベンチがAcceptのリクエストを送って成功した時間
	BenchMatchedAt time.Time

	// Evaluated ライドの評価が完了しているかどうか
	Evaluated atomic.Bool

	Statuses RideStatuses
}

func (r *Ride) String() string {
	chairID := "<nil>"
	if r.Chair != nil {
		chairID = strconv.Itoa(int(r.Chair.ID))
	}
	return fmt.Sprintf(
		"Ride{id=%d,status=%s,user=%d,from=%s,to=%s,chair=%s,time=%s}",
		r.ID, r.Statuses.String(), r.User.ID, r.PickupPoint, r.DestinationPoint, chairID, r.timelineString(),
	)
}

func (r *Ride) SetID(id RideID) {
	r.ID = id
}

// Sales 売り上げ
func (r *Ride) Sales() int {
	return InitialFare + r.PickupPoint.DistanceTo(r.DestinationPoint)*FarePerDistance
}

// Fare ユーザーが支払う料金
func (r *Ride) Fare() int {
	return InitialFare + max(r.PickupPoint.DistanceTo(r.DestinationPoint)*FarePerDistance-r.Discount, 0)
}

// ActualDiscount 実際に割り引いた価格
func (r *Ride) ActualDiscount() int {
	return r.Sales() - r.Fare()
}

// CalculateEvaluation 送迎の評価値を計算する
func (r *Ride) CalculateEvaluation() Evaluation {
	if !(r.MatchedAt > 0 && r.DispatchedAt > 0 && r.PickedUpAt > 0 && r.ArrivedAt > 0) {
		panic("計算に必要な時間情報が足りていない状況なのに評価値を計算しようとしている")
	}

	// TODO: いい感じにする
	result := Evaluation{}
	{
		// マッチング待ち時間評価
		if r.MatchedAt-r.RequestedAt < 100 {
			// 100ticks以内ならOK
			result.Matching = true
		}
	}
	{
		// 乗車待ち時間評価
		if r.StartPoint.V.DistanceTo(r.PickupPoint) < 25 {
			// 割り当てられた椅子が自分の場所から距離25以内
			result.Dispatch = true
		}
	}
	{
		// 乗車待ち時間誤差評価
		idealTime := neededTime(r.StartPoint.V.DistanceTo(r.PickupPoint), r.Chair.Model.Speed)
		actualTime := int(r.PickedUpAt - r.MatchedAt)
		if actualTime-idealTime < 15 {
			// 理想時間との誤差が15ticks以内ならOK
			result.Pickup = true
		}
	}
	{
		// 乗車時間誤差評価
		idealTime := neededTime(r.PickupPoint.DistanceTo(r.DestinationPoint), r.Chair.Model.Speed)
		actualTime := int(r.ArrivedAt - r.PickedUpAt)
		if actualTime-idealTime < 5 {
			// 理想時間との誤差が5ticks以内ならOK
			result.Drive = true
		}
	}

	return result
}

func (r *Ride) Intervals() [3]int64 {
	return [3]int64{
		max(0, r.MatchedAt-r.RequestedAt),
		max(0, r.DispatchedAt-r.MatchedAt),
		max(0, r.ArrivedAt-r.DispatchedAt),
	}
}

func (r *Ride) timelineString() string {
	baseTime := r.RequestedAt
	matchTime := max(0, r.MatchedAt-r.RequestedAt)
	dispatchTime := max(0, r.DispatchedAt-r.RequestedAt)
	pickedUpTime := max(0, r.PickedUpAt-r.RequestedAt)
	arrivedTime := max(0, r.ArrivedAt-r.RequestedAt)
	completedTime := max(0, r.CompletedAt-r.RequestedAt)
	return fmt.Sprintf("[0(base=%d),%d,%d,%d,%d,%d]", baseTime, matchTime, dispatchTime, pickedUpTime, arrivedTime, completedTime)
}

const ForwardingScoreDenominator = 10

func (r *Ride) Score() int {
	return r.Sales() + r.StartPoint.V.DistanceTo(r.PickupPoint)*FarePerDistance/ForwardingScoreDenominator
}

func (r *Ride) PartialScore() int {
	switch r.Statuses.Desired {
	case RideStatusMatching:
		return 0
	case RideStatusEnRoute:
		return r.StartPoint.V.DistanceTo(r.Chair.Location.Current()) * FarePerDistance / ForwardingScoreDenominator
	case RideStatusPickup:
		return r.StartPoint.V.DistanceTo(r.PickupPoint) * FarePerDistance / ForwardingScoreDenominator
	case RideStatusCarrying:
		return r.StartPoint.V.DistanceTo(r.PickupPoint)*FarePerDistance/ForwardingScoreDenominator + r.PickupPoint.DistanceTo(r.Chair.Location.Current())*FarePerDistance
	case RideStatusArrived:
		return r.Score() - InitialFare
	case RideStatusCompleted:
		return r.Score()
	default:
		panic("unknown status")
	}
}

type Evaluation struct {
	Matching bool
	Dispatch bool
	Pickup   bool
	Drive    bool
}

func (e Evaluation) String() string {
	return fmt.Sprintf("score: %d (matching:%v, dispath:%v, pickup:%v, drive:%v)", e.Score(), e.Matching, e.Dispatch, e.Pickup, e.Drive)
}

func (e Evaluation) Score() int {
	result := 1
	if e.Matching {
		result++
	}
	if e.Dispatch {
		result++
	}
	if e.Pickup {
		result++
	}
	if e.Drive {
		result++
	}
	return result
}

type RideStatuses struct {
	// Desired 現在の想定されるステータス
	Desired RideStatus
	// Chair 現在椅子が認識しているステータス
	Chair RideStatus
	// User 現在ユーザーが認識しているステータス
	User RideStatus

	m sync.RWMutex
}

func (s *RideStatuses) String() string {
	d, c, u := s.Get()
	return fmt.Sprintf("(%v,%v,%v)", d, c, u)
}

func (s *RideStatuses) Get() (desired, chair, user RideStatus) {
	return s.Desired, s.Chair, s.User
}

// Lock DesiredのみをWrite Lockします
// MEMO: ロックを取らなけらばならないところ以外はとってない
func (s *RideStatuses) Lock() { s.m.Lock() }

// Unlock DesiredのみをWrite Unlockします
// MEMO: ロックを取らなけらばならないところ以外はとってない
func (s *RideStatuses) Unlock() { s.m.Unlock() }

// RLock DesiredのみをRead Lockします
// MEMO: ロックを取らなけらばならないところ以外はとってない
func (s *RideStatuses) RLock() { s.m.RLock() }

// RUnlock DesiredのみをRead Unlockします
// MEMO: ロックを取らなけらばならないところ以外はとってない
func (s *RideStatuses) RUnlock() { s.m.RUnlock() }
