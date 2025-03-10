-- +goose Up
CREATE TABLE refresh_tokens (
  id              UUID PRIMARY KEY,
  user_id         UUID REFERENCES users(id) ON DELETE CASCADE,

  token_hash      TEXT NOT NULL,

  issued_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL
);

-- +goose Down 
DROP TABLE refresh_tokens;
