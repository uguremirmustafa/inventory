-- name: ListItemTypes :many
SELECT * FROM item_type where deleted_at is null;

-- name: CreateItemType :exec
insert into item_type ("name", "description") values ($1, $2);



