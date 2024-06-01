-- name: ListLocationsOfGroup :many
SELECT * FROM location WHERE group_id = $1 AND deleted_at is null;

-- name: GetLocation :one
SELECT * FROM location WHERE id = $1 AND deleted_at is null LIMIT 1;


-- name: InsertLocation :one
INSERT INTO location (name, image_url, description, group_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: UpdateLocation :one
UPDATE location
SET name = $2,
    image_url = $3,
    description = $4,
    group_id = $5
WHERE id = $1
RETURNING id;

-- name: DeleteLocation :one
UPDATE location
SET deleted_at = $2
WHERE id = $1
RETURNING id;