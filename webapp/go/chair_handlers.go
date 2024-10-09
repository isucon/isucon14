package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
)

type postChairChairsRequest struct {
	Name  string `json:"name"`
	Model string `json:"model"`
}

type postChairChairsResponse struct {
	AccessToken string `json:"access_token"`
	ID          string `json:"id"`
}

func postChairChairs(w http.ResponseWriter, r *http.Request) {
	provider := r.Context().Value("provider").(*Provider)

	req := &postChairChairsRequest{}
	if err := bindJSON(r, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chairID := ulid.Make().String()

	if req.Name == "" || req.Model == "" {
		writeError(w, http.StatusBadRequest, errors.New("some of required fields(name, model) are empty"))
		return
	}

	accessToken := secureRandomStr(32)
	_, err := db.Exec(
		"INSERT INTO chairs (id, provider_id, name, model, is_active, access_token, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, isu_now(), isu_now())",
		chairID, provider.ID, req.Name, req.Model, false, accessToken,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, &postChairChairsResponse{
		AccessToken: accessToken,
		ID:          chairID,
	})
}

type postChairActivityRequest struct {
	IsActive bool `json:"is_active"`
}

func postChairActivity(w http.ResponseWriter, r *http.Request) {
	chair := r.Context().Value("chair").(*Chair)

	req := &postChairActivityRequest{}
	if err := bindJSON(r, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE chairs SET is_active = ?, updated_at = isu_now() WHERE id = ?", req.IsActive, chair.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func postChairCoordinate(w http.ResponseWriter, r *http.Request) {
	req := &Coordinate{}
	if err := bindJSON(r, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chair := r.Context().Value("chair").(*Chair)

	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()
	chairLocationID := ulid.Make().String()
	if _, err := tx.Exec(
		`INSERT INTO chair_locations (id, chair_id, latitude, longitude, created_at) VALUES (?, ?, ?, ?, isu_now())`,
		chairLocationID, chair.ID, req.Latitude, req.Longitude,
	); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	ride := &Ride{}
	if err := tx.Get(ride, `SELECT * FROM rides WHERE chair_id = ? AND status NOT IN ('COMPLETED', 'CANCELED')`, chair.ID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		if req.Latitude == ride.PickupLatitude && req.Longitude == ride.PickupLongitude {
			if _, err := tx.Exec("UPDATE rides SET status = 'PICKUP', pickup_arrived_at = isu_now(), updated_at = isu_now() WHERE id = ? AND status = 'ENROUTE'", ride.ID); err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
		}

		if req.Latitude == ride.DestinationLatitude && req.Longitude == ride.DestinationLongitude {
			if _, err := tx.Exec("UPDATE rides SET status = 'ARRIVED', destination_arrived_at = isu_now(), updated_at = isu_now() WHERE id = ? AND status = 'CARRYING'", ride.ID); err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getChairNotification(w http.ResponseWriter, r *http.Request) {
	chair := r.Context().Value("chair").(*Chair)
	found := true
	ride := &Ride{}
	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	if _, err := tx.Exec("SELECT * FROM chairs WHERE id = ? FOR UPDATE", chair.ID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Get(ride, `SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC LIMIT 1`, chair.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			found = false
		} else {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}

	if !found || ride.Status == "COMPLETED" || ride.Status == "CANCELED" {
		rideToMatch := &Ride{}
		if err := tx.Get(rideToMatch, `SELECT * FROM rides WHERE status = 'MATCHING' AND chair_id IS NULL ORDER BY RAND() LIMIT 1 FOR UPDATE`); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		if _, err := tx.Exec("UPDATE rides SET chair_id = ?, matched_at = isu_now(), updated_at = isu_now() WHERE id = ?", chair.ID, rideToMatch.ID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		if !found {
			ride = rideToMatch
		}
	}

	user := &User{}
	err = tx.Get(user, "SELECT * FROM users WHERE id = ?", ride.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, &getChairRequestResponse{
		RideID: ride.ID,
		User: simpleUser{
			ID:   user.ID,
			Name: fmt.Sprintf("%s %s", user.Firstname, user.Lastname),
		},
		DestinationCoordinate: Coordinate{
			Latitude:  ride.DestinationLatitude,
			Longitude: ride.DestinationLongitude,
		},
		Status: ride.Status,
	})
}

func chairGetNotificationSSE(w http.ResponseWriter, r *http.Request) {
	chair := r.Context().Value("chair").(*Chair)

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
			err := func() error {
				found := true
				ride := &Ride{}
				tx, err := db.Beginx()
				if err != nil {
					return err
				}
				defer tx.Rollback()

				if _, err := tx.Exec("SELECT * FROM chairs WHERE id = ? FOR UPDATE", chair.ID); err != nil {
					return err
				}

				if err := tx.Get(ride, `SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC LIMIT 1`, chair.ID); err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						found = false
					} else {
						return err
					}
				}

				if !found || ride.Status == "COMPLETED" || ride.Status == "CANCELED" {
					rideToMatch := &Ride{}
					if err := tx.Get(rideToMatch, `SELECT * FROM rides WHERE status = 'MATCHING' AND chair_id IS NULL ORDER BY RAND() LIMIT 1 FOR UPDATE`); err != nil {
						if errors.Is(err, sql.ErrNoRows) {
							return nil
						}
						return err
					}

					if _, err := tx.Exec("UPDATE rides SET chair_id = ?, matched_at = isu_now(), updated_at = isu_now() WHERE id = ?", chair.ID, rideToMatch.ID); err != nil {
						return err
					}

					if !found {
						ride = rideToMatch
					}
				}

				if lastRide != nil && ride.ID == lastRide.ID && ride.Status == lastRide.Status {
					return nil
				}

				user := &User{}
				err = tx.Get(user, "SELECT * FROM users WHERE id = ?", ride.UserID)
				if err != nil {
					return err
				}

				if err := tx.Commit(); err != nil {
					return err
				}

				if err := writeSSE(w, "matched", &getChairRequestResponse{
					RideID: ride.ID,
					User: simpleUser{
						ID:   user.ID,
						Name: fmt.Sprintf("%s %s", user.Firstname, user.Lastname),
					},
					DestinationCoordinate: Coordinate{
						Latitude:  ride.DestinationLatitude,
						Longitude: ride.DestinationLongitude,
					},
					Status: ride.Status,
				}); err != nil {
					return err
				}
				lastRide = ride

				return nil
			}()

			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

type simpleUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type getChairRequestResponse struct {
	RideID                string     `json:"ride_id"`
	User                  simpleUser `json:"user"`
	DestinationCoordinate Coordinate `json:"destination_coordinate"`
	Status                string     `json:"status"`
}

func getChairRidesRideID(w http.ResponseWriter, r *http.Request) {
	rideID := r.PathValue("ride_id")

	ride := &Ride{}
	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	if err := tx.Get(ride, "SELECT * FROM rides WHERE id = ?", rideID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("ride not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	user := &User{}
	if err := tx.Get(user, "SELECT * FROM users WHERE id = ?", ride.UserID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, &getChairRequestResponse{
		RideID: ride.ID,
		User: simpleUser{
			ID:   user.ID,
			Name: fmt.Sprintf("%s %s", user.Firstname, user.Lastname),
		},
		DestinationCoordinate: Coordinate{
			Latitude:  ride.DestinationLatitude,
			Longitude: ride.DestinationLongitude,
		},
		Status: ride.Status,
	})
}

type postChairRidesRideIDStatusRequest struct {
	Status string `json:"status"`
}

func postChairRidesRideIDStatus(w http.ResponseWriter, r *http.Request) {
	rideID := r.PathValue("ride_id")

	req := &postChairRidesRideIDStatusRequest{}
	if err := bindJSON(r, req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	chair := r.Context().Value("chair").(*Chair)

	ride := &Ride{}
	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	if err := tx.Get(ride, "SELECT * FROM rides WHERE id = ?", rideID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("ride not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if ride.ChairID != nil && *ride.ChairID != chair.ID {
		writeError(w, http.StatusBadRequest, errors.New("not assigned to this ride"))
		return
	}

	switch req.Status {
	case "MATCHING":
		if _, err := tx.Exec("UPDATE rides SET chair_id = NULL, status = 'MATCHING', matched_at = NULL, updated_at = isu_now() WHERE id = ?", rideID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	case "ENROUTE":
		if _, err := tx.Exec("UPDATE rides SET status = 'ENROUTE', updated_at = isu_now() WHERE id = ?", rideID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	case "CARRYING":
		if ride.Status != "PICKUP" {
			writeError(w, http.StatusBadRequest, errors.New("invalid status transition"))
			return
		}
		if _, err := tx.Exec("UPDATE rides SET status = 'CARRYING', updated_at = isu_now() WHERE id = ?", rideID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	default:
		writeError(w, http.StatusBadRequest, errors.New("invalid status"))
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func postChairRidesRideIDPayment(w http.ResponseWriter, r *http.Request) {
	rideID := r.PathValue("ride_id")

	tx, err := db.Beginx()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	ride := &Ride{}
	if err := tx.Get(ride, `SELECT * FROM rides WHERE id = ?`, rideID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("request not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	paymentToken := &PaymentToken{}
	if err := tx.Get(paymentToken, `SELECT * FROM payment_tokens WHERE user_id = ?`, ride.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusBadRequest, errors.New("payment token not registered"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if ride.Status == "COMPLETED" {
		writeError(w, http.StatusBadRequest, errors.New("already paid"))
		return
	}
	if ride.Status != "ARRIVED" {
		writeError(w, http.StatusBadRequest, errors.New("not arrived yet"))
		return
	}

	paymentGatewayRequest := &paymentGatewayPostPaymentRequest{
		Token: paymentToken.Token,
		// TODO: calculate payment amount
		Amount: 100,
	}
	if err := requestPaymentGatewayPostPayment(paymentGatewayRequest); err != nil {
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

	w.WriteHeader(http.StatusNoContent)
}
