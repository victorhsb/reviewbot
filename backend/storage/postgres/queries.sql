-- name: ListMessagesByParticipant :many
SELECT * FROM messages WHERE sender_id = $1 OR receiver_id = $1 ORDER BY created_at;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: SaveMessage :exec
INSERT INTO messages (receiver_id, sender_id, message) VALUES ($1, $2, $3);

-- name: SaveUser :one
INSERT INTO users (username, role) VALUES ($1, $2) RETURNING *;