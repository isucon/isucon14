package world

import (
	"time"

	"github.com/guregu/null/v5"
)

type WorldClient interface {
	// RegisterUser サーバーにユーザーを登録する
	RegisterUser(ctx *Context, data *RegisterUserRequest) (*RegisterUserResponse, error)
	// RegisterProvider サーバーにプロバイダーを登録する
	RegisterProvider(ctx *Context, data *RegisterProviderRequest) (*RegisterProviderResponse, error)
	// RegisterChair サーバーにユーザーを登録する
	RegisterChair(ctx *Context, provider *Provider, data *RegisterChairRequest) (*RegisterChairResponse, error)
}

type UserClient interface {
	// SendCreateRequest サーバーにリクエスト作成を送信する
	SendCreateRequest(ctx *Context, req *Request) (*SendCreateRequestResponse, error)
	// GetRequests サーバーからリクエスト一覧を取得する
	GetRequests(ctx *Context) (*GetRequestsResponse, error)
	// SendEvaluation サーバーに今回の送迎の評価を送信する
	SendEvaluation(ctx *Context, req *Request, score int) (*SendEvaluationResponse, error)
	// RegisterPaymentMethods サーバーにユーザーの支払い情報を登録する
	RegisterPaymentMethods(ctx *Context, user *User) error
	// ConnectUserNotificationStream ユーザー用の通知ストリームに接続する
	ConnectUserNotificationStream(ctx *Context, user *User, receiver NotificationReceiverFunc) (NotificationStream, error)
}

type ProviderClient interface {
	// GetProviderSales サーバーからプロバイダーの売り上げ情報を取得する
	GetProviderSales(ctx *Context, args *GetProviderSalesRequest) (*GetProviderSalesResponse, error)
	// GetProviderChairs サーバーからプロバイダーの椅子一覧を取得する
	GetProviderChairs(ctx *Context, args *GetProviderChairsRequest) (*GetProviderChairsResponse, error)
}

type ChairClient interface {
	// SendChairCoordinate サーバーに椅子の座標を送信する
	SendChairCoordinate(ctx *Context, chair *Chair) (*SendChairCoordinateResponse, error)
	// SendAcceptRequest サーバーに配椅子要求を受理することを報告する
	SendAcceptRequest(ctx *Context, chair *Chair, req *Request) error
	// SendDenyRequest サーバーに配椅子要求を受理することを報告する
	SendDenyRequest(ctx *Context, chair *Chair, serverRequestID string) error
	// SendDepart サーバーに客が搭乗完了して出発することを報告する
	SendDepart(ctx *Context, req *Request) error
	// SendActivate サーバーにリクエストの受付開始を通知する
	SendActivate(ctx *Context, chair *Chair) error
	// SendDeactivate サーバーにリクエストの受付停止を通知する
	SendDeactivate(ctx *Context, chair *Chair) error
	// GetRequestByChair サーバーからリクエストの詳細を取得する(椅子側)
	GetRequestByChair(ctx *Context, chair *Chair, serverRequestID string) (*GetRequestByChairResponse, error)
	// ConnectChairNotificationStream 椅子用の通知ストリームに接続する
	ConnectChairNotificationStream(ctx *Context, chair *Chair, receiver NotificationReceiverFunc) (NotificationStream, error)
}

type SendCreateRequestResponse struct {
	ServerRequestID string
}

type GetRequestsResponse struct {
	Requests []*RequestHistory
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

type GetRequestByChairResponse struct{}

type GetProviderSalesRequest struct {
	Since time.Time
	Until time.Time
}

type GetProviderSalesResponse struct {
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

type GetProviderChairsRequest struct{}

type GetProviderChairsResponse struct {
	Chairs []*ProviderChair
}

type ProviderChair struct {
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
	Fare        int
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

type RegisterProviderRequest struct {
	Name string
}

type RegisterProviderResponse struct {
	ServerProviderID     string
	ChairRegisteredToken string

	Client ProviderClient
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
