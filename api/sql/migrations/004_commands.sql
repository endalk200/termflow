-- +goose Up
CREATE TABLE commands (
  id          UUID PRIMARY KEY,
  user_id     UUID NOT NULL,
  command     TEXT NOT NULL,
  description TEXT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE tags (
  id          UUID PRIMARY KEY,
  user_id     UUID NOT NULL,
  name        VARCHAR(255) NOT NULL,
  description TEXT,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  UNIQUE (name), -- Ensures a user cannot have duplicate tag names
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE command_tags (
  command_id  UUID NOT NULL,
  tag_id      UUID NOT NULL,

  PRIMARY KEY (command_id, tag_id),
  FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- CREATE TRIGGER update_commands_updated_at
-- BEFORE UPDATE ON commands
-- FOR EACH ROW
-- EXECUTE FUNCTION update_updated_at_column();
--
-- CREATE TRIGGER update_tags_updated_at
-- BEFORE UPDATE ON tags
-- FOR EACH ROW
-- EXECUTE FUNCTION update_updated_at_column();

-- +goose Down 
-- DROP TRIGGER IF EXISTS update_commands_updated_at ON commands;
-- DROP TRIGGER IF EXISTS update_tags_updated_at ON tags;

DROP TABLE command_tags;
DROP TABLE tags;
DROP TABLE commands;
