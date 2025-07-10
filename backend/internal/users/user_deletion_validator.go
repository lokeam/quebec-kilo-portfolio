package users

import (
	"fmt"
	"strings"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/shared/logger"
	validationErrors "github.com/lokeam/qko-beta/internal/shared/validation"
)

const (
	MinDeletionReasonLength = 10
	MaxDeletionReasonLength = 500
	GracePeriodDays         = 30
)

type UserDeletionValidator struct {
	sanitizer  interfaces.Sanitizer
	timeSource func() time.Time
	logger     logger.Logger
}

func NewUserDeletionValidator(sanitizer interfaces.Sanitizer) (*UserDeletionValidator, error) {
	log, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &UserDeletionValidator{
		sanitizer:  sanitizer,
		timeSource: time.Now,
		logger:     *log,
	}, nil
}

// ValidateDeletionRequest validates the deletion reason
func (udv *UserDeletionValidator) ValidateDeletionRequest(reason string) (string, error) {
	udv.logger.Debug("Validating deletion request", map[string]any{
		"reason": reason,
	})

	// Check if reason is empty
	if reason == "" {
		return "", &validationErrors.ValidationError{
			Field:   "deletion_reason",
			Message: "Deletion reason cannot be empty",
		}
	}

	// Check minimum length
	if len(reason) < MinDeletionReasonLength {
		return "", &validationErrors.ValidationError{
			Field:   "deletion_reason",
			Message: fmt.Sprintf("Deletion reason must be at least %d characters", MinDeletionReasonLength),
		}
	}

	// Check maximum length
	if len(reason) > MaxDeletionReasonLength {
		return "", &validationErrors.ValidationError{
			Field:   "deletion_reason",
			Message: fmt.Sprintf("Deletion reason must be less than %d characters", MaxDeletionReasonLength),
		}
	}

	// Sanitize the reason
	sanitizedReason, err := udv.sanitizer.SanitizeString(reason)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "deletion_reason",
			Message: fmt.Sprintf("Invalid deletion reason content: %v", err),
		}
	}

	udv.logger.Debug("Deletion request validation successful", map[string]any{
		"reason": sanitizedReason,
	})

	return sanitizedReason, nil
}

// ValidateUserID validates the user ID
func (udv *UserDeletionValidator) ValidateUserID(userID string) (string, error) {
	udv.logger.Debug("Validating user ID", map[string]any{
		"userID": userID,
	})

	// Check if user ID is empty
	if userID == "" {
		return "", &validationErrors.ValidationError{
			Field:   "user_id",
			Message: "User ID cannot be empty",
		}
	}

	// Auth0 user IDs are typically in format: auth0|1234567890abcdef
	if !strings.Contains(userID, "|") {
		return "", &validationErrors.ValidationError{
			Field:   "user_id",
			Message: "Invalid Auth0 user ID format",
		}
	}

	// Sanitize the user ID
	sanitizedUserID, err := udv.sanitizer.SanitizeString(userID)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "user_id",
			Message: fmt.Sprintf("Invalid user ID content: %v", err),
		}
	}

	udv.logger.Debug("User ID validation successful", map[string]any{
		"userID": sanitizedUserID,
	})

	return sanitizedUserID, nil
}

// ValidateGracePeriod validates if a user is eligible for deletion based on grace period
func (udv *UserDeletionValidator) ValidateGracePeriod(user models.User) error {
	udv.logger.Debug("Validating grace period", map[string]any{
		"userID": user.UserID,
		"deletionRequestedAt": user.DeletionRequestedAt,
	})

	// Check if user has requested deletion
	if user.DeletionRequestedAt == nil {
		return &validationErrors.ValidationError{
			Field:   "grace_period",
			Message: "User has not requested deletion",
		}
	}

	// Check if user is already deleted
	if user.DeletedAt != nil {
		return &validationErrors.ValidationError{
			Field:   "grace_period",
			Message: "User is already permanently deleted",
		}
	}

	// Calculate grace period end
	gracePeriodEnd := user.DeletionRequestedAt.AddDate(0, 0, GracePeriodDays)
	now := udv.timeSource()

	// Check if grace period has expired
	if now.Before(gracePeriodEnd) {
		return &validationErrors.ValidationError{
			Field:   "grace_period",
			Message: fmt.Sprintf("Grace period has not expired. Deletion will be available on %s", gracePeriodEnd.Format("2006-01-02")),
		}
	}

	udv.logger.Debug("Grace period validation successful", map[string]any{
		"userID": user.UserID,
		"gracePeriodEnd": gracePeriodEnd,
	})

	return nil
}

// ValidateUserExists validates that a user exists and is not already deleted
func (udv *UserDeletionValidator) ValidateUserExists(user models.User) error {
	udv.logger.Debug("Validating user exists", map[string]any{
		"userID": user.UserID,
	})

	// Check if user is already permanently deleted
	if user.DeletedAt != nil {
		return &validationErrors.ValidationError{
			Field:   "user_status",
			Message: "User is already permanently deleted",
		}
	}

	udv.logger.Debug("User exists validation successful", map[string]any{
		"userID": user.UserID,
	})

	return nil
}

// ValidateDeletionRequestExists validates that a user has a pending deletion request
func (udv *UserDeletionValidator) ValidateDeletionRequestExists(user models.User) error {
	udv.logger.Debug("Validating deletion request exists", map[string]any{
		"userID": user.UserID,
		"deletionRequestedAt": user.DeletionRequestedAt,
	})

	// Check if user has requested deletion
	if user.DeletionRequestedAt == nil {
		return &validationErrors.ValidationError{
			Field:   "deletion_request",
			Message: "No pending deletion request found",
		}
	}

	// Check if user is already deleted
	if user.DeletedAt != nil {
		return &validationErrors.ValidationError{
			Field:   "deletion_request",
			Message: "User is already permanently deleted",
		}
	}

	udv.logger.Debug("Deletion request exists validation successful", map[string]any{
		"userID": user.UserID,
	})

	return nil
}