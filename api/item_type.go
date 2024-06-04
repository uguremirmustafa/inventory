package api

import (
	"net/http"
	"time"

	"github.com/uguremirmustafa/inventory/db"
	"github.com/uguremirmustafa/inventory/utils"
)

type ItemTypeService struct {
	q *db.Queries
}

func NewItemTypeService(q *db.Queries) *ItemTypeService {
	return &ItemTypeService{
		q: q,
	}
}

func (s *ItemTypeService) HandleListItemTypes(w http.ResponseWriter, r *http.Request) error {
	itemTypes, err := s.q.ListItemTypes(r.Context())
	if err != nil {
		return NotFound()
	}
	var itemTypeJsonList []ItemType
	for _, itemType := range itemTypes {
		itemTypeJsonList = append(itemTypeJsonList, *getItemTypeJson(itemType))
	}
	return writeJson(w, http.StatusOK, itemTypeJsonList)
}

type CreateItemTypeParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *ItemTypeService) HandleCreateItemType(w http.ResponseWriter, r *http.Request) error {
	// TODO: request validation - empty string
	var body CreateItemTypeParams
	err := decode(r, body)
	if err != nil {
		return err
	}
	err = s.q.CreateItemType(r.Context(), db.CreateItemTypeParams{
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, "Success")
}

type ItemType struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IconClass   string     `json:"icon_class"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

func getItemTypeJson(it db.ItemType) *ItemType {
	return &ItemType{
		ID:          it.ID,
		Name:        it.Name,
		Description: it.Description,
		IconClass:   it.IconClass,
		CreatedAt:   utils.GetNilTime(&it.CreatedAt),
		UpdatedAt:   utils.GetNilTime(&it.UpdatedAt),
		DeletedAt:   utils.GetNilTime(&it.DeletedAt),
	}
}
