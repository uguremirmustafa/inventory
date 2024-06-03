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
	// var reqBody struct {
	// 	ParentID *int64 `json:"parentID"`
	// }
	// err := decode(r, &reqBody)
	// if err != nil {
	// 	return InvalidJSON()
	// }

	// var itemTypes []db.ItemType
	// if reqBody.ParentID != nil {
	// 	itemTypes, err = s.q.ListItemTypes(r.Context(), sql.NullInt64{Int64: *reqBody.ParentID, Valid: true})
	// 	if err != nil {
	// 		return NotFound()
	// 	}
	// } else {
	itemTypes, err := s.q.ListAllItemTypes(r.Context())
	if err != nil {
		return NotFound()
	}
	// }

	var itemTypeJsonList []ItemType
	for _, itemType := range itemTypes {
		itemTypeJsonList = append(itemTypeJsonList, *getItemTypeJson(itemType))
	}
	writeJson(w, http.StatusOK, itemTypeJsonList)
	return nil
}

type CreateItemTypeParams struct {
	Name string `json:"name"`
}

func (s *ItemTypeService) HandleCreateItemType(w http.ResponseWriter, r *http.Request) error {
	// TODO: request validation - empty string
	var body CreateItemTypeParams
	err := decode(r, body)
	if err != nil {
		return err
	}
	err = s.q.CreateItemType(r.Context(), body.Name)
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, "Success")
}

type ItemType struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	ParentID  *int64     `json:"parent_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

func getItemTypeJson(it db.ItemType) *ItemType {
	return &ItemType{
		ID:        it.ID,
		Name:      it.Name,
		ParentID:  utils.GetNilInt64(&it.ParentID),
		CreatedAt: utils.GetNilTime(&it.CreatedAt),
		UpdatedAt: utils.GetNilTime(&it.UpdatedAt),
		DeletedAt: utils.GetNilTime(&it.DeletedAt),
	}
}
