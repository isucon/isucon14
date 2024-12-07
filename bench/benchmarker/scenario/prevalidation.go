package scenario

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/isucon/isucon14/bench/benchmarker/webapp"
	"github.com/isucon/isucon14/bench/benchmarker/webapp/api"
	"github.com/isucon/isucon14/bench/benchrun"
)

func (s *Scenario) validateFrontendFiles(ctx context.Context) error {
	client, err := webapp.NewClient(webapp.ClientConfig{
		TargetBaseURL:         s.target,
		TargetAddr:            s.addr,
		ClientIdleConnTimeout: 10 * time.Second,
	})
	if err != nil {
		return err
	}

	frontendHashes := benchrun.FrontendHashesMap
	indexHtmlHash := frontendHashes["index.html"]

	{
		actualHash, err := client.StaticGetFileHash(ctx, "/client")
		if err != nil {
			return err
		}
		if actualHash != indexHtmlHash {
			return errors.New("/の内容が期待したものと一致しません")
		}
	}

	for path, expectedHash := range frontendHashes {
		// check separately
		if path == "/index.html" {
			continue
		}

		actualHash, err := client.StaticGetFileHash(ctx, path)
		if err != nil {
			return err
		}
		if actualHash != expectedHash {
			return errors.New(path + "の内容が期待したものと一致しません")
		}
	}

	// check index.html for other paths
	{
		actualHash, err := client.StaticGetFileHash(ctx, "/owner")
		if err != nil {
			return err
		}
		if actualHash != indexHtmlHash {
			return errors.New("/ownerの内容が期待したものと一致しません")
		}
	}

	return nil
}

