-- name: ListLocationsOfUser :many
SELECT * FROM location WHERE user_id = $1;

-- name: GetLocation :one
SELECT * FROM location WHERE id = $1 LIMIT 1;


