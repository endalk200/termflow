-- +goose Up
CREATE TABLE history (
  id   SERIAL,

  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,

  command TEXT NOT NULL,

  PRIMARY KEY(id)
);

-- +goose Down 
DROP TABLE history; 
