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

func (c *client) SaveMessage(_ context.Context, _ service.Message) error {
	panic("not implemented") // TODO: Implement
}

func (c *client) ListMessagesByParticipant(ctx context.Context, id uuid.UUID) ([]service.Message, error) {
	queries := sqlc.New(c.conn)

	queries.ListMessagesByParticipant(ctx, id)
	panic("not implemented") // TODO: Implement
}

func New(ctx context.Context, connUrl string) (service.MessageRepository, error) {
	conn, err := pgxpool.New(ctx, connUrl)
	if err != nil {
		return nil, fmt.Errorf("could not establish connection to postgres; %w", err)
	}

	repo := &client{
		conn: conn,
	}

	return repo, nil
}
