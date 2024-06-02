// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: groups.sql

package db

import (
	"context"
	"database/sql"
)

const connectUserAndGroup = `-- name: ConnectUserAndGroup :exec
INSERT INTO
    user_groups (user_id, group_id)
VALUES
    ($1, $2)
`

type ConnectUserAndGroupParams struct {
	UserID  int64 `db:"user_id" json:"user_id"`
	GroupID int64 `db:"group_id" json:"group_id"`
}

func (q *Queries) ConnectUserAndGroup(ctx context.Context, arg ConnectUserAndGroupParams) error {
	_, err := q.db.ExecContext(ctx, connectUserAndGroup, arg.UserID, arg.GroupID)
	return err
}

const createGroup = `-- name: CreateGroup :one
INSERT INTO
    groups (name, description, group_owner_id)
VALUES 
    ($1, $2, $3)
RETURNING id, name, description, group_owner_id, created_at, updated_at, deleted_at
`

type CreateGroupParams struct {
	Name         string         `db:"name" json:"name"`
	Description  sql.NullString `db:"description" json:"description"`
	GroupOwnerID int64          `db:"group_owner_id" json:"group_owner_id"`
}

func (q *Queries) CreateGroup(ctx context.Context, arg CreateGroupParams) (Group, error) {
	row := q.db.QueryRowContext(ctx, createGroup, arg.Name, arg.Description, arg.GroupOwnerID)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.GroupOwnerID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getGroupsOfUser = `-- name: GetGroupsOfUser :many
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
    ug.user_id = $1 and g.deleted_at is null
`

type GetGroupsOfUserRow struct {
	GroupID      int64          `db:"group_id" json:"group_id"`
	GroupName    string         `db:"group_name" json:"group_name"`
	GroupDesc    sql.NullString `db:"group_desc" json:"group_desc"`
	GroupOwnerID int64          `db:"group_owner_id" json:"group_owner_id"`
	CreatedAt    sql.NullTime   `db:"created_at" json:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at" json:"updated_at"`
}

func (q *Queries) GetGroupsOfUser(ctx context.Context, userID int64) ([]GetGroupsOfUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getGroupsOfUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGroupsOfUserRow
	for rows.Next() {
		var i GetGroupsOfUserRow
		if err := rows.Scan(
			&i.GroupID,
			&i.GroupName,
			&i.GroupDesc,
			&i.GroupOwnerID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMembersOfGroup = `-- name: GetMembersOfGroup :many
SELECT
	u.id,
	u."name",
	u.email,
	u.avatar,
	g.id as group_id,
	g.name as group_name
FROM
	user_groups ug
	join users u on ug.user_id = u.id
	join groups g on ug.group_id = g.id
WHERE ug.group_id = $1
`

type GetMembersOfGroupRow struct {
	ID        int64          `db:"id" json:"id"`
	Name      string         `db:"name" json:"name"`
	Email     string         `db:"email" json:"email"`
	Avatar    sql.NullString `db:"avatar" json:"avatar"`
	GroupID   int64          `db:"group_id" json:"group_id"`
	GroupName string         `db:"group_name" json:"group_name"`
}

func (q *Queries) GetMembersOfGroup(ctx context.Context, groupID int64) ([]GetMembersOfGroupRow, error) {
	rows, err := q.db.QueryContext(ctx, getMembersOfGroup, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMembersOfGroupRow
	for rows.Next() {
		var i GetMembersOfGroupRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Avatar,
			&i.GroupID,
			&i.GroupName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
