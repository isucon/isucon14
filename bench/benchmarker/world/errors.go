package world

import (
	"errors"
	"fmt"
	"maps"
	"sync"
)

const (
	ErrorLimit = 200
)

type ErrorCode int

const (
	// ErrorCodeUnknown 不明なエラー
	ErrorCodeUnknown ErrorCode = iota
	// ErrorCodeFailedToSendChairCoordinate 椅子の座標送信リクエストが失敗した
	ErrorCodeFailedToSendChairCoordinate
	// ErrorCodeFailedToDepart 椅子が出発しようとしたが、departリクエストが失敗した
	ErrorCodeFailedToDepart
	// ErrorCodeFailedToAcceptRequest 椅子がリクエストを受理しようとしたが失敗した
	ErrorCodeFailedToAcceptRequest
	_
	// ErrorCodeFailedToEvaluate ユーザーが送迎の評価をしようとしたが失敗した
	ErrorCodeFailedToEvaluate
	// ErrorCodeEvaluateTimeout ユーザーが送迎の評価をしようとしたがタイムアウトした
	ErrorCodeEvaluateTimeout
	// ErrorCodeFailedToCheckRequestHistory ユーザーがリクエスト履歴を確認しようとしたが失敗した
	ErrorCodeFailedToCheckRequestHistory
	// ErrorCodeFailedToCreateRequest ユーザーがリクエストを作成しようとしたが失敗した
	ErrorCodeFailedToCreateRequest
	// ErrorCodeUserNotRequestingButStatusChanged リクエストしていないユーザーのリクエストステータスが更新された
	ErrorCodeUserNotRequestingButStatusChanged
	// ErrorCodeChairNotAssignedButStatusChanged 椅子にリクエストが割り当てられていないのに、椅子のステータスが更新された
	ErrorCodeChairNotAssignedButStatusChanged
	// ErrorCodeUnexpectedUserRequestStatusTransitionOccurred 想定されていないユーザーのRequestStatusの遷移が発生した
	ErrorCodeUnexpectedUserRequestStatusTransitionOccurred
	// ErrorCodeUnexpectedChairRequestStatusTransitionOccurred 想定されていない椅子のRequestStatusの遷移が発生した
	ErrorCodeUnexpectedChairRequestStatusTransitionOccurred
	// ErrorCodeFailedToActivate 椅子がリクエストの受付を開始しようとしたが失敗した
	ErrorCodeFailedToActivate
	_
	// ErrorCodeChairAlreadyHasRequest 既にリクエストが割り当てられている椅子に、別のリクエストが割り当てられた
	ErrorCodeChairAlreadyHasRequest
	_
	// ErrorCodeFailedToRegisterUser ユーザー登録に失敗した
	ErrorCodeFailedToRegisterUser
	// ErrorCodeFailedToRegisterOwner オーナー登録に失敗した
	ErrorCodeFailedToRegisterOwner
	// ErrorCodeFailedToRegisterChair 椅子登録に失敗した
	ErrorCodeFailedToRegisterChair
	// ErrorCodeFailedToConnectNotificationStream 通知ストリームへの接続に失敗した
	ErrorCodeFailedToConnectNotificationStream
	// ErrorCodeFailedToRegisterPaymentMethods ユーザーの支払い情報の登録に失敗した
	ErrorCodeFailedToRegisterPaymentMethods
	// ErrorCodeFailedToGetOwnerSales オーナーの売り上げ情報の取得に失敗した
	ErrorCodeFailedToGetOwnerSales
	// ErrorCodeIncorrectAmountOfFareCharged ユーザーのリクエストに対して誤った金額が請求されました
	ErrorCodeIncorrectAmountOfFareCharged
	// ErrorCodeSalesMismatched 取得したオーナーの売り上げ情報が想定しているものとズレています
	ErrorCodeSalesMismatched
	// ErrorCodeFailedToGetOwnerChairs オーナーの椅子一覧の取得に失敗した
	ErrorCodeFailedToGetOwnerChairs
	// ErrorCodeIncorrectOwnerChairsData 取得したオーナーの椅子一覧の情報が合ってない
	ErrorCodeIncorrectOwnerChairsData
	// ErrorCodeTooOldNearbyChairsResponse 取得した付近の椅子情報が古すぎます
	ErrorCodeTooOldNearbyChairsResponse
	// ErrorCodeUncontrollableRequestReceived アサインされたライドがベンチマーカー外で作られたものであるため、処理できない
	ErrorCodeUncontrollableRequestReceived
	// ErrorCodeChairReceivedDataIsWrong 椅子が通知から受け取ったデータが想定と異なります
	ErrorCodeChairReceivedDataIsWrong
	// ErrorCodeWrongNearbyChairs 取得した付近の椅子情報に不備があります
	ErrorCodeWrongNearbyChairs
	// ErrorCodeLackOfNearbyChairs 取得した付近の椅子情報が足りません
	ErrorCodeLackOfNearbyChairs
	// ErrorCodeMatchingTimeout マッチングに時間がかかりすぎです
	ErrorCodeMatchingTimeout
	// ErrorCodeUserReceivedDataIsWrong ユーザーが通知から受け取ったデータが想定と異なります
	ErrorCodeUserReceivedDataIsWrong
)

var CriticalErrorCodes = map[ErrorCode]bool{
	ErrorCodeUserNotRequestingButStatusChanged:              true,
	ErrorCodeChairNotAssignedButStatusChanged:               true,
	ErrorCodeUnexpectedUserRequestStatusTransitionOccurred:  true,
	ErrorCodeUnexpectedChairRequestStatusTransitionOccurred: true,
	ErrorCodeChairAlreadyHasRequest:                         true,
	ErrorCodeIncorrectAmountOfFareCharged:                   true,
	ErrorCodeUncontrollableRequestReceived:                  true,
	ErrorCodeMatchingTimeout:                                true,
	ErrorCodeEvaluateTimeout:                                true,
}

