// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO
    users (name, email, avatar)
VALUES
    ($1, $2, $3) RETURNING id, name, email, avatar, active_group_id, created_at, updated_at, deleted_at
`

type CreateUserParams struct {
	Name   string         `db:"name" json:"name"`
	Email  string         `db:"email" json:"email"`
	Avatar sql.NullString `db:"avatar" json:"avatar"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Name, arg.Email, arg.Avatar)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Avatar,
		&i.ActiveGroupID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM
    users
WHERE
    id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT
    id, name, email, avatar, active_group_id, created_at, updated_at, deleted_at
FROM
    users
WHERE
    id = $1
LIMIT
    1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Avatar,
		&i.ActiveGroupID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
    id, name, email, avatar, active_group_id, created_at, updated_at, deleted_at
FROM
    users
WHERE
    email = $1
LIMIT
    1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Avatar,
		&i.ActiveGroupID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateUserActiveGroupID = `-- name: UpdateUserActiveGroupID :one
UPDATE users
SET active_group_id = $2
WHERE id = $1
RETURNING id, name, email, avatar, active_group_id, created_at, updated_at, deleted_at
`

type UpdateUserActiveGroupIDParams struct {
	ID            int64         `db:"id" json:"id"`
	ActiveGroupID sql.NullInt64 `db:"active_group_id" json:"active_group_id"`
}

func (q *Queries) UpdateUserActiveGroupID(ctx context.Context, arg UpdateUserActiveGroupIDParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserActiveGroupID, arg.ID, arg.ActiveGroupID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Avatar,
		&i.ActiveGroupID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const upsertUser = `-- name: UpsertUser :one
INSERT INTO users (name, email, avatar, active_group_id) 
VALUES ($1, $2, $3, $4)
ON CONFLICT (email) 
DO UPDATE SET 
    name = EXCLUDED.name,
    avatar = EXCLUDED.avatar,
    active_group_id = EXCLUDED.active_group_id
RETURNING id, name, email, avatar, active_group_id, created_at, updated_at, deleted_at
`

type UpsertUserParams struct {
	Name          string         `db:"name" json:"name"`
	Email         string         `db:"email" json:"email"`
	Avatar        sql.NullString `db:"avatar" json:"avatar"`
	ActiveGroupID sql.NullInt64  `db:"active_group_id" json:"active_group_id"`
}

func (q *Queries) UpsertUser(ctx context.Context, arg UpsertUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, upsertUser,
		arg.Name,
		arg.Email,
		arg.Avatar,
		arg.ActiveGroupID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Avatar,
		&i.ActiveGroupID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
