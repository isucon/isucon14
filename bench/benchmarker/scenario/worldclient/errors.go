package worldclient

import (
	"errors"
	"fmt"
)

type ErrorCode int

const (
	// ErrorCodeFailedToPostCoordinate 座標送信に失敗したエラー
	ErrorCodeFailedToPostCoordinate = iota + 10000
	// ErrorCodeFailedToPostAccept リクエスト受諾に失敗したエラー
	ErrorCodeFailedToPostAccept
	// ErrorCodeFailedToPostDeny リクエスト拒否に失敗したエラー
	ErrorCodeFailedToPostDeny
	// ErrorCodeFailedToPostDepart 出発通知に失敗したエラー
	ErrorCodeFailedToPostDepart
	// ErrorCodeFailedToPostEvaluate 評価送信に失敗したエラー
	ErrorCodeFailedToPostEvaluate
	// ErrorCodeFailedToPostActivate 配車受付の開始に失敗したエラー
	ErrorCodeFailedToPostActivate
	// ErrorCodeFailedToPostDeactivate 配車受付の停止に失敗したエラー
	ErrorCodeFailedToPostDeactivate
	// ErrorCodeFailedToGetChairRequest 椅子のリクエスト取得に失敗したエラー
	ErrorCodeFailedToGetChairRequest
	// ErrorCodeFailedToCreateWebappClient WebappClientの作成に失敗したエラー
	ErrorCodeFailedToCreateWebappClient
	// ErrorCodeFailedToRegisterUser ユーザー登録に失敗したエラー
	ErrorCodeFailedToRegisterUser
	// ErrorCodeFailedToRegisterProvider プロバイダー登録に失敗したエラー
	ErrorCodeFailedToRegisterProvider
	// ErrorCodeFailedToRegisterChair 椅子の登録に失敗したエラー
	ErrorCodeFailedToRegisterChair
	// ErrorCodeFailedToPostRequest リクエスト送信に失敗したエラー
	ErrorCodeFailedToPostRequest
	// ErrorCodeFailedToGetRequests リクエスト一覧の取得に失敗したエラー
	ErrorCodeFailedToGetRequests
	// ErrorCodeFailedToPostPaymentMethods ユーザー支払い情報登録に失敗したエラー
	ErrorCodeFailedToPostPaymentMethods
	// ErrorCodeFailedToGetProviderSales プロバイダーの売り上げ情報の取得に失敗したエラー
	ErrorCodeFailedToGetProviderSales
	// ErrorCodeFailedToPostRidesEstimatedFare ライド料金の見積もりの取得に失敗したエラー
	ErrorCodeFailedToPostRidesEstimatedFare
	// ErrorCodeFailedToGetNearbyChairs 近くの椅子情報の取得に失敗したエラー
	ErrorCodeFailedToGetNearbyChairs
)

type codeError struct {
	code ErrorCode
	err  error
}

func (e *codeError) Error() string {
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
