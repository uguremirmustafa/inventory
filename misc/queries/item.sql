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

-- name: InsertUserItem :one
INSERT INTO item (
    name,
    description,
    user_id,
    item_type_id,
    manufacturer_id
) VALUES (
    $1, -- name
	$2, -- description
	$3, -- user_id
	$4, -- item_type_id
	$5  -- manufacturer_id
) RETURNING id;

-- name: InsertItemInfo :one
INSERT INTO item_info (
    item_id,
    purchase_date,
    purchase_location,
    price,
    expiration_date,
    last_used,
    location_id
) VALUES (
    $1, -- item_id
	$2, -- purchase_date
	$3, -- purchase_location
	$4, -- price
	$5, -- expiration_date
	$6, -- last_used
	$7  -- location_id
) RETURNING id;