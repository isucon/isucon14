package main

import (
	"time"
)

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

type Chair struct {
	ID          string    `db:"id"`
	ProviderID  string    `db:"provider_id"`
	Name        string    `db:"name"`
	AccessToken string    `db:"access_token"`
	Model       string    `db:"model"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Ride struct {
	ID                   string     `db:"id"`
	UserID               string     `db:"user_id"`
	DriverID             string     `db:"driver_id"`
	ChairID              *string    `db:"chair_id"`
	Status               string     `db:"status"`
	PickupLatitude       int        `db:"pickup_latitude"`
	PickupLongitude      int        `db:"pickup_longitude"`
	DestinationLatitude  int        `db:"destination_latitude"`
	DestinationLongitude int        `db:"destination_longitude"`
	Evaluation           *int       `db:"evaluation"`
	RequestedAt          time.Time  `db:"requested_at"`
	MatchedAt            *time.Time `db:"matched_at"`
	PickupArrivedAt      *time.Time `db:"pickup_arrived_at"`
	PickuppedAt          *time.Time `db:"pickupped_at"`
	DestinationArrivedAt *time.Time `db:"destination_arrived_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
}

type ChairLocation struct {
	ChairID   string    `db:"chair_id"`
	Latitude  int       `db:"latitude"`
	Longitude int       `db:"longitude"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Provider struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	AccessToken string    `db:"access_token"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
