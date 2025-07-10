package users

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
)

type UserDeletionService struct {
	appCtx     *appcontext.AppContext
	dbAdapter  interfaces.UserDeletionDbAdapter
	sanitizer  interfaces.Sanitizer
	validator  interfaces.UserDeletionValidator
}

func NewUserDeletionService(appCtx *appcontext.AppContext) (*UserDeletionService, error) {
	if appCtx == nil {
		return nil, fmt.Errorf("app context cannot be nil")
	}

	// Create DB adapter
	dbAdapter, err := NewUserDeletionDbAdapter(appCtx)
	if err != nil {
		appCtx.Logger.Error("Failed to create user deletion db adapter", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to create user deletion db adapter: %w", err)
	}
	appCtx.Logger.Info("User deletion db adapter created successfully", nil)

	// Create sanitizer
	sanitizer, err := security.NewSanitizer()
	if err != nil {
		appCtx.Logger.Error("Failed to create sanitizer", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to create sanitizer: %w", err)
	}
	appCtx.Logger.Info("User deletion sanitizer created successfully", nil)

	// Create validator
	validator, err := NewUserDeletionValidator(sanitizer)
	if err != nil {
		appCtx.Logger.Error("Failed to create user deletion validator", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to create user deletion validator: %w", err)
	}
	appCtx.Logger.Info("User deletion validator created successfully", nil)

	return &UserDeletionService{
		appCtx:    appCtx,
		dbAdapter: dbAdapter,
		sanitizer: sanitizer,
		validator: validator,
	}, nil
}

// RequestDeletion marks a user for deletion (soft delete)
func (uds *UserDeletionService) RequestDeletion(
	ctx context.Context,
	userID string,
	reason string,
) error {
	uds.appCtx.Logger.Debug("RequestDeletion called", map[string]any{
		"userID": userID,
		"reason": reason,
	})

	// Validate user ID
	validatedUserID, err := uds.validator.ValidateUserID(userID)
	if err != nil {
		uds.appCtx.Logger.Error("User ID validation failed", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return fmt.Errorf("validation failed: %w", err)
	}

	// Validate deletion reason
	validatedReason, err := uds.validator.ValidateDeletionRequest(reason)
	if err != nil {
		uds.appCtx.Logger.Error("Deletion reason validation failed", map[string]any{
			"reason": reason,
			"error":  err,
		})
		return fmt.Errorf("validation failed: %w", err)
	}

	// Get user to validate they exist and aren't already deleted
	user, err := uds.dbAdapter.GetUser(ctx, validatedUserID)
	if err != nil {
		uds.appCtx.Logger.Error("Failed to get user for deletion request", map[string]any{
			"userID": validatedUserID,
			"error":  err,
		})
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Validate user exists and isn't already deleted
	if err := uds.validator.ValidateGracePeriod(user); err != nil {
		uds.appCtx.Logger.Error("User validation failed", map[string]any{
			"userID": validatedUserID,
			"error":  err,
		})
		return fmt.Errorf("user validation failed: %w", err)
	}

	// Request deletion through DB adapter
	err = uds.dbAdapter.RequestDeletion(
		ctx,
		validatedUserID,
		validatedReason,
	)
	if err != nil {
		uds.appCtx.Logger.Error("Failed to request deletion", map[string]any{
			"userID": validatedUserID,
			"reason": validatedReason,
			"error":  err,
		})
		return fmt.Errorf("failed to request deletion: %w", err)
	}

	uds.appCtx.Logger.Info("User deletion requested successfully", map[string]any{
		"userID": validatedUserID,
		"reason": validatedReason,
	})

	return nil
}

// CancelDeletionRequest cancels a pending deletion request
func (uds *UserDeletionService) CancelDeletionRequest(ctx context.Context, userID string) error {
	uds.appCtx.Logger.Debug("CancelDeletionRequest called", map[string]any{
		"userID": userID,
	})

	// Validate user ID
	validatedUserID, err := uds.validator.ValidateUserID(userID)
	if err != nil {
		uds.appCtx.Logger.Error("User ID validation failed", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return fmt.Errorf("validation failed: %w", err)
	}

	// Get user to validate they have a pending deletion request
	user, err := uds.dbAdapter.GetUser(ctx, validatedUserID)
	if err != nil {
		uds.appCtx.Logger.Error("Failed to get user for cancellation", map[string]any{
			"userID": validatedUserID,
			"error":  err,
		})
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Validate user has a pending deletion request (simplified validation)
	if user.DeletionRequestedAt == nil {
		uds.appCtx.Logger.Error("No pending deletion request found", map[string]any{
			"userID": validatedUserID,
		})
		return fmt.Errorf("no pending deletion request found")
	}

	// Cancel deletion request through DB adapter
	err = uds.dbAdapter.CancelDeletionRequest(ctx, validatedUserID)
	if err != nil {
		uds.appCtx.Logger.Error("Failed to cancel deletion request", map[string]any{
			"userID": validatedUserID,
			"error":  err,
		})
		return fmt.Errorf("failed to cancel deletion request: %w", err)
	}

	uds.appCtx.Logger.Info("User deletion request cancelled successfully", map[string]any{
		"userID": validatedUserID,
	})

	return nil
}

// GetUser retrieves a user by ID
func (uds *UserDeletionService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	uds.appCtx.Logger.Debug("GetUser called", map[string]any{
		"userID": userID,
	})

	// Validate user ID
	validatedUserID, err := uds.validator.ValidateUserID(userID)
	if err != nil {
		uds.appCtx.Logger.Error("User ID validation failed", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get user through DB adapter
	user, err := uds.dbAdapter.GetUser(ctx, validatedUserID)
	if err != nil {
		uds.appCtx.Logger.Error("Failed to get user", map[string]any{
			"userID": validatedUserID,
			"error":  err,
		})
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	uds.appCtx.Logger.Debug("User retrieved successfully", map[string]any{
		"userID": validatedUserID,
	})

	return &user, nil
}

// GetUsersPendingDeletion returns users who have requested deletion and are past the grace period
func (uds *UserDeletionService) GetUsersPendingDeletion(ctx context.Context) ([]string, error) {
	uds.appCtx.Logger.Debug("GetUsersPendingDeletion called", map[string]any{})

	// Get users pending deletion through DB adapter
	userIDs, err := uds.dbAdapter.GetUsersPendingDeletion(ctx)
	if err != nil {
		uds.appCtx.Logger.Error("Failed to get users pending deletion", map[string]any{
			"error": err,
		})
		return nil, fmt.Errorf("failed to get users pending deletion: %w", err)
	}

	uds.appCtx.Logger.Info("Users pending deletion retrieved", map[string]any{
		"count": len(userIDs),
	})

	return userIDs, nil
}

// PermanentlyDeleteUser marks a user as permanently deleted
func (uds *UserDeletionService) PermanentlyDeleteUser(ctx context.Context, userID string) error {
	uds.appCtx.Logger.Debug("PermanentlyDeleteUser called", map[string]any{
		"userID": userID,
	})

	// Validate user ID
	validatedUserID, err := uds.validator.ValidateUserID(userID)
	if err != nil {
		uds.appCtx.Logger.Error("User ID validation failed", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return fmt.Errorf("validation failed: %w", err)
	}

	// Get user to validate they are eligible for permanent deletion
	user, err := uds.dbAdapter.GetUser(ctx, validatedUserID)
	if err != nil {
		uds.appCtx.Logger.Error("Failed to get user for permanent deletion", map[string]any{
			"userID": validatedUserID,
			"error":  err,
		})
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Validate grace period has expired
	if err := uds.validator.ValidateGracePeriod(user); err != nil {
		uds.appCtx.Logger.Error("Grace period validation failed", map[string]any{
			"userID": validatedUserID,
			"error":  err,
		})
		return fmt.Errorf("grace period validation failed: %w", err)
	}

	// Permanently delete user through DB adapter
	err = uds.dbAdapter.PermanentlyDeleteUser(ctx, validatedUserID)
	if err != nil {
		uds.appCtx.Logger.Error("Failed to permanently delete user", map[string]any{
			"userID": validatedUserID,
			"error":  err,
		})
		return fmt.Errorf("failed to permanently delete user: %w", err)
	}

	uds.appCtx.Logger.Info("User permanently deleted successfully", map[string]any{
		"userID": validatedUserID,
	})

	return nil
}