-- +goose Up
CREATE TABLE chirpy(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NUll,
  updated_at TIMESTAMP NOT NUll,
  body TEXT NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirpy;
