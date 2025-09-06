-- +goose Up
CREATE TABLE messages (
  id UUID NOT NULL PRIMARY KEY,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE messages;
