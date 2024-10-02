// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  first_name, last_name, email, is_email_verified, is_active, github_handle, password
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, first_name, last_name, password, refresh_token, email, is_email_verified, is_active, github_handle
`

type CreateUserParams struct {
	FirstName       string
	LastName        string
	Email           string
	IsEmailVerified pgtype.Bool
	IsActive        pgtype.Bool
	GithubHandle    pgtype.Text
	Password        string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.IsEmailVerified,
		arg.IsActive,
		arg.GithubHandle,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Password,
		&i.RefreshToken,
		&i.Email,
		&i.IsEmailVerified,
		&i.IsActive,
		&i.GithubHandle,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, first_name, last_name, password, refresh_token, email, is_email_verified, is_active, github_handle FROM users
WHERE (id = $1 OR github_handle = $2 OR email = $3)
LIMIT 1
`

type GetUserParams struct {
	ID           int32
	GithubHandle pgtype.Text
	Email        string
}

func (q *Queries) GetUser(ctx context.Context, arg GetUserParams) (User, error) {
	row := q.db.QueryRow(ctx, getUser, arg.ID, arg.GithubHandle, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Password,
		&i.RefreshToken,
		&i.Email,
		&i.IsEmailVerified,
		&i.IsActive,
		&i.GithubHandle,
	)
	return i, err
}

const getUserWithRefreshTokens = `-- name: GetUserWithRefreshTokens :one
SELECT u.id, u.first_name, u.last_name, u.password, u.refresh_token, u.email, u.is_email_verified, u.is_active, u.github_handle, rt.id, rt.user_id, rt.token_hash, rt.issued_at, rt.expires_at, rt.revoked, rt.revoked_at
FROM users u
LEFT JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE u.id = $1
`

type GetUserWithRefreshTokensRow struct {
	ID              int32
	FirstName       string
	LastName        string
	Password        string
	RefreshToken    pgtype.Text
	Email           string
	IsEmailVerified pgtype.Bool
	IsActive        pgtype.Bool
	GithubHandle    pgtype.Text
	ID_2            pgtype.Int4
	UserID          pgtype.Int8
	TokenHash       pgtype.Text
	IssuedAt        pgtype.Timestamp
	ExpiresAt       pgtype.Timestamp
	Revoked         pgtype.Bool
	RevokedAt       pgtype.Timestamp
}

func (q *Queries) GetUserWithRefreshTokens(ctx context.Context, id int32) (GetUserWithRefreshTokensRow, error) {
	row := q.db.QueryRow(ctx, getUserWithRefreshTokens, id)
	var i GetUserWithRefreshTokensRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Password,
		&i.RefreshToken,
		&i.Email,
		&i.IsEmailVerified,
		&i.IsActive,
		&i.GithubHandle,
		&i.ID_2,
		&i.UserID,
		&i.TokenHash,
		&i.IssuedAt,
		&i.ExpiresAt,
		&i.Revoked,
		&i.RevokedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, first_name, last_name, password, refresh_token, email, is_email_verified, is_active, github_handle FROM users
WHERE ($1 IS NULL OR is_active = $1)
AND ($2 IS NULL OR is_email_verified = $2)
ORDER BY first_name
`

type ListUsersParams struct {
	Column1 interface{}
	Column2 interface{}
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.Column1, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Password,
			&i.RefreshToken,
			&i.Email,
			&i.IsEmailVerified,
			&i.IsActive,
			&i.GithubHandle,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users 
  set first_name = $2,
  last_name = $3
WHERE id = $1
`

type UpdateUserParams struct {
	ID        int32
	FirstName string
	LastName  string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser, arg.ID, arg.FirstName, arg.LastName)
	return err
}