package worldclient

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/guregu/null/v5"
	"github.com/samber/lo"

	"github.com/isucon/isucon14/bench/benchmarker/webapp"
	"github.com/isucon/isucon14/bench/benchmarker/webapp/api"
	"github.com/isucon/isucon14/bench/benchmarker/world"
	"github.com/isucon/isucon14/bench/benchrun"
)

type userClient struct {
	ctx                       context.Context
	client                    *webapp.Client
	skipStaticFileSanityCheck bool
}

type ownerClient struct {
	ctx                       context.Context
	client                    *webapp.Client
	webappClientConfig        webapp.ClientConfig
	skipStaticFileSanityCheck bool
}

type chairClient struct {
	ctx    context.Context
	client *webapp.Client
}

type WorldClient struct {
	ctx                       context.Context
	webappClientConfig        webapp.ClientConfig
	skipStaticFileSanityCheck bool
}

func NewWorldClient(ctx context.Context, webappClientConfig webapp.ClientConfig, skipStaticFileSanityCheck bool) *WorldClient {
	return &WorldClient{
		ctx:                       ctx,
		webappClientConfig:        webappClientConfig,
		skipStaticFileSanityCheck: skipStaticFileSanityCheck,
	}
}

func (c *WorldClient) RegisterUser(ctx *world.Context, data *world.RegisterUserRequest, beforeRequest func(client world.UserClient) error) (*world.RegisterUserResponse, error) {
	client, err := webapp.NewClient(c.webappClientConfig)
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToCreateWebappClient, err)
	}
	userClient := &userClient{
		ctx:                       c.ctx,
		client:                    client,
		skipStaticFileSanityCheck: c.skipStaticFileSanityCheck,
	}

	err = beforeRequest(userClient)
	if err != nil {
		return nil, err
	}

	response, err := client.AppPostRegister(c.ctx, &api.AppPostUsersReq{
		Username:       data.UserName,
		Firstname:      data.FirstName,
		Lastname:       data.LastName,
		DateOfBirth:    data.DateOfBirth,
		InvitationCode: api.OptString{Set: len(data.InvitationCode) > 0, Value: data.InvitationCode},
	})
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToRegisterUser, err)
	}

	return &world.RegisterUserResponse{
		ServerUserID:   response.ID,
		InvitationCode: response.InvitationCode,
		Client:         userClient,
	}, nil
}

func (c *WorldClient) RegisterOwner(ctx *world.Context, data *world.RegisterOwnerRequest, beforeRequest func(client world.OwnerClient) error) (*world.RegisterOwnerResponse, error) {
	client, err := webapp.NewClient(c.webappClientConfig)
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToCreateWebappClient, err)
	}
	ownerClient := &ownerClient{
		ctx:                       c.ctx,
		client:                    client,
		webappClientConfig:        c.webappClientConfig,
		skipStaticFileSanityCheck: c.skipStaticFileSanityCheck,
	}

	err = beforeRequest(ownerClient)
	if err != nil {
		return nil, err
	}

	response, err := client.OwnerPostRegister(c.ctx, &api.OwnerPostOwnersReq{
		Name: data.Name,
	})
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToRegisterOwner, err)
	}

	return &world.RegisterOwnerResponse{
		ServerOwnerID:        response.ID,
		ChairRegisteredToken: response.ChairRegisterToken,
		Client:               ownerClient,
	}, nil
}

func (c *WorldClient) RegisterChair(ctx *world.Context, owner *world.Owner, data *world.RegisterChairRequest) (*world.RegisterChairResponse, error) {
	client, err := webapp.NewClient(c.webappClientConfig)
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToCreateWebappClient, err)
	}

	response, err := client.ChairPostRegister(c.ctx, &api.ChairPostChairsReq{
		Name:               data.Name,
		Model:              data.Model,
		ChairRegisterToken: owner.RegisteredData.ChairRegisterToken,
	})
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToRegisterChair, err)
	}

	return &world.RegisterChairResponse{
		ServerChairID: response.ID,
		ServerOwnerID: response.OwnerID,
		Client: &chairClient{
			ctx:    c.ctx,
			client: client,
		},
	}, nil
}

