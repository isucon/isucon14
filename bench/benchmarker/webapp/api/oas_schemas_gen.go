// Code generated by ogen, DO NOT EDIT.

package api

import (
	"github.com/go-faster/errors"
)

// AppGetNotificationNoContent is response for AppGetNotification operation.
type AppGetNotificationNoContent struct{}

func (*AppGetNotificationNoContent) appGetNotificationRes() {}

// AppPostInquiryNoContent is response for AppPostInquiry operation.
type AppPostInquiryNoContent struct{}

type AppPostInquiryReq struct {
	// 件名.
	Subject string `json:"subject"`
	// 問い合わせ内容.
	Body string `json:"body"`
}

// GetSubject returns the value of Subject.
func (s *AppPostInquiryReq) GetSubject() string {
	return s.Subject
}

// GetBody returns the value of Body.
func (s *AppPostInquiryReq) GetBody() string {
	return s.Body
}

// SetSubject sets the value of Subject.
func (s *AppPostInquiryReq) SetSubject(val string) {
	s.Subject = val
}

// SetBody sets the value of Body.
func (s *AppPostInquiryReq) SetBody(val string) {
	s.Body = val
}

// AppPostPaymentMethodsNoContent is response for AppPostPaymentMethods operation.
type AppPostPaymentMethodsNoContent struct{}

type AppPostPaymentMethodsReq struct {
	// 決済トークン.
	Token string `json:"token"`
}

// GetToken returns the value of Token.
func (s *AppPostPaymentMethodsReq) GetToken() string {
	return s.Token
}

// SetToken sets the value of Token.
func (s *AppPostPaymentMethodsReq) SetToken(val string) {
	s.Token = val
}

type AppPostRegisterOK struct {
	// アクセストークン.
	AccessToken string `json:"access_token"`
	// ユーザーID.
	ID string `json:"id"`
}

// GetAccessToken returns the value of AccessToken.
func (s *AppPostRegisterOK) GetAccessToken() string {
	return s.AccessToken
}

// GetID returns the value of ID.
func (s *AppPostRegisterOK) GetID() string {
	return s.ID
}

// SetAccessToken sets the value of AccessToken.
func (s *AppPostRegisterOK) SetAccessToken(val string) {
	s.AccessToken = val
}

// SetID sets the value of ID.
func (s *AppPostRegisterOK) SetID(val string) {
	s.ID = val
}

func (*AppPostRegisterOK) appPostRegisterRes() {}

type AppPostRegisterReq struct {
	// ユーザー名.
	Username string `json:"username"`
	// 名前.
	Firstname string `json:"firstname"`
	// 名字.
	Lastname string `json:"lastname"`
	// 生年月日.
	DateOfBirth string `json:"date_of_birth"`
}

// GetUsername returns the value of Username.
func (s *AppPostRegisterReq) GetUsername() string {
	return s.Username
}

// GetFirstname returns the value of Firstname.
func (s *AppPostRegisterReq) GetFirstname() string {
	return s.Firstname
}

// GetLastname returns the value of Lastname.
func (s *AppPostRegisterReq) GetLastname() string {
	return s.Lastname
}

// GetDateOfBirth returns the value of DateOfBirth.
func (s *AppPostRegisterReq) GetDateOfBirth() string {
	return s.DateOfBirth
}

// SetUsername sets the value of Username.
func (s *AppPostRegisterReq) SetUsername(val string) {
	s.Username = val
}

// SetFirstname sets the value of Firstname.
func (s *AppPostRegisterReq) SetFirstname(val string) {
	s.Firstname = val
}

// SetLastname sets the value of Lastname.
func (s *AppPostRegisterReq) SetLastname(val string) {
	s.Lastname = val
}

// SetDateOfBirth sets the value of DateOfBirth.
func (s *AppPostRegisterReq) SetDateOfBirth(val string) {
	s.DateOfBirth = val
}

type AppPostRequestAccepted struct {
	// 配車要求ID.
	RequestID string `json:"request_id"`
}

// GetRequestID returns the value of RequestID.
func (s *AppPostRequestAccepted) GetRequestID() string {
	return s.RequestID
}

// SetRequestID sets the value of RequestID.
func (s *AppPostRequestAccepted) SetRequestID(val string) {
	s.RequestID = val
}

type AppPostRequestEvaluateBadRequest Error

func (*AppPostRequestEvaluateBadRequest) appPostRequestEvaluateRes() {}

// AppPostRequestEvaluateNoContent is response for AppPostRequestEvaluate operation.
type AppPostRequestEvaluateNoContent struct{}

