package api

import (
	"net/http"
	"time"

	"github.com/uguremirmustafa/inventory/db"
	"github.com/uguremirmustafa/inventory/utils"
)

type Manufacturer struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	LogoUrl     *string    `json:"logo_url"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

func handleListManufacturer(q *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manufacturers, err := q.ListManufacturers(r.Context())
		if err != nil {
			writeJson(w, http.StatusNotFound, "no items found")
		}

		var list []Manufacturer
		for _, item := range manufacturers {
			list = append(list, *getManufacturerJson(item))
		}
		writeJson(w, http.StatusOK, list)
	})
}

func getManufacturerJson(it db.Manufacturer) *Manufacturer {
	return &Manufacturer{
		ID:          it.ID,
		Name:        it.Name,
		LogoUrl:     utils.GetNilString(&it.LogoUrl),
		Description: utils.GetNilString(&it.Description),
		CreatedAt:   utils.GetNilTime(&it.CreatedAt),
		UpdatedAt:   utils.GetNilTime(&it.UpdatedAt),
		DeletedAt:   utils.GetNilTime(&it.DeletedAt),
	}
}
