-- name: CreateMessage :one
INSERT INTO messages (content)
VALUES ($1)
RETURNING *;

-- name: ListMessages :many
SELECT * FROM messages
ORDER BY created_at DESC;