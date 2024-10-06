package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

type postAppRegisterRequest struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	DateOfBirth string `json:"date_of_birth"`
}

type postAppRegisterResponse struct {
	AccessToken string `json:"access_token"`
	ID          string `json:"id"`
}

func appPostRegister(w http.ResponseWriter, r *http.Request) {
	req := &postAppRegisterRequest{}
	if err := bindJSON(r, req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	userID := ulid.Make().String()

	if req.Username == "" || req.FirstName == "" || req.LastName == "" || req.DateOfBirth == "" {
		respondError(w, http.StatusBadRequest, errors.New("required fields(username, firstname, lastname, date_of_birth) are empty"))
		return
	}
	accessToken := secureRandomStr(32)
	_, err := db.Exec(
		"INSERT INTO users (id, username, firstname, lastname, date_of_birth, access_token, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, isu_now(), isu_now())",
		userID, req.Username, req.FirstName, req.LastName, req.DateOfBirth, accessToken,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusCreated, &postAppRegisterResponse{
		AccessToken: accessToken,
		ID:          userID,
	})
}

func appAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		if accessToken == "" {
			respondError(w, http.StatusUnauthorized, errors.New("access token is required"))
			return
		}

		user := &User{}
		err := db.Get(user, "SELECT * FROM users WHERE access_token = ?", accessToken)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondError(w, http.StatusUnauthorized, errors.New("invalid access token"))
				return
			}
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type postAppPaymentMethodsRequest struct {
	Token string `json:"token"`
}