func (*AppPostRequestEvaluateNoContent) appPostRequestEvaluateRes() {}

type AppPostRequestEvaluateNotFound Error

func (*AppPostRequestEvaluateNotFound) appPostRequestEvaluateRes() {}

type AppPostRequestEvaluateReq struct {
	// 椅子の評価.
	Evaluation int `json:"evaluation"`
}

// GetEvaluation returns the value of Evaluation.
func (s *AppPostRequestEvaluateReq) GetEvaluation() int {
	return s.Evaluation
}

// SetEvaluation sets the value of Evaluation.
func (s *AppPostRequestEvaluateReq) SetEvaluation(val int) {
	s.Evaluation = val
}

type AppPostRequestReq struct {
	// 配車位置.
	PickupCoordinate Coordinate `json:"pickup_coordinate"`
	// 目的地.
	DestinationCoordinate Coordinate `json:"destination_coordinate"`
}

// GetPickupCoordinate returns the value of PickupCoordinate.
func (s *AppPostRequestReq) GetPickupCoordinate() Coordinate {
	return s.PickupCoordinate
}

// GetDestinationCoordinate returns the value of DestinationCoordinate.
func (s *AppPostRequestReq) GetDestinationCoordinate() Coordinate {
	return s.DestinationCoordinate
}

// SetPickupCoordinate sets the value of PickupCoordinate.
func (s *AppPostRequestReq) SetPickupCoordinate(val Coordinate) {
	s.PickupCoordinate = val
}

// SetDestinationCoordinate sets the value of DestinationCoordinate.
func (s *AppPostRequestReq) SetDestinationCoordinate(val Coordinate) {
	s.DestinationCoordinate = val
}

// App向け配車要求情報.
// Ref: #/components/schemas/AppRequest
type AppRequest struct {
	// 配車要求ID.
	RequestID             string        `json:"request_id"`
	PickupCoordinate      Coordinate    `json:"pickup_coordinate"`
	DestinationCoordinate Coordinate    `json:"destination_coordinate"`
	Status                RequestStatus `json:"status"`
	Chair                 OptChair      `json:"chair"`
	// 配車要求日時.
	CreatedAt float64 `json:"created_at"`
	// 配車要求更新日時.
	UpdatedAt float64 `json:"updated_at"`
}

// GetRequestID returns the value of RequestID.
func (s *AppRequest) GetRequestID() string {
	return s.RequestID
}

// GetPickupCoordinate returns the value of PickupCoordinate.
func (s *AppRequest) GetPickupCoordinate() Coordinate {
	return s.PickupCoordinate
}

// GetDestinationCoordinate returns the value of DestinationCoordinate.
func (s *AppRequest) GetDestinationCoordinate() Coordinate {
	return s.DestinationCoordinate
}

// GetStatus returns the value of Status.
func (s *AppRequest) GetStatus() RequestStatus {
	return s.Status
}

// GetChair returns the value of Chair.
func (s *AppRequest) GetChair() OptChair {
	return s.Chair
}

// GetCreatedAt returns the value of CreatedAt.
func (s *AppRequest) GetCreatedAt() float64 {
	return s.CreatedAt
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *AppRequest) GetUpdatedAt() float64 {
	return s.UpdatedAt
}

// SetRequestID sets the value of RequestID.
func (s *AppRequest) SetRequestID(val string) {
	s.RequestID = val
}

// SetPickupCoordinate sets the value of PickupCoordinate.
func (s *AppRequest) SetPickupCoordinate(val Coordinate) {
	s.PickupCoordinate = val
}

// SetDestinationCoordinate sets the value of DestinationCoordinate.
func (s *AppRequest) SetDestinationCoordinate(val Coordinate) {
	s.DestinationCoordinate = val
}

// SetStatus sets the value of Status.
func (s *AppRequest) SetStatus(val RequestStatus) {
	s.Status = val
}

// SetChair sets the value of Chair.
func (s *AppRequest) SetChair(val OptChair) {
	s.Chair = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *AppRequest) SetCreatedAt(val float64) {
	s.CreatedAt = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *AppRequest) SetUpdatedAt(val float64) {
	s.UpdatedAt = val
}

func (*AppRequest) appGetNotificationRes() {}
func (*AppRequest) appGetRequestRes()      {}

// 簡易椅子情報.
// Ref: #/components/schemas/Chair
type Chair struct {
	// 椅子ID.
	ID string `json:"id"`
	// 椅子名.
	Name string `json:"name"`
	// 車種.
	ChairModel string `json:"chair_model"`
	// カーナンバー.
	ChairNo string `json:"chair_no"`
}

