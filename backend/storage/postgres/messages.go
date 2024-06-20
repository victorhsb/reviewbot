package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/victorhsb/review-bot/backend/service"
	"github.com/victorhsb/review-bot/backend/storage/postgres/sqlc"
)

// SaveMessage stores a message in the database
func (c *client) SaveMessage(ctx context.Context, msg service.Message) error {
	return sqlc.New(c.conn).SaveMessage(ctx, sqlc.SaveMessageParams{
		// ID is expected to be automatically set by the database
		UserID:  &msg.UserID,
		Message: msg.Message,
		// CreatedAt is expected to be automatically set by the database
	})
}

var directionMap = map[sqlc.Direction]service.Direction{
	sqlc.DirectionSent:     service.DirectionSent,
	sqlc.DirectionReceived: service.DirectionReceived,
}

// ListMessagesByUserID retrieves all messages sent or received by a participant from the database
func (c *client) ListMessagesByUserID(ctx context.Context, id uuid.UUID) ([]service.Message, error) {
	queries := sqlc.New(c.conn)

	messages, err := queries.ListMessagesByUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not list messages by participant; %w", err)
	}

	result := make([]service.Message, 0)
	for _, m := range messages {
		result = append(result, service.Message{
			Message:   m.Message,
			Direction: directionMap[m.Direction],
			Timestamp: m.CreatedAt.Time,
		})
	}

	return result, nil
}
