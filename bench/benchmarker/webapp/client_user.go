package webapp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/isucon/isucon14/bench/benchmarker/webapp/api"
)

func (c *Client) AppPostRegister(ctx context.Context, reqBody *api.AppPostUsersReq) (*api.AppPostUsersCreated, error) {
	reqBodyBuf, err := reqBody.MarshalJSON()
	if err != nil {
		return nil, err
	}

	req, err := c.agent.NewRequest(http.MethodPost, "/api/app/users", bytes.NewReader(reqBodyBuf))
	if err != nil {
		return nil, err
	}

	for _, modifier := range c.requestModifiers {
		modifier(req)
	}

	resp, err := c.agent.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("POST /api/app/usersのリクエストが失敗しました: %w", err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("POST /api/app/usersへのリクエストに対して、期待されたHTTPステータスコードが確認できませんでした (expected:%d, actual:%d)", http.StatusCreated, resp.StatusCode)
	}

	resBody := &api.AppPostUsersCreated{}
	if err := json.NewDecoder(resp.Body).Decode(resBody); err != nil {
		return nil, fmt.Errorf("POST /api/app/usersのJSONのdecodeに失敗しました: %w", err)
	}

	return resBody, nil
}

func (c *Client) AppGetRequests(ctx context.Context) (*api.AppGetRidesOK, error) {
	req, err := c.agent.NewRequest(http.MethodGet, "/api/app/rides", nil)
	if err != nil {
		return nil, err
	}

	for _, modifier := range c.requestModifiers {
		modifier(req)
	}

	resp, err := c.agent.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GET /app/rides のリクエストが失敗しました: %w", err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /app/rides へのリクエストに対して、期待されたHTTPステータスコードが確認できませんでした (expected:%d, actual:%d)", http.StatusOK, resp.StatusCode)
	}

	resBody := &api.AppGetRidesOK{}
	if err := json.NewDecoder(resp.Body).Decode(resBody); err != nil {
		return nil, fmt.Errorf("GET /app/ridesのJSONのdecodeに失敗しました: %w", err)
	}

	return resBody, nil
}

func (c *Client) AppPostRidesEstimatedFare(ctx context.Context, reqBody *api.AppPostRidesEstimatedFareReq) (*api.AppPostRidesEstimatedFareOK, error) {
	reqBodyBuf, err := reqBody.MarshalJSON()
	if err != nil {
		return nil, err
	}

	req, err := c.agent.NewRequest(http.MethodPost, "/api/app/rides/estimated-fare", bytes.NewReader(reqBodyBuf))
	if err != nil {
		return nil, err
	}

	for _, modifier := range c.requestModifiers {
		modifier(req)
	}

	resp, err := c.agent.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("POST /app/rides/estimated-fareのリクエストが失敗しました: %w", err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("POST /app/rides/estimated-fareへのリクエストに対して、期待されたHTTPステータスコードが確認できませんでした (expected:%d, actual:%d)", http.StatusOK, resp.StatusCode)
	}

	resBody := &api.AppPostRidesEstimatedFareOK{}
	if err := json.NewDecoder(resp.Body).Decode(resBody); err != nil {
		return nil, fmt.Errorf("POST /app/rides/estimated-fareのJSONのdecodeに失敗しました: %w", err)
	}

	return resBody, nil
}

func (c *Client) AppPostRequest(ctx context.Context, reqBody *api.AppPostRidesReq) (*api.AppPostRidesAccepted, error) {
	reqBodyBuf, err := reqBody.MarshalJSON()
	if err != nil {
		return nil, err
	}

	req, err := c.agent.NewRequest(http.MethodPost, "/api/app/rides", bytes.NewReader(reqBodyBuf))
	if err != nil {
		return nil, err
	}

	for _, modifier := range c.requestModifiers {
		modifier(req)
	}

	resp, err := c.agent.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("POST /api/app/ridesのリクエストが失敗しました: %w", err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("POST /api/app/ridesへのリクエストに対して、期待されたHTTPステータスコードが確認できませんでした (expected:%d, actual:%d)", http.StatusAccepted, resp.StatusCode)
	}

	resBody := &api.AppPostRidesAccepted{}
	if err := json.NewDecoder(resp.Body).Decode(resBody); err != nil {
		return nil, fmt.Errorf("POST /api/app/ridesのJSONのdecodeに失敗しました: %w", err)
	}

	return resBody, nil
}

