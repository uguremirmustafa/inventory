-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    id = $1
LIMIT
    1;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1
LIMIT
    1;

-- name: CreateUser :one
INSERT INTO
    users (name, email, avatar)
VALUES
    ($1, $2, $3) RETURNING *;

-- name: DeleteUser :exec
DELETE FROM
    users
WHERE
    id = $1;


-- name: UpsertUser :one
INSERT INTO users (name, email, avatar, active_group_id) 
VALUES ($1, $2, $3, $4)
ON CONFLICT (email) 
DO UPDATE SET 
    name = EXCLUDED.name,
    avatar = EXCLUDED.avatar,
    active_group_id = EXCLUDED.active_group_id
RETURNING *;

-- name: UpdateUserActiveGroupID :one
UPDATE users
SET active_group_id = $2
WHERE id = $1
RETURNING *;


