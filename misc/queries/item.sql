-- name: ListUserItems :many
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
WHERE u.id = $1 AND i.deleted_at is null;
