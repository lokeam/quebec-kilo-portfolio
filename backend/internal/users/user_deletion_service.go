package users

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
)

type UserDeletionService struct {
	appCtx *appcontext.AppContext
	db     *sqlx.DB
}

func NewUserDeletionService(appCtx *appcontext.AppContext) (*UserDeletionService, error) {
	if appCtx == nil {
		return nil, fmt.Errorf("app context cannot be nil")
	}

	// Create sqlx connection following the pattern of other services
	db, err := sqlx.Connect("pgx", appCtx.Config.Postgres.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	return &UserDeletionService{
		appCtx: appCtx,
		db:     db,
	}, nil
}

// RequestDeletion marks a user for deletion (soft delete)
func (uds *UserDeletionService) RequestDeletion(ctx context.Context, userID string, reason string) error {
	query := `
		UPDATE users
		SET deletion_requested_at = NOW(),
		    deletion_reason = $2,
		    updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := uds.db.ExecContext(ctx, query, userID, reason)
	if err != nil {
		return fmt.Errorf("failed to request deletion: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	uds.appCtx.Logger.Info("User deletion requested", map[string]any{
		"user_id": userID,
		"reason":  reason,
	})

	return nil
}

// CancelDeletionRequest cancels a pending deletion request
func (uds *UserDeletionService) CancelDeletionRequest(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET deletion_requested_at = NULL,
		    deletion_reason = NULL,
		    updated_at = NOW()
		WHERE id = $1 AND deletion_requested_at IS NOT NULL AND deleted_at IS NULL
	`

	result, err := uds.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to cancel deletion request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no pending deletion request found")
	}

	uds.appCtx.Logger.Info("User deletion request cancelled", map[string]any{
		"user_id": userID,
	})

	return nil
}

// GetUser retrieves a user by ID
func (uds *UserDeletionService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, user_id, email, user_type, created_at, updated_at,
		       deleted_at, deletion_requested_at, deletion_reason
		FROM users
		WHERE id = $1
	`

	err := uds.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUsersPendingDeletion returns users who have requested deletion and are past the grace period
func (uds *UserDeletionService) GetUsersPendingDeletion(ctx context.Context) ([]string, error) {
	var userIDs []string
	query := `
		SELECT id
		FROM users
		WHERE deletion_requested_at IS NOT NULL
		  AND deletion_requested_at < NOW() - INTERVAL '30 days'
		  AND deleted_at IS NULL
	`

	err := uds.db.SelectContext(ctx, &userIDs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users pending deletion: %w", err)
	}

	return userIDs, nil
}

// PermanentlyDeleteUser marks a user as permanently deleted
func (uds *UserDeletionService) PermanentlyDeleteUser(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET deleted_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1 AND deletion_requested_at IS NOT NULL
	`

	result, err := uds.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to permanently delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or not eligible for permanent deletion")
	}

	uds.appCtx.Logger.Info("User permanently deleted", map[string]any{
		"user_id": userID,
	})

	return nil
}