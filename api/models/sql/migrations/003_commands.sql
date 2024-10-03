-- +goose Up
CREATE TABLE commands (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  command TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE tags (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE (user_id, name), -- Ensures a user cannot have duplicate tag names
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE command_tags (
  command_id INT NOT NULL,
  tag_id INT NOT NULL,
  PRIMARY KEY (command_id, tag_id),
  FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- +goose Down 
DROP TABLE tags;
DROP TABLE command_tags;