func (c *ownerClient) GetOwnerSales(ctx *world.Context, args *world.GetOwnerSalesRequest) (*world.GetOwnerSalesResponse, error) {
	params := api.OwnerGetSalesParams{}
	if !args.Since.IsZero() {
		params.Since.SetTo(args.Since.UnixMilli())
	}
	if !args.Until.IsZero() {
		params.Until.SetTo(args.Until.UnixMilli())
	}

	response, err := c.client.OwnerGetSales(c.ctx, &params)
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToGetOwnerSales, err)
	}

	return &world.GetOwnerSalesResponse{
		Total: response.TotalSales,
		Chairs: lo.Map(response.Chairs, func(v api.OwnerGetSalesOKChairsItem, _ int) *world.ChairSales {
			return &world.ChairSales{
				ID:    v.ID,
				Name:  v.Name,
				Sales: v.Sales,
			}
		}),
		Models: lo.Map(response.Models, func(v api.OwnerGetSalesOKModelsItem, _ int) *world.ChairSalesPerModel {
			return &world.ChairSalesPerModel{
				Model: v.Model,
				Sales: v.Sales,
			}
		}),
	}, nil
}

func (c *ownerClient) GetOwnerChairs(ctx *world.Context, args *world.GetOwnerChairsRequest) (*world.GetOwnerChairsResponse, error) {
	response, err := c.client.OwnerGetChairs(c.ctx)
	if err != nil {
		return nil, err
	}

	return &world.GetOwnerChairsResponse{Chairs: lo.Map(response.Chairs, func(v api.OwnerGetChairsOKChairsItem, _ int) *world.OwnerChair {
		return &world.OwnerChair{
			ID:                     v.ID,
			Name:                   v.Name,
			Model:                  v.Model,
			Active:                 v.Active,
			RegisteredAt:           time.UnixMilli(v.RegisteredAt),
			TotalDistance:          v.TotalDistance,
			TotalDistanceUpdatedAt: null.NewTime(time.UnixMilli(v.TotalDistanceUpdatedAt.Value), v.TotalDistanceUpdatedAt.Set),
		}
	})}, nil
}

func (c *ownerClient) BrowserAccess(ctx *world.Context, scenario benchrun.FrontendPathScenario) error {
	if c.skipStaticFileSanityCheck {
		return nil
	}
	return browserAccess(c.ctx, c.client, scenario)
}

func (c *chairClient) SendChairCoordinate(ctx *world.Context, chair *world.Chair) (*world.SendChairCoordinateResponse, error) {
	response, err := c.client.ChairPostCoordinate(c.ctx, &api.Coordinate{
		Latitude:  chair.Location.Current().X,
		Longitude: chair.Location.Current().Y,
	})
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToPostCoordinate, err)
	}

	return &world.SendChairCoordinateResponse{RecordedAt: time.UnixMilli(response.RecordedAt)}, nil
}

func (c *chairClient) SendAcceptRequest(ctx *world.Context, chair *world.Chair, req *world.Request) error {
	_, err := c.client.ChairPostRideStatus(c.ctx, req.ServerID, &api.ChairPostRideStatusReq{
		Status: api.ChairPostRideStatusReqStatusENROUTE,
	})
	if err != nil {
		return WrapCodeError(ErrorCodeFailedToPostAccept, err)
	}

	return nil
}

func (c *chairClient) SendDepart(ctx *world.Context, req *world.Request) error {
	_, err := c.client.ChairPostRideStatus(c.ctx, req.ServerID, &api.ChairPostRideStatusReq{
		Status: api.ChairPostRideStatusReqStatusCARRYING,
	})
	if err != nil {
		return WrapCodeError(ErrorCodeFailedToPostDepart, err)
	}

	return nil
}

