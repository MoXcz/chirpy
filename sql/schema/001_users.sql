-- +goose Up
CREATE TABLE users(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NUll,
  updated_at TIMESTAMP NOT NUll,
  email TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
