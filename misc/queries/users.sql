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
INSERT INTO users (name, email, avatar) 
VALUES ($1, $2, $3)
ON CONFLICT (email) 
DO UPDATE SET 
    name = EXCLUDED.name,
    avatar = EXCLUDED.avatar
RETURNING *;