package api

import (
	"net/http"
	"time"

	"github.com/uguremirmustafa/inventory/db"
	"github.com/uguremirmustafa/inventory/utils"
)

type ItemService struct {
	q *db.Queries
}

func NewItemService(q *db.Queries) *ItemService {
	return &ItemService{
		q: q,
	}
}

func (s *ItemService) HandleListUserItem(w http.ResponseWriter, r *http.Request) error {
	userID := getUserID(w, r)
	userItems, err := s.q.ListUserItems(r.Context(), userID)
	if err != nil {
		return NotFound()
	}
	var list []UserItemRow
	for _, item := range userItems {
		list = append(list, *getUserItemJson(item))
	}
	encode(w, http.StatusOK, list)
	return nil
}

type UserItemRow struct {
	ItemID              int64      `json:"item_id"`
	ItemName            string     `json:"item_name"`
	ItemDescription     *string    `json:"item_description"`
	UserName            *string    `json:"user_name"`
	UserEmail           *string    `json:"user_email"`
	ItemTypeID          *int64     `json:"item_type_id"`
	ItemTypeName        *string    `json:"item_type_name"`
	ManufacturerID      *int64     `json:"manufacturer_id"`
	ManufacturerName    *string    `json:"manufacturer_name"`
	ItemInfoID          *int64     `json:"item_info_id"`
	PurchaseDate        *time.Time `json:"purchase_date"`
	PurchaseLocation    *string    `json:"purchase_location"`
	Price               *int64     `json:"price"`
	ExpirationDate      *time.Time `json:"expiration_date"`
	LastUsed            *time.Time `json:"last_used"`
	LocationID          *int64     `json:"location_id"`
	LocationName        *string    `json:"location_name"`
	LocationDescription *string    `json:"location_description"`
	LocationImg         *string    `json:"location_img"`
}

func getUserItemJson(l db.ListUserItemsRow) *UserItemRow {
	return &UserItemRow{
		ItemID:              l.ItemID,
		ItemName:            l.ItemName,
		ItemDescription:     utils.GetNilString(&l.ItemDescription),
		UserName:            utils.GetNilString(&l.UserName),
		UserEmail:           utils.GetNilString(&l.UserEmail),
		ItemTypeID:          utils.GetNilInt64(&l.ItemTypeID),
		ItemTypeName:        utils.GetNilString(&l.ItemTypeName),
		ManufacturerID:      utils.GetNilInt64(&l.ManufacturerID),
		ManufacturerName:    utils.GetNilString(&l.ManufacturerName),
		ItemInfoID:          utils.GetNilInt64(&l.ItemInfoID),
		PurchaseDate:        utils.GetNilTime(&l.PurchaseDate),
		PurchaseLocation:    utils.GetNilString(&l.PurchaseLocation),
		Price:               utils.GetNilInt64(&l.Price),
		ExpirationDate:      utils.GetNilTime(&l.ExpirationDate),
		LastUsed:            utils.GetNilTime(&l.LastUsed),
		LocationID:          utils.GetNilInt64(&l.LocationID),
		LocationName:        utils.GetNilString(&l.LocationName),
		LocationDescription: utils.GetNilString(&l.LocationDescription),
		LocationImg:         utils.GetNilString(&l.LocationImg),
	}
}
