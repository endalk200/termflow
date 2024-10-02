-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3);

-- name: GetRefreshTokenByUserID :one
SELECT *
FROM refresh_tokens
WHERE user_id = $1 AND revoked = FALSE
LIMIT 1;

-- name: GetAllRefreshTokensByUserID :many
SELECT *
FROM refresh_tokens
WHERE user_id = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked = TRUE, revoked_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE id = $1;

-- name: CleanupExpiredTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW();
