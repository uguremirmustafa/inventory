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


-- name: GetGroupsOfUser :many
SELECT 
    g.id as group_id,
    g.name as group_name,
    g.description as group_desc,
    g.group_owner_id,
    g.created_at,
    g.updated_at
FROM 
    groups g
JOIN 
    user_groups ug ON g.id = ug.group_id
WHERE 
    ug.user_id = $1 and g.deleted_at is null;

