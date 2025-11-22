// Package database handles db connection and migrations.
package database

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx" // database implementation for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"  // source implementation for migrations
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect applies migrations and creates connection pool to database.
func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	connString := os.Getenv("DB_URL")

	err := migrateUp(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, connString)
	return pool, err
}

func migrateUp(connString string) error {
	m, err := migrate.New("file://migrations/", connString)
	if err != nil {
		return err
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("Database is up-to-date, no migrations applied")
		return nil
	} else if err != nil {
		return err
	}

	slog.Info("Migrations applied successfully")
	return nil
}
