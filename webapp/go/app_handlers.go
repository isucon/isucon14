package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
)

type postAppUsersRequest struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	DateOfBirth string `json:"date_of_birth"`
}

type postAppUsersResponse struct {
	AccessToken string `json:"access_token"`
	ID          string `json:"id"`
}

func postAppUsers(w http.ResponseWriter, r *http.Request) {
	req := &postAppUsersRequest{}
	if err := bindJSON(r, req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	userID := ulid.Make().String()

	if req.Username == "" || req.FirstName == "" || req.LastName == "" || req.DateOfBirth == "" {
		writeError(w, http.StatusBadRequest, errors.New("required fields(username, firstname, lastname, date_of_birth) are empty"))
		return
	}
	accessToken := secureRandomStr(32)
	_, err := db.Exec(
		"INSERT INTO users (id, username, firstname, lastname, date_of_birth, access_token, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, isu_now(), isu_now())",
		userID, req.Username, req.FirstName, req.LastName, req.DateOfBirth, accessToken,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, &postAppUsersResponse{
		AccessToken: accessToken,
		ID:          userID,
	})
}

type postAppPaymentMethodsRequest struct {
	Token string `json:"token"`
}

func postAppPaymentMethods(w http.ResponseWriter, r *http.Request) {
	req := &postAppPaymentMethodsRequest{}
	if err := bindJSON(r, req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user := r.Context().Value("user").(*User)

	_, err := db.Exec(
		`INSERT INTO payment_tokens (user_id, token, created_at) VALUES (?, ?, isu_now())`,
		user.ID,
		req.Token,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type postAppRidesRequest struct {
	PickupCoordinate      *Coordinate `json:"pickup_coordinate"`
	DestinationCoordinate *Coordinate `json:"destination_coordinate"`
}

type postAppRidesResponse struct {
	RideID string `json:"ride_id"`
}

func postAppRides(w http.ResponseWriter, r *http.Request) {
	req := &postAppRidesRequest{}
	if err := bindJSON(r, req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user := r.Context().Value("user").(*User)

	if req.PickupCoordinate == nil || req.DestinationCoordinate == nil {
		writeError(w, http.StatusBadRequest, errors.New("required fields(pickup_coordinate, destination_coordinate) are empty"))
		return
	}
	rideID := ulid.Make().String()

	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	requestCount := 0
	if err := tx.Get(&requestCount, `SELECT COUNT(*) FROM rides WHERE user_id = ? AND status NOT IN ('COMPLETED', 'CANCELED')`, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if requestCount > 0 {
		writeError(w, http.StatusConflict, errors.New("ride already exists"))
		return
	}

	if _, err := tx.Exec(
		`INSERT INTO rides (id, user_id, status, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude, requested_at, updated_at)
				  VALUES (?, ?, ?, ?, ?, ?, ?, isu_now(), isu_now())`,
		rideID, user.ID, "MATCHING", req.PickupCoordinate.Latitude, req.PickupCoordinate.Longitude, req.DestinationCoordinate.Latitude, req.DestinationCoordinate.Longitude,
	); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusAccepted, &postAppRidesResponse{
		RideID: rideID,
	})
}

type simpleChair struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Model string `json:"model"`
}

type getAppRidesRideIDResponse struct {
	RideID                string      `json:"ride_id"`
	PickupCoordinate      Coordinate  `json:"pickup_coordinate"`
	DestinationCoordinate Coordinate  `json:"destination_coordinate"`
	Status                string      `json:"status"`
	Chair                 simpleChair `json:"chair"`
	CreatedAt             int64       `json:"created_at"`
	UpdateAt              int64       `json:"updated_at"`
}

func getAppRidesRideID(w http.ResponseWriter, r *http.Request) {
	rideID := r.PathValue("ride_id")

	ride := &Ride{}
	err := db.Get(
		ride,
		`SELECT * FROM ridess WHERE id = ?`,
		rideID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("request not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	res := &getAppRidesRideIDResponse{
		RideID:                ride.ID,
		PickupCoordinate:      Coordinate{Latitude: ride.PickupLatitude, Longitude: ride.PickupLongitude},
		DestinationCoordinate: Coordinate{Latitude: ride.DestinationLatitude, Longitude: ride.DestinationLongitude},
		Status:                ride.Status,
		CreatedAt:             ride.RequestedAt.Unix(),
		UpdateAt:              ride.UpdatedAt.Unix(),
	}

	chair := &Chair{}
	if ride.ChairID != nil {
		err := db.Get(
			chair,
			`SELECT * FROM chairs WHERE id = ?`,
			*ride.ChairID,
		)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		res.Chair = simpleChair{
			ID:    chair.ID,
			Name:  chair.Name,
			Model: chair.Model,
		}
	}

	writeJSON(w, http.StatusOK, res)
}

type postAppRidesRideIDEvaluationRequest struct {
	Evaluation int `json:"evaluation"`
}

func postAppRidesRideIDEvaluation(w http.ResponseWriter, r *http.Request) {
	rideID := r.PathValue("ride_id")

	req := &postAppRidesRideIDEvaluationRequest{}
	if err := bindJSON(r, req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	result, err := db.Exec(
		`UPDATE ride SET evaluation = ?, status = ?, updated_at = isu_now() WHERE id = ?`,
		req.Evaluation, "COMPLETED", rideID)
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

	w.WriteHeader(http.StatusNoContent)
}

func getAppNotification(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)

	ride := &Ride{}
	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()
	if err := tx.Get(ride, `SELECT * FROM rides WHERE user_id = ? ORDER BY requested_at DESC LIMIT 1`, user.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	chair := &Chair{}
	if ride.ChairID != nil {
		if err := tx.Get(chair, `SELECT * FROM chairs WHERE id = ?`, *ride.ChairID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}

	writeJSON(w, http.StatusOK, &getAppRidesRideIDResponse{
		RideID: ride.ID,
		PickupCoordinate: Coordinate{
			Latitude:  ride.PickupLatitude,
			Longitude: ride.PickupLongitude,
		},
		DestinationCoordinate: Coordinate{
			Latitude:  ride.DestinationLatitude,
			Longitude: ride.DestinationLongitude,
		},
		Status: ride.Status,
		Chair: simpleChair{
			ID:    chair.ID,
			Name:  chair.Name,
			Model: chair.Model,
		},
		CreatedAt: ride.RequestedAt.Unix(),
		UpdateAt:  ride.UpdatedAt.Unix(),
	})
}

func appGetNotificationSSE(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)

	// Server Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	var lastRide *Ride
	for {
		select {
		case <-r.Context().Done():
			w.WriteHeader(http.StatusOK)
			return

		default:
			ride := &Ride{}
			err := db.Get(ride, `SELECT * FROM rides WHERE user_id = ? ORDER BY requested_at DESC LIMIT 1`, user.ID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					time.Sleep(100 * time.Millisecond)
					continue
				}
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			if lastRide != nil && ride.ID == lastRide.ID && ride.Status == lastRide.Status {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			chair := &Chair{}
			if ride.ChairID != nil {
				if err := db.Get(chair, `SELECT * FROM chairs WHERE id = ?`, *ride.ChairID); err != nil {
					writeError(w, http.StatusInternalServerError, err)
					return
				}
			}

			if err := writeSSE(w, "matched", &getAppRidesRideIDResponse{
				RideID: ride.ID,
				PickupCoordinate: Coordinate{
					Latitude:  ride.PickupLatitude,
					Longitude: ride.PickupLongitude,
				},
				DestinationCoordinate: Coordinate{
					Latitude:  ride.DestinationLatitude,
					Longitude: ride.DestinationLongitude,
				},
				Status: ride.Status,
				Chair: simpleChair{
					ID:    chair.ID,
					Name:  chair.Name,
					Model: chair.Model,
				},
				CreatedAt: ride.RequestedAt.Unix(),
				UpdateAt:  ride.UpdatedAt.Unix(),
			}); err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			lastRide = ride
		}
	}
}
