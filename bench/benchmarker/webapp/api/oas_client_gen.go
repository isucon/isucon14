// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/uri"
)

// Invoker invokes operations described by OpenAPI v3 specification.
type Invoker interface {
	// AppGetNotification invokes app-get-notification operation.
	//
	// ポーリング方式にしない場合に、ユーザーのアプリに配車要求の各種状態遷移を通知するなどに使う想定.
	//
	// GET /app/notification
	AppGetNotification(ctx context.Context) error
	// AppGetRequest invokes app-get-request operation.
	//
	// ユーザーが配車要求の状態を確認する.
	//
	// GET /app/requests/{request_id}
	AppGetRequest(ctx context.Context, params AppGetRequestParams) (AppGetRequestRes, error)
	// AppPostInquiry invokes app-post-inquiry operation.
	//
	// ユーザーが問い合わせを送信する.
	//
	// POST /app/inquiry
	AppPostInquiry(ctx context.Context, request OptAppPostInquiryReq) error
	// AppPostRegister invokes app-post-register operation.
	//
	// ユーザーが会員登録を行う.
	//
	// POST /app/register
	AppPostRegister(ctx context.Context, request OptAppPostRegisterReq) (AppPostRegisterRes, error)
	// AppPostRequest invokes app-post-request operation.
	//
	// ユーザーが配車要求を行う.
	//
	// POST /app/requests
	AppPostRequest(ctx context.Context, request OptAppPostRequestReq) (*AppPostRequestAccepted, error)
	// AppPostRequestEvaluate invokes app-post-request-evaluate operation.
	//
	// ユーザーが椅子を評価する.
	//
	// POST /app/requests/{request_id}/evaluate
	AppPostRequestEvaluate(ctx context.Context, request OptAppPostRequestEvaluateReq, params AppPostRequestEvaluateParams) (AppPostRequestEvaluateRes, error)
	// ChairGetInquiries invokes chair-get-inquiries operation.
	//
	// 問い合わせの一覧を取得する.
	//
	// GET /admin/inquiries
	ChairGetInquiries(ctx context.Context, params ChairGetInquiriesParams) (*ChairGetInquiriesOK, error)
	// ChairGetInquiry invokes chair-get-inquiry operation.
	//
	// 指定したIDの問い合わせ内容を取得.
	//
	// GET /admin/inquiries/{inquiry_id}
	ChairGetInquiry(ctx context.Context, params ChairGetInquiryParams) (ChairGetInquiryRes, error)
	// ChairGetNotification invokes chair-get-notification operation.
	//
	// 椅子に配車要求を通知するなどで使う想定.
	//
	// GET /chair/notification
	ChairGetNotification(ctx context.Context) error
	// ChairGetRequest invokes chair-get-request operation.
	//
	// 椅子向け通知エンドポイントから通知されたidの情報を取得する想定.
	//
	// GET /chair/requests/{request_id}
	ChairGetRequest(ctx context.Context, params ChairGetRequestParams) (ChairGetRequestRes, error)
	// ChairPostActivate invokes chair-post-activate operation.
	//
	// 椅子が配車受付を開始する.
	//
	// POST /chair/activate
	ChairPostActivate(ctx context.Context, request *ChairPostActivateReq) error
	// ChairPostCoordinate invokes chair-post-coordinate operation.
	//
	// 椅子が位置情報を送信する.
	//
	// POST /chair/coordinate
	ChairPostCoordinate(ctx context.Context, request OptCoordinate) error
	// ChairPostDeactivate invokes chair-post-deactivate operation.
	//
	// 椅子が配車受付を停止する.
	//
	// POST /chair/deactivate
	ChairPostDeactivate(ctx context.Context, request *ChairPostDeactivateReq) error
	// ChairPostRegister invokes chair-post-register operation.
	//
	// 椅子登録を行う.
	//
	// POST /chair/register
	ChairPostRegister(ctx context.Context, request OptChairPostRegisterReq) (*ChairPostRegisterCreated, error)
	// ChairPostRequestAccept invokes chair-post-request-accept operation.
	//
	// 椅子が配車要求を受理する.
	//
	// POST /chair/requests/{request_id}/accept
	ChairPostRequestAccept(ctx context.Context, params ChairPostRequestAcceptParams) (ChairPostRequestAcceptRes, error)
	// ChairPostRequestDeny invokes chair-post-request-deny operation.
	//
	// 椅子が配車要求を拒否する.
	//
	// POST /chair/requests/{request_id}/deny
	ChairPostRequestDeny(ctx context.Context, params ChairPostRequestDenyParams) (ChairPostRequestDenyRes, error)
	// ChairPostRequestDepart invokes chair-post-request-depart operation.
	//
	// 椅子が配車位置から出発する(ユーザーが乗車完了した).
	//
	// POST /chair/requests/{request_id}/depart
	ChairPostRequestDepart(ctx context.Context, params ChairPostRequestDepartParams) (ChairPostRequestDepartRes, error)
	// Initialize invokes initialize operation.
	//
	// サービスを初期化する.
	//
	// POST /initialize
	Initialize(ctx context.Context) (*InitializeOK, error)
}