var ErrorTexts = map[ErrorCode]string{
	ErrorCodeFailedToSendChairCoordinate:                    "椅子の座標送信に失敗しました",
	ErrorCodeFailedToDepart:                                 "椅子が出発できませんでした",
	ErrorCodeFailedToAcceptRequest:                          "椅子がライドを受理できませんでした",
	ErrorCodeFailedToEvaluate:                               "ユーザーのライド評価に失敗しました",
	ErrorCodeEvaluateTimeout:                                "ユーザーのライド評価がタイムアウトしました",
	ErrorCodeFailedToCheckRequestHistory:                    "ユーザーがライド履歴の取得に失敗しました",
	ErrorCodeFailedToCreateRequest:                          "ユーザーが新しくライドを作成できませんでした",
	ErrorCodeUserNotRequestingButStatusChanged:              "ユーザーが想定していない通知を受け取りました",
	ErrorCodeChairNotAssignedButStatusChanged:               "椅子が想定していない通知を受け取りました",
	ErrorCodeUnexpectedUserRequestStatusTransitionOccurred:  "ユーザーに想定していないライドの状態遷移の通知がありました",
	ErrorCodeUnexpectedChairRequestStatusTransitionOccurred: "椅子に想定していないライドの状態遷移の通知がありました",
	ErrorCodeFailedToActivate:                               "椅子がアクティベートに失敗しました",
	ErrorCodeChairAlreadyHasRequest:                         "椅子がライドの完了通知を受け取る前に、別の新しいライドの通知を受け取りました",
	ErrorCodeFailedToRegisterUser:                           "ユーザー登録に失敗しました",
	ErrorCodeFailedToRegisterOwner:                          "オーナー登録に失敗しました",
	ErrorCodeFailedToRegisterChair:                          "椅子登録に失敗しました",
	ErrorCodeFailedToConnectNotificationStream:              "通知APIの接続に失敗しました",
	ErrorCodeFailedToRegisterPaymentMethods:                 "ユーザーの支払い情報の登録に失敗しました",
	ErrorCodeFailedToGetOwnerSales:                          "オーナーの売り上げ情報の取得に失敗しました",
	ErrorCodeIncorrectAmountOfFareCharged:                   "ユーザーに誤った金額が請求されました",
	ErrorCodeSalesMismatched:                                "取得したオーナーの売り上げ情報が想定しているものと異なります",
	ErrorCodeFailedToGetOwnerChairs:                         "オーナーの椅子一覧の取得に失敗しました",
	ErrorCodeIncorrectOwnerChairsData:                       "取得したオーナーの椅子一覧の情報が想定しているものと異なります",
	ErrorCodeTooOldNearbyChairsResponse:                     "取得した付近の椅子情報が古すぎます",
	ErrorCodeUncontrollableRequestReceived:                  "アサインされたライドがベンチマーカー外で作られたものであるため処理できません",
	ErrorCodeChairReceivedDataIsWrong:                       "椅子が受け取った通知の内容が想定と異なります",
	ErrorCodeWrongNearbyChairs:                              "取得した付近の椅子情報に不備があります",
	ErrorCodeLackOfNearbyChairs:                             "付近の椅子情報が想定よりも足りていません",
	ErrorCodeMatchingTimeout:                                "ライドが長時間マッチングされませんでした",
	ErrorCodeUserReceivedDataIsWrong:                        "ユーザーが受け取った通知の内容が想定と異なります",
}

type codeError struct {
	code ErrorCode
	err  error
}

func (e *codeError) Error() string {
	text, ok := ErrorTexts[e.code]
	if ok {
		if e.err == nil {
			return fmt.Sprintf("%s (CODE=%d)", text, e.code)
		}
		return fmt.Sprintf("%s (CODE=%d): %s", text, e.code, e.err)
	}
	if e.err == nil {
		return fmt.Sprintf("CODE=%d", e.code)
	}
	return fmt.Sprintf("CODE=%d: %s", e.code, e.err)
}

func (e *codeError) Unwrap() error {
	return e.err
}

func (e *codeError) Code() ErrorCode {
	return e.code
}

func (e *codeError) Is(target error) bool {
	var t *codeError
	if errors.As(target, &t) {
		return t.code == e.code
	}
	return false
}

func WrapCodeError(code ErrorCode, err error) error {
	return &codeError{code, err}
}

func CodeError(code ErrorCode) error {
	return &codeError{code, nil}
}

func IsCriticalError(err error) bool {
	return CriticalErrorCodes[GetErrorCode(err)]
}

func GetErrorCode(err error) ErrorCode {
	var t *codeError
	if errors.As(err, &t) {
		return t.code
	}
	return ErrorCodeUnknown
}

type ErrorCounter struct {
	counter map[ErrorCode]int
	total   int
	m       sync.Mutex
}

func NewErrorCounter() *ErrorCounter {
	return &ErrorCounter{
		counter: make(map[ErrorCode]int),
	}
}

func (c *ErrorCounter) Add(err error) error {
	c.m.Lock()
	defer c.m.Unlock()
	c.total++
	c.counter[GetErrorCode(err)]++
	if c.total > ErrorLimit {
		return errors.New("発生しているエラーが多すぎます")
	}
	return nil
}

func (c *ErrorCounter) Total() int {
	c.m.Lock()
	defer c.m.Unlock()
	return c.total
}

func (c *ErrorCounter) Count() map[ErrorCode]int {
	c.m.Lock()
	defer c.m.Unlock()
	return maps.Clone(c.counter)
}
