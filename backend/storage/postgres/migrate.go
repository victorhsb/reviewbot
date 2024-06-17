package postgres

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratedb "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	pgxstdlib "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*
var migrations embed.FS

// Migrate will perform the database migrations embedded in the source code. This means that the server will be able to run without the need of any external interference in terms of database migrations.
// It also means that in case of failure it is left for the developer to decide how to handle it. (presumably the server will not start)
func (c *client) Migrate() error {
	directory, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("could not read migrations embedded directory; %w", err)
	}

	db := pgxstdlib.OpenDBFromPool(c.conn)

	driver, err := migratedb.WithInstance(db, &migratedb.Config{})
	if err != nil {
		return fmt.Errorf("could not get driver instance; %w", err)
	}

	migration, err := migrate.NewWithInstance("iofs", directory, "program-api", driver)
	if err != nil {
		return fmt.Errorf("could not instantiate migrate instance; %w", err)
	}

	// This will perform all the migrations necessary in order to get to the latest possible version.
	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not perform migration; %w", err)
	}

	return nil
}
