-- name: CreateInvitation :one
INSERT INTO invitations
(email, token, group_id, invitor_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetInvitationByToken :one
SELECT * FROM invitations
WHERE token = $1
LIMIT 1;


