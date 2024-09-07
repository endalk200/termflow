-- name: GetCommand :one
SELECT * FROM Command
WHERE id = ? LIMIT 1;

-- name: ListCommands :many
SELECT * FROM Command
ORDER BY name;

-- name: AddCommand :exec
INSERT INTO Command (
  name, description, command
) VALUES (
  ?, ?, ?
);

-- name: UpdateCommand :exec
UPDATE Command
set name = ?,
description = ?,
command = ?
WHERE id = ?;

-- name: DeleteCommand :exec
DELETE FROM Command
WHERE id = ?;