// 実装の検証を行う
func (s *Scenario) prevalidation(ctx context.Context, client *webapp.Client) error {
	clientConfig := webapp.ClientConfig{
		TargetBaseURL:         s.target,
		TargetAddr:            s.addr,
		ClientIdleConnTimeout: 10 * time.Second,
	}

	if err := validateInitialData(ctx, clientConfig); err != nil {
		s.contestantLogger.Error("初期データのチェックに失敗しました", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func validateInitialData(ctx context.Context, clientConfig webapp.ClientConfig) error {
	validationData := LoadData()

	cmpOptions := []cmp.Option{
		cmpopts.SortSlices(func(i, j api.OwnerGetChairsOKChairsItem) bool {
			return i.ID < j.ID
		}),
		cmpopts.SortSlices(func(i, j api.OwnerGetSalesOKChairsItem) bool {
			return i.ID < j.ID
		}),
		cmpopts.SortSlices(func(i, j api.OwnerGetSalesOKModelsItem) bool {
			return i.Model < j.Model
		}),
		cmpopts.SortSlices(func(i, j api.AppGetRidesOKRidesItem) bool { return i.ID < j.ID }),
	}

	{
		ownerClient, err := webapp.NewClient(clientConfig)
		if err != nil {
			return err
		}
		ownerClient.SetCookie(&http.Cookie{Name: "owner_session", Value: "0811617de5c97aea5ddb433f085c3d1e"})

		chairs, err := ownerClient.OwnerGetChairs(ctx)
		if err != nil {
			return err
		}
		if !cmp.Equal(chairs, &validationData.Owner01JDFEDF00B09BNMV8MP0RB34G.Chairs, cmpOptions...) {
			return errors.New("GET /api/owner/chairs のレスポンスの内容が期待したものと一致しません")
		}

		sales, err := ownerClient.OwnerGetSales(ctx, &api.OwnerGetSalesParams{})
		if err != nil {
			return err
		}
		if !cmp.Equal(sales, &validationData.Owner01JDFEDF00B09BNMV8MP0RB34G.Sales, cmpOptions...) {
			return errors.New("GET /api/owner/sales のレスポンスの内容が期待したものと一致しません")
		}

		sales2, err := ownerClient.OwnerGetSales(ctx, &api.OwnerGetSalesParams{Since: api.NewOptInt64(1732579200000), Until: api.NewOptInt64(1732622400000)})
		if err != nil {
			return err
		}
		if !cmp.Equal(sales2, &validationData.Owner01JDFEDF00B09BNMV8MP0RB34G.Sales1732579200000to1732622400000, cmpOptions...) {
			return errors.New("GET /api/owner/sales のレスポンスの内容が期待したものと一致しません")
		}
	}
	{
		userClient, err := webapp.NewClient(clientConfig)
		if err != nil {
			return err
		}
		userClient.SetCookie(&http.Cookie{Name: "app_session", Value: "21e9562de048ee9b34da840296509fa9"})

		rides, err := userClient.AppGetRequests(ctx)
		if err != nil {
			return err
		}
		if !cmp.Equal(rides, &validationData.User01JDM0N9W89PK57C7XEVGD5C80.Rides, cmpOptions...) {
			return errors.New("GET /api/app/rides のレスポンスの内容が期待したものと一致しません")
		}
	}
	{
		userClient, err := webapp.NewClient(clientConfig)
		if err != nil {
			return err
		}
		userClient.SetCookie(&http.Cookie{Name: "app_session", Value: "c9e15fd57545f43105ace9088f1c467e"})

		rides, err := userClient.AppGetRequests(ctx)
		if err != nil {
			return err
		}
		if !cmp.Equal(rides, &validationData.User01JDK5EFNGT8ZHMTQXQ4BNH8NQ.Rides, cmpOptions...) {
			return errors.New("GET /api/app/rides のレスポンスの内容が期待したものと一致しません")
		}

		estimated1, err := userClient.AppPostRidesEstimatedFare(ctx, &api.AppPostRidesEstimatedFareReq{
			PickupCoordinate:      api.Coordinate{Latitude: 0 + 10, Longitude: 0 + 10},
			DestinationCoordinate: api.Coordinate{Latitude: 3 + 10, Longitude: 10 + 10},
		})
		if err != nil {
			return err
		}
		if !cmp.Equal(estimated1, &validationData.User01JDK5EFNGT8ZHMTQXQ4BNH8NQ.Estimated_3_10, cmpOptions...) {
			return errors.New("POST /api/app/rides/estimated-fare のレスポンスの内容が期待したものと一致しません")
		}

		estimated2, err := userClient.AppPostRidesEstimatedFare(ctx, &api.AppPostRidesEstimatedFareReq{
			PickupCoordinate:      api.Coordinate{Latitude: 0 - 10, Longitude: 0 - 10},
			DestinationCoordinate: api.Coordinate{Latitude: -11 - 10, Longitude: 10 - 10},
		})
		if err != nil {
			return err
		}
		if !cmp.Equal(estimated2, &validationData.User01JDK5EFNGT8ZHMTQXQ4BNH8NQ.Estimated_m11_10, cmpOptions...) {
			return errors.New("POST /api/app/rides/estimated-fare のレスポンスの内容が期待したものと一致しません")
		}
	}
	{
		userClient, err := webapp.NewClient(clientConfig)
		if err != nil {
			return err
		}
		userClient.SetCookie(&http.Cookie{Name: "app_session", Value: "a8b21d78f143c3facdece4dffba964cc"})

		rides, err := userClient.AppGetRequests(ctx)
		if err != nil {
			return err
		}
		if !cmp.Equal(rides, &validationData.User01JDJ4XN10E2CRZ37RNZ5GAFW6.Rides, cmpOptions...) {
			return errors.New("GET /api/app/rides のレスポンスの内容が期待したものと一致しません")
		}

		estimated1, err := userClient.AppPostRidesEstimatedFare(ctx, &api.AppPostRidesEstimatedFareReq{
			PickupCoordinate:      api.Coordinate{Latitude: 0 + 10, Longitude: 0 + 10},
			DestinationCoordinate: api.Coordinate{Latitude: 3 + 10, Longitude: 10 + 10},
		})
		if err != nil {
			return err
		}
		if !cmp.Equal(estimated1, &validationData.User01JDJ4XN10E2CRZ37RNZ5GAFW6.Estimated_3_10, cmpOptions...) {
			return errors.New("POST /api/app/rides/estimated-fare のレスポンスの内容が期待したものと一致しません")
		}

		estimated2, err := userClient.AppPostRidesEstimatedFare(ctx, &api.AppPostRidesEstimatedFareReq{
			PickupCoordinate:      api.Coordinate{Latitude: 0 - 10, Longitude: 0 - 10},
			DestinationCoordinate: api.Coordinate{Latitude: -11 - 10, Longitude: 10 - 10},
		})
		if err != nil {
			return err
		}
		if !cmp.Equal(estimated2, &validationData.User01JDJ4XN10E2CRZ37RNZ5GAFW6.Estimated_m11_10, cmpOptions...) {
			return errors.New("POST /api/app/rides/estimated-fare のレスポンスの内容が期待したものと一致しません")
		}
	}

	return nil
}

func validateSuccessFlow(ctx context.Context, clientConfig webapp.ClientConfig) error {
	userClient, err := webapp.NewClient(clientConfig)
	if err != nil {
		return err
	}
	ownerClient, err := webapp.NewClient(clientConfig)
	if err != nil {
		return err
	}
	chairClient, err := webapp.NewClient(clientConfig)
	if err != nil {
		return err
	}

	userID := ""
	// POST /api/app/register
	{
		result, err := userClient.AppPostRegister(ctx, &api.AppPostUsersReq{
			Username:    "hoge",
			Firstname:   "hoge",
			Lastname:    "hoge",
			DateOfBirth: "2000-01-01",
		})
		if err != nil {
			return err
		}
		if result.ID == "" {
			return errors.New("POST /api/app/register の返却するIDは、空であってはいけません")
		}
		userID = result.ID
	}

	paymentToken := "token"
	// POST /api/app/payment-methods
	{
		_, err := userClient.AppPostPaymentMethods(ctx, &api.AppPostPaymentMethodsReq{
			Token: paymentToken,
		})
		if err != nil {
			return err
		}
	}

	// POST /api/app/requests
	requestID := ""
	{
		result, err := userClient.AppPostRequest(ctx, &api.AppPostRidesReq{
			PickupCoordinate: api.Coordinate{
				Latitude:  0,
				Longitude: 0,
			},
			DestinationCoordinate: api.Coordinate{
				Latitude:  10,
				Longitude: 10,
			},
		})
		if err != nil {
			return err
		}
		if result.RideID == "" {
			return errors.New("POST /api/app/requests の返却するIDは、空であってはいけません")
		}
		requestID = result.RideID
	}

	// TODO: 登録された直後の椅子の状態を確認する

	// GET /api/app/notifications
	{
		for result, err := range userClient.AppGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateAppNotification(result.Data.V, requestID, api.RideStatusMATCHING); err != nil {
				return err
			}
			if result.Data.V.Chair.Set {
				return errors.New("GET /api/app/requests/:requestID の返却するchairがセットされているべきではありません")
			}
			break
		}
	}

	// GET /api/app/nearby-chairs
	{
		result, err := userClient.AppGetNearbyChairs(ctx, &api.AppGetNearbyChairsParams{
			Latitude:  0,
			Longitude: 0,
		})
		if err != nil {
			return err
		}
		if len(result.Chairs) != 0 {
			return fmt.Errorf("GET /api/app/nearby-chairs の返却するchairsの数が異なります (expected:%d, actual:%d)", 0, len(result.Chairs))
		}
	}

	chairRegisterToken := ""
	// POST /api/owner/register
	{
		result, err := ownerClient.OwnerPostRegister(ctx, &api.OwnerPostOwnersReq{
			Name: "hoge",
		})
		if err != nil {
			return err
		}
		if result.ID == "" {
			return errors.New("POST /api/owner/register の返却するIDは、空であってはいけません")
		}
		if result.ChairRegisterToken == "" {
			return errors.New("POST /api/owner/register の返却するchair_register_tokenは、空であってはいけません")
		}
		chairRegisterToken = result.ChairRegisterToken
	}

	chairID := ""
	// POST /api/chair/register
	{
		result, err := chairClient.ChairPostRegister(ctx, &api.ChairPostChairsReq{
			Name:               "hoge",
			Model:              "A",
			ChairRegisterToken: chairRegisterToken,
		})
		if err != nil {
			return err
		}
		if result.ID == "" {
			return errors.New("POST /api/chair/register の返却するIDは、空であってはいけません")
		}
		chairID = result.ID
	}

	// POST /api/chair/activate
	{
		_, err := chairClient.ChairPostActivity(ctx, &api.ChairPostActivityReq{
			IsActive: true,
		})
		if err != nil {
			return err
		}
	}

	// POST /api/chair/coordinate
	{
		_, err := chairClient.ChairPostCoordinate(ctx, &api.Coordinate{
			Latitude:  1,
			Longitude: 1,
		})
		if err != nil {
			return err
		}
	}

	// GET /api/chair/notification
	{
		for result, err := range chairClient.ChairGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateChairNotification(result.Data.V, requestID, userID, api.RideStatusMATCHING); err != nil {
				return err
			}
			break
		}
	}

	// GET /api/app/notifications
	{
		for result, err := range userClient.AppGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateAppNotificationWithChair(result.Data.V, requestID, api.RideStatusMATCHING, chairID); err != nil {
				return err
			}
			break
		}
	}

	// POST /api/chair/requests/accept
	{
		_, err := chairClient.ChairPostRideStatus(ctx, requestID, &api.ChairPostRideStatusReq{
			Status: api.ChairPostRideStatusReqStatusENROUTE,
		})
		if err != nil {
			return err
		}
	}

	// GET /api/chair/notification
	{
		for result, err := range chairClient.ChairGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateChairNotification(result.Data.V, requestID, userID, api.RideStatusENROUTE); err != nil {
				return err
			}
			break
		}
	}

	// GET /api/app/notifications
	{
		for result, err := range userClient.AppGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateAppNotificationWithChair(result.Data.V, requestID, api.RideStatusENROUTE, chairID); err != nil {
				return err
			}
			break
		}
	}

	// POST /api/chair/coordinate
	{
		_, err := chairClient.ChairPostCoordinate(ctx, &api.Coordinate{
			Latitude:  0,
			Longitude: 0,
		})
		if err != nil {
			return err
		}
	}

	// GET /api/chair/notification
	{
		for result, err := range chairClient.ChairGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateChairNotification(result.Data.V, requestID, userID, api.RideStatusPICKUP); err != nil {
				return err
			}
			break
		}
	}

	// GET /api/app/notifications
	{
		for result, err := range userClient.AppGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateAppNotificationWithChair(result.Data.V, requestID, api.RideStatusPICKUP, chairID); err != nil {
				return err
			}
			break
		}
	}

	// POST /api/chair/requests/depart
	{
		_, err := chairClient.ChairPostRideStatus(ctx, requestID, &api.ChairPostRideStatusReq{
			Status: api.ChairPostRideStatusReqStatusCARRYING,
		})
		if err != nil {
			return err
		}
	}

	// GET /api/chair/notification
	{
		for result, err := range chairClient.ChairGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateChairNotification(result.Data.V, requestID, userID, api.RideStatusCARRYING); err != nil {
				return err
			}
			break
		}
	}

	// GET /api/app/notifications
	{
		for result, err := range userClient.AppGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateAppNotificationWithChair(result.Data.V, requestID, api.RideStatusCARRYING, chairID); err != nil {
				return err
			}
			break
		}
	}

	// POST /api/chair/coordinate
	{
		_, err := chairClient.ChairPostCoordinate(ctx, &api.Coordinate{
			Latitude:  10,
			Longitude: 10,
		})
		if err != nil {
			return err
		}
	}

	// GET /api/chair/notification
	{
		for result, err := range chairClient.ChairGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateChairNotification(result.Data.V, requestID, userID, api.RideStatusARRIVED); err != nil {
				return err
			}
			break
		}
	}

	// GET /api/app/notifications
	{
		for result, err := range userClient.AppGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateAppNotificationWithChair(result.Data.V, requestID, api.RideStatusARRIVED, chairID); err != nil {
				return err
			}
			break
		}
	}

	// POST /api/app/request/:requestID/evaluate
	{
		_, err := userClient.AppPostRequestEvaluate(ctx, requestID, &api.AppPostRideEvaluationReq{
			Evaluation: 5,
		})
		if err != nil {
			return err
		}
	}

	// GET /api/app/nearby-chairs
	{
		result, err := userClient.AppGetNearbyChairs(ctx, &api.AppGetNearbyChairsParams{
			Latitude:  0,
			Longitude: 0,
		})
		if err != nil {
			return err
		}
		if len(result.Chairs) != 1 {
			return fmt.Errorf("GET /api/app/nearby-chairs の返却するchairsの数が異なります (expected:%d, actual:%d)", 1, len(result.Chairs))
		}
		if result.Chairs[0].ID != chairID {
			return fmt.Errorf("GET /api/app/nearby-chairs の返却するchairのIDが異なります (expected:%s, actual:%s)", chairID, result.Chairs[0].ID)
		}
		if result.Chairs[0].Name != "hoge" {
			return fmt.Errorf("GET /api/app/nearby-chairs の返却するchairのnameが異なります (expected:%s, actual:%s)", "hoge", result.Chairs[0].Name)
		}
		if result.Chairs[0].Model != "A" {
			return fmt.Errorf("GET /api/app/nearby-chairs の返却するchairのmodelが異なります (expected:%s, actual:%s)", "A", result.Chairs[0].Model)
		}
		if result.Chairs[0].CurrentCoordinate.Latitude != 10 {
			return fmt.Errorf("GET /api/app/nearby-chairs の返却するchairのcurrent_coordinateのlatitudeが異なります (expected:%d, actual:%d)", 10, result.Chairs[0].CurrentCoordinate.Latitude)
		}
		if result.Chairs[0].CurrentCoordinate.Longitude != 10 {
			return fmt.Errorf("GET /api/app/nearby-chairs の返却するchairのcurrent_coordinateのlongitudeが異なります (expected:%d, actual:%d)", 10, result.Chairs[0].CurrentCoordinate.Longitude)
		}
	}

	// GET /api/app/notifications
	{
		for result, err := range userClient.AppGetNotification(ctx) {
			if err != nil {
				return err
			}
			if err := validateAppNotification(result.Data.V, requestID, api.RideStatusCOMPLETED); err != nil {
				return err
			}
			if result.Data.V.Chair.Value.Stats.TotalEvaluationAvg != 5 {
				return fmt.Errorf("GET /api/app/nearby-chairs の返却するchairのstatsのtotal_evaluation_avgが異なります (expected:%f, actual:%f)", 5.0, result.Data.V.Chair.Value.Stats.TotalEvaluationAvg)
			}
			if result.Data.V.Chair.Value.Stats.TotalRidesCount != 1 {
				return fmt.Errorf("GET /api/app/nearby-chairs の返却するchairのstatsのtotal_rides_countが異なります (expected:%d, actual:%d)", 1, result.Data.V.Chair.Value.Stats.TotalRidesCount)
			}
			break
		}
	}

	return nil
}

