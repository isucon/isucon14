package verify

// import (
// 	"context"
// 	"fmt"
// 	"math/rand/v2"
// 	"net/http"
// 	"time"

// 	"go.uber.org/zap"

// 	"github.com/isucon/isucon14/bench/benchmarker/webapp"
// 	"github.com/isucon/isucon14/bench/benchmarker/webapp/api"
// )

// type Agent struct {
// 	target           string
// 	contestantLogger *zap.Logger
// }

// func NewAgent(target string, contestantLogger *zap.Logger) (*Agent, error) {
// 	return &Agent{
// 		target:           target,
// 		contestantLogger: contestantLogger,
// 	}, nil
// }

// func randomCoordinate() api.Coordinate {
// 	return api.Coordinate{
// 		rand.Float64()*180 - 90,
// 		rand.Float64()*360 - 180,
// 	}
// }

// func (a *Agent) Run() error {
// 	userClient, err := webapp.NewClient(webapp.ClientConfig{
// 		TargetBaseURL:         a.target,
// 		DefaultClientTimeout:  5 * time.Second,
// 		ClientIdleConnTimeout: 10 * time.Second,
// 		InsecureSkipVerify:    true,
// 		ContestantLogger:      a.contestantLogger,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	ctx := context.Background()
// 	// ユーザーの登録
// 	{
// 		registerRes, err := userClient.AppPostRegister(ctx)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to register user", zap.Error(err))
// 			return err
// 		}

// 		userClient.AddRequestModifier(func(req *http.Ride) {
// 			req.Header.Set("Authorization", "Bearer "+registerRes.AccessToken)
// 		})
// 	}

// 	chairClient, err := webapp.NewClient(webapp.ClientConfig{
// 		TargetBaseURL:         a.target,
// 		DefaultClientTimeout:  5 * time.Second,
// 		ClientIdleConnTimeout: 10 * time.Second,
// 		InsecureSkipVerify:    true,
// 		ContestantLogger:      a.contestantLogger,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// ドライバーの登録
// 	{
// 		registerRes, err := chairClient.ChairPostRegister(ctx)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to register driver", zap.Error(err))
// 			return err
// 		}

// 		chairClient.AddRequestModifier(func(req *http.Ride) {
// 			req.Header.Set("Authorization", "Bearer "+registerRes.AccessToken)
// 		})
// 	}

// 	pickupCoordinate := randomCoordinate()
// 	destinationCoordinate := randomCoordinate()

// 	receivedRequestCh := make(chan *api.GetRequestOK)
// 	// ドライバーが受信開始
// 	go func() {
// 		ctx := context.Background()
// 		res, result, err := chairClient.ChairGetNotification(ctx)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to receive notifications", zap.Error(err))
// 			return
// 		}
// 		for receivedRequest := range res {
// 			receivedRequestCh <- receivedRequest
// 		}

// 		if err := result(); err != nil {
// 			a.contestantLogger.Error("Failed to receive notifications", zap.Error(err))
// 			return
// 		}
// 	}()

// 	// ドライバーが待機開始
// 	{
// 		_, err := chairClient.ChairPostActivate(ctx)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to post activate", zap.Error(err))
// 			return err
// 		}
// 	}

// 	// ドライバーが現在位置を送信
// 	{
// 		cord := randomCoordinate()
// 		_, err := chairClient.ChairPostCoordinate(ctx, &cord)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to post coordinate", zap.Error(err))
// 			return err
// 		}
// 	}

// 	var requestID string
// 	// 配車要求
// 	{
// 		req, err := userClient.AppPostRequest(ctx, &api.PostRequestReq{
// 			PickupCoordinate:      pickupCoordinate,
// 			DestinationCoordinate: destinationCoordinate,
// 		})
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to post request", zap.Error(err))
// 			return err
// 		}
// 		requestID = req.RequestID
// 	}

// 	// 配車要求の取得
// 	{
// 		req, err := userClient.AppGetRequest(ctx, requestID)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to get request", zap.Error(err))
// 			return err
// 		}
// 		if req.Status != api.RequestStatusDispatching {
// 			a.contestantLogger.Error("Ride status is not matching")
// 			return nil
// 		}
// 	}

// 	// マッチングしたらaccept
// 	{
// 		receivedRequest := <-receivedRequestCh
// 		_, err := chairClient.ChairPostRequestAccept(ctx, receivedRequest.RequestID)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to post accept", zap.Error(err))
// 			return err
// 		}
// 	}

// 	// 配車要求の状態を確認
// 	{
// 		req, err := userClient.AppGetRequest(ctx, requestID)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to get request", zap.Error(err))
// 			return err
// 		}
// 		if req.Status != api.RequestStatusDispatched {
// 			a.contestantLogger.Error("Ride status is not dispatched")
// 			return fmt.Errorf("request status is not dispatched")
// 		}
// 	}

// 	// ドライバーが配車位置に到着
// 	{
// 		_, err := chairClient.ChairPostCoordinate(ctx, &pickupCoordinate)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to post coordinate", zap.Error(err))
// 			return err
// 		}
// 	}

// 	// ユーザーが乗車した
// 	{
// 		_, err := chairClient.ChairPostRequestDepart(ctx, requestID)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to post depart", zap.Error(err))
// 			return err
// 		}
// 	}

// 	// 配車要求の状態を確認
// 	{
// 		req, err := userClient.AppGetRequest(ctx, requestID)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to get request", zap.Error(err))
// 			return err
// 		}

// 		if req.Status != api.RequestStatusCarrying {
// 			a.contestantLogger.Error("Ride status is not carrying")
// 			return fmt.Errorf("request status is not carrying")
// 		}
// 	}

// 	// ドライバーが目的地に到着
// 	{
// 		_, err := chairClient.ChairPostCoordinate(ctx, &destinationCoordinate)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to post coordinate", zap.Error(err))
// 			return err
// 		}
// 	}

// 	// 配車要求の状態を確認
// 	//{
// 	//	req, err := userClient.AppGetRequest(ctx, requestID)
// 	//	if err != nil {
// 	//		a.contestantLogger.Error("Failed to get request", zap.Error(err))
// 	//		return err
// 	//	}
// 	//
// 	//	if req.Status != api.RequestStatusArrived {
// 	//		a.contestantLogger.Error("Ride status is not arrived")
// 	//		return fmt.Errorf("request status is not arrived")
// 	//	}
// 	//}

// 	// ユーザーが降車・評価
// 	{
// 		_, err := userClient.AppPostRequestEvaluate(ctx, requestID, &api.EvaluateReq{
// 			rand.IntN(5) + 1,
// 		})
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to post evaluate", zap.Error(err))
// 			return err
// 		}
// 	}

// 	// 配車要求の状態を確認
// 	{
// 		req, err := userClient.AppGetRequest(ctx, requestID)
// 		if err != nil {
// 			a.contestantLogger.Error("Failed to get request", zap.Error(err))
// 			return err
// 		}

// 		if req.Status != api.RequestStatusCompleted {
// 			a.contestantLogger.Error("Ride status is not completed")
// 			return fmt.Errorf("request status is not completed")
// 		}
// 	}
// 	return nil
// }
