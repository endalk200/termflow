-- name: GetHistory :one
SELECT * FROM history 
WHERE id = $1 LIMIT 1;

-- name: ListHistory :many
SELECT * FROM history 
ORDER BY name;
