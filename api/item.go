package api

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/uguremirmustafa/inventory/db"
	"github.com/uguremirmustafa/inventory/utils"
)

type ItemService struct {
	q  *db.Queries
	db *sql.DB
}

func NewItemService(q *db.Queries, db *sql.DB) *ItemService {
	return &ItemService{
		q:  q,
		db: db,
	}
}

type ItemRow struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Images      []string   `json:"images"`
}

func getItemRowJson(i db.ListItemsRow, images []string) *ItemRow {
	return &ItemRow{
		ID:          i.ItemID,
		Name:        i.ItemName,
		Description: *utils.GetNilString(&i.ItemDescription),
		CreatedAt:   utils.GetNilTime(&i.CreatedAt),
		UpdatedAt:   utils.GetNilTime(&i.UpdatedAt),
		Images:      images,
	}
}

func (s *ItemService) HandleListItems(w http.ResponseWriter, r *http.Request) error {
	groupID := getUserActiveGroupID(w, r)
	groupItems, err := s.q.ListItems(r.Context(), groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return writeJson(w, http.StatusOK, []ItemRow{})
		}
		return err
	}
	var list []ItemRow
	for _, item := range groupItems {
		images, err := s.q.ListItemImages(r.Context(), db.ListItemImagesParams{
			ItemID: item.ItemID,
			Limit:  3,
		})
		var imageUrls []string
		if err != nil {
			if err == sql.ErrNoRows {
				// No images found, proceed with an empty image list
				imageUrls = []string{}
			} else {
				slog.Error("Error fetching images", slog.Int64("itemID", item.ItemID))
				continue
			}
		} else {
			for _, image := range images {
				imageUrls = append(imageUrls, image.ImageUrl)
			}
		}
		list = append(list, *getItemRowJson(item, imageUrls))
	}
	return writeJson(w, http.StatusOK, list)
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
	return writeJson(w, http.StatusOK, list)
}

func (s *ItemService) HandleInsertUserItem(w http.ResponseWriter, r *http.Request) error {
	userID := getUserID(w, r)
	groupID := getUserActiveGroupID(w, r)
	var reqBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ItemTypeID  int64  `json:"item_type_id"`
		ImgUrl      string `json:"img_url"`
	}
	err := decode(r, &reqBody)
	if err != nil {
		return InvalidJSON()
	}
	tx, err := s.db.BeginTx(r.Context(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := s.q.WithTx(tx)

	// Insert a new user item within the transaction
	itemParams := db.InsertUserItemParams{
		Name:           reqBody.Name,
		Description:    sql.NullString{String: reqBody.Description, Valid: true},
		UserID:         userID,
		ItemTypeID:     1,
		ManufacturerID: sql.NullInt64{Int64: 1, Valid: true},
		GroupID:        groupID,
	}
	itemID, err := qtx.InsertUserItem(r.Context(), itemParams)
	if err != nil {
		return err
	}
	err = qtx.InsertItemImage(r.Context(), db.InsertItemImageParams{
		ItemID:   itemID,
		ImageUrl: reqBody.ImgUrl,
	})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("transaction error")
		return err
	}
	return writeJson(w, http.StatusOK, itemID)
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
