package world

import (
	"fmt"
	"math"
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

type RequestStatus int

func (r RequestStatus) String() string {
	switch r {
	case RequestStatusMatching:
		return "MATCHING"
	case RequestStatusDispatching:
		return "ENROUTE"
	case RequestStatusDispatched:
		return "PICKUP"
	case RequestStatusCarrying:
		return "CARRYING"
	case RequestStatusArrived:
		return "ARRIVED"
	case RequestStatusCompleted:
		return "COMPLETED"
	default:
		return "UNKNOWN"
	}
}

const (
	RequestStatusMatching RequestStatus = iota
	RequestStatusDispatching
	RequestStatusDispatched
	RequestStatusCarrying
	RequestStatusArrived
	RequestStatusCompleted
)

type RequestID int

type Request struct {
	// ID ベンチマーカー内部リクエストID
	ID RequestID
	// ServerID サーバー上でのリクエストID
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

	// RequestedAt リクエストを行った時間
	RequestedAt int64
	// MatchedAt マッチングが完了した時間。割り当てられるまで0
	MatchedAt int64
	// DispatchedAt 配椅子位置についた時間。割り当てられるまで0
	DispatchedAt int64
	// PickedUpAt ピックアップされ出発された時間。割り当てられるまで0
	PickedUpAt int64
	// ArrivedAt 目的地に到着した時間。割り当てられるまで0
	ArrivedAt int64
	// CompletedAt リクエストが正常に完了した時間。割り当てられるまで0
	CompletedAt int64
	// ServerCompletedAt サーバー側で記録されている完了時間
	ServerCompletedAt time.Time
	// BenchMatchedAt ベンチがAcceptのリクエストを送って成功した時間
	BenchMatchedAt time.Time

	// Evaluated リクエストの評価が完了しているかどうか
	Evaluated atomic.Bool

	Statuses RequestStatuses
}

func (r *Request) String() string {
	chairID := "<nil>"
	if r.Chair != nil {
		chairID = strconv.Itoa(int(r.Chair.ID))
	}
	return fmt.Sprintf(
		"Request{id=%d,status=%s,user=%d,from=%s,to=%s,chair=%s,time=%s}",
		r.ID, r.Statuses.String(), r.User.ID, r.PickupPoint, r.DestinationPoint, chairID, r.timelineString(),
	)
}

func (r *Request) SetID(id RequestID) {
	r.ID = id
}

// Sales 売り上げ
func (r *Request) Sales() int {
	return InitialFare + r.PickupPoint.DistanceTo(r.DestinationPoint)*FarePerDistance
}

// Fare ユーザーが支払う料金
func (r *Request) Fare() int {
	return InitialFare + max(r.PickupPoint.DistanceTo(r.DestinationPoint)*FarePerDistance-r.Discount, 0)
}

// ActualDiscount 実際に割り引いた価格
func (r *Request) ActualDiscount() int {
	return r.Sales() - r.Fare()
}

// CalculateEvaluation 送迎の評価値を計算する
func (r *Request) CalculateEvaluation() Evaluation {
	if !(r.MatchedAt > 0 && r.DispatchedAt > 0 && r.PickedUpAt > 0 && r.ArrivedAt > 0) {
		panic("計算に必要な時間情報が足りていない状況なのに評価値を計算しようとしている")
	}

	// TODO: いい感じにする
	result := Evaluation{}
	{
		// マッチング待ち時間評価
		diff := r.MatchedAt - r.RequestedAt
		if diff < 100 {
			// 100ticks以内なら満点
			result.Matching = 1
		} else {
			result.Matching = f1_x2(float64(100), float64(diff))
		}
	}
	{
		// 乗車待ち時間評価
		distance := r.StartPoint.V.DistanceTo(r.PickupPoint)
		if distance < 25 {
			// 割り当てられた椅子が自分の場所から距離25以内なら満点
			result.Dispatch = 1
		} else {
			result.Dispatch = f1_x2(float64(25), float64(distance))
		}
	}
	{
		// 乗車待ち時間誤差評価
		idealTime := neededTime(r.StartPoint.V.DistanceTo(r.PickupPoint), r.Chair.Model.Speed)
		actualTime := int(r.PickedUpAt - r.MatchedAt)
		diff := actualTime - idealTime
		if diff < 15 {
			// 理想時間との誤差が15ticks以内なら満点
			result.Pickup = 1
		} else {
			result.Pickup = f1_x2(float64(15), float64(diff))
		}
	}
	{
		// 乗車時間誤差評価
		idealTime := neededTime(r.PickupPoint.DistanceTo(r.DestinationPoint), r.Chair.Model.Speed)
		actualTime := int(r.ArrivedAt - r.PickedUpAt)
		diff := actualTime - idealTime
		if diff < 5 {
			// 理想時間との誤差が5ticks以内なら満点
			result.Drive = 1
		} else {
			result.Drive = f1_x2(float64(5), float64(diff))
		}
	}

	return result
}

func (r *Request) Intervals() [3]int64 {
	return [3]int64{
		max(0, r.MatchedAt-r.RequestedAt),
		max(0, r.DispatchedAt-r.MatchedAt),
		max(0, r.ArrivedAt-r.DispatchedAt),
	}
}

func (r *Request) timelineString() string {
	baseTime := r.RequestedAt
	matchTime := max(0, r.MatchedAt-r.RequestedAt)
	dispatchTime := max(0, r.DispatchedAt-r.RequestedAt)
	pickedUpTime := max(0, r.PickedUpAt-r.RequestedAt)
	arrivedTime := max(0, r.ArrivedAt-r.RequestedAt)
	completedTime := max(0, r.CompletedAt-r.RequestedAt)
	return fmt.Sprintf("[0(base=%d),%d,%d,%d,%d,%d]", baseTime, matchTime, dispatchTime, pickedUpTime, arrivedTime, completedTime)
}

const ForwardingScoreDenominator = 10

func (r *Request) Score() int {
	return r.Sales() + r.StartPoint.V.DistanceTo(r.PickupPoint)*FarePerDistance/ForwardingScoreDenominator
}

func (r *Request) PartialScore() int {
	switch r.Statuses.Desired {
	case RequestStatusMatching:
		return 0
	case RequestStatusDispatching:
		return r.StartPoint.V.DistanceTo(r.Chair.Location.Current()) * FarePerDistance / ForwardingScoreDenominator
	case RequestStatusDispatched:
		return r.StartPoint.V.DistanceTo(r.PickupPoint) * FarePerDistance / ForwardingScoreDenominator
	case RequestStatusCarrying:
		return r.StartPoint.V.DistanceTo(r.PickupPoint)*FarePerDistance/ForwardingScoreDenominator + r.PickupPoint.DistanceTo(r.Chair.Location.Current())*FarePerDistance
	case RequestStatusArrived:
		return r.Score() - InitialFare
	case RequestStatusCompleted:
		return r.Score()
	default:
		panic("unknown status")
	}
}

type Evaluation struct {
	Matching float64
	Dispatch float64
	Pickup   float64
	Drive    float64
}

func (e Evaluation) String() string {
	return fmt.Sprintf("score: %d (matching:%v, dispath:%v, pickup:%v, drive:%v)", e.Score(), e.Matching, e.Dispatch, e.Pickup, e.Drive)
}

func (e Evaluation) Map() [4]float64 {
	return [4]float64{e.Matching, e.Dispatch, e.Pickup, e.Drive}
}

func (e Evaluation) Score() int {
	total := int(math.Round(e.Matching + e.Dispatch + e.Pickup + e.Drive))
	if total <= 0 {
		total = 1
	}
	return total
}

func f1_x2(upper float64, lower float64) float64 {
	return (upper / lower) * (upper / lower)
}

type RequestStatuses struct {
	// Desired 現在の想定されるリクエストステータス
	Desired RequestStatus
	// Chair 現在椅子が認識しているステータス
	Chair RequestStatus
	// User 現在ユーザーが認識しているステータス
	User RequestStatus

	m sync.RWMutex
}

func (s *RequestStatuses) String() string {
	d, c, u := s.Get()
	return fmt.Sprintf("(%v,%v,%v)", d, c, u)
}

func (s *RequestStatuses) Get() (desired, chair, user RequestStatus) {
	return s.Desired, s.Chair, s.User
}

// Lock DesiredのみをWrite Lockします
// MEMO: ロックを取らなけらばならないところ以外はとってない
func (s *RequestStatuses) Lock() { s.m.Lock() }

// Unlock DesiredのみをWrite Unlockします
// MEMO: ロックを取らなけらばならないところ以外はとってない
func (s *RequestStatuses) Unlock() { s.m.Unlock() }

// RLock DesiredのみをRead Lockします
// MEMO: ロックを取らなけらばならないところ以外はとってない
func (s *RequestStatuses) RLock() { s.m.RLock() }

// RUnlock DesiredのみをRead Unlockします
// MEMO: ロックを取らなけらばならないところ以外はとってない
func (s *RequestStatuses) RUnlock() { s.m.RUnlock() }