func (c *Client) AppGetRequest(ctx context.Context, rideID string) (*api.AppRide, error) {
	req, err := c.agent.NewRequest(http.MethodGet, fmt.Sprintf("/api/app/rides/%s", rideID), nil)
	if err != nil {
		return nil, err
	}

	for _, modifier := range c.requestModifiers {
		modifier(req)
	}

	resp, err := c.agent.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GET /api/app/rides/{ride_id}のリクエストが失敗しました: %w", err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /api/app/rides/{ride_id}へのリクエストに対して、期待されたHTTPステータスコードが確認できませんでした (expected:%d, actual:%d)", http.StatusOK, resp.StatusCode)
	}

	resBody := &api.AppRide{}
	if err := json.NewDecoder(resp.Body).Decode(resBody); err != nil {
		return nil, fmt.Errorf("GET /api/app/rides/{ride_id}のJSONのdecodeに失敗しました: %w", err)
	}

	return resBody, nil
}

func (c *Client) AppPostRequestEvaluate(ctx context.Context, rideID string, reqBody *api.AppPostRideEvaluationReq) (*api.AppPostRideEvaluationOK, error) {
	reqBodyBuf, err := reqBody.MarshalJSON()
	if err != nil {
		return nil, err
	}

	req, err := c.agent.NewRequest(http.MethodPost, fmt.Sprintf("/api/app/rides/%s/evaluation", rideID), bytes.NewReader(reqBodyBuf))
	if err != nil {
		return nil, err
	}

	for _, modifier := range c.requestModifiers {
		modifier(req)
	}

	resp, err := c.agent.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("POST /api/app/rides/{ride_id}/evaluationのリクエストが失敗しました: %w", err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("POST /api/app/rides/{ride_id}/evaluationへのリクエストに対して、期待されたHTTPステータスコードが確認できませんでした (expected:%d, actual:%d)", http.StatusOK, resp.StatusCode)
	}

	resBody := &api.AppPostRideEvaluationOK{}
	if err := json.NewDecoder(resp.Body).Decode(resBody); err != nil {
		return nil, fmt.Errorf("POST /api/app/rides/{ride_id}/evaluationのJSONのdecodeに失敗しました: %w", err)
	}

	return resBody, nil
}

func (c *Client) AppPostPaymentMethods(ctx context.Context, reqBody *api.AppPostPaymentMethodsReq) (*api.AppPostPaymentMethodsNoContent, error) {
	reqBodyBuf, err := reqBody.MarshalJSON()
	if err != nil {
		return nil, err
	}

	req, err := c.agent.NewRequest(http.MethodPost, "/api/app/payment-methods", bytes.NewReader(reqBodyBuf))
	if err != nil {
		return nil, err
	}

	for _, modifier := range c.requestModifiers {
		modifier(req)
	}

	resp, err := c.agent.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("POST /api/app/payment-methodsのリクエストが失敗しました: %w", err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("POST /api/app/payment-methodsへのリクエストに対して、期待されたHTTPステータスコードが確認できませんでした (expected:%d, actual:%d)", http.StatusNoContent, resp.StatusCode)
	}

	resBody := &api.AppPostPaymentMethodsNoContent{}
	return resBody, nil
}

func (c *Client) AppGetNotification(ctx context.Context) iter.Seq2[*api.AppRide, error] {
	return func(yield func(*api.AppRide, error) bool) {
		for notification, err := range c.appGetNotification(ctx, false) {
			if notification == nil {
				if !yield(nil, err) {
					return
				}
			} else {
				if !yield(&api.AppRide{
					ID:                    notification.RideID,
					PickupCoordinate:      notification.PickupCoordinate,
					DestinationCoordinate: notification.DestinationCoordinate,
					Status:                notification.Status,
					Chair:                 notification.Chair,
					CreatedAt:             notification.CreatedAt,
					UpdatedAt:             notification.UpdatedAt,
				}, err) {
					return
				}
			}
		}
	}
}

