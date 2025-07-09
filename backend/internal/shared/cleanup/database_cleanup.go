package cleanup

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
)

// DatabaseCleanupService handles cleanup of expired data
type DatabaseCleanupService struct {
	appCtx *appcontext.AppContext
	db     *sqlx.DB
}

// Close closes the database connection
func (dcs *DatabaseCleanupService) Close() error {
	if dcs.db != nil {
		return dcs.db.Close()
	}
	return nil
}

// NewDatabaseCleanupService creates a new database cleanup service
func NewDatabaseCleanupService(appCtx *appcontext.AppContext) (*DatabaseCleanupService, error) {
	// Create sqlx connection
	db, err := sqlx.Connect("pgx", appCtx.Config.Postgres.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	return &DatabaseCleanupService{
		appCtx: appCtx,
		db:     db,
	}, nil
}

// CleanupExpiredData cleans up expired users and their data
func (dcs *DatabaseCleanupService) CleanupExpiredData(ctx context.Context) error {
	dcs.appCtx.Logger.Info("Starting database cleanup", map[string]any{
		"timestamp": time.Now(),
	})

	// Add timeout to prevent hanging
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Clean up users who requested deletion and are past the grace period
	if err := dcs.cleanupExpiredUsers(ctx); err != nil {
		dcs.appCtx.Logger.Error("Failed to cleanup expired users", map[string]any{
			"error": err.Error(),
		})
		// Don't return error, continue with other cleanup tasks
	}

	dcs.appCtx.Logger.Info("Database cleanup completed", map[string]any{
		"timestamp": time.Now(),
	})
	return nil
}

// cleanupExpiredUsers removes users who requested deletion and are past the grace period
func (dcs *DatabaseCleanupService) cleanupExpiredUsers(ctx context.Context) error {
	// Find users who requested deletion 30+ days ago
	expiredUsers, err := dcs.getExpiredUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired users: %w", err)
	}

	if len(expiredUsers) == 0 {
		dcs.appCtx.Logger.Debug("No expired users found", map[string]any{})
		return nil
	}

	dcs.appCtx.Logger.Info("Found expired users", map[string]any{
		"count": len(expiredUsers),
	})

	// Clean up each expired user
	for _, userID := range expiredUsers {
		if err := dcs.cleanupUserData(ctx, userID); err != nil {
			dcs.appCtx.Logger.Error("Failed to cleanup user data", map[string]any{
				"user_id": userID,
				"error":   err,
			})
			// Continue with other users even if one fails
			continue
		}
	}

	return nil
}

// getExpiredUsers gets a list of user IDs who requested deletion 30+ days ago
func (dcs *DatabaseCleanupService) getExpiredUsers(ctx context.Context) ([]string, error) {
	var userIDs []string
	query := `
		SELECT id
		FROM users
		WHERE deletion_requested_at IS NOT NULL
		AND deletion_requested_at < NOW() - INTERVAL '30 days'
		AND deleted_at IS NULL
	`
	err := dcs.db.SelectContext(ctx, &userIDs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query expired users: %w", err)
	}

	return userIDs, nil
}

// cleanupUserData removes all data associated with a user
func (dcs *DatabaseCleanupService) cleanupUserData(ctx context.Context, userID string) error {
	dcs.appCtx.Logger.Debug("Cleaning up user data", map[string]any{
		"user_id": userID,
	})

	// Start a transaction with timeout
	tx, err := dcs.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			dcs.appCtx.Logger.Error("Failed to rollback transaction", map[string]any{
				"user_id": userID,
				"error":   err,
			})
		}
	}()

	// Delete user data in order (respecting foreign key constraints)
	if err := dcs.deleteUserGames(ctx, tx, userID); err != nil {
		return fmt.Errorf("failed to delete user games: %w", err)
	}

	if err := dcs.deletePhysicalLocations(ctx, tx, userID); err != nil {
		return fmt.Errorf("failed to delete physical locations: %w", err)
	}

	if err := dcs.deleteDigitalLocations(ctx, tx, userID); err != nil {
		return fmt.Errorf("failed to delete digital locations: %w", err)
	}

	if err := dcs.deleteSpendingData(ctx, tx, userID); err != nil {
		return fmt.Errorf("failed to delete spending data: %w", err)
	}

	if err := dcs.deleteWishlistItems(ctx, tx, userID); err != nil {
		return fmt.Errorf("failed to delete wishlist items: %w", err)
	}

	// Finally, delete the user
	if err := dcs.deleteUser(ctx, tx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	dcs.appCtx.Logger.Debug("User data cleanup completed", map[string]any{
		"user_id": userID,
	})

	return nil
}

// deleteUserGames deletes all games for a user
func (dcs *DatabaseCleanupService) deleteUserGames(ctx context.Context, tx *sqlx.Tx, userID string) error {
	query := `DELETE FROM user_games WHERE user_id = $1`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

// deletePhysicalLocations deletes all physical locations for a user
func (dcs *DatabaseCleanupService) deletePhysicalLocations(ctx context.Context, tx *sqlx.Tx, userID string) error {
	query := `DELETE FROM physical_locations WHERE user_id = $1`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

// deleteDigitalLocations deletes all digital locations for a user
func (dcs *DatabaseCleanupService) deleteDigitalLocations(ctx context.Context, tx *sqlx.Tx, userID string) error {
	query := `DELETE FROM digital_locations WHERE user_id = $1`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

// deleteSpendingData deletes all spending data for a user
func (dcs *DatabaseCleanupService) deleteSpendingData(ctx context.Context, tx *sqlx.Tx, userID string) error {
	query := `DELETE FROM spending_data WHERE user_id = $1`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

// deleteWishlistItems deletes all wishlist items for a user
func (dcs *DatabaseCleanupService) deleteWishlistItems(ctx context.Context, tx *sqlx.Tx, userID string) error {
	query := `DELETE FROM wishlist_items WHERE user_id = $1`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

// deleteUser deletes a user
func (dcs *DatabaseCleanupService) deleteUser(ctx context.Context, tx *sqlx.Tx, userID string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}