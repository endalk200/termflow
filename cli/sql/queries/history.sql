-- name: GetHistory :one
SELECT * FROM history 
WHERE id = ? LIMIT 1;

-- name: ListHistory :many
SELECT * FROM history 
ORDER BY name;
