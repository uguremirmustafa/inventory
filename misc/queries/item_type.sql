-- name: ListItemTypes :many
SELECT * FROM item_type;

-- name: CreateItemType :exec
insert into item_type ("name") values ($1);