// Client implements OAS client.
type Client struct {
	serverURL *url.URL
	baseClient
}

func trimTrailingSlashes(u *url.URL) {
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawPath = strings.TrimRight(u.RawPath, "/")
}

// NewClient initializes new Client defined by OAS.
func NewClient(serverURL string, opts ...ClientOption) (*Client, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	trimTrailingSlashes(u)

	c, err := newClientConfig(opts...).baseClient()
	if err != nil {
		return nil, err
	}
	return &Client{
		serverURL:  u,
		baseClient: c,
	}, nil
}

type serverURLKey struct{}

// WithServerURL sets context key to override server URL.
func WithServerURL(ctx context.Context, u *url.URL) context.Context {
	return context.WithValue(ctx, serverURLKey{}, u)
}

func (c *Client) requestURL(ctx context.Context) *url.URL {
	u, ok := ctx.Value(serverURLKey{}).(*url.URL)
	if !ok {
		return c.serverURL
	}
	return u
}

// AppGetNotification invokes app-get-notification operation.
//
// ポーリング方式にしない場合に、ユーザーのアプリに配車要求の各種状態遷移を通知するなどに使う想定.
//
// GET /app/notification
func (c *Client) AppGetNotification(ctx context.Context) error {
	_, err := c.sendAppGetNotification(ctx)
	return err
}

