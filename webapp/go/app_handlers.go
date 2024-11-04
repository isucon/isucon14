package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type appPostRegisterRequest struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	DateOfBirth string `json:"date_of_birth"`
}

type appPostRegisterResponse struct {
	ID string `json:"id"`
}

func appPostRegister(w http.ResponseWriter, r *http.Request) {
	req := &appPostRegisterRequest{}
	if err := bindJSON(r, req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if req.Username == "" || req.FirstName == "" || req.LastName == "" || req.DateOfBirth == "" {
		writeError(w, http.StatusBadRequest, errors.New("required fields(username, firstname, lastname, date_of_birth) are empty"))
		return
	}

	userID := ulid.Make().String()
	accessToken := secureRandomStr(32)

	_, err := db.Exec(
		"INSERT INTO users (id, username, firstname, lastname, date_of_birth, access_token) VALUES (?, ?, ?, ?, ?, ?)",
		userID, req.Username, req.FirstName, req.LastName, req.DateOfBirth, accessToken,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     "app_session",
		Value:    accessToken,
		HttpOnly: true,
	})

	writeJSON(w, http.StatusCreated, &appPostRegisterResponse{
		ID: userID,
	})
}

type appPostPaymentMethodsRequest struct {
	Token string `json:"token"`
}

func appPostPaymentMethods(w http.ResponseWriter, r *http.Request) {
	req := &appPostPaymentMethodsRequest{}
	if err := bindJSON(r, req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if req.Token == "" {
		writeError(w, http.StatusBadRequest, errors.New("token is required but was empty"))
		return
	}

	user := r.Context().Value("user").(*User)

	_, err := db.Exec(
		`INSERT INTO payment_tokens (user_id, token) VALUES (?, ?)`,
		user.ID,
		req.Token,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type appPostRequestsRequest struct {
	PickupCoordinate      *Coordinate `json:"pickup_coordinate"`
	DestinationCoordinate *Coordinate `json:"destination_coordinate"`
}

type appPostRequestsResponse struct {
	RequestID string `json:"request_id"`
}

func appPostRequests(w http.ResponseWriter, r *http.Request) {
	req := &appPostRequestsRequest{}
	if err := bindJSON(r, req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if req.PickupCoordinate == nil || req.DestinationCoordinate == nil {
		writeError(w, http.StatusBadRequest, errors.New("required fields(pickup_coordinate, destination_coordinate) are empty"))
		return
	}

	user := r.Context().Value("user").(*User)
	requestID := ulid.Make().String()

	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	requestCount := 0
	if err := tx.Get(&requestCount, `SELECT COUNT(*) FROM ride_requests WHERE user_id = ? AND status NOT IN ('COMPLETED', 'CANCELED')`, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if requestCount > 0 {
		writeError(w, http.StatusConflict, errors.New("request already exists"))
		return
	}

	if _, err := tx.Exec(
		`INSERT INTO ride_requests (id, user_id, status, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude)
				  VALUES (?, ?, ?, ?, ?, ?, ?)`,
		requestID, user.ID, "MATCHING", req.PickupCoordinate.Latitude, req.PickupCoordinate.Longitude, req.DestinationCoordinate.Latitude, req.DestinationCoordinate.Longitude,
	); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusAccepted, &appPostRequestsResponse{
		RequestID: requestID,
	})
}

type recentRide struct {
	ID                    string        `json:"id"`
	PickupCoordinate      Coordinate    `json:"pickup_coordinate"`
	DestinationCoordinate Coordinate    `json:"destination_coordinate"`
	Distance              int           `json:"distance"`
	Duration              time.Duration `json:"duration"`
	Evaluation            int           `json:"evaluation"`
}

type appChairStats struct {
	// 最近の乗車履歴
	RecentRides []recentRide `json:"recent_rides"`

	// 累計の情報
	TotalRidesCount    int     `json:"total_rides_count"`
	TotalEvaluationAvg float64 `json:"total_evaluation_avg"`
}

type appChair struct {
	ID    string        `json:"id"`
	Name  string        `json:"name"`
	Model string        `json:"model"`
	Stats appChairStats `json:"stats"`
}

type appGetRequestResponse struct {
	RequestID             string     `json:"request_id"`
	PickupCoordinate      Coordinate `json:"pickup_coordinate"`
	DestinationCoordinate Coordinate `json:"destination_coordinate"`
	Status                string     `json:"status"`
	Chair                 *appChair  `json:"chair,omitempty"`
	CreatedAt             int64      `json:"created_at"`
	UpdateAt              int64      `json:"updated_at"`
}

func appGetRequest(w http.ResponseWriter, r *http.Request) {
	requestID := r.PathValue("request_id")

	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	rideRequest := &RideRequest{}
	err = tx.Get(
		rideRequest,
		`SELECT * FROM ride_requests WHERE id = ?`,
		requestID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("request not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	response := &appGetRequestResponse{
		RequestID:             rideRequest.ID,
		PickupCoordinate:      Coordinate{Latitude: rideRequest.PickupLatitude, Longitude: rideRequest.PickupLongitude},
		DestinationCoordinate: Coordinate{Latitude: rideRequest.DestinationLatitude, Longitude: rideRequest.DestinationLongitude},
		Status:                rideRequest.Status,
		CreatedAt:             rideRequest.RequestedAt.Unix(),
		UpdateAt:              rideRequest.UpdatedAt.Unix(),
	}

	if rideRequest.ChairID.Valid {
		chair := &Chair{}
		err := tx.Get(
			chair,
			`SELECT * FROM chairs WHERE id = ?`,
			rideRequest.ChairID,
		)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		stats, err := getChairStats(tx, chair.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		response.Chair = &appChair{
			ID:    chair.ID,
			Name:  chair.Name,
			Model: chair.Model,
			Stats: stats,
		}
	}

	writeJSON(w, http.StatusOK, response)
}

func getChairStats(tx *sqlx.Tx, chairID string) (appChairStats, error) {
	stats := appChairStats{RecentRides: make([]recentRide, 0)}

	// 最近の乗車履歴
	rideRequests := []RideRequest{}
	err := tx.Select(
		&rideRequests,
		`SELECT *
		 FROM ride_requests
		 WHERE ride_requests.chair_id = ? AND ride_requests.status = 'COMPLETED'
		 ORDER BY ride_requests.updated_at DESC`,
		chairID,
	)
	if err != nil {
		return stats, err
	}

	totalRideCount := len(rideRequests)
	totalEvaluation := 0.0
	for _, rideRequest := range rideRequests {
		chairLocations := []ChairLocation{}
		err := tx.Select(
			&chairLocations,
			`SELECT * FROM chair_locations WHERE chair_id = ? AND created_at BETWEEN ? AND ? ORDER BY created_at`,
			chairID, rideRequest.RequestedAt, rideRequest.UpdatedAt,
		)
		if err != nil {
			return stats, err
		}

		distance := 0
		lastLocation := ChairLocation{
			Latitude:  rideRequest.PickupLatitude,
			Longitude: rideRequest.PickupLongitude,
		}
		for _, location := range chairLocations {
			distance += calculateDistance(lastLocation.Latitude, lastLocation.Longitude, location.Latitude, location.Longitude)
			lastLocation = location
		}
		distance += calculateDistance(lastLocation.Latitude, lastLocation.Longitude, rideRequest.DestinationLatitude, rideRequest.DestinationLongitude)

		stats.RecentRides = append(stats.RecentRides, recentRide{
			ID:                    rideRequest.ID,
			PickupCoordinate:      Coordinate{Latitude: rideRequest.PickupLatitude, Longitude: rideRequest.PickupLongitude},
			DestinationCoordinate: Coordinate{Latitude: rideRequest.DestinationLatitude, Longitude: rideRequest.DestinationLongitude},
			Distance:              distance,
			Duration:              rideRequest.ArrivedAt.Sub(*rideRequest.RodeAt),
			Evaluation:            *rideRequest.Evaluation,
		})

		totalEvaluation += float64(*rideRequest.Evaluation)
	}

	// 5件以上の履歴がある場合は5件までにする
	if totalRideCount > 5 {
		stats.RecentRides = stats.RecentRides[:5]
	}

	stats.TotalRidesCount = totalRideCount
	if totalRideCount > 0 {
		stats.TotalEvaluationAvg = totalEvaluation / float64(totalRideCount)
	}

	return stats, nil
}

// マンハッタン距離を求める
func calculateDistance(aLatitude, aLongitude, bLatitude, bLongitude int) int {
	return abs(aLatitude-bLatitude) + abs(aLongitude-bLongitude)
}
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type appPostEvaluateRequest struct {
	Evaluation int `json:"evaluation"`
}

type appPostEvaluateResponse struct {
	Fare        int       `json:"fare"`
	CompletedAt time.Time `json:"completed_at"`
}

func appPostRequestEvaluate(w http.ResponseWriter, r *http.Request) {
	requestID := r.PathValue("request_id")

	postAppEvaluateRequest := &appPostEvaluateRequest{}
	if err := bindJSON(r, postAppEvaluateRequest); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if postAppEvaluateRequest.Evaluation < 1 || postAppEvaluateRequest.Evaluation > 5 {
		writeError(w, http.StatusBadRequest, errors.New("evaluation must be between 1 and 5"))
		return
	}

	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	rideRequest := &RideRequest{}
	if err := tx.Get(rideRequest, `SELECT * FROM ride_requests WHERE id = ?`, requestID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("request not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if rideRequest.Status != "ARRIVED" {
		writeError(w, http.StatusBadRequest, errors.New("not arrived yet"))
		return
	}

	result, err := tx.Exec(
		`UPDATE ride_requests SET evaluation = ?, status = ? WHERE id = ?`,
		postAppEvaluateRequest.Evaluation, "COMPLETED", requestID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if count, err := result.RowsAffected(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	} else if count == 0 {
		writeError(w, http.StatusNotFound, errors.New("request not found"))
		return
	}

	if err := tx.Get(rideRequest, `SELECT * FROM ride_requests WHERE id = ?`, requestID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("request not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	paymentToken := &PaymentToken{}
	if err := tx.Get(paymentToken, `SELECT * FROM payment_tokens WHERE user_id = ?`, rideRequest.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusBadRequest, errors.New("payment token not registered"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	paymentGatewayRequest := &paymentGatewayPostPaymentRequest{
		Amount: calculateSale(*rideRequest),
	}
	if err := requestPaymentGatewayPostPayment(paymentToken.Token, paymentGatewayRequest, func() ([]RideRequest, error) {
		rideRequests := []RideRequest{}
		if err := tx.Select(&rideRequests, `SELECT * FROM ride_requests WHERE user_id = ? ORDER BY requested_at ASC`, rideRequest.UserID); err != nil {
			return nil, err
		}
		return rideRequests, nil
	}); err != nil {
		if errors.Is(err, erroredUpstream) {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, &appPostEvaluateResponse{
		Fare:        calculateSale(*rideRequest),
		CompletedAt: rideRequest.UpdatedAt,
	})
}

type appGetNotificationResponse struct {
	RequestID             string     `json:"request_id"`
	PickupCoordinate      Coordinate `json:"pickup_coordinate"`
	DestinationCoordinate Coordinate `json:"destination_coordinate"`
	Status                string     `json:"status"`
	Chair                 *appChair  `json:"chair,omitempty"`
	CreatedAt             int64      `json:"created_at"`
	UpdateAt              int64      `json:"updated_at"`
}

func appGetNotification(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)

	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	rideRequest := &RideRequest{}
	if err := tx.Get(rideRequest, `SELECT * FROM ride_requests WHERE user_id = ? ORDER BY requested_at DESC LIMIT 1`, user.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	response := &appGetNotificationResponse{
		RequestID: rideRequest.ID,
		PickupCoordinate: Coordinate{
			Latitude:  rideRequest.PickupLatitude,
			Longitude: rideRequest.PickupLongitude,
		},
		DestinationCoordinate: Coordinate{
			Latitude:  rideRequest.DestinationLatitude,
			Longitude: rideRequest.DestinationLongitude,
		},
		Status:    rideRequest.Status,
		CreatedAt: rideRequest.RequestedAt.Unix(),
		UpdateAt:  rideRequest.UpdatedAt.Unix(),
	}

	if rideRequest.ChairID.Valid {
		chair := &Chair{}
		if err := tx.Get(chair, `SELECT * FROM chairs WHERE id = ?`, rideRequest.ChairID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		stats, err := getChairStats(tx, chair.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		response.Chair = &appChair{
			ID:    chair.ID,
			Name:  chair.Name,
			Model: chair.Model,
			Stats: stats,
		}
	}

	writeJSON(w, http.StatusOK, response)
}

func appGetNotificationSSE(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)

	// Server Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	var lastRideRequest *RideRequest
	for {
		select {
		case <-r.Context().Done():
			w.WriteHeader(http.StatusOK)
			return

		default:
			err := func() error {
				tx, err := db.Beginx()
				if err != nil {
					return err
				}
				defer tx.Rollback()

				rideRequest := &RideRequest{}
				err = tx.Get(rideRequest, `SELECT * FROM ride_requests WHERE user_id = ? ORDER BY requested_at DESC LIMIT 1`, user.ID)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						time.Sleep(100 * time.Millisecond)
						return nil
					}
					return err
				}
				if lastRideRequest != nil && rideRequest.ID == lastRideRequest.ID && rideRequest.Status == lastRideRequest.Status {
					time.Sleep(100 * time.Millisecond)
					return nil
				}

				chair := &Chair{}
				stats := appChairStats{}
				if rideRequest.ChairID.Valid {
					if err := tx.Get(chair, `SELECT * FROM chairs WHERE id = ?`, rideRequest.ChairID); err != nil {
						return err
					}
					stats, err = getChairStats(tx, chair.ID)
					if err != nil {
						return err
					}
				}

				if err := writeSSE(w, "matched", &appGetNotificationResponse{
					RequestID: rideRequest.ID,
					PickupCoordinate: Coordinate{
						Latitude:  rideRequest.PickupLatitude,
						Longitude: rideRequest.PickupLongitude,
					},
					DestinationCoordinate: Coordinate{
						Latitude:  rideRequest.DestinationLatitude,
						Longitude: rideRequest.DestinationLongitude,
					},
					Status: rideRequest.Status,
					Chair: &appChair{
						ID:    chair.ID,
						Name:  chair.Name,
						Model: chair.Model,
						Stats: stats,
					},
					CreatedAt: rideRequest.RequestedAt.Unix(),
					UpdateAt:  rideRequest.UpdatedAt.Unix(),
				}); err != nil {
					return err
				}
				lastRideRequest = rideRequest

				return nil
			}()
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}

		}
	}
}

type appGetNearbyChairsResponse struct {
	Chairs      []appChair `json:"chairs"`
	RetrievedAt int64      `json:"retrieved_at"`
}

func appGetNearbyChairs(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("latitude")
	lonStr := r.URL.Query().Get("longitude")
	distanceStr := r.URL.Query().Get("distance")
	if latStr == "" || lonStr == "" {
		writeError(w, http.StatusBadRequest, errors.New("latitude or longitude is empty"))
		return
	}

	lat, err := strconv.Atoi(latStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("latitude is invalid"))
		return
	}

	lon, err := strconv.Atoi(lonStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("longitude is invalid"))
		return
	}

	distance := 50
	if distanceStr != "" {
		distance, err = strconv.Atoi(distanceStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, errors.New("distance is invalid"))
			return
		}
	}

	coordinate := Coordinate{Latitude: lat, Longitude: lon}

	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	chairs := []Chair{}
	err = tx.Select(
		&chairs,
		`SELECT * FROM chairs`,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	nearbyChairs := []appChair{}
	for _, chair := range chairs {
		// 現在進行中のリクエストがある場合はスキップ
		rideRequest := &RideRequest{}
		err := tx.Get(
			rideRequest,
			`SELECT * FROM ride_requests WHERE chair_id = ? ORDER BY requested_at DESC LIMIT 1`,
			chair.ID,
		)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
		}
		if rideRequest.Status != "COMPLETED" {
			continue
		}

		// 5分以内に更新されている最新の位置情報を取得
		chairLocation := &ChairLocation{}
		err = tx.Get(
			chairLocation,
			`SELECT * FROM chair_locations WHERE chair_id = ? AND created_at > DATE_SUB(CURRENT_TIMESTAMP(6), INTERVAL 5 MINUTE) ORDER BY created_at DESC LIMIT 1`,
			chair.ID,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		if calculateDistance(coordinate.Latitude, coordinate.Longitude, chairLocation.Latitude, chairLocation.Longitude) <= distance {
			stats, err := getChairStats(tx, chair.ID)
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}

			nearbyChairs = append(nearbyChairs, appChair{
				ID:    chair.ID,
				Name:  chair.Name,
				Model: chair.Model,
				Stats: stats,
			})
		}
	}

	retrievedAt := &time.Time{}
	err = tx.Get(
		retrievedAt,
		`SELECT CURRENT_TIMESTAMP(6)`,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, &appGetNearbyChairsResponse{
		Chairs:      nearbyChairs,
		RetrievedAt: retrievedAt.Unix(),
	})
}