func (c *Client) appGetNotification(ctx context.Context, nested bool) iter.Seq2[*api.AppGetNotificationOK, error] {
	req, err := c.agent.NewRequest(http.MethodGet, "/api/app/notification", nil)
	if err != nil {
		return func(yield func(*api.AppGetNotificationOK, error) bool) { yield(nil, err) }
	}

	for _, modifier := range c.requestModifiers {
		modifier(req)
	}

	httpClient := &http.Client{
		Transport:     c.agent.HttpClient.Transport,
		CheckRedirect: c.agent.HttpClient.CheckRedirect,
		Jar:           c.agent.HttpClient.Jar,
		Timeout:       60 * time.Second,
	}

	resp, err := httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return func(yield func(*api.AppGetNotificationOK, error) bool) {
			yield(nil, fmt.Errorf("GET /api/app/notificationsのリクエストが失敗しました: %w", err))
		}
	}

	if strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream") {
		return func(yield func(*api.AppGetNotificationOK, error) bool) {
			defer closeBody(resp)

			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				request := &api.AppGetNotificationOK{}
				line := scanner.Text()
				if strings.HasPrefix(line, "data:") {
					err := json.Unmarshal([]byte(line[5:]), request)
					if !yield(request, err) || err != nil {
						return
					}
				}
			}
		}
	}

	defer closeBody(resp)
	request := &api.AppGetNotificationOK{}
	if resp.StatusCode == http.StatusOK {
		if err = json.NewDecoder(resp.Body).Decode(request); err != nil {
			err = fmt.Errorf("requestのJSONのdecodeに失敗しました: %w", err)
		}
	} else if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("GET /api/app/notificationsへのリクエストに対して、期待されたHTTPステータスコードが確認できませんでした (expected:%d or %d, actual:%d)", http.StatusOK, http.StatusNoContent, resp.StatusCode)
	}

	if nested {
		return func(yield func(*api.AppGetNotificationOK, error) bool) { yield(request, err) }
	} else {
		return func(yield func(*api.AppGetNotificationOK, error) bool) {
			if !yield(request, err) {
				return
			}

			const defaultWaitTime = 30 * time.Millisecond
			waitTime := defaultWaitTime
			for {
				select {
				// こちらから切断
				case <-ctx.Done():
					return

				default:
					time.Sleep(waitTime)
					for notification, err := range c.appGetNotification(ctx, true) {
						if !yield(notification, err) {
							return
						}
						if notification != nil && notification.RetryAfterMs.IsSet() {
							waitTime = time.Duration(notification.RetryAfterMs.Value) * time.Millisecond
						} else {
							waitTime = defaultWaitTime
						}
					}
				}
			}
		}
	}
}

func (c *Client) AppGetNearbyChairs(ctx context.Context, params *api.AppGetNearbyChairsParams) (*api.AppGetNearbyChairsOK, error) {
	queryParams := url.Values{}
	queryParams.Set("latitude", strconv.Itoa(params.Latitude))
	queryParams.Set("longitude", strconv.Itoa(params.Longitude))
	if params.Distance.Set {
		queryParams.Set("distance", strconv.Itoa(params.Distance.Value))
	}

	req, err := c.agent.NewRequest(http.MethodGet, "/api/app/nearby-chairs?"+queryParams.Encode(), nil)
	if err != nil {
		return nil, err
	}

	for _, modifier := range c.requestModifiers {
		modifier(req)
	}

	resp, err := c.agent.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GET /api/app/requests/nearby-chairsのリクエストが失敗しました: %w", err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /api/app/requests/nearby-chairsへのリクエストに対して、期待されたHTTPステータスコードが確認できませんでした (expected:%d, actual:%d)", http.StatusOK, resp.StatusCode)
	}

	resBody := &api.AppGetNearbyChairsOK{}
	if err := json.NewDecoder(resp.Body).Decode(resBody); err != nil {
		return nil, fmt.Errorf("GET /api/app/requests/nearby-chairsのJSONのdecodeに失敗しました: %w", err)
	}

	return resBody, nil
}
