-- name: InsertItemImage :exec
INSERT INTO item_image
("item_id","image_url")
values
($1, $2);
 