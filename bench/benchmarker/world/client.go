package world

import (
	"time"

	"github.com/guregu/null/v5"
	"github.com/isucon/isucon14/bench/benchrun"
)

type WorldClient interface {
	// RegisterUser サーバーにユーザーを登録する
	RegisterUser(ctx *Context, data *RegisterUserRequest, beforeRequest func(client UserClient) error) (*RegisterUserResponse, error)
	// RegisterOwner サーバーにオーナーを登録する
	RegisterOwner(ctx *Context, data *RegisterOwnerRequest, beforeRequest func(client OwnerClient) error) (*RegisterOwnerResponse, error)
	// RegisterChair サーバーにユーザーを登録する
	RegisterChair(ctx *Context, owner *Owner, data *RegisterChairRequest) (*RegisterChairResponse, error)
}

type UserClient interface {
	// SendCreateRequest サーバーにリクエスト作成を送信する
	SendCreateRequest(ctx *Context, req *Ride) (*SendCreateRequestResponse, error)
	// GetRequests サーバーからリクエスト一覧を取得する
	GetRequests(ctx *Context) (*GetRequestsResponse, error)
	// GetNearbyChairs サーバーから近くの椅子の情報を取得する
	GetNearbyChairs(ctx *Context, current Coordinate, distance int) (*GetNearbyChairsResponse, error)
	// GetEstimatedFare サーバーから料金の見積もりを取る
	GetEstimatedFare(ctx *Context, pickup Coordinate, dest Coordinate) (*GetEstimatedFareResponse, error)
	// SendEvaluation サーバーに今回の送迎の評価を送信する
	SendEvaluation(ctx *Context, req *Ride, score int) (*SendEvaluationResponse, error)
	// RegisterPaymentMethods サーバーにユーザーの支払い情報を登録する
	RegisterPaymentMethods(ctx *Context, user *User) error
	// ConnectUserNotificationStream ユーザー用の通知ストリームに接続する
	ConnectUserNotificationStream(ctx *Context, user *User, receiver NotificationReceiverFunc) (NotificationStream, error)
	// BrowserAccess ブラウザでアクセスしたときのリクエストを送信する
	BrowserAccess(ctx *Context, scenario benchrun.FrontendPathScenario) error
}

type OwnerClient interface {
	// GetOwnerSales サーバーからオーナーの売り上げ情報を取得する
	GetOwnerSales(ctx *Context, args *GetOwnerSalesRequest) (*GetOwnerSalesResponse, error)
	// GetOwnerChairs サーバーからオーナーの椅子一覧を取得する
	GetOwnerChairs(ctx *Context, args *GetOwnerChairsRequest) (*GetOwnerChairsResponse, error)
	// BrowserAccess ブラウザでアクセスしたときのリクエストを送信する
	BrowserAccess(ctx *Context, scenario benchrun.FrontendPathScenario) error
}

type ChairClient interface {
	// SendChairCoordinate サーバーに椅子の座標を送信する
	SendChairCoordinate(ctx *Context, chair *Chair) (*SendChairCoordinateResponse, error)
	// SendAcceptRequest サーバーに配椅子要求を受理することを報告する
	SendAcceptRequest(ctx *Context, chair *Chair, req *Ride) error
	// SendDepart サーバーに客が搭乗完了して出発することを報告する
	SendDepart(ctx *Context, req *Ride) error
	// SendActivate サーバーにリクエストの受付開始を通知する
	SendActivate(ctx *Context, chair *Chair) error
	// ConnectChairNotificationStream 椅子用の通知ストリームに接続する
	ConnectChairNotificationStream(ctx *Context, chair *Chair, receiver NotificationReceiverFunc) (NotificationStream, error)
}

type SendCreateRequestResponse struct {
	ServerRequestID string
}

type GetRequestsResponse struct {
	Requests []*RequestHistory
}

type GetEstimatedFareResponse struct {
	Fare     int
	Discount int
}

type GetNearbyChairsResponse struct {
	RetrievedAt time.Time
	Chairs      []*AppChair
}

type AppChair struct {
	ID         string
	Name       string
	Model      string
	Coordinate Coordinate
}

type RequestHistory struct {
	ID                    string
	PickupCoordinate      Coordinate
	DestinationCoordinate Coordinate
	Chair                 RequestHistoryChair
	Fare                  int
	Evaluation            int
	RequestedAt           time.Time
	CompletedAt           time.Time
}

type RequestHistoryChair struct {
	ID    string
	Owner string
	Name  string
	Model string
}

type GetOwnerSalesRequest struct {
	Since time.Time
	Until time.Time
}

type GetOwnerSalesResponse struct {
	Total  int
	Chairs []*ChairSales
	Models []*ChairSalesPerModel
}

type ChairSales struct {
	ID    string
	Name  string
	Sales int
}

type ChairSalesPerModel struct {
	Model string
	Sales int
}

type GetOwnerChairsRequest struct{}

type GetOwnerChairsResponse struct {
	Chairs []*OwnerChair
}

type OwnerChair struct {
	ID                     string
	Name                   string
	Model                  string
	Active                 bool
	RegisteredAt           time.Time
	TotalDistance          int
	TotalDistanceUpdatedAt null.Time
}

type SendChairCoordinateResponse struct {
	RecordedAt time.Time
}

type SendEvaluationResponse struct {
	CompletedAt time.Time
}

type RegisterUserRequest struct {
	UserName       string
	FirstName      string
	LastName       string
	DateOfBirth    string
	InvitationCode string
}

type RegisterUserResponse struct {
	ServerUserID   string
	InvitationCode string

	Client UserClient
}

type RegisterOwnerRequest struct {
	Name string
}

type RegisterOwnerResponse struct {
	ServerOwnerID        string
	ChairRegisteredToken string

	Client OwnerClient
}

type RegisterChairRequest struct {
	Name  string
	Model string
}

type RegisterChairResponse struct {
	ServerChairID string
	ServerOwnerID string

	Client ChairClient
}

type NotificationReceiverFunc func(event NotificationEvent)

type NotificationStream interface {
	Close()
}
