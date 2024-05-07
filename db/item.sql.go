// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: item.sql

package db

import (
	"context"
	"database/sql"
)

const listUserItems = `-- name: ListUserItems :many
SELECT 
    i.id AS item_id,
	i.name AS item_name,
	i.description AS item_description,
	u.name AS user_name,
	u.email AS user_email,
	it.id AS item_type_id,
	it.name AS item_type_name,
	man.id AS manufacturer_id,
	man.name AS manufacturer_name,
	ii.id AS item_info_id,
	ii.purchase_date AS purchase_date,
	ii.purchase_location AS purchase_location,
	ii.price AS price,
	ii.expiration_date AS expiration_date,
	ii.last_used AS last_used,
	loc.id AS location_id,
	loc.name AS location_name,
	loc.description AS location_description,
	loc.image_url AS location_img
FROM item i
LEFT JOIN users u ON i.user_id = u.id AND u.deleted_at is null
LEFT JOIN item_type it ON i.item_type_id = it.id AND it.deleted_at is null
LEFT JOIN manufacturer man ON i.manufacturer_id = man.id AND man.deleted_at is null
LEFT JOIN item_info ii ON ii.item_id = i.id AND ii.deleted_at is null
LEFT JOIN location loc ON ii.location_id = loc.id AND loc.deleted_at is null
WHERE u.id = $1 AND i.deleted_at is null
`

type ListUserItemsRow struct {
	ItemID              int64          `db:"item_id" json:"item_id"`
	ItemName            string         `db:"item_name" json:"item_name"`
	ItemDescription     sql.NullString `db:"item_description" json:"item_description"`
	UserName            sql.NullString `db:"user_name" json:"user_name"`
	UserEmail           sql.NullString `db:"user_email" json:"user_email"`
	ItemTypeID          sql.NullInt64  `db:"item_type_id" json:"item_type_id"`
	ItemTypeName        sql.NullString `db:"item_type_name" json:"item_type_name"`
	ManufacturerID      sql.NullInt64  `db:"manufacturer_id" json:"manufacturer_id"`
	ManufacturerName    sql.NullString `db:"manufacturer_name" json:"manufacturer_name"`
	ItemInfoID          sql.NullInt64  `db:"item_info_id" json:"item_info_id"`
	PurchaseDate        sql.NullTime   `db:"purchase_date" json:"purchase_date"`
	PurchaseLocation    sql.NullString `db:"purchase_location" json:"purchase_location"`
	Price               sql.NullInt64  `db:"price" json:"price"`
	ExpirationDate      sql.NullTime   `db:"expiration_date" json:"expiration_date"`
	LastUsed            sql.NullTime   `db:"last_used" json:"last_used"`
	LocationID          sql.NullInt64  `db:"location_id" json:"location_id"`
	LocationName        sql.NullString `db:"location_name" json:"location_name"`
	LocationDescription sql.NullString `db:"location_description" json:"location_description"`
	LocationImg         sql.NullString `db:"location_img" json:"location_img"`
}

func (q *Queries) ListUserItems(ctx context.Context, id int64) ([]ListUserItemsRow, error) {
	rows, err := q.db.QueryContext(ctx, listUserItems, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUserItemsRow
	for rows.Next() {
		var i ListUserItemsRow
		if err := rows.Scan(
			&i.ItemID,
			&i.ItemName,
			&i.ItemDescription,
			&i.UserName,
			&i.UserEmail,
			&i.ItemTypeID,
			&i.ItemTypeName,
			&i.ManufacturerID,
			&i.ManufacturerName,
			&i.ItemInfoID,
			&i.PurchaseDate,
			&i.PurchaseLocation,
			&i.Price,
			&i.ExpirationDate,
			&i.LastUsed,
			&i.LocationID,
			&i.LocationName,
			&i.LocationDescription,
			&i.LocationImg,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
