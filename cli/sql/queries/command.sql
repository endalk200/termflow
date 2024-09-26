-- name: GetCommand :one
SELECT * FROM Command
WHERE id = ? LIMIT 1;

-- name: ListCommands :many
SELECT * FROM Command
ORDER BY id;

-- name: ListCommandsForTagByName :many
SELECT c.id, c.command, c.description
FROM Command c
JOIN CommandTag ct ON c.id = ct.commandId
JOIN Tag t ON ct.tagId = t.id
WHERE t.name = ?;

-- name: ListCommandsWithTagsByTagName :many
SELECT 
  c.id AS command_id, 
  c.command, 
  c.description AS command_description,
  t.id AS tag_id, 
  t.name AS tag_name, 
  t.description AS tag_description
FROM Command c
LEFT JOIN CommandTag ct ON c.id = ct.commandId
LEFT JOIN Tag t ON ct.tagId = t.id
WHERE t.name = ?
ORDER BY c.id, t.name;

-- name: ListCommandsWithTags :many
SELECT 
    c.id AS command_id, 
    c.command, 
    c.description AS command_description,
    t.id AS tag_id, 
    t.name AS tag_name, 
    t.description AS tag_description
FROM Command c
LEFT JOIN CommandTag ct ON c.id = ct.commandId
LEFT JOIN Tag t ON ct.tagId = t.id
ORDER BY c.id, t.name;

-- name: AddCommand :one
INSERT INTO Command (
  command, description
) VALUES (
  ?, ?
) RETURNING *;

-- name: UpdateCommand :exec
UPDATE Command
set command = ?,
description = ?
WHERE id = ?;

-- name: DeleteCommand :exec
DELETE FROM Command
WHERE id = ?;

-- name: GetTag :one
SELECT * FROM Tag
WHERE id = ? LIMIT 1;

-- name: GetTagByName :one
SELECT * FROM Tag
WHERE name = ? LIMIT 1;

-- name: ListTags :many
SELECT * FROM Tag
ORDER BY name;

-- name: AddTag :one
INSERT INTO Tag (
  name, description
) VALUES (
  ?, ?
) RETURNING *;

-- name: UpdateTag :exec
UPDATE Tag
SET name = ?,
description = ?
WHERE id = ?;

-- name: DeleteTag :exec
DELETE FROM Tag
WHERE id = ?;

-- name: AddCommandTag :exec
INSERT INTO CommandTag (
  commandId, tagId
) VALUES (
  ?, ?
);

-- name: RemoveCommandTag :exec
DELETE FROM CommandTag
WHERE commandId = ? AND tagId = ?;

-- name: GetTagsForCommand :many
SELECT t.id, t.name, t.description
FROM Tag t
JOIN CommandTag ct ON t.id = ct.tagId
WHERE ct.commandId = ?;



