-- name: GetCollection :one
SELECT * FROM Collection
WHERE id = ? LIMIT 1;

-- name: ListCollections :many
SELECT * FROM Collection
ORDER BY name;

-- name: CreateCollection :one
INSERT INTO Collection (
  name, description
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateCollection :exec
UPDATE Collection
set name = ?,
description = ?
WHERE id = ?;

-- name: DeleteCollection :exec
DELETE FROM Collection
WHERE id = ?;