// GetID returns the value of ID.
func (s *Chair) GetID() string {
	return s.ID
}

// GetName returns the value of Name.
func (s *Chair) GetName() string {
	return s.Name
}

// GetChairModel returns the value of ChairModel.
func (s *Chair) GetChairModel() string {
	return s.ChairModel
}

// GetChairNo returns the value of ChairNo.
func (s *Chair) GetChairNo() string {
	return s.ChairNo
}

// SetID sets the value of ID.
func (s *Chair) SetID(val string) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *Chair) SetName(val string) {
	s.Name = val
}

// SetChairModel sets the value of ChairModel.
func (s *Chair) SetChairModel(val string) {
	s.ChairModel = val
}

// SetChairNo sets the value of ChairNo.
func (s *Chair) SetChairNo(val string) {
	s.ChairNo = val
}

type ChairGetInquiriesOK struct {
	Inquiries []ChairGetInquiriesOKInquiriesItem `json:"inquiries"`
}

// GetInquiries returns the value of Inquiries.
func (s *ChairGetInquiriesOK) GetInquiries() []ChairGetInquiriesOKInquiriesItem {
	return s.Inquiries
}

// SetInquiries sets the value of Inquiries.
func (s *ChairGetInquiriesOK) SetInquiries(val []ChairGetInquiriesOKInquiriesItem) {
	s.Inquiries = val
}

type ChairGetInquiriesOKInquiriesItem struct {
	// 問い合わせID.
	ID string `json:"id"`
	// 件名.
	Subject string `json:"subject"`
	// 問い合わせ日時.
	CreatedAt float64 `json:"created_at"`
}

// GetID returns the value of ID.
func (s *ChairGetInquiriesOKInquiriesItem) GetID() string {
	return s.ID
}

// GetSubject returns the value of Subject.
func (s *ChairGetInquiriesOKInquiriesItem) GetSubject() string {
	return s.Subject
}

// GetCreatedAt returns the value of CreatedAt.
func (s *ChairGetInquiriesOKInquiriesItem) GetCreatedAt() float64 {
	return s.CreatedAt
}

// SetID sets the value of ID.
func (s *ChairGetInquiriesOKInquiriesItem) SetID(val string) {
	s.ID = val
}

// SetSubject sets the value of Subject.
func (s *ChairGetInquiriesOKInquiriesItem) SetSubject(val string) {
	s.Subject = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *ChairGetInquiriesOKInquiriesItem) SetCreatedAt(val float64) {
	s.CreatedAt = val
}

// ChairGetNotificationNoContent is response for ChairGetNotification operation.
type ChairGetNotificationNoContent struct{}

func (*ChairGetNotificationNoContent) chairGetNotificationRes() {}

// ChairPostActivateNoContent is response for ChairPostActivate operation.
type ChairPostActivateNoContent struct{}

type ChairPostActivateReq struct{}

// ChairPostCoordinateNoContent is response for ChairPostCoordinate operation.
type ChairPostCoordinateNoContent struct{}

// ChairPostDeactivateNoContent is response for ChairPostDeactivate operation.
type ChairPostDeactivateNoContent struct{}

type ChairPostDeactivateReq struct{}

type ChairPostRegisterCreated struct {
	// アクセストークン.
	AccessToken string `json:"access_token"`
	// 椅子ID.
	ID string `json:"id"`
}

// GetAccessToken returns the value of AccessToken.
func (s *ChairPostRegisterCreated) GetAccessToken() string {
	return s.AccessToken
}

// GetID returns the value of ID.
func (s *ChairPostRegisterCreated) GetID() string {
	return s.ID
}

// SetAccessToken sets the value of AccessToken.
func (s *ChairPostRegisterCreated) SetAccessToken(val string) {
	s.AccessToken = val
}

// SetID sets the value of ID.
func (s *ChairPostRegisterCreated) SetID(val string) {
	s.ID = val
}

type ChairPostRegisterReq struct {
	// 椅子名.
	Username string `json:"username"`
	// 名前.
	Firstname string `json:"firstname"`
	// 名字.
	Lastname string `json:"lastname"`
	// 生年月日.
	DateOfBirth string `json:"date_of_birth"`
	// 車種.
	ChairModel string `json:"chair_model"`
	// カーナンバー.
	ChairNo string `json:"chair_no"`
}

// GetUsername returns the value of Username.
func (s *ChairPostRegisterReq) GetUsername() string {
	return s.Username
}

// GetFirstname returns the value of Firstname.
func (s *ChairPostRegisterReq) GetFirstname() string {
	return s.Firstname
}

