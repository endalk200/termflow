-- +goose Up
CREATE TABLE Collection (
  id   INTEGER PRIMARY KEY,
  name text    NOT NULL,
  description  text
);

CREATE TABLE Command (
  id INTEGER PRIMARY KEY,
  name text NOT NULL,
  description text,
  command text
)

-- +goose Down 
DROP TABLE Collection;
DROP TABLE Command;