func (c *chairClient) SendActivate(ctx *world.Context, chair *world.Chair) error {
	_, err := c.client.ChairPostActivity(c.ctx, &api.ChairPostActivityReq{
		IsActive: true,
	})
	if err != nil {
		return WrapCodeError(ErrorCodeFailedToPostActivate, err)
	}

	return nil
}

func (c *chairClient) ConnectChairNotificationStream(ctx *world.Context, chair *world.Chair, receiver world.NotificationReceiverFunc) (world.NotificationStream, error) {
	sseContext, cancel := context.WithCancel(c.ctx)

	go func() {
		for r, err := range c.client.ChairGetNotification(sseContext) {
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					// TODO ロギング
					slog.Debug(err.Error())
				}
				continue
			}
			if r.Data.Valid {
				data := r.Data.V
				var event world.NotificationEvent
				notificationEvent := world.ChairNotificationEvent{
					User: world.ChairNotificationEventUserPayload{
						ID:   data.User.ID,
						Name: data.User.Name,
					},
					Pickup:      world.C(data.PickupCoordinate.Latitude, data.PickupCoordinate.Longitude),
					Destination: world.C(data.DestinationCoordinate.Latitude, data.DestinationCoordinate.Longitude),
				}
				switch data.Status {
				case api.RideStatusMATCHING:
					event = &world.ChairNotificationEventMatched{
						ServerRequestID:        data.RideID,
						ChairNotificationEvent: notificationEvent,
					}
				case api.RideStatusENROUTE:
					event = &world.ChairNotificationEventDispatching{
						ServerRequestID:        data.RideID,
						ChairNotificationEvent: notificationEvent,
					}
				case api.RideStatusPICKUP:
					event = &world.ChairNotificationEventDispatched{
						ServerRequestID:        data.RideID,
						ChairNotificationEvent: notificationEvent,
					}
				case api.RideStatusCARRYING:
					event = &world.ChairNotificationEventCarrying{
						ServerRequestID:        data.RideID,
						ChairNotificationEvent: notificationEvent,
					}
				case api.RideStatusARRIVED:
					event = &world.ChairNotificationEventArrived{
						ServerRequestID:        data.RideID,
						ChairNotificationEvent: notificationEvent,
					}
				case api.RideStatusCOMPLETED:
					event = &world.ChairNotificationEventCompleted{
						ServerRequestID:        data.RideID,
						ChairNotificationEvent: notificationEvent,
					}
				}
				if event == nil {
					// 意図しない通知の種類は無視する
					continue
				}
				receiver(event)
			}
		}
	}()

	return &notificationConnectionImpl{
		close: cancel,
	}, nil
}

func (c *userClient) getInternalClient() *webapp.Client {
	return c.client
}

func (c *userClient) SendEvaluation(ctx *world.Context, req *world.Request, score int) (*world.SendEvaluationResponse, error) {
	res, err := c.client.AppPostRequestEvaluate(c.ctx, req.ServerID, &api.AppPostRideEvaluationReq{
		Evaluation: score,
	})
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToPostEvaluate, err)
	}

	return &world.SendEvaluationResponse{
		CompletedAt: time.UnixMilli(res.CompletedAt),
	}, nil
}

func (c *userClient) SendCreateRequest(ctx *world.Context, req *world.Request) (*world.SendCreateRequestResponse, error) {
	pickup := req.PickupPoint
	destination := req.DestinationPoint
	response, err := c.client.AppPostRequest(c.ctx, &api.AppPostRidesReq{
		PickupCoordinate: api.Coordinate{
			Latitude:  pickup.X,
			Longitude: pickup.Y,
		},
		DestinationCoordinate: api.Coordinate{
			Latitude:  destination.X,
			Longitude: destination.Y,
		},
	})
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToPostRequest, err)
	}

	return &world.SendCreateRequestResponse{ServerRequestID: response.RideID}, nil
}

