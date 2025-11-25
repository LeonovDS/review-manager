// Package database handles db connection and migrations.
package database

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx" // database implementation for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"  // source implementation for migrations
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect applies migrations and creates connection pool to database.
func Connect(ctx context.Context, connStr, migrationSrc string) (*pgxpool.Pool, error) {
	err := migrateUp(connStr, migrationSrc)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, connStr)
	return pool, err
}

func migrateUp(connString, migrationsSrc string) error {
	connString = strings.Replace(connString, "postgres", "pgx", 1)
	m, err := migrate.New(migrationsSrc, connString)
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