func appPostPaymentMethods(w http.ResponseWriter, r *http.Request) {
	req := &postAppPaymentMethodsRequest{}
	if err := bindJSON(r, req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	user := r.Context().Value("user").(*User)

	_, err := db.Exec(
		`INSERT INTO payment_tokens (user_id, token, created_at) VALUES (?, ?, isu_now())`,
		user.ID,
		req.Token,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type postAppRequestsRequest struct {
	PickupCoordinate      *Coordinate `json:"pickup_coordinate"`
	DestinationCoordinate *Coordinate `json:"destination_coordinate"`
}

type postAppRequestsResponse struct {
	RequestID string `json:"request_id"`
}

func appPostRequests(w http.ResponseWriter, r *http.Request) {
	req := &postAppRequestsRequest{}
	if err := bindJSON(r, req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	user := r.Context().Value("user").(*User)

	if req.PickupCoordinate == nil || req.DestinationCoordinate == nil {
		respondError(w, http.StatusBadRequest, errors.New("required fields(pickup_coordinate, destination_coordinate) are empty"))
		return
	}
	requestID := ulid.Make().String()

	tx, err := db.Beginx()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	requestCount := 0
	if err := tx.Get(&requestCount, `SELECT COUNT(*) FROM ride_requests WHERE user_id = ? AND status NOT IN ('COMPLETED', 'CANCELED')`, user.ID); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if requestCount > 0 {
		respondError(w, http.StatusConflict, errors.New("request already exists"))
		return
	}

	if _, err := tx.Exec(
		`INSERT INTO ride_requests (id, user_id, status, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude, requested_at, updated_at)
				  VALUES (?, ?, ?, ?, ?, ?, ?, isu_now(), isu_now())`,
		requestID, user.ID, "MATCHING", req.PickupCoordinate.Latitude, req.PickupCoordinate.Longitude, req.DestinationCoordinate.Latitude, req.DestinationCoordinate.Longitude,
	); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusAccepted, &postAppRequestsResponse{
		RequestID: requestID,
	})
}

type simpleChair struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ChairModel string `json:"chair_model"`
	ChairNo    string `json:"chair_no"`
}

type getAppRequestResponse struct {
	RequestID                string      `json:"request_id"`
	PickupCoordinate         Coordinate  `json:"pickup_coordinate"`
	DestinationCoordinate    Coordinate  `json:"destination_coordinate"`
	Status                   string      `json:"status"`
	Chair                    simpleChair `json:"chair"`
	LastChairLocation        *Coordinate `json:"last_chair_location"`
	EstimatedPickupDuration  int         `json:"estimated_pickup_duration"`
	EstimatedArrivalDuration int         `json:"estimated_arrival_duration"`
	CreatedAt                int64       `json:"created_at"`
	UpdateAt                 int64       `json:"updated_at"`
}

func appGetRequest(w http.ResponseWriter, r *http.Request) {
	requestID := r.PathValue("request_id")

	rideRequest := &RideRequest{}
	err := db.Get(
		rideRequest,
		`SELECT * FROM ride_requests WHERE id = ?`,
		requestID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, errors.New("request not found"))
			return
		}
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	response := &getAppRequestResponse{
		RequestID:             rideRequest.ID,
		PickupCoordinate:      Coordinate{Latitude: rideRequest.PickupLatitude, Longitude: rideRequest.PickupLongitude},
		DestinationCoordinate: Coordinate{Latitude: rideRequest.DestinationLatitude, Longitude: rideRequest.DestinationLongitude},
		Status:                rideRequest.Status,
		CreatedAt:             rideRequest.RequestedAt.Unix(),
		UpdateAt:              rideRequest.UpdatedAt.Unix(),
	}

	chair := &Chair{}
	if rideRequest.ChairID.Valid {
		err := db.Get(
			chair,
			`SELECT * FROM chairs WHERE id = ?`,
			rideRequest.ChairID,
		)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		response.Chair = simpleChair{
			ID:         chair.ID,
			Name:       chair.Firstname + " " + chair.Lastname,
			ChairModel: chair.ChairModel,
			ChairNo:    chair.ChairNo,
		}
	}

	respondJSON(w, http.StatusOK, response)
}

type postAppEvaluateRequest struct {
	Evaluation int `json:"evaluation"`
}

func appPostRequestEvaluate(w http.ResponseWriter, r *http.Request) {
	requestID := r.PathValue("request_id")
	postAppEvaluateRequest := &postAppEvaluateRequest{}
	if err := bindJSON(r, postAppEvaluateRequest); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	result, err := db.Exec(
		`UPDATE ride_requests SET evaluation = ?, status = ?, updated_at = isu_now() WHERE id = ?`,
		postAppEvaluateRequest.Evaluation, "COMPLETED", requestID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if count, err := result.RowsAffected(); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	} else if count == 0 {
		respondError(w, http.StatusNotFound, errors.New("request not found"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func appGetNotification(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)

	rideRequest := &RideRequest{}
	tx, err := db.Beginx()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()
	if err := tx.Get(rideRequest, `SELECT * FROM ride_requests WHERE user_id = ? ORDER BY requested_at DESC LIMIT 1`, user.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		respondError(w, http.StatusInternalServerError, err)
	}

	pickupCoordinate := Coordinate{
		Latitude:  rideRequest.PickupLatitude,
		Longitude: rideRequest.PickupLongitude,
	}
	destinationCoordinate := Coordinate{
		Latitude:  rideRequest.DestinationLatitude,
		Longitude: rideRequest.DestinationLongitude,
	}

	var lastChairLocation *Coordinate = nil
	estimatedPickupDuration := -1
	estimatedArrivalDuration := -1

	chair := &Chair{}
	if rideRequest.ChairID.Valid {
		if err := tx.Get(chair, `SELECT * FROM chairs WHERE id = ?`, rideRequest.ChairID); err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		chairLocations := []ChairLocation{}
		if err := tx.Select(chairLocations, `SELECT * FROM chair_locations WHERE chair_id = ? ORDER BY created_at`, rideRequest.ChairID); err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		if len(chairLocations) > 0 {
			lastChairLocation = &Coordinate{
				Latitude:  chairLocations[0].Latitude,
				Longitude: chairLocations[0].Longitude,
			}

			// マンハッタン距離
			totalDistance := int64(0)
			totalSeconds := int64(0)
			for i, chairLocation := range chairLocations {
				if i == 0 {
					continue
				}

				previousChairLocation := chairLocations[i-1]
				// TODO: オーバーフローに対処する
				totalDistance += int64(calculateDistanceFromChairLocation(previousChairLocation, chairLocation))
				totalSeconds += chairLocation.CreatedAt.Sub(previousChairLocation.CreatedAt).Milliseconds()
			}
			averageSpeed := float64(totalDistance) / float64(totalSeconds)

			pickupDistance := calculateDistance(*lastChairLocation, pickupCoordinate)
			rideDistance := calculateDistance(pickupCoordinate, destinationCoordinate)

			estimatedPickupDuration = int(float64(pickupDistance) / averageSpeed)
			estimatedArrivalDuration = int(float64(rideDistance) / averageSpeed)
		}
	}

	// 最後の椅子の位置と椅子の速度を元に、配車と到着にどれだけかかるかを求める

	respondJSON(w, http.StatusOK, &getAppRequestResponse{
		RequestID:             rideRequest.ID,
		PickupCoordinate:      pickupCoordinate,
		DestinationCoordinate: destinationCoordinate,
		Status:                rideRequest.Status,
		Chair: simpleChair{
			ID:         chair.ID,
			Name:       chair.Firstname + " " + chair.Lastname,
			ChairModel: chair.ChairModel,
			ChairNo:    chair.ChairNo,
		},
		LastChairLocation:        lastChairLocation,
		EstimatedPickupDuration:  estimatedPickupDuration,
		EstimatedArrivalDuration: estimatedArrivalDuration,
		CreatedAt:                rideRequest.RequestedAt.Unix(),
		UpdateAt:                 rideRequest.UpdatedAt.Unix(),
	})
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
			rideRequest := &RideRequest{}
			err := db.Get(rideRequest, `SELECT * FROM ride_requests WHERE user_id = ? ORDER BY requested_at DESC LIMIT 1`, user.ID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				respondError(w, http.StatusInternalServerError, err)
				return
			}
			if lastRideRequest != nil && rideRequest.ID == lastRideRequest.ID && rideRequest.Status == lastRideRequest.Status {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			chair := &Chair{}
			if rideRequest.ChairID.Valid {
				if err := db.Get(chair, `SELECT * FROM chairs WHERE id = ?`, rideRequest.ChairID); err != nil {
					respondError(w, http.StatusInternalServerError, err)
					return
				}
			}

			if err := writeSSE(w, "matched", &getAppRequestResponse{
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
				Chair: simpleChair{
					ID:         chair.ID,
					Name:       chair.Firstname + " " + chair.Lastname,
					ChairModel: chair.ChairModel,
					ChairNo:    chair.ChairNo,
				},
				CreatedAt: rideRequest.RequestedAt.Unix(),
				UpdateAt:  rideRequest.UpdatedAt.Unix(),
			}); err != nil {
				respondError(w, http.StatusInternalServerError, err)
				return
			}
			lastRideRequest = rideRequest
		}
	}
}

type postAppInquiryRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func appPostInquiry(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)

	req := &postAppInquiryRequest{}
	if err := bindJSON(r, req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	_, err := db.Exec(
		`INSERT INTO inquiries (user_id, subject, body, created_at) VALUES (?, ?, ?, isu_now())`,
		user.ID, req.Subject, req.Body,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// マンハッタン距離を求める
func calculateDistance(coordinate1, coordinate2 Coordinate) int {
	return (coordinate1.Latitude - coordinate2.Latitude) + (coordinate1.Longitude - coordinate2.Longitude)
}

func calculateDistanceFromChairLocation(coordinate1, coordinate2 ChairLocation) int {
	return (coordinate1.Latitude - coordinate2.Latitude) + (coordinate1.Longitude - coordinate2.Longitude)
}
