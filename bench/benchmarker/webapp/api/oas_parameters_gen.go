// Code generated by ogen, DO NOT EDIT.

package api

// AppGetNearbyChairsParams is parameters of app-get-nearby-chairs operation.
type AppGetNearbyChairsParams struct {
	// 緯度.
	Latitude int
	// 経度.
	Longitude int
	// 検索距離.
	Distance OptInt
}

// AppGetRequestParams is parameters of app-get-request operation.
type AppGetRequestParams struct {
	// 配車要求ID.
	RequestID string
}

// AppPostRequestEvaluateParams is parameters of app-post-request-evaluate operation.
type AppPostRequestEvaluateParams struct {
	// 配車要求ID.
	RequestID string
}

// ChairGetRequestParams is parameters of chair-get-request operation.
type ChairGetRequestParams struct {
	// 配車要求ID.
	RequestID string
}

// ChairPostRequestAcceptParams is parameters of chair-post-request-accept operation.
type ChairPostRequestAcceptParams struct {
	// 配車要求ID.
	RequestID string
}

// ChairPostRequestDenyParams is parameters of chair-post-request-deny operation.
type ChairPostRequestDenyParams struct {
	// 配車要求ID.
	RequestID string
}

// ChairPostRequestDepartParams is parameters of chair-post-request-depart operation.
type ChairPostRequestDepartParams struct {
	// 配車要求ID.
	RequestID string
}

// OwnerGetChairDetailParams is parameters of owner-get-chair-detail operation.
type OwnerGetChairDetailParams struct {
	// 椅子ID.
	ChairID string
}

// OwnerGetSalesParams is parameters of owner-get-sales operation.
type OwnerGetSalesParams struct {
	// 開始日時（含む）.
	Since OptInt64
	// 終了日時（含む）.
	Until OptInt64
}
