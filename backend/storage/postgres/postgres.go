package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/victorhsb/review-bot/backend/service"
	"github.com/victorhsb/review-bot/backend/storage/postgres/sqlc"
)

type client struct {
	conn *pgxpool.Pool
}

// New creates a new instance of the postgres client
func New(ctx context.Context, connUrl string) (service.Storage, error) {
	conn, err := pgxpool.New(ctx, connUrl)
	if err != nil {
		return nil, fmt.Errorf("could not establish connection to postgres; %w", err)
	}

	repo := &client{
		conn: conn,
	}

	return repo, nil
}

// SaveMessage stores a message in the database
func (c *client) SaveMessage(ctx context.Context, msg service.Message) error {
	return sqlc.New(c.conn).SaveMessage(ctx, sqlc.SaveMessageParams{
		// ID is expected to be automatically set by the database
		ReceiverID: msg.Target,
		SenderID:   msg.Sender,
		Message:    msg.Message,
		// CreatedAt is expected to be automatically set by the database
	})
}

// ListMessagesByParticipant retrieves all messages sent or received by a participant from the database
func (c *client) ListMessagesByParticipant(ctx context.Context, id uuid.UUID) ([]service.Message, error) {
	queries := sqlc.New(c.conn)

	messages, err := queries.ListMessagesByParticipant(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not list messages by participant; %w", err)
	}

	result := make([]service.Message, 0)
	for _, m := range messages {
		result = append(result, service.Message{
			Sender:    m.SenderID,
			Target:    m.ReceiverID,
			Message:   m.Message,
			Timestamp: m.CreatedAt.Time,
		})
	}

	return result, nil
}
