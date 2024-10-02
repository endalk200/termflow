-- +goose Up
CREATE TABLE refresh_tokens (
  id              SERIAL PRIMARY KEY,
  user_id         BIGINT REFERENCES users(id) ON DELETE CASCADE,

  token_hash      TEXT NOT NULL,

  issued_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at      TIMESTAMP NOT NULL,

  revoked         BOOLEAN DEFAULT FALSE,
  revoked_at      TIMESTAMP
);

-- +goose Down 
DROP TABLE refresh_tokens;
