-- name: GetUser :one
SELECT * FROM users
WHERE (id = $1 OR email = $2)
LIMIT 1;

-- name: GetUserWithRefreshTokens :one
SELECT u.*, rt.*
FROM users u
LEFT JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE u.id = $1;

-- name: ListUsers :many
SELECT * FROM users
WHERE (sqlc.arg(is_active) IS NULL OR is_active = sqlc.arg(is_active))
AND (sqlc.arg(is_email_verified) IS NULL OR is_email_verified = sqlc.arg(is_email_verified))
ORDER BY first_name;

-- name: InsertUser :one
INSERT INTO users (
  id, first_name, last_name, email, is_email_verified, is_active, password
) VALUES (
  uuid_generate_v4(), $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users 
  set first_name = $2,
  last_name = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;
