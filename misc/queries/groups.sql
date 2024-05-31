-- name: CreateGroup :one
INSERT INTO
    groups (name, description, group_owner_id)
VALUES 
    ($1, $2, $3)
RETURNING *;

-- name: ConnectUserAndGroup :exec
INSERT INTO
    user_groups (user_id, group_id)
VALUES
    ($1, $2);


