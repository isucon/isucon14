package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
)

type postOwnerRegisterRequest struct {
	Name string `json:"name"`
}

type postOwnerRegisterResponse struct {
	ID string `json:"id"`
}

func ownerPostRegister(w http.ResponseWriter, r *http.Request) {
	req := &postOwnerRegisterRequest{}
	if err := bindJSON(r, req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ownerID := ulid.Make().String()

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, errors.New("some of required fields(name) are empty"))
		return
	}

	accessToken := secureRandomStr(32)
	_, err := db.Exec(
		"INSERT INTO owners (id, name, access_token, created_at, updated_at) VALUES (?, ?, ?, isu_now(), isu_now())",
		ownerID, req.Name, accessToken,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Path:     "/",
		Name:     "owner_session",
		Value:    accessToken,
		HttpOnly: true,
	})

	writeJSON(w, http.StatusCreated, &postOwnerRegisterResponse{
		ID: ownerID,
	})
}

type getOwnerSalesResponse struct {
	TotalSales int          `json:"total_sales"`
	Chairs     []ChairSales `json:"chairs"`
	Models     []ModelSales `json:"models"`
}

type ChairSales struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Sales int    `json:"sales"`
}

type ModelSales struct {
	Model string `json:"model"`
	Sales int    `json:"sales"`
}

func ownerGetSales(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value("owner").(*Owner)

	since := time.Unix(0, 0)
	until := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	if r.URL.Query().Get("since") != "" {
		parsed, err := time.Parse(time.RFC3339Nano, r.URL.Query().Get("since"))
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
		}
		since = parsed
	}
	if r.URL.Query().Get("until") != "" {
		parsed, err := time.Parse(time.RFC3339Nano, r.URL.Query().Get("until"))
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
		}
		until = parsed
	}

	chairs := []Chair{}
	if err := db.Select(&chairs, "SELECT * FROM chairs WHERE owner_id = ?", owner.ID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	res := getOwnerSalesResponse{
		TotalSales: 0,
	}

	modelSalesByModel := map[string]int{}

	for _, chair := range chairs {
		reqs := []RideRequest{}
		if err := db.Select(&reqs, "SELECT * FROM ride_requests WHERE chair_id = ? AND status = 'COMPLETED' AND updated_at BETWEEN ? AND ?", chair.ID, since, until); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		chairSales := sumSales(reqs)
		res.TotalSales += chairSales

		res.Chairs = append(res.Chairs, ChairSales{
			ID:    chair.ID,
			Name:  chair.Name,
			Sales: chairSales,
		})

		modelSalesByModel[chair.Model] += chairSales
	}

	modelSales := []ModelSales{}
	for model, sales := range modelSalesByModel {
		modelSales = append(modelSales, ModelSales{
			Model: model,
			Sales: sales,
		})
	}

	res.Models = modelSales

	writeJSON(w, http.StatusOK, res)
}

const (
	initialFare     = 500
	farePerDistance = 100
)

func sumSales(requests []RideRequest) int {
	sale := 0
	for _, req := range requests {
		sale += calculateSale(req)
	}
	return sale
}

func calculateSale(req RideRequest) int {
	latDiff := max(req.DestinationLatitude-req.PickupLatitude, req.PickupLatitude-req.DestinationLatitude)
	lonDiff := max(req.DestinationLongitude-req.PickupLongitude, req.PickupLongitude-req.DestinationLongitude)
	return initialFare + farePerDistance*(latDiff+lonDiff)
}

type ChairWithDetail struct {
	ID                     string       `db:"id"`
	OwnerID                string       `db:"owner_id"`
	Name                   string       `db:"name"`
	AccessToken            string       `db:"access_token"`
	Model                  string       `db:"model"`
	IsActive               bool         `db:"is_active"`
	CreatedAt              time.Time    `db:"created_at"`
	UpdatedAt              time.Time    `db:"updated_at"`
	TotalDistance          int          `db:"total_distance"`
	TotalDistanceUpdatedAt sql.NullTime `db:"total_distance_updated_at"`
}

