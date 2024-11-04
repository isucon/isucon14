package main

import (
	"database/sql"
	"time"
)

type ChairModel struct {
	Name  string `db:"name"`
	Speed int    `db:"speed"`
}

type Chair struct {
	ID          string    `db:"id"`
	OwnerID     string    `db:"owner_id"`
	Name        string    `db:"name"`
	Model       string    `db:"model"`
	IsActive    bool      `db:"is_active"`
	AccessToken string    `db:"access_token"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type ChairLocation struct {
	ID        string    `db:"id"`
	ChairID   string    `db:"chair_id"`
	Latitude  int       `db:"latitude"`
	Longitude int       `db:"longitude"`
	CreatedAt time.Time `db:"created_at"`
}

type User struct {
	ID          string    `db:"id"`
	Username    string    `db:"username"`
	Firstname   string    `db:"firstname"`
	Lastname    string    `db:"lastname"`
	DateOfBirth string    `db:"date_of_birth"`
	AccessToken string    `db:"access_token"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type PaymentToken struct {
	UserID    string    `db:"user_id"`
	Token     string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
}

type RideRequest struct {
	ID                   string         `db:"id"`
	UserID               string         `db:"user_id"`
	ChairID              sql.NullString `db:"chair_id"`
	Status               string         `db:"status"`
	PickupLatitude       int            `db:"pickup_latitude"`
	PickupLongitude      int            `db:"pickup_longitude"`
	DestinationLatitude  int            `db:"destination_latitude"`
	DestinationLongitude int            `db:"destination_longitude"`
	Evaluation           *int           `db:"evaluation"`
	RequestedAt          time.Time      `db:"requested_at"`
	MatchedAt            *time.Time     `db:"matched_at"`
	DispatchedAt         *time.Time     `db:"dispatched_at"`
	RodeAt               *time.Time     `db:"rode_at"`
	ArrivedAt            *time.Time     `db:"arrived_at"`
	UpdatedAt            time.Time      `db:"updated_at"`
}

type Owner struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	AccessToken string    `db:"access_token"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