func (c *userClient) RegisterPaymentMethods(ctx *world.Context, user *world.User) error {
	_, err := c.client.AppPostPaymentMethods(c.ctx, &api.AppPostPaymentMethodsReq{Token: user.PaymentToken})
	if err != nil {
		return WrapCodeError(ErrorCodeFailedToPostPaymentMethods, err)
	}
	return nil
}

func (c *userClient) GetRequests(ctx *world.Context) (*world.GetRequestsResponse, error) {
	res, err := c.client.AppGetRequests(c.ctx)
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToGetRequests, err)
	}

	requests := make([]*world.RequestHistory, len(res.Rides))
	for i, r := range res.Rides {
		requests[i] = &world.RequestHistory{
			ID: r.ID,
			PickupCoordinate: world.Coordinate{
				X: r.PickupCoordinate.Latitude,
				Y: r.PickupCoordinate.Longitude,
			},
			DestinationCoordinate: world.Coordinate{
				X: r.DestinationCoordinate.Latitude,
				Y: r.DestinationCoordinate.Longitude,
			},
			Chair: world.RequestHistoryChair{
				ID:    r.Chair.ID,
				Owner: r.Chair.Owner,
				Name:  r.Chair.Name,
				Model: r.Chair.Model,
			},
			Fare:        r.Fare,
			Evaluation:  r.Evaluation,
			RequestedAt: time.UnixMilli(r.RequestedAt),
			CompletedAt: time.UnixMilli(r.CompletedAt),
		}
	}

	return &world.GetRequestsResponse{
		Requests: requests,
	}, nil
}

func (c *userClient) GetEstimatedFare(ctx *world.Context, pickup world.Coordinate, dest world.Coordinate) (*world.GetEstimatedFareResponse, error) {
	res, err := c.client.AppPostRidesEstimatedFare(c.ctx, &api.AppPostRidesEstimatedFareReq{
		PickupCoordinate: api.Coordinate{
			Latitude:  pickup.X,
			Longitude: pickup.Y,
		},
		DestinationCoordinate: api.Coordinate{
			Latitude:  dest.X,
			Longitude: dest.Y,
		},
	})
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToPostRidesEstimatedFare, err)
	}
	return &world.GetEstimatedFareResponse{
		Fare:     res.Fare,
		Discount: res.Discount,
	}, nil
}

func (c *userClient) GetNearbyChairs(ctx *world.Context, current world.Coordinate, distance int) (*world.GetNearbyChairsResponse, error) {
	res, err := c.client.AppGetNearbyChairs(c.ctx, &api.AppGetNearbyChairsParams{
		Latitude:  current.X,
		Longitude: current.Y,
		Distance:  api.NewOptInt(distance),
	})
	if err != nil {
		return nil, WrapCodeError(ErrorCodeFailedToGetNearbyChairs, err)
	}
	return &world.GetNearbyChairsResponse{
		RetrievedAt: time.UnixMilli(res.RetrievedAt),
		Chairs: lo.Map(res.Chairs, func(chair api.AppGetNearbyChairsOKChairsItem, _ int) *world.AppChair {
			return &world.AppChair{
				ID:         chair.ID,
				Name:       chair.Name,
				Model:      chair.Model,
				Coordinate: world.C(chair.CurrentCoordinate.Latitude, chair.CurrentCoordinate.Longitude),
			}
		}),
	}, nil
}

