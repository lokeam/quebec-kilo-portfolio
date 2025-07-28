package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
)

type UserDeletionDbAdapter struct {
	db     *sqlx.DB
	logger interfaces.Logger
}

func NewUserDeletionDbAdapter(appContext *appcontext.AppContext) (*UserDeletionDbAdapter, error) {
	appContext.Logger.Debug("Creating UserDeletionDbAdapter", map[string]any{"appContext": appContext})

	// Use shared DB pool
	db := appContext.DB

	return &UserDeletionDbAdapter{
		db:     db,
		logger: appContext.Logger,
	}, nil
}

// RequestDeletion marks a user for deletion (soft delete)
func (uda *UserDeletionDbAdapter) RequestDeletion(
	ctx context.Context,
	userID string,
	reason string,
) error {
	uda.logger.Debug("Requesting user deletion", map[string]any{
		"userID": userID,
		"reason": reason,
	})

	result, err := uda.db.ExecContext(
		ctx,
		RequestDeletionQuery,
		userID,
		reason,
	)
	if err != nil {
		uda.logger.Error("Failed to request deletion", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return fmt.Errorf("failed to request deletion: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		uda.logger.Error("User not found or already deleted", map[string]any{
			"userID": userID,
		})
		return fmt.Errorf("user not found or already deleted")
	}

	uda.logger.Info("User deletion requested successfully", map[string]any{
		"userID": userID,
		"reason": reason,
	})

	return nil
}

// CancelDeletionRequest cancels a pending deletion request
func (uda *UserDeletionDbAdapter) CancelDeletionRequest(ctx context.Context, userID string) error {
	uda.logger.Debug("Canceling user deletion request", map[string]any{
		"userID": userID,
	})

	result, err := uda.db.ExecContext(
		ctx,
		CancelDeletionRequestQuery,
		userID,
	)
	if err != nil {
		uda.logger.Error("Failed to cancel deletion request", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return fmt.Errorf("failed to cancel deletion request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		uda.logger.Error("No pending deletion request found", map[string]any{
			"userID": userID,
		})
		return fmt.Errorf("no pending deletion request found")
	}

	uda.logger.Info("User deletion request cancelled successfully", map[string]any{
		"userID": userID,
	})

	return nil
}

// GetUser retrieves a user by ID
func (uda *UserDeletionDbAdapter) GetUser(ctx context.Context, userID string) (models.User, error) {
	uda.logger.Debug("Getting user for deletion", map[string]any{"userID": userID})

	var user models.User
	err := uda.db.GetContext(
		ctx,
		&user,
		GetUserQuery,
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			uda.logger.Error("User not found", map[string]any{"userID": userID})
			return models.User{}, fmt.Errorf("user not found: %w", err)
		}
		uda.logger.Error("Failed to get user", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	uda.logger.Debug("User retrieved successfully", map[string]any{
		"userID": userID,
		"email":  user.Email,
	})

	return user, nil
}

// GetUsersPendingDeletion returns users who have requested deletion and are past the grace period
func (uda *UserDeletionDbAdapter) GetUsersPendingDeletion(ctx context.Context) ([]string, error) {
	uda.logger.Debug("Getting users pending deletion", map[string]any{})

	var userIDs []string
	err := uda.db.SelectContext(
		ctx,
		&userIDs,
		GetUsersPendingDeletionQuery,
	)
	if err != nil {
		uda.logger.Error("Failed to get users pending deletion", map[string]any{
			"error": err,
		})
		return nil, fmt.Errorf("failed to get users pending deletion: %w", err)
	}

	uda.logger.Info("Users pending deletion retrieved", map[string]any{
		"count": len(userIDs),
	})

	return userIDs, nil
}

// PermanentlyDeleteUser marks a user as permanently deleted
// This is a soft delete and is required by GDPR Article 17 (Right to be forgotten)
// Method needed for audit purposes
func (uda *UserDeletionDbAdapter) PermanentlyDeleteUser(ctx context.Context, userID string) error {
	uda.logger.Debug("Permanently deleting user", map[string]any{"userID": userID})

	result, err := uda.db.ExecContext(
		ctx,
		PermanentlyDeleteUserQuery,
		userID,
	)
	if err != nil {
		uda.logger.Error("Failed to permanently delete user", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return fmt.Errorf("failed to permanently delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		uda.logger.Error("User not found or not eligible for permanent deletion", map[string]any{
			"userID": userID,
		})
		return fmt.Errorf("user not found or not eligible for permanent deletion")
	}

	uda.logger.Info("User permanently deleted successfully", map[string]any{
		"userID": userID,
	})

	return nil
}