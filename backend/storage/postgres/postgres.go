package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/victorhsb/review-bot/backend/service"
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

var _ service.Storage = (*client)(nil)

func (c *client) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}
