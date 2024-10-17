-- name: FindTags :many
SELECT * FROM tags
WHERE (user_id = $1);

-- name: FindTagById :one
SELECT * FROM tags
WHERE (id = $1);

-- name: FindTagByName :one
SELECT * FROM tags
WHERE (name = $1);

-- name: FindTagsWithCommands :many
SELECT 
    t.id AS tag_id, t.name AS tag_name, t.description AS tag_description,
    c.id AS command_id, c.command AS command_name, c.description AS command_description
FROM tags t
LEFT JOIN command_tags ct ON t.id = ct.tag_id
LEFT JOIN commands c ON ct.command_id = c.id
WHERE t.user_id = $1;
-- SELECT 
--     c.id AS command_id, 
--     c.command AS command_name, 
--     c.description AS command_description, 
--     c.created_at AS command_created_at, 
--     c.updated_at AS command_updated_at,
--     array_agg(
--         json_build_object(
--             'id', t.id,
--             'name', t.name,
--             'description', t.description
--         )
--     ) AS tags
-- FROM commands c
-- LEFT JOIN command_tags ct ON c.id = ct.command_id
-- LEFT JOIN tags t ON ct.tag_id = t.id
-- WHERE c.user_id = $1
-- GROUP BY c.id;

-- name: InsertTag :one
INSERT INTO tags (
  id, user_id, name, description
) VALUES (
  uuid_generate_v4(), $1, $2, sqlc.narg(description)::Text
)
RETURNING *;

-- name: UpdateTag :one
UPDATE tags
SET name = COALESCE(NULLIF($2, ''), name),
    description = COALESCE(NULLIF($3, ''), description)
WHERE id = $1
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE id = $1;

-- name: FindCommands :many
SELECT * FROM commands
WHERE (user_id = $1);

-- name: FindCommandsWithTags :many
SELECT 
    c.id AS command_id, c.command AS command_name, c.description AS command_description, c.created_at AS command_created_at, c.updated_at AS command_updated_at,
    t.id AS tag_id, t.name AS tag_name, t.description AS tag_description, t.created_at AS tag_created_at, t.updated_at AS tag_updated_at
FROM commands c
LEFT JOIN command_tags ct ON c.id = ct.command_id
LEFT JOIN tags t ON ct.tag_id = t.id
WHERE c.user_id = $1;

-- name: FindCommandsByTagId :many
SELECT 
    c.id AS command_id, c.command AS command_name, c.description AS command_description, c.created_at AS command_created_at, c.updated_at AS command_updated_at,
    t.id AS tag_id, t.name AS tag_name, t.description AS tag_description, t.created_at AS tag_created_at, t.updated_at AS tag_updated_at
FROM commands c
LEFT JOIN command_tags ct ON c.id = ct.command_id
LEFT JOIN tags t ON ct.tag_id = t.id
WHERE t.id = $1;

-- name: InsertCommands :one
INSERT INTO commands (
  id, user_id, command, description
) VALUES (
  uuid_generate_v4(), $1, $2, sqlc.narg(description)::Text
)
RETURNING *;

-- name: UpdateCommand :one
UPDATE commands
SET command = COALESCE(NULLIF($2, ''), command),
    description = COALESCE(NULLIF($3, ''), description)
WHERE id = $1
RETURNING *;

-- name: DeleteCommand :exec
DELETE FROM commands
WHERE id = $1;

-- name: DeleteCommandTagRelationByCommandId :exec
DELETE FROM command_tags
WHERE command_id = $1;

-- name: DeleteCommandTagRelationByTagId :exec
DELETE FROM command_tags
WHERE tag_id = $1;

-- name: AttachCommandToTag :one
INSERT INTO command_tags (command_id, tag_id)
VALUES ($1, $2)
RETURNING *;

-- name: ReplaceCommandTag :exec
WITH deleted AS (
    DELETE FROM command_tags
    WHERE command_id = $1
)
INSERT INTO command_tags (command_id, tag_id)
VALUES ($1, $2);
