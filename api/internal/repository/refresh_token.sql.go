// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: refresh_token.sql

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const cleanupExpiredTokens = `-- name: CleanupExpiredTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW()
`

func (q *Queries) CleanupExpiredTokens(ctx context.Context) error {
	_, err := q.db.Exec(ctx, cleanupExpiredTokens)
	return err
}

const createRefreshToken = `-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
VALUES (uuid_generate_v4(), $1, $2, $3)
`

type CreateRefreshTokenParams struct {
	UserID    pgtype.UUID        `json:"user_id"`
	TokenHash string             `json:"token_hash"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) error {
	_, err := q.db.Exec(ctx, createRefreshToken, arg.UserID, arg.TokenHash, arg.ExpiresAt)
	return err
}

const deleteRefreshToken = `-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE id = $1
`

func (q *Queries) DeleteRefreshToken(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteRefreshToken, id)
	return err
}

const getAllRefreshTokensByUserID = `-- name: GetAllRefreshTokensByUserID :many
SELECT id, user_id, token_hash, issued_at, expires_at
FROM refresh_tokens
WHERE user_id = $1
`

func (q *Queries) GetAllRefreshTokensByUserID(ctx context.Context, userID pgtype.UUID) ([]RefreshToken, error) {
	rows, err := q.db.Query(ctx, getAllRefreshTokensByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RefreshToken
	for rows.Next() {
		var i RefreshToken
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.TokenHash,
			&i.IssuedAt,
			&i.ExpiresAt,
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

const getRefreshTokenByUserID = `-- name: GetRefreshTokenByUserID :one
SELECT id, user_id, token_hash, issued_at, expires_at
FROM refresh_tokens
WHERE user_id = $1 AND revoked = FALSE
LIMIT 1
`

func (q *Queries) GetRefreshTokenByUserID(ctx context.Context, userID pgtype.UUID) (RefreshToken, error) {
	row := q.db.QueryRow(ctx, getRefreshTokenByUserID, userID)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TokenHash,
		&i.IssuedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked = TRUE, revoked_at = CURRENT_TIMESTAMP
WHERE id = $1
`

func (q *Queries) RevokeRefreshToken(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, revokeRefreshToken, id)
	return err
}