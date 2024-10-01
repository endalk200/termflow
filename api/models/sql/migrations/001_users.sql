-- +goose Up
CREATE TABLE users (
  id   SERIAL,

  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,

  password VARCHAR(255) NOT NULL,

  email VARCHAR(320) UNIQUE NOT NULL,
	is_email_verified BOOLEAN DEFAULT FALSE,
	is_active BOOLEAN DEFAULT TRUE,
	github_handle VARCHAR(255),

  PRIMARY KEY(id)
);

-- +goose Down 
DROP TABLE users; 
