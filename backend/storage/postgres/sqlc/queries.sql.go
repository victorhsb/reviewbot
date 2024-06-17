// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const getUser = `-- name: GetUser :one
SELECT id, username, role, created_at FROM users WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

const listMessagesByParticipant = `-- name: ListMessagesByParticipant :many
SELECT id, receiver_id, sender_id, message, created_at FROM messages WHERE sender_id = $1 OR receiver_id = $1 ORDER BY created_at
`

func (q *Queries) ListMessagesByParticipant(ctx context.Context, senderID pgtype.UUID) ([]Message, error) {
	rows, err := q.db.Query(ctx, listMessagesByParticipant, senderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.ReceiverID,
			&i.SenderID,
			&i.Message,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveMessage = `-- name: SaveMessage :exec
INSERT INTO messages (receiver_id, sender_id, message) VALUES ($1, $2, $3)
`

type SaveMessageParams struct {
	ReceiverID pgtype.UUID
	SenderID   pgtype.UUID
	Message    string
}

func (q *Queries) SaveMessage(ctx context.Context, arg SaveMessageParams) error {
	_, err := q.db.Exec(ctx, saveMessage, arg.ReceiverID, arg.SenderID, arg.Message)
	return err
}

const saveUser = `-- name: SaveUser :one
INSERT INTO users (username, role) VALUES ($1, $2) RETURNING id, username, role, created_at
`

type SaveUserParams struct {
	Username string
	Role     string
}

func (q *Queries) SaveUser(ctx context.Context, arg SaveUserParams) (User, error) {
	row := q.db.QueryRow(ctx, saveUser, arg.Username, arg.Role)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}