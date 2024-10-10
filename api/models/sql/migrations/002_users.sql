-- +goose Up
CREATE TABLE users (
  id   SERIAL,

  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,

  password VARCHAR(255) NOT NULL,
  refresh_token TEXT,

  email VARCHAR(320) UNIQUE NOT NULL,
	is_email_verified BOOLEAN DEFAULT FALSE,

	is_active BOOLEAN DEFAULT TRUE,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  PRIMARY KEY(id)
);

-- Apply trigger for `users` table
-- CREATE TRIGGER update_users_updated_at
-- BEFORE UPDATE ON users
-- FOR EACH ROW
-- EXECUTE FUNCTION update_updated_at_column();

-- +goose Down 
-- DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TABLE users;
