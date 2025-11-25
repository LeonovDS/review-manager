package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionManager abstracts database transactions for use in use cases.
type TransactionManager interface {
	WithTransaction(ctx context.Context, transaction func(context.Context) error) error
}

// DBTransactionManager implements TransactionManager for pgx connection pool.
type DBTransactionManager struct {
	Pool *pgxpool.Pool
}

// WithTransaction wraps actions in postgres transaction.
func (db *DBTransactionManager) WithTransaction(
	ctx context.Context, transaction func(context.Context) error,
) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = transaction(ctx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
