package main

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"

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
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Sprintf("failed to convert DB port number from DB_PORT environment variable into int: %v", err))
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "isucon"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "isucon"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "isucon"
	}

	dbConfig := &mysql.Config{
		User:      user,
		Passwd:    password,
		Net:       "tcp",
		Addr:      net.JoinHostPort(host, port),
		DBName:    dbname,
		ParseTime: true,
	}

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
		mux.HandleFunc("POST /app/register", appPostRegister)

		authedMux := mux.With(appAuthMiddleware)
		authedMux.HandleFunc("POST /app/payment-methods", appPostPaymentMethods)
		authedMux.HandleFunc("POST /app/requests", appPostRequests)
		authedMux.HandleFunc("GET /app/requests/{request_id}", appGetRequest)
		authedMux.HandleFunc("POST /app/requests/{request_id}/evaluate", appPostRequestEvaluate)
		//authedMux.HandleFunc("GET /app/notification", appGetNotificationSSE)
		authedMux.HandleFunc("GET /app/notification", appGetNotification)
		authedMux.HandleFunc("POST /app/inquiry", appPostInquiry)
	}

	// chair handlers
	{
		mux.HandleFunc("POST /chair/register", chairPostRegister)

		authedMux := mux.With(chairAuthMiddleware)
		authedMux.HandleFunc("POST /chair/activate", chairPostActivate)
		authedMux.HandleFunc("POST /chair/deactivate", chairPostDeactivate)
		authedMux.HandleFunc("POST /chair/coordinate", chairPostCoordinate)
		//authedMux.HandleFunc("GET /chair/notification", chairGetNotificationSSE)
		authedMux.HandleFunc("GET /chair/notification", chairGetNotification)
		authedMux.HandleFunc("GET /chair/requests/{request_id}", chairGetRequest)
		authedMux.HandleFunc("POST /chair/requests/{request_id}/accept", chairPostRequestAccept)
		authedMux.HandleFunc("POST /chair/requests/{request_id}/deny", chairPostRequestDeny)
		authedMux.HandleFunc("POST /chair/requests/{request_id}/depart", chairPostRequestDepart)
		authedMux.HandleFunc("POST /chair/requests/{request_id}/payment", chairPostRequestPayment)
	}

	// admin
	{
		mux.HandleFunc("GET /admin/inquiries", adminGetInquiries)
		mux.HandleFunc("GET /admin/inquiries/{inquiry_id}", adminGetInquiry)
	}

	return mux
}

type postInitializeRequest struct {
	PaymentServer string `json:"payment_server"`
}

func postInitialize(w http.ResponseWriter, r *http.Request) {
	req := &postInitializeRequest{}
	if err := bindJSON(r, req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	tables := []string{
		"chair_locations",
		"ride_requests",
		"inquiries",
		"users",
		"chairs",
	}
	tx, err := db.Beginx()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	tx.MustExec("SET FOREIGN_KEY_CHECKS = 0")
	for _, table := range tables {
		_, err := tx.Exec("TRUNCATE TABLE " + table)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	tx.MustExec("SET FOREIGN_KEY_CHECKS = 1")
	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	paymentURL = req.PaymentServer

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(`{"language":"golang"}`))
}

type Coordinate struct {
	Latitude  int `json:"latitude"`
	Longitude int `json:"longitude"`
}

func bindJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func respondJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	buf, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(buf)
	return
}

func writeSSE(w http.ResponseWriter, event string, data interface{}) error {
	_, err := w.Write([]byte("event: " + event + "\n"))
	if err != nil {
		return err
	}

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

func respondError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	buf, marshalError := json.Marshal(map[string]string{"message": err.Error()})
	if marshalError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"marshaling error failed"}`))
		return
	}
	w.Write(buf)
	return
}

func secureRandomStr(b int) string {
	k := make([]byte, b)
	if _, err := crand.Read(k); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", k)
}
