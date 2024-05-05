package api

import (
	"net/http"
	"time"

	"github.com/uguremirmustafa/inventory/db"
	"github.com/uguremirmustafa/inventory/utils"
)

type ItemType struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func handleItemTypeList(q *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itemTypes, err := q.ListItemTypes(r.Context())
		if err != nil {
			encode(w, http.StatusNotFound, "no items found")
		}
		var itemTypeJsonList []ItemType
		for _, itemType := range itemTypes {
			itemTypeJsonList = append(itemTypeJsonList, *getItemTypeJson(itemType))
		}
		encode(w, http.StatusOK, itemTypeJsonList)
	})
}

func getItemTypeJson(it db.ItemType) *ItemType {
	return &ItemType{
		ID:        it.ID,
		Name:      it.Name,
		CreatedAt: utils.GetNilTime(&it.CreatedAt),
		UpdatedAt: utils.GetNilTime(&it.UpdatedAt),
		DeletedAt: utils.GetNilTime(&it.DeletedAt),
	}
}