type getOwnerChairResponse struct {
	Chairs []ownerChair `json:"chairs"`
}

type ownerChair struct {
	ID                     string     `json:"id"`
	Name                   string     `json:"name"`
	Model                  string     `json:"model"`
	Active                 bool       `json:"active"`
	RegisteredAt           time.Time  `json:"registered_at"`
	TotalDistance          int        `json:"total_distance"`
	TotalDistanceUpdatedAt *time.Time `json:"total_distance_updated_at,omitempty"`
}

func ownerGetChairs(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value("owner").(*Owner)

	chairs := []ChairWithDetail{}
	if err := db.Select(&chairs, `SELECT id,
       owner_id,
       name,
       access_token,
       model,
       is_active,
       created_at,
       updated_at,
       IFNULL(total_distance, 0) AS total_distance,
       total_distance_updated_at
FROM chairs
       LEFT JOIN (SELECT chair_id,
                          SUM(IFNULL(distance, 0)) AS total_distance,
                          MAX(created_at)          AS total_distance_updated_at
                   FROM (SELECT chair_id,
                                created_at,
                                ABS(latitude - LAG(latitude) OVER (PARTITION BY chair_id ORDER BY created_at)) +
                                ABS(longitude - LAG(longitude) OVER (PARTITION BY chair_id ORDER BY created_at)) AS distance
                         FROM chair_locations) tmp
                   GROUP BY chair_id) distance_table ON distance_table.chair_id = chairs.id
WHERE owner_id = ?
`, owner.ID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	res := getOwnerChairResponse{}
	for _, chair := range chairs {
		c := ownerChair{
			ID:            chair.ID,
			Name:          chair.Name,
			Model:         chair.Model,
			Active:        chair.IsActive,
			RegisteredAt:  chair.CreatedAt,
			TotalDistance: chair.TotalDistance,
		}
		if chair.TotalDistanceUpdatedAt.Valid {
			c.TotalDistanceUpdatedAt = &chair.TotalDistanceUpdatedAt.Time
		}
		res.Chairs = append(res.Chairs, c)
	}
	writeJSON(w, http.StatusOK, res)
}

type ownerChairDetail struct {
	ID                     string     `json:"id"`
	Name                   string     `json:"name"`
	Model                  string     `json:"model"`
	Active                 bool       `json:"active"`
	RegisteredAt           time.Time  `json:"registered_at"`
	TotalDistance          int        `json:"total_distance"`
	TotalDistanceUpdatedAt *time.Time `json:"total_distance_updated_at"`
}

func ownerGetChairDetail(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value("owner").(*Owner)
	chairID := r.PathValue("chair_id")

	chair := ChairWithDetail{}
	if err := db.Get(&chair, `SELECT id,
       owner_id,
       name,
       access_token,
       model,
       is_active,
       created_at,
       updated_at,
       IFNULL(total_distance, 0) AS total_distance,
       total_distance_updated_at
FROM chairs
       LEFT JOIN (SELECT chair_id,
                          SUM(IFNULL(distance, 0)) AS total_distance,
                          MAX(created_at)          AS total_distance_updated_at
                   FROM (SELECT chair_id,
                                created_at,
                                ABS(latitude - LAG(latitude) OVER (PARTITION BY chair_id ORDER BY created_at)) +
                                ABS(longitude - LAG(longitude) OVER (PARTITION BY chair_id ORDER BY created_at)) AS distance
                         FROM chair_locations) tmp
                   GROUP BY chair_id) distance_table ON distance_table.chair_id = chairs.id
WHERE owner_id = ? AND id = ?`, owner.ID, chairID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("chair not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	resp := ownerChairDetail{
		ID:            chair.ID,
		Name:          chair.Name,
		Model:         chair.Model,
		Active:        chair.IsActive,
		RegisteredAt:  chair.CreatedAt,
		TotalDistance: chair.TotalDistance,
	}
	if chair.TotalDistanceUpdatedAt.Valid {
		resp.TotalDistanceUpdatedAt = &chair.TotalDistanceUpdatedAt.Time
	}
	writeJSON(w, http.StatusOK, resp)
}