func (c *Client) sendAppGetNotification(ctx context.Context) (res *AppGetNotificationOK, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/app/notification"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeAppGetNotificationResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// AppGetRequest invokes app-get-request operation.
//
// ユーザーが配車要求の状態を確認する.
//
// GET /app/requests/{request_id}
func (c *Client) AppGetRequest(ctx context.Context, params AppGetRequestParams) (AppGetRequestRes, error) {
	res, err := c.sendAppGetRequest(ctx, params)
	return res, err
}

func (c *Client) sendAppGetRequest(ctx context.Context, params AppGetRequestParams) (res AppGetRequestRes, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [2]string
	pathParts[0] = "/app/requests/"
	{
		// Encode "request_id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "request_id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.StringToString(params.RequestID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeAppGetRequestResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// AppPostInquiry invokes app-post-inquiry operation.
//
// ユーザーが問い合わせを送信する.
//
// POST /app/inquiry
func (c *Client) AppPostInquiry(ctx context.Context, request OptAppPostInquiryReq) error {
	_, err := c.sendAppPostInquiry(ctx, request)
	return err
}

func (c *Client) sendAppPostInquiry(ctx context.Context, request OptAppPostInquiryReq) (res *AppPostInquiryNoContent, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/app/inquiry"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeAppPostInquiryRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeAppPostInquiryResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// AppPostRegister invokes app-post-register operation.
//
// ユーザーが会員登録を行う.
//
// POST /app/register
func (c *Client) AppPostRegister(ctx context.Context, request OptAppPostRegisterReq) (AppPostRegisterRes, error) {
	res, err := c.sendAppPostRegister(ctx, request)
	return res, err
}

func (c *Client) sendAppPostRegister(ctx context.Context, request OptAppPostRegisterReq) (res AppPostRegisterRes, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/app/register"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeAppPostRegisterRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeAppPostRegisterResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// AppPostRequest invokes app-post-request operation.
//
// ユーザーが配車要求を行う.
//
// POST /app/requests
func (c *Client) AppPostRequest(ctx context.Context, request OptAppPostRequestReq) (*AppPostRequestAccepted, error) {
	res, err := c.sendAppPostRequest(ctx, request)
	return res, err
}

func (c *Client) sendAppPostRequest(ctx context.Context, request OptAppPostRequestReq) (res *AppPostRequestAccepted, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/app/requests"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeAppPostRequestRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeAppPostRequestResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// AppPostRequestEvaluate invokes app-post-request-evaluate operation.
//
// ユーザーが椅子を評価する.
//
// POST /app/requests/{request_id}/evaluate
func (c *Client) AppPostRequestEvaluate(ctx context.Context, request OptAppPostRequestEvaluateReq, params AppPostRequestEvaluateParams) (AppPostRequestEvaluateRes, error) {
	res, err := c.sendAppPostRequestEvaluate(ctx, request, params)
	return res, err
}

func (c *Client) sendAppPostRequestEvaluate(ctx context.Context, request OptAppPostRequestEvaluateReq, params AppPostRequestEvaluateParams) (res AppPostRequestEvaluateRes, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [3]string
	pathParts[0] = "/app/requests/"
	{
		// Encode "request_id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "request_id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.StringToString(params.RequestID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	pathParts[2] = "/evaluate"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeAppPostRequestEvaluateRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeAppPostRequestEvaluateResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// ChairGetInquiries invokes chair-get-inquiries operation.
//
// 問い合わせの一覧を取得する.
//
// GET /admin/inquiries
func (c *Client) ChairGetInquiries(ctx context.Context, params ChairGetInquiriesParams) (*ChairGetInquiriesOK, error) {
	res, err := c.sendChairGetInquiries(ctx, params)
	return res, err
}

func (c *Client) sendChairGetInquiries(ctx context.Context, params ChairGetInquiriesParams) (res *ChairGetInquiriesOK, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/admin/inquiries"
	uri.AddPathParts(u, pathParts[:]...)

	q := uri.NewQueryEncoder()
	{
		// Encode "limit" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "limit",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := params.Limit.Get(); ok {
				return e.EncodeValue(conv.Float64ToString(val))
			}
			return nil
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "cursor" parameter.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "cursor",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}

		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := params.Cursor.Get(); ok {
				return e.EncodeValue(conv.StringToString(val))
			}
			return nil
		}); err != nil {
			return res, errors.Wrap(err, "encode query")
		}
	}
	u.RawQuery = q.Values().Encode()

	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeChairGetInquiriesResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// ChairGetInquiry invokes chair-get-inquiry operation.
//
// 指定したIDの問い合わせ内容を取得.
//
// GET /admin/inquiries/{inquiry_id}
func (c *Client) ChairGetInquiry(ctx context.Context, params ChairGetInquiryParams) (ChairGetInquiryRes, error) {
	res, err := c.sendChairGetInquiry(ctx, params)
	return res, err
}

func (c *Client) sendChairGetInquiry(ctx context.Context, params ChairGetInquiryParams) (res ChairGetInquiryRes, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [2]string
	pathParts[0] = "/admin/inquiries/"
	{
		// Encode "inquiry_id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "inquiry_id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.StringToString(params.InquiryID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeChairGetInquiryResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// ChairGetNotification invokes chair-get-notification operation.
//
// 椅子に配車要求を通知するなどで使う想定.
//
// GET /chair/notification
func (c *Client) ChairGetNotification(ctx context.Context) error {
	_, err := c.sendChairGetNotification(ctx)
	return err
}

func (c *Client) sendChairGetNotification(ctx context.Context) (res *ChairGetNotificationOK, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/chair/notification"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeChairGetNotificationResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// ChairGetRequest invokes chair-get-request operation.
//
// 椅子向け通知エンドポイントから通知されたidの情報を取得する想定.
//
// GET /chair/requests/{request_id}
func (c *Client) ChairGetRequest(ctx context.Context, params ChairGetRequestParams) (ChairGetRequestRes, error) {
	res, err := c.sendChairGetRequest(ctx, params)
	return res, err
}

func (c *Client) sendChairGetRequest(ctx context.Context, params ChairGetRequestParams) (res ChairGetRequestRes, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [2]string
	pathParts[0] = "/chair/requests/"
	{
		// Encode "request_id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "request_id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.StringToString(params.RequestID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "GET", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeChairGetRequestResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// ChairPostActivate invokes chair-post-activate operation.
//
// 椅子が配車受付を開始する.
//
// POST /chair/activate
func (c *Client) ChairPostActivate(ctx context.Context, request *ChairPostActivateReq) error {
	_, err := c.sendChairPostActivate(ctx, request)
	return err
}

func (c *Client) sendChairPostActivate(ctx context.Context, request *ChairPostActivateReq) (res *ChairPostActivateNoContent, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/chair/activate"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeChairPostActivateRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeChairPostActivateResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// ChairPostCoordinate invokes chair-post-coordinate operation.
//
// 椅子が位置情報を送信する.
//
// POST /chair/coordinate
func (c *Client) ChairPostCoordinate(ctx context.Context, request OptCoordinate) error {
	_, err := c.sendChairPostCoordinate(ctx, request)
	return err
}

func (c *Client) sendChairPostCoordinate(ctx context.Context, request OptCoordinate) (res *ChairPostCoordinateNoContent, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/chair/coordinate"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeChairPostCoordinateRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeChairPostCoordinateResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// ChairPostDeactivate invokes chair-post-deactivate operation.
//
// 椅子が配車受付を停止する.
//
// POST /chair/deactivate
func (c *Client) ChairPostDeactivate(ctx context.Context, request *ChairPostDeactivateReq) error {
	_, err := c.sendChairPostDeactivate(ctx, request)
	return err
}

func (c *Client) sendChairPostDeactivate(ctx context.Context, request *ChairPostDeactivateReq) (res *ChairPostDeactivateNoContent, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/chair/deactivate"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeChairPostDeactivateRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeChairPostDeactivateResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// ChairPostRegister invokes chair-post-register operation.
//
// 椅子登録を行う.
//
// POST /chair/register
func (c *Client) ChairPostRegister(ctx context.Context, request OptChairPostRegisterReq) (*ChairPostRegisterCreated, error) {
	res, err := c.sendChairPostRegister(ctx, request)
	return res, err
}

func (c *Client) sendChairPostRegister(ctx context.Context, request OptChairPostRegisterReq) (res *ChairPostRegisterCreated, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/chair/register"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeChairPostRegisterRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeChairPostRegisterResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// ChairPostRequestAccept invokes chair-post-request-accept operation.
//
// 椅子が配車要求を受理する.
//
// POST /chair/requests/{request_id}/accept
func (c *Client) ChairPostRequestAccept(ctx context.Context, params ChairPostRequestAcceptParams) (ChairPostRequestAcceptRes, error) {
	res, err := c.sendChairPostRequestAccept(ctx, params)
	return res, err
}

func (c *Client) sendChairPostRequestAccept(ctx context.Context, params ChairPostRequestAcceptParams) (res ChairPostRequestAcceptRes, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [3]string
	pathParts[0] = "/chair/requests/"
	{
		// Encode "request_id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "request_id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.StringToString(params.RequestID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	pathParts[2] = "/accept"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeChairPostRequestAcceptResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// ChairPostRequestDeny invokes chair-post-request-deny operation.
//
// 椅子が配車要求を拒否する.
//
// POST /chair/requests/{request_id}/deny
func (c *Client) ChairPostRequestDeny(ctx context.Context, params ChairPostRequestDenyParams) (ChairPostRequestDenyRes, error) {
	res, err := c.sendChairPostRequestDeny(ctx, params)
	return res, err
}

func (c *Client) sendChairPostRequestDeny(ctx context.Context, params ChairPostRequestDenyParams) (res ChairPostRequestDenyRes, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [3]string
	pathParts[0] = "/chair/requests/"
	{
		// Encode "request_id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "request_id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.StringToString(params.RequestID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	pathParts[2] = "/deny"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeChairPostRequestDenyResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// ChairPostRequestDepart invokes chair-post-request-depart operation.
//
// 椅子が配車位置から出発する(ユーザーが乗車完了した).
//
// POST /chair/requests/{request_id}/depart
func (c *Client) ChairPostRequestDepart(ctx context.Context, params ChairPostRequestDepartParams) (ChairPostRequestDepartRes, error) {
	res, err := c.sendChairPostRequestDepart(ctx, params)
	return res, err
}

func (c *Client) sendChairPostRequestDepart(ctx context.Context, params ChairPostRequestDepartParams) (res ChairPostRequestDepartRes, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [3]string
	pathParts[0] = "/chair/requests/"
	{
		// Encode "request_id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "request_id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.StringToString(params.RequestID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	pathParts[2] = "/depart"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeChairPostRequestDepartResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// Initialize invokes initialize operation.
//
// サービスを初期化する.
//
// POST /initialize
func (c *Client) Initialize(ctx context.Context) (*InitializeOK, error) {
	res, err := c.sendInitialize(ctx)
	return res, err
}

func (c *Client) sendInitialize(ctx context.Context) (res *InitializeOK, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [1]string
	pathParts[0] = "/initialize"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeInitializeResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}
