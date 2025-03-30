package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/interfaces"
)

// WithTransaction executes the given function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func WithTransaction(
	ctx context.Context,
	db *sqlx.DB,
	logger interfaces.Logger,
	fn func(*sqlx.Tx) error,
) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if panicValue := recover(); panicValue != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logger.Error("failed to rollback transaction after panic", map[string]any{
					"rollback_error": rollbackErr,
					"panic_value": panicValue,
				})
				panic(panicValue)
			} else {
				logger.Error("transaction rolled back after panic", map[string]any{
					"panic_value": panicValue,
				})
			}
			panic(panicValue) // re-throw panic after rollback so that it may be handled further up callstack
		}
	}()

	if err := fn(tx); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logger.Error("failed to rollback transaction after error", map[string]any{
				"rollback_error": rollbackErr,
				"original_error": err,
			})
		}
		return err
	}

	return tx.Commit()
}