func (c *userClient) ConnectUserNotificationStream(ctx *world.Context, user *world.User, receiver world.NotificationReceiverFunc) (world.NotificationStream, error) {
	sseContext, cancel := context.WithCancel(c.ctx)

	go func() {
		for r, err := range c.client.AppGetNotification(sseContext) {
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					// TODO ロギング
					slog.Debug(err.Error())
				}
				continue
			}
			if r.Data.Valid {
				data := r.Data.V
				var event world.NotificationEvent
				userNotificationEvent := world.UserNotificationEvent{
					Pickup:      world.C(data.PickupCoordinate.Latitude, data.PickupCoordinate.Longitude),
					Destination: world.C(data.DestinationCoordinate.Latitude, data.DestinationCoordinate.Longitude),
					Fare:        data.Fare,
					Chair:       nil,
				}
				if data.Chair.Set {
					userNotificationEvent.Chair = &world.UserNotificationEventChairPayload{
						ID:    data.Chair.Value.ID,
						Name:  data.Chair.Value.Name,
						Model: data.Chair.Value.Model,
						Stats: world.UserNotificationEventChairStatsPayload{
							TotalRidesCount:    data.Chair.Value.Stats.TotalRidesCount,
							TotalEvaluationAvg: data.Chair.Value.Stats.TotalEvaluationAvg,
						},
					}
				}
				switch data.Status {
				case api.RideStatusMATCHING:
					event = &world.UserNotificationEventMatching{
						ServerRequestID:       data.RideID,
						UserNotificationEvent: userNotificationEvent,
					}
				case api.RideStatusENROUTE:
					event = &world.UserNotificationEventDispatching{
						ServerRequestID:       data.RideID,
						UserNotificationEvent: userNotificationEvent,
					}
				case api.RideStatusPICKUP:
					event = &world.UserNotificationEventDispatched{
						ServerRequestID:       data.RideID,
						UserNotificationEvent: userNotificationEvent,
					}
				case api.RideStatusCARRYING:
					event = &world.UserNotificationEventCarrying{
						ServerRequestID:       data.RideID,
						UserNotificationEvent: userNotificationEvent,
					}
				case api.RideStatusARRIVED:
					event = &world.UserNotificationEventArrived{
						ServerRequestID:       data.RideID,
						UserNotificationEvent: userNotificationEvent,
					}
				case api.RideStatusCOMPLETED:
					event = &world.UserNotificationEventCompleted{
						ServerRequestID:       data.RideID,
						UserNotificationEvent: userNotificationEvent,
					}
				}
				if event == nil {
					// 意図しない通知の種類は無視する
					continue
				}
				receiver(event)
			}
		}
	}()

	return &notificationConnectionImpl{
		close: cancel,
	}, nil
}

type notificationConnectionImpl struct {
	close func()
}

func (c *notificationConnectionImpl) Close() {
	c.close()
}

func (c *userClient) BrowserAccess(ctx *world.Context, scenario benchrun.FrontendPathScenario) error {
	if c.skipStaticFileSanityCheck {
		return nil
	}
	return browserAccess(c.ctx, c.client, scenario)
}

func browserAccess(ctx context.Context, client *webapp.Client, scenario benchrun.FrontendPathScenario) error {
	paths := benchrun.FRONTEND_PATH_SCENARIOS[scenario]
	path := paths[len(paths)-1]

	// ハードナビゲーション
	if len(paths) == 1 {
		hash, err := client.StaticGetFileHash(ctx, path)
		if err != nil {
			return WrapCodeError(ErrorCodeFailedToGetStaticFile, err)
		}
		if hash != benchrun.FrontendHashesMap["index.html"] {
			return WrapCodeError(ErrorCodeInvalidContent, errors.New(path+"の内容が期待したものと一致しません"))
		}
	}

	filePaths := benchrun.FrontendPathScenarioFiles[scenario]
	for _, filePath := range filePaths {
		hash, err := client.StaticGetFileHash(ctx, filePath)
		if err != nil {
			return WrapCodeError(ErrorCodeFailedToGetStaticFile, err)
		}
		if hash != benchrun.FrontendHashesMap[filePath[1:]] {
			return WrapCodeError(ErrorCodeInvalidContent, errors.New(filePath+"の内容が期待したものと一致しません"))
		}
	}

	return nil
}
