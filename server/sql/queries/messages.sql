-- name: CreateMessage :one
INSERT INTO messages (user_id, content)
VALUES ($1, $2)
RETURNING *;

-- name: ListMessages :many
SELECT * FROM messages
ORDER BY created_at DESC;