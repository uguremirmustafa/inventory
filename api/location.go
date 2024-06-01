package api

import (
	"database/sql"
	"fmt"
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
	GroupID     int64      `json:"group_id"`
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
	groupID := getUserActiveGroupID(w, r)
	locations, err := s.q.ListLocationsOfGroup(r.Context(), groupID)
	if err != nil {
		return NotFound()
	}

	var list []Location
	for _, item := range locations {
		list = append(list, *getLocationJson(item))
	}
	return writeJson(w, http.StatusOK, list)
}

type SaveLocationParams struct {
	Name        string `json:"name"`
	ImageUrl    string `json:"image_url"`
	Description string `json:"description"`
}

func (s *LocationService) HandleInsertUserLocation(w http.ResponseWriter, r *http.Request) error {
	groupID := getUserActiveGroupID(w, r)

	var reqBody SaveLocationParams
	err := decode(r, &reqBody)
	if err != nil {
		fmt.Println("here")
		return InvalidJSON()
	}
	locationID, err := s.q.InsertLocation(r.Context(), db.InsertLocationParams{
		Name:        reqBody.Name,
		ImageUrl:    sql.NullString{String: reqBody.ImageUrl, Valid: true},
		Description: sql.NullString{String: reqBody.Description, Valid: true},
		GroupID:     groupID,
	})
	if err != nil {
		return FailedInsert()
	}
	return writeJson(w, http.StatusOK, locationID)
}

func (s *LocationService) HandleUpdateUserLocation(w http.ResponseWriter, r *http.Request) error {
	groupID := getUserActiveGroupID(w, r)
	id, err := utils.GetPathID(r)
	if err != nil {
		return err
	}
	var reqBody SaveLocationParams
	err = decode(r, &reqBody)
	if err != nil {
		fmt.Println("here")
		return InvalidJSON()
	}
	locationID, err := s.q.UpdateLocation(r.Context(), db.UpdateLocationParams{
		ID:          id,
		Name:        reqBody.Name,
		ImageUrl:    sql.NullString{String: reqBody.ImageUrl, Valid: true},
		Description: sql.NullString{String: reqBody.Description, Valid: true},
		GroupID:     groupID,
	})
	if err != nil {
		return FailedUpdate()
	}
	return writeJson(w, http.StatusOK, locationID)
}

func (s *LocationService) HandleDeleteUserLocation(w http.ResponseWriter, r *http.Request) error {
	id, err := utils.GetPathID(r)
	if err != nil {
		return err
	}
	locationID, err := s.q.DeleteLocation(r.Context(), db.DeleteLocationParams{
		ID:        id,
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return FailedUpdate()
	}
	return writeJson(w, http.StatusOK, locationID)
}

func getLocationJson(l db.Location) *Location {
	return &Location{
		ID:          l.ID,
		Name:        l.Name,
		ImageUrl:    utils.GetNilString(&l.ImageUrl),
		Description: utils.GetNilString(&l.Description),
		GroupID:     l.GroupID,
		CreatedAt:   utils.GetNilTime(&l.CreatedAt),
		UpdatedAt:   utils.GetNilTime(&l.UpdatedAt),
		DeletedAt:   utils.GetNilTime(&l.DeletedAt),
	}
}
