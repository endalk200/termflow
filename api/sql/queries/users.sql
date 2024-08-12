-- name: GetUser :one
SELECT * FROM users 
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users 
ORDER BY first_name;

-- name: CreateUser :one
INSERT INTO users (
  first_name, last_name, email, is_email_verified, is_active, github_handle, password
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
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
