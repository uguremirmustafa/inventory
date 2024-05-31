package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/uguremirmustafa/inventory/db"
	"github.com/uguremirmustafa/inventory/utils"
)

type LocationService struct {
	q  *db.Queries
	db *sql.DB
}

func NewLocationService(q *db.Queries, db *sql.DB) *LocationService {
	return &LocationService{
		q:  q,
		db: db,
	}
}

type Location struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	ImageUrl    *string    `json:"image_url"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

func (s *LocationService) HandleGetLocation(w http.ResponseWriter, r *http.Request) error {
	locationID, err := utils.GetPathID(r)
	if err != nil {
		return err
	}
	location, err := s.q.GetLocation(r.Context(), locationID)
	if err != nil {
		return NotFound()
	}
	return writeJson(w, http.StatusOK, getLocationJson(location))
}

func (s *LocationService) HandleListUserLocations(w http.ResponseWriter, r *http.Request) error {
	userID := getUserID(w, r)
	locations, err := s.q.ListLocationsOfUser(r.Context(), userID)
	if err != nil {
		return NotFound()
	}

	var list []Location
	for _, item := range locations {
		list = append(list, *getLocationJson(item))
	}
	return writeJson(w, http.StatusOK, list)
}

func getLocationJson(l db.Location) *Location {
	return &Location{
		ID:          l.ID,
		Name:        l.Name,
		ImageUrl:    utils.GetNilString(&l.ImageUrl),
		Description: utils.GetNilString(&l.Description),
		CreatedAt:   utils.GetNilTime(&l.CreatedAt),
		UpdatedAt:   utils.GetNilTime(&l.UpdatedAt),
		DeletedAt:   utils.GetNilTime(&l.DeletedAt),
	}
}