func validateAppNotification(req webapp.UserNotificationData, requestID string, status api.RideStatus) error {
	if req.RideID != requestID {
		return fmt.Errorf("GET /api/app/notification の返却するIDが、リクエストIDと一致しません (expected:%s, actual:%s)", requestID, req.RideID)
	}
	if req.PickupCoordinate.Latitude != 0 {
		return fmt.Errorf("GET /api/app/notification の返却するpickup_coordinateのlatitudeが異なります (expected:%d, actual:%d)", 0, req.PickupCoordinate.Latitude)
	}
	if req.PickupCoordinate.Longitude != 0 {
		return fmt.Errorf("GET /api/app/notification の返却するpickup_coordinateのlongitudeが異なります (expected:%d, actual:%d)", 0, req.PickupCoordinate.Longitude)
	}
	if req.DestinationCoordinate.Latitude != 10 {
		return fmt.Errorf("GET /api/app/notification の返却するdestination_coordinateのlatitudeが異なります (expected:%d, actual:%d)", 10, req.DestinationCoordinate.Latitude)
	}
	if req.DestinationCoordinate.Longitude != 10 {
		return fmt.Errorf("GET /api/app/notification の返却するdestination_coordinateのlongitudeが異なります (expected:%d, actual:%d)", 10, req.DestinationCoordinate.Longitude)
	}

	if req.Status != status {
		return fmt.Errorf("GET /api/app/notification の返却するstatusが異なります (expected:%s, actual:%s)", status, req.Status)
	}

	return nil
}

