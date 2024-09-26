-- +goose Up
CREATE TABLE Tag (
  id   INTEGER PRIMARY KEY,
  name text    NOT NULL UNIQUE,
  description  text
);

CREATE TABLE Command (
  id INTEGER PRIMARY KEY,
  command text,
  description text
);

CREATE TABLE CommandTag (
  commandId INTEGER,
  tagId INTEGER,

  PRIMARY KEY (commandId, tagId),
  FOREIGN KEY (commandId) REFERENCES Command(id) ON DELETE CASCADE,
  FOREIGN KEY (tagId) REFERENCES Tag(id) ON DELETE CASCADE
);

-- +goose Down 
DROP TABLE Tag;
DROP TABLE Command;
DROP TABLE CommandTag;
