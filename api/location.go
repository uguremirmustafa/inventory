package api

import (
	"net/http"
	"time"

	"github.com/uguremirmustafa/inventory/db"
	"github.com/uguremirmustafa/inventory/utils"
)

type Location struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	ImageUrl    *string    `json:"image_url"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

func handleListLocation(q *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locations, err := q.ListLocations(r.Context())
		if err != nil {
			writeJson(w, http.StatusNotFound, "no items found")
		}

		var list []Location
		for _, item := range locations {
			list = append(list, *getLocationJson(item))
		}
		writeJson(w, http.StatusOK, list)
	})
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