func validateAppNotificationWithChair(req webapp.UserNotificationData, requestID string, status api.RideStatus, chairID string) error {
	if err := validateAppNotification(req, requestID, status); err != nil {
		return err
	}
	if !req.Chair.Set {
		return errors.New("GET /api/app/notification の返却するchairが、返却されるべきです")
	}
	if req.Chair.Value.ID != chairID {
		return fmt.Errorf("GET /api/app/notification の返却するchair.idが異なります (expected:%s, actual:%s)", chairID, req.Chair.Value.ID)
	}
	if req.Chair.Value.Name != "hoge" {
		return fmt.Errorf("GET /api/app/notification の返却するchair.nameが異なります (expected:%s, actual:%s)", "hoge", req.Chair.Value.Name)
	}
	if req.Chair.Value.Model != "A" {
		return fmt.Errorf("GET /api/app/notification の返却するchair.modelが異なります (expected:%s, actual:%s)", "A", req.Chair.Value.Model)
	}
	return nil
}

func validateChairNotification(req webapp.ChairNotificationData, requestID string, userID string, status api.RideStatus) error {
	if req.RideID != requestID {
		return fmt.Errorf("GET /api/chair/notification の返却するIDが、リクエストIDと一致しません (expected:%s, actual:%s)", requestID, req.RideID)
	}
	if req.User.ID != userID {
		return fmt.Errorf("GET /api/chair/notification の返却するuser.idが、ユーザーIDと一致しません (expected:%s, actual:%s)", userID, req.User.ID)
	}
	if req.User.Name != "hoge hoge" {
		return fmt.Errorf("GET /api/chair/notification の返却するuser.nameが異なります (expected:%s, actual:%s)", "hoge hoge", req.User.Name)
	}
	if req.PickupCoordinate.Latitude != 0 {
		return fmt.Errorf("GET /api/chair/notification の返却するpickup_coordinateのlatitudeが異なります (expected:%d, actual:%d)", 0, req.PickupCoordinate.Latitude)
	}
	if req.PickupCoordinate.Longitude != 0 {
		return fmt.Errorf("GET /api/chair/notification の返却するpickup_coordinateのlongitudeが異なります (expected:%d, actual:%d)", 0, req.PickupCoordinate.Longitude)
	}
	if req.DestinationCoordinate.Latitude != 10 {
		return fmt.Errorf("GET /api/chair/notification の返却するdestination_coordinateのlatitudeが異なります (expected:%d, actual:%d)", 10, req.DestinationCoordinate.Latitude)
	}
	if req.DestinationCoordinate.Longitude != 10 {
		return fmt.Errorf("GET /api/chair/notification の返却するdestination_coordinateのlongitudeが異なります (expected:%d, actual:%d)", 10, req.DestinationCoordinate.Longitude)
	}
	if req.Status != status {
		return fmt.Errorf("GET /api/chair/notification の返却するstatusが異なります (expected:%s, actual:%s)", status, req.Status)
	}
	return nil
}
