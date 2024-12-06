package main

import (
	crand "crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {
	mux := setup()
	slog.Info("Listening on :8080")
	http.ListenAndServe(":8080", mux)
}

func setup() http.Handler {
	host := os.Getenv("ISUCON_DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("ISUCON_DB_PORT")
	if port == "" {
		port = "3306"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("failed to convert DB port number from ISUCON_DB_PORT environment variable into int: %v", err))
	}
	user := os.Getenv("ISUCON_DB_USER")
	if user == "" {
		user = "isucon"
	}
	password := os.Getenv("ISUCON_DB_PASSWORD")
	if password == "" {
		password = "isucon"
	}
	dbname := os.Getenv("ISUCON_DB_NAME")
	if dbname == "" {
		dbname = "isuride"
	}

	dbConfig := mysql.NewConfig()
	dbConfig.User = user
	dbConfig.Passwd = password
	dbConfig.Addr = net.JoinHostPort(host, port)
	dbConfig.Net = "tcp"
	dbConfig.DBName = dbname
	dbConfig.ParseTime = true

	_db, err := sqlx.Connect("mysql", dbConfig.FormatDSN())
	if err != nil {
		panic(err)
	}
	db = _db

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.HandleFunc("POST /api/initialize", postInitialize)

	// app handlers
	{
		mux.HandleFunc("POST /api/app/users", appPostUsers)

		authedMux := mux.With(appAuthMiddleware)
		authedMux.HandleFunc("POST /api/app/payment-methods", appPostPaymentMethods)
		authedMux.HandleFunc("GET /api/app/rides", appGetRides)
		authedMux.HandleFunc("POST /api/app/rides", appPostRides)
		authedMux.HandleFunc("POST /api/app/rides/estimated-fare", appPostRidesEstimatedFare)
		authedMux.HandleFunc("POST /api/app/rides/{ride_id}/evaluation", appPostRideEvaluatation)
		authedMux.HandleFunc("GET /api/app/notification", appGetNotificationSSE)
		authedMux.HandleFunc("GET /api/app/nearby-chairs", appGetNearbyChairs)
	}

	// owner handlers
	{
		mux.HandleFunc("POST /api/owner/owners", ownerPostOwners)

		authedMux := mux.With(ownerAuthMiddleware)
		authedMux.HandleFunc("GET /api/owner/sales", ownerGetSales)
		authedMux.HandleFunc("GET /api/owner/chairs", ownerGetChairs)
	}

	// chair handlers
	{
		mux.HandleFunc("POST /api/chair/chairs", chairPostChairs)

		authedMux := mux.With(chairAuthMiddleware)
		authedMux.HandleFunc("POST /api/chair/activity", chairPostActivity)
		authedMux.HandleFunc("POST /api/chair/coordinate", chairPostCoordinate)
		authedMux.HandleFunc("GET /api/chair/notification", chairGetNotificationSSE)
		authedMux.HandleFunc("POST /api/chair/rides/{ride_id}/status", chairPostRideStatus)
	}

	// internal handlers
	{
		mux.HandleFunc("GET /api/internal/matching", internalGetMatching)
	}

	return mux
}

type postInitializeRequest struct {
	PaymentServer string `json:"payment_server"`
}

type postInitializeResponse struct {
	Language string `json:"language"`
}

func postInitialize(w http.ResponseWriter, r *http.Request) {
	req := &postInitializeRequest{}
	if err := bindJSON(r, req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if out, err := exec.Command("../sql/init.sh").CombinedOutput(); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to initialize: %s: %w", string(out), err))
		return
	}

	if _, err := db.Exec("UPDATE settings SET value = ? WHERE name = 'payment_gateway_url'", req.PaymentServer); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, postInitializeResponse{Language: "go"})
}

type Coordinate struct {
	Latitude  int `json:"latitude"`
	Longitude int `json:"longitude"`
}

func bindJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	buf, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(buf)
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	buf, marshalError := json.Marshal(map[string]string{"message": err.Error()})
	if marshalError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"marshaling error failed"}`))
		return
	}
	w.Write(buf)

	fmt.Fprintln(os.Stderr, err)
}

func secureRandomStr(b int) string {
	k := make([]byte, b)
	if _, err := crand.Read(k); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", k)
}

func writeSSE(w http.ResponseWriter, data interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("data: " + string(buf) + "\n\n"))
	if err != nil {
		return err
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}

func appGetNotificationSSE(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*User)

	// Server Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	var lastRide *Ride
	var lastRideStatus string
	f := func() (respond bool, err error) {
		tx, err := db.Beginx()
		if err != nil {
			return false, err
		}
		defer tx.Rollback()

		ride := &Ride{}
		err = tx.Get(ride, `SELECT * FROM rides WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`, user.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			}
			return false, err

		}
		status, err := getLatestRideStatus(tx, ride.ID)
		if err != nil {
			return false, err

		}
		if lastRide != nil && ride.ID == lastRide.ID && status == lastRideStatus {
			return false, nil
		}

		fare, err := calculateDiscountedFare(tx, user.ID, ride, ride.PickupLatitude, ride.PickupLongitude, ride.DestinationLatitude, ride.DestinationLongitude)
		if err != nil {
			return false, err
		}

		chair := &Chair{}
		stats := appGetNotificationResponseChairStats{}
		if ride.ChairID.Valid {
			if err := tx.Get(chair, `SELECT * FROM chairs WHERE id = ?`, ride.ChairID); err != nil {
				return false, err
			}
			stats, err = getChairStats(tx, chair.ID)
			if err != nil {
				return false, err
			}
		}

		if err := writeSSE(w, &appGetNotificationResponseData{
			RideID: ride.ID,
			PickupCoordinate: Coordinate{
				Latitude:  ride.PickupLatitude,
				Longitude: ride.PickupLongitude,
			},
			DestinationCoordinate: Coordinate{
				Latitude:  ride.DestinationLatitude,
				Longitude: ride.DestinationLongitude,
			},
			Fare:   fare,
			Status: status,
			Chair: &appGetNotificationResponseChair{
				ID:    chair.ID,
				Name:  chair.Name,
				Model: chair.Model,
				Stats: stats,
			},
			CreatedAt: ride.CreatedAt.UnixMilli(),
			UpdateAt:  ride.UpdatedAt.UnixMilli(),
		}); err != nil {
			return false, err
		}
		lastRide = ride
		lastRideStatus = status

		return true, nil
	}

	// 初回送信を必ず行う
	respond, err := f()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if !respond {
		if err := writeSSE(w, nil); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}

	for {
		select {
		case <-r.Context().Done():
			w.WriteHeader(http.StatusOK)
			return

		default:
			respond, err := f()
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			if !respond {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func chairGetNotificationSSE(w http.ResponseWriter, r *http.Request) {
	chair := r.Context().Value("chair").(*Chair)

	// Server Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	var lastRide *Ride
	var lastRideStatus string
	f := func() (respond bool, err error) {
		found := true
		ride := &Ride{}
		tx, err := db.Beginx()
		if err != nil {
			return false, err
		}
		defer tx.Rollback()

		if _, err := tx.Exec("SELECT * FROM chairs WHERE id = ? FOR UPDATE", chair.ID); err != nil {
			return false, err
		}

		if err := tx.Get(ride, `SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC LIMIT 1`, chair.ID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				found = false
			} else {
				return false, err
			}
		}

		var status string
		if found {
			status, err = getLatestRideStatus(tx, ride.ID)
			if err != nil {
				return false, err
			}
		}

		if !found || status == "COMPLETED" {
			matched := &Ride{}
			if err := tx.Get(matched, `SELECT * FROM rides WHERE chair_id IS NULL ORDER BY created_at LIMIT 1 FOR UPDATE`); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return false, nil
				}
				return false, err
			}

			if _, err := tx.Exec("UPDATE rides SET chair_id = ? WHERE id = ?", chair.ID, matched.ID); err != nil {
				return false, err
			}

			if !found {
				ride = matched
			}
		}

		if lastRide != nil && ride.ID == lastRide.ID && status == lastRideStatus {
			return false, nil
		}

		user := &User{}
		err = tx.Get(user, "SELECT * FROM users WHERE id = ?", ride.UserID)
		if err != nil {
			return false, err
		}

		if err := tx.Commit(); err != nil {
			return false, err
		}

		if err := writeSSE(w, &chairGetNotificationResponseData{
			RideID: ride.ID,
			User: simpleUser{
				ID:   user.ID,
				Name: fmt.Sprintf("%s %s", user.Firstname, user.Lastname),
			},
			PickupCoordinate: Coordinate{
				Latitude:  ride.PickupLatitude,
				Longitude: ride.PickupLongitude,
			},
			DestinationCoordinate: Coordinate{
				Latitude:  ride.DestinationLatitude,
				Longitude: ride.DestinationLongitude,
			},
			Status: status,
		}); err != nil {
			return false, err
		}
		lastRide = ride
		lastRideStatus = status

		return true, nil
	}

	// 初回送信を必ず行う
	respond, err := f()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if !respond {
		if err := writeSSE(w, nil); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}

	for {
		select {
		case <-r.Context().Done():
			w.WriteHeader(http.StatusOK)
			return

		default:
			respond, err := f()
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			if !respond {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
