-- name: ListMainItemTypes :many
SELECT * FROM item_type where deleted_at is null AND parent_id is null;

-- name: ListItemTypes :many
SELECT * FROM item_type where deleted_at is null AND parent_id = $1;

-- name: ListAllItemTypes :many
SELECT * FROM item_type where deleted_at is null;

-- name: CreateItemType :exec
insert into item_type ("name") values ($1);