// GetLastname returns the value of Lastname.
func (s *ChairPostRegisterReq) GetLastname() string {
	return s.Lastname
}

// GetDateOfBirth returns the value of DateOfBirth.
func (s *ChairPostRegisterReq) GetDateOfBirth() string {
	return s.DateOfBirth
}

// GetChairModel returns the value of ChairModel.
func (s *ChairPostRegisterReq) GetChairModel() string {
	return s.ChairModel
}

// GetChairNo returns the value of ChairNo.
func (s *ChairPostRegisterReq) GetChairNo() string {
	return s.ChairNo
}

// SetUsername sets the value of Username.
func (s *ChairPostRegisterReq) SetUsername(val string) {
	s.Username = val
}

// SetFirstname sets the value of Firstname.
func (s *ChairPostRegisterReq) SetFirstname(val string) {
	s.Firstname = val
}

// SetLastname sets the value of Lastname.
func (s *ChairPostRegisterReq) SetLastname(val string) {
	s.Lastname = val
}

// SetDateOfBirth sets the value of DateOfBirth.
func (s *ChairPostRegisterReq) SetDateOfBirth(val string) {
	s.DateOfBirth = val
}

// SetChairModel sets the value of ChairModel.
func (s *ChairPostRegisterReq) SetChairModel(val string) {
	s.ChairModel = val
}

// SetChairNo sets the value of ChairNo.
func (s *ChairPostRegisterReq) SetChairNo(val string) {
	s.ChairNo = val
}

// ChairPostRequestAcceptNoContent is response for ChairPostRequestAccept operation.
type ChairPostRequestAcceptNoContent struct{}

func (*ChairPostRequestAcceptNoContent) chairPostRequestAcceptRes() {}

// ChairPostRequestDenyNoContent is response for ChairPostRequestDeny operation.
type ChairPostRequestDenyNoContent struct{}

func (*ChairPostRequestDenyNoContent) chairPostRequestDenyRes() {}

type ChairPostRequestDepartBadRequest Error

func (*ChairPostRequestDepartBadRequest) chairPostRequestDepartRes() {}

// ChairPostRequestDepartNoContent is response for ChairPostRequestDepart operation.
type ChairPostRequestDepartNoContent struct{}

func (*ChairPostRequestDepartNoContent) chairPostRequestDepartRes() {}

type ChairPostRequestDepartNotFound Error

func (*ChairPostRequestDepartNotFound) chairPostRequestDepartRes() {}

type ChairPostRequestPaymentBadRequest Error

func (*ChairPostRequestPaymentBadRequest) chairPostRequestPaymentRes() {}

// ChairPostRequestPaymentNoContent is response for ChairPostRequestPayment operation.
type ChairPostRequestPaymentNoContent struct{}

func (*ChairPostRequestPaymentNoContent) chairPostRequestPaymentRes() {}

type ChairPostRequestPaymentNotFound Error

func (*ChairPostRequestPaymentNotFound) chairPostRequestPaymentRes() {}

// Chair向け配車要求情報.
// Ref: #/components/schemas/ChairRequest
type ChairRequest struct {
	// 配車要求ID.
	RequestID             string           `json:"request_id"`
	User                  User             `json:"user"`
	DestinationCoordinate Coordinate       `json:"destination_coordinate"`
	Status                OptRequestStatus `json:"status"`
}

// GetRequestID returns the value of RequestID.
func (s *ChairRequest) GetRequestID() string {
	return s.RequestID
}

// GetUser returns the value of User.
func (s *ChairRequest) GetUser() User {
	return s.User
}

// GetDestinationCoordinate returns the value of DestinationCoordinate.
func (s *ChairRequest) GetDestinationCoordinate() Coordinate {
	return s.DestinationCoordinate
}

// GetStatus returns the value of Status.
func (s *ChairRequest) GetStatus() OptRequestStatus {
	return s.Status
}

// SetRequestID sets the value of RequestID.
func (s *ChairRequest) SetRequestID(val string) {
	s.RequestID = val
}

// SetUser sets the value of User.
func (s *ChairRequest) SetUser(val User) {
	s.User = val
}

// SetDestinationCoordinate sets the value of DestinationCoordinate.
func (s *ChairRequest) SetDestinationCoordinate(val Coordinate) {
	s.DestinationCoordinate = val
}

// SetStatus sets the value of Status.
func (s *ChairRequest) SetStatus(val OptRequestStatus) {
	s.Status = val
}

func (*ChairRequest) chairGetNotificationRes() {}
func (*ChairRequest) chairGetRequestRes()      {}

// 座標情報.
// Ref: #/components/schemas/Coordinate
type Coordinate struct {
	// 経度.
	Latitude int `json:"latitude"`
	// 緯度.
	Longitude int `json:"longitude"`
}

// GetLatitude returns the value of Latitude.
func (s *Coordinate) GetLatitude() int {
	return s.Latitude
}

// GetLongitude returns the value of Longitude.
func (s *Coordinate) GetLongitude() int {
	return s.Longitude
}

// SetLatitude sets the value of Latitude.
func (s *Coordinate) SetLatitude(val int) {
	s.Latitude = val
}

// SetLongitude sets the value of Longitude.
func (s *Coordinate) SetLongitude(val int) {
	s.Longitude = val
}

// Ref: #/components/schemas/Error
type Error struct {
	Message string `json:"message"`
}

// GetMessage returns the value of Message.
func (s *Error) GetMessage() string {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *Error) SetMessage(val string) {
	s.Message = val
}

func (*Error) appGetRequestRes()          {}
func (*Error) appPostRegisterRes()        {}
func (*Error) chairGetInquiryRes()        {}
func (*Error) chairGetRequestRes()        {}
func (*Error) chairPostRequestAcceptRes() {}
func (*Error) chairPostRequestDenyRes()   {}

// 問い合わせ内容.
// Ref: #/components/schemas/InquiryContent
type InquiryContent struct {
	// 問い合わせID.
	ID string `json:"id"`
	// 件名.
	Subject string `json:"subject"`
	// 問い合わせ内容.
	Body string `json:"body"`
	// 問い合わせ日時.
	CreatedAt float64 `json:"created_at"`
}

// GetID returns the value of ID.
func (s *InquiryContent) GetID() string {
	return s.ID
}

// GetSubject returns the value of Subject.
func (s *InquiryContent) GetSubject() string {
	return s.Subject
}

// GetBody returns the value of Body.
func (s *InquiryContent) GetBody() string {
	return s.Body
}

// GetCreatedAt returns the value of CreatedAt.
func (s *InquiryContent) GetCreatedAt() float64 {
	return s.CreatedAt
}

// SetID sets the value of ID.
func (s *InquiryContent) SetID(val string) {
	s.ID = val
}

// SetSubject sets the value of Subject.
func (s *InquiryContent) SetSubject(val string) {
	s.Subject = val
}

// SetBody sets the value of Body.
func (s *InquiryContent) SetBody(val string) {
	s.Body = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *InquiryContent) SetCreatedAt(val float64) {
	s.CreatedAt = val
}

func (*InquiryContent) chairGetInquiryRes() {}

// NewOptAppPostInquiryReq returns new OptAppPostInquiryReq with value set to v.
func NewOptAppPostInquiryReq(v AppPostInquiryReq) OptAppPostInquiryReq {
	return OptAppPostInquiryReq{
		Value: v,
		Set:   true,
	}
}

// OptAppPostInquiryReq is optional AppPostInquiryReq.
type OptAppPostInquiryReq struct {
	Value AppPostInquiryReq
	Set   bool
}

// IsSet returns true if OptAppPostInquiryReq was set.
func (o OptAppPostInquiryReq) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptAppPostInquiryReq) Reset() {
	var v AppPostInquiryReq
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptAppPostInquiryReq) SetTo(v AppPostInquiryReq) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptAppPostInquiryReq) Get() (v AppPostInquiryReq, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptAppPostInquiryReq) Or(d AppPostInquiryReq) AppPostInquiryReq {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptAppPostPaymentMethodsReq returns new OptAppPostPaymentMethodsReq with value set to v.
func NewOptAppPostPaymentMethodsReq(v AppPostPaymentMethodsReq) OptAppPostPaymentMethodsReq {
	return OptAppPostPaymentMethodsReq{
		Value: v,
		Set:   true,
	}
}

// OptAppPostPaymentMethodsReq is optional AppPostPaymentMethodsReq.
type OptAppPostPaymentMethodsReq struct {
	Value AppPostPaymentMethodsReq
	Set   bool
}

// IsSet returns true if OptAppPostPaymentMethodsReq was set.
func (o OptAppPostPaymentMethodsReq) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptAppPostPaymentMethodsReq) Reset() {
	var v AppPostPaymentMethodsReq
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptAppPostPaymentMethodsReq) SetTo(v AppPostPaymentMethodsReq) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptAppPostPaymentMethodsReq) Get() (v AppPostPaymentMethodsReq, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptAppPostPaymentMethodsReq) Or(d AppPostPaymentMethodsReq) AppPostPaymentMethodsReq {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptAppPostRegisterReq returns new OptAppPostRegisterReq with value set to v.
func NewOptAppPostRegisterReq(v AppPostRegisterReq) OptAppPostRegisterReq {
	return OptAppPostRegisterReq{
		Value: v,
		Set:   true,
	}
}

// OptAppPostRegisterReq is optional AppPostRegisterReq.
type OptAppPostRegisterReq struct {
	Value AppPostRegisterReq
	Set   bool
}

// IsSet returns true if OptAppPostRegisterReq was set.
func (o OptAppPostRegisterReq) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptAppPostRegisterReq) Reset() {
	var v AppPostRegisterReq
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptAppPostRegisterReq) SetTo(v AppPostRegisterReq) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptAppPostRegisterReq) Get() (v AppPostRegisterReq, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptAppPostRegisterReq) Or(d AppPostRegisterReq) AppPostRegisterReq {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptAppPostRequestEvaluateReq returns new OptAppPostRequestEvaluateReq with value set to v.
func NewOptAppPostRequestEvaluateReq(v AppPostRequestEvaluateReq) OptAppPostRequestEvaluateReq {
	return OptAppPostRequestEvaluateReq{
		Value: v,
		Set:   true,
	}
}

// OptAppPostRequestEvaluateReq is optional AppPostRequestEvaluateReq.
type OptAppPostRequestEvaluateReq struct {
	Value AppPostRequestEvaluateReq
	Set   bool
}

// IsSet returns true if OptAppPostRequestEvaluateReq was set.
func (o OptAppPostRequestEvaluateReq) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptAppPostRequestEvaluateReq) Reset() {
	var v AppPostRequestEvaluateReq
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptAppPostRequestEvaluateReq) SetTo(v AppPostRequestEvaluateReq) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptAppPostRequestEvaluateReq) Get() (v AppPostRequestEvaluateReq, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptAppPostRequestEvaluateReq) Or(d AppPostRequestEvaluateReq) AppPostRequestEvaluateReq {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptAppPostRequestReq returns new OptAppPostRequestReq with value set to v.
func NewOptAppPostRequestReq(v AppPostRequestReq) OptAppPostRequestReq {
	return OptAppPostRequestReq{
		Value: v,
		Set:   true,
	}
}

// OptAppPostRequestReq is optional AppPostRequestReq.
type OptAppPostRequestReq struct {
	Value AppPostRequestReq
	Set   bool
}

// IsSet returns true if OptAppPostRequestReq was set.
func (o OptAppPostRequestReq) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptAppPostRequestReq) Reset() {
	var v AppPostRequestReq
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptAppPostRequestReq) SetTo(v AppPostRequestReq) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptAppPostRequestReq) Get() (v AppPostRequestReq, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptAppPostRequestReq) Or(d AppPostRequestReq) AppPostRequestReq {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptChair returns new OptChair with value set to v.
func NewOptChair(v Chair) OptChair {
	return OptChair{
		Value: v,
		Set:   true,
	}
}

// OptChair is optional Chair.
type OptChair struct {
	Value Chair
	Set   bool
}

// IsSet returns true if OptChair was set.
func (o OptChair) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptChair) Reset() {
	var v Chair
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptChair) SetTo(v Chair) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptChair) Get() (v Chair, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptChair) Or(d Chair) Chair {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptChairPostRegisterReq returns new OptChairPostRegisterReq with value set to v.
func NewOptChairPostRegisterReq(v ChairPostRegisterReq) OptChairPostRegisterReq {
	return OptChairPostRegisterReq{
		Value: v,
		Set:   true,
	}
}

// OptChairPostRegisterReq is optional ChairPostRegisterReq.
type OptChairPostRegisterReq struct {
	Value ChairPostRegisterReq
	Set   bool
}

// IsSet returns true if OptChairPostRegisterReq was set.
func (o OptChairPostRegisterReq) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptChairPostRegisterReq) Reset() {
	var v ChairPostRegisterReq
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptChairPostRegisterReq) SetTo(v ChairPostRegisterReq) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptChairPostRegisterReq) Get() (v ChairPostRegisterReq, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptChairPostRegisterReq) Or(d ChairPostRegisterReq) ChairPostRegisterReq {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptCoordinate returns new OptCoordinate with value set to v.
func NewOptCoordinate(v Coordinate) OptCoordinate {
	return OptCoordinate{
		Value: v,
		Set:   true,
	}
}

// OptCoordinate is optional Coordinate.
type OptCoordinate struct {
	Value Coordinate
	Set   bool
}

// IsSet returns true if OptCoordinate was set.
func (o OptCoordinate) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptCoordinate) Reset() {
	var v Coordinate
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptCoordinate) SetTo(v Coordinate) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptCoordinate) Get() (v Coordinate, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptCoordinate) Or(d Coordinate) Coordinate {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptFloat64 returns new OptFloat64 with value set to v.
func NewOptFloat64(v float64) OptFloat64 {
	return OptFloat64{
		Value: v,
		Set:   true,
	}
}

// OptFloat64 is optional float64.
type OptFloat64 struct {
	Value float64
	Set   bool
}

// IsSet returns true if OptFloat64 was set.
func (o OptFloat64) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptFloat64) Reset() {
	var v float64
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptFloat64) SetTo(v float64) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptFloat64) Get() (v float64, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptFloat64) Or(d float64) float64 {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptPostInitializeReq returns new OptPostInitializeReq with value set to v.
func NewOptPostInitializeReq(v PostInitializeReq) OptPostInitializeReq {
	return OptPostInitializeReq{
		Value: v,
		Set:   true,
	}
}

// OptPostInitializeReq is optional PostInitializeReq.
type OptPostInitializeReq struct {
	Value PostInitializeReq
	Set   bool
}

// IsSet returns true if OptPostInitializeReq was set.
func (o OptPostInitializeReq) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptPostInitializeReq) Reset() {
	var v PostInitializeReq
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptPostInitializeReq) SetTo(v PostInitializeReq) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptPostInitializeReq) Get() (v PostInitializeReq, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptPostInitializeReq) Or(d PostInitializeReq) PostInitializeReq {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptRequestStatus returns new OptRequestStatus with value set to v.
func NewOptRequestStatus(v RequestStatus) OptRequestStatus {
	return OptRequestStatus{
		Value: v,
		Set:   true,
	}
}

// OptRequestStatus is optional RequestStatus.
type OptRequestStatus struct {
	Value RequestStatus
	Set   bool
}

// IsSet returns true if OptRequestStatus was set.
func (o OptRequestStatus) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptRequestStatus) Reset() {
	var v RequestStatus
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptRequestStatus) SetTo(v RequestStatus) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptRequestStatus) Get() (v RequestStatus, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptRequestStatus) Or(d RequestStatus) RequestStatus {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

type PostInitializeOK struct {
	// 実装言語.
	Language PostInitializeOKLanguage `json:"language"`
}

// GetLanguage returns the value of Language.
func (s *PostInitializeOK) GetLanguage() PostInitializeOKLanguage {
	return s.Language
}

// SetLanguage sets the value of Language.
func (s *PostInitializeOK) SetLanguage(val PostInitializeOKLanguage) {
	s.Language = val
}

// 実装言語.
type PostInitializeOKLanguage string

const (
	PostInitializeOKLanguageGo     PostInitializeOKLanguage = "go"
	PostInitializeOKLanguagePerl   PostInitializeOKLanguage = "perl"
	PostInitializeOKLanguagePhp    PostInitializeOKLanguage = "php"
	PostInitializeOKLanguagePython PostInitializeOKLanguage = "python"
	PostInitializeOKLanguageRuby   PostInitializeOKLanguage = "ruby"
	PostInitializeOKLanguageRust   PostInitializeOKLanguage = "rust"
	PostInitializeOKLanguageNode   PostInitializeOKLanguage = "node"
)

// AllValues returns all PostInitializeOKLanguage values.
func (PostInitializeOKLanguage) AllValues() []PostInitializeOKLanguage {
	return []PostInitializeOKLanguage{
		PostInitializeOKLanguageGo,
		PostInitializeOKLanguagePerl,
		PostInitializeOKLanguagePhp,
		PostInitializeOKLanguagePython,
		PostInitializeOKLanguageRuby,
		PostInitializeOKLanguageRust,
		PostInitializeOKLanguageNode,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s PostInitializeOKLanguage) MarshalText() ([]byte, error) {
	switch s {
	case PostInitializeOKLanguageGo:
		return []byte(s), nil
	case PostInitializeOKLanguagePerl:
		return []byte(s), nil
	case PostInitializeOKLanguagePhp:
		return []byte(s), nil
	case PostInitializeOKLanguagePython:
		return []byte(s), nil
	case PostInitializeOKLanguageRuby:
		return []byte(s), nil
	case PostInitializeOKLanguageRust:
		return []byte(s), nil
	case PostInitializeOKLanguageNode:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *PostInitializeOKLanguage) UnmarshalText(data []byte) error {
	switch PostInitializeOKLanguage(data) {
	case PostInitializeOKLanguageGo:
		*s = PostInitializeOKLanguageGo
		return nil
	case PostInitializeOKLanguagePerl:
		*s = PostInitializeOKLanguagePerl
		return nil
	case PostInitializeOKLanguagePhp:
		*s = PostInitializeOKLanguagePhp
		return nil
	case PostInitializeOKLanguagePython:
		*s = PostInitializeOKLanguagePython
		return nil
	case PostInitializeOKLanguageRuby:
		*s = PostInitializeOKLanguageRuby
		return nil
	case PostInitializeOKLanguageRust:
		*s = PostInitializeOKLanguageRust
		return nil
	case PostInitializeOKLanguageNode:
		*s = PostInitializeOKLanguageNode
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

type PostInitializeReq struct {
	// 決済サーバーアドレス.
	PaymentServer string `json:"payment_server"`
}

// GetPaymentServer returns the value of PaymentServer.
func (s *PostInitializeReq) GetPaymentServer() string {
	return s.PaymentServer
}

// SetPaymentServer sets the value of PaymentServer.
func (s *PostInitializeReq) SetPaymentServer(val string) {
	s.PaymentServer = val
}

// 配車要求ステータス
// MATCHING:
// サービス上でマッチング処理を行なっていて椅子が確定していない
// DISPATCHING: 椅子が確定し、乗車位置に向かっている
// DISPATCHED: 椅子が乗車位置に到着して、ユーザーの乗車を待機している
// CARRYING: ユーザーが乗車し、椅子が目的地に向かっている
// ARRIVED: 目的地に到着した
// COMPLETED: ユーザーの決済・椅子評価が完了した
// CANCELED:
// 何らかの理由により途中でキャンセルされた(一定時間待ったが椅子を割り当てられなかった場合などを想定).
// Ref: #/components/schemas/RequestStatus
type RequestStatus string

const (
	RequestStatusMATCHING    RequestStatus = "MATCHING"
	RequestStatusDISPATCHING RequestStatus = "DISPATCHING"
	RequestStatusDISPATCHED  RequestStatus = "DISPATCHED"
	RequestStatusCARRYING    RequestStatus = "CARRYING"
	RequestStatusARRIVED     RequestStatus = "ARRIVED"
	RequestStatusCOMPLETED   RequestStatus = "COMPLETED"
	RequestStatusCANCELED    RequestStatus = "CANCELED"
)

// AllValues returns all RequestStatus values.
func (RequestStatus) AllValues() []RequestStatus {
	return []RequestStatus{
		RequestStatusMATCHING,
		RequestStatusDISPATCHING,
		RequestStatusDISPATCHED,
		RequestStatusCARRYING,
		RequestStatusARRIVED,
		RequestStatusCOMPLETED,
		RequestStatusCANCELED,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s RequestStatus) MarshalText() ([]byte, error) {
	switch s {
	case RequestStatusMATCHING:
		return []byte(s), nil
	case RequestStatusDISPATCHING:
		return []byte(s), nil
	case RequestStatusDISPATCHED:
		return []byte(s), nil
	case RequestStatusCARRYING:
		return []byte(s), nil
	case RequestStatusARRIVED:
		return []byte(s), nil
	case RequestStatusCOMPLETED:
		return []byte(s), nil
	case RequestStatusCANCELED:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *RequestStatus) UnmarshalText(data []byte) error {
	switch RequestStatus(data) {
	case RequestStatusMATCHING:
		*s = RequestStatusMATCHING
		return nil
	case RequestStatusDISPATCHING:
		*s = RequestStatusDISPATCHING
		return nil
	case RequestStatusDISPATCHED:
		*s = RequestStatusDISPATCHED
		return nil
	case RequestStatusCARRYING:
		*s = RequestStatusCARRYING
		return nil
	case RequestStatusARRIVED:
		*s = RequestStatusARRIVED
		return nil
	case RequestStatusCOMPLETED:
		*s = RequestStatusCOMPLETED
		return nil
	case RequestStatusCANCELED:
		*s = RequestStatusCANCELED
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// 簡易ユーザー情報.
// Ref: #/components/schemas/User
type User struct {
	// ユーザーID.
	ID string `json:"id"`
	// ユーザー名.
	Name string `json:"name"`
}

// GetID returns the value of ID.
func (s *User) GetID() string {
	return s.ID
}

// GetName returns the value of Name.
func (s *User) GetName() string {
	return s.Name
}

// SetID sets the value of ID.
func (s *User) SetID(val string) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *User) SetName(val string) {
	s.Name = val
}
