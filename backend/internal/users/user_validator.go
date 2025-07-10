package users

import (
	"fmt"
	"strings"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/shared/logger"
	validationErrors "github.com/lokeam/qko-beta/internal/shared/validation"
	"github.com/lokeam/qko-beta/internal/types"
)

const (
	MaxNameLength = 50
	MinNameLength = 1
	MaxEmailLength = 255
)

type UserValidator struct {
	sanitizer  interfaces.Sanitizer
	timeSource  func() time.Time
	logger      logger.Logger
}

func NewUserValidator(sanitizer interfaces.Sanitizer) (*UserValidator, error) {
	log, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &UserValidator{
		sanitizer:   sanitizer,
		timeSource:  time.Now,
		logger:      *log,
	}, nil
}

func (uv *UserValidator) ValidateUserProfile(user models.User) (models.User, error) {
	uv.logger.Debug("Validating user profile", map[string]any{
		"userID": user.UserID,
		"email": user.Email,
	})

	var validatedUser models.User
	var violations []string

	// Copy ID and timestamps - these are required and thus don't need validation
	validatedUser.ID = user.ID
	validatedUser.UserID = user.UserID
	validatedUser.Email = user.Email
	validatedUser.CreatedAt = user.CreatedAt
	validatedUser.UpdatedAt = user.UpdatedAt

	// Validate first name
	if sanitizedFirstName, err := uv.validateFirstName(user.FirstName); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedUser.FirstName = sanitizedFirstName
	}

	// Validate last name
	if sanitizedLastName, err := uv.validateLastName(user.LastName); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedUser.LastName = sanitizedLastName
	}

	// Copy deletion fields (don't need validation)
	validatedUser.DeletionRequestedAt = user.DeletionRequestedAt
	validatedUser.DeletionReason = user.DeletionReason
	validatedUser.DeletedAt = user.DeletedAt

	if len(violations) > 0 {
		uv.logger.Debug("User profile validation failed", map[string]any{
			"violations": violations,
		})
		return models.User{}, &validationErrors.ValidationError{
			Field:   "user_profile",
			Message: fmt.Sprintf("User profile validation failed: %v", violations),
		}
	}

	uv.logger.Debug("User profile validation successful", map[string]any{
		"userID": validatedUser.UserID,
	})
	return validatedUser, nil
}

func (uv *UserValidator) ValidateCreateUserRequest(request types.CreateUserRequest) (types.CreateUserRequest, error) {
	uv.logger.Debug("Validating create user request", map[string]any{
		"auth0UserID": request.Auth0UserID,
		"email": request.Email,
	})

	var validatedRequest types.CreateUserRequest
	var violations []string

	// Validate Auth0 User ID
	if sanitizedAuth0UserID, err := uv.validateAuth0UserID(request.Auth0UserID); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedRequest.Auth0UserID = sanitizedAuth0UserID
	}

	// Validate email
	if sanitizedEmail, err := uv.validateEmail(request.Email); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedRequest.Email = sanitizedEmail
	}

	// Validate first name
	if sanitizedFirstName, err := uv.validateFirstName(request.FirstName); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedRequest.FirstName = sanitizedFirstName
	}

	// Validate last name
	if sanitizedLastName, err := uv.validateLastName(request.LastName); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedRequest.LastName = sanitizedLastName
	}

	if len(violations) > 0 {
		uv.logger.Debug("Create user request validation failed", map[string]any{
			"violations": violations,
		})
		return types.CreateUserRequest{}, &validationErrors.ValidationError{
			Field:   "create_user_request",
			Message: fmt.Sprintf("Create user request validation failed: %v", violations),
		}
	}

	uv.logger.Debug("Create user request validation successful", map[string]any{
		"auth0UserID": validatedRequest.Auth0UserID,
	})
	return validatedRequest, nil
}

func (uv *UserValidator) ValidateUpdateUserProfileRequest(request types.UpdateUserProfileRequest) (types.UpdateUserProfileRequest, error) {
	uv.logger.Debug("Validating update user profile request", map[string]any{
		"firstName": request.FirstName,
		"lastName": request.LastName,
	})

	var validatedRequest types.UpdateUserProfileRequest
	var violations []string

	// Validate first name
	if sanitizedFirstName, err := uv.validateFirstName(request.FirstName); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedRequest.FirstName = sanitizedFirstName
	}

	// Validate last name
	if sanitizedLastName, err := uv.validateLastName(request.LastName); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedRequest.LastName = sanitizedLastName
	}

	if len(violations) > 0 {
		uv.logger.Debug("Update user profile request validation failed", map[string]any{
			"violations": violations,
		})
		return types.UpdateUserProfileRequest{}, &validationErrors.ValidationError{
			Field:   "update_user_profile_request",
			Message: fmt.Sprintf("Update user profile request validation failed: %v", violations),
		}
	}

	uv.logger.Debug("Update user profile request validation successful", map[string]any{
		"firstName": validatedRequest.FirstName,
		"lastName": validatedRequest.LastName,
	})
	return validatedRequest, nil
}

// Helper validation methods
func (uv *UserValidator) validateAuth0UserID(userID string) (string, error) {
	if userID == "" {
		return "", &validationErrors.ValidationError{
			Field:   "auth0_user_id",
			Message: "Auth0 user ID cannot be empty",
		}
	}

	// Auth0 user IDs are typically in format: auth0|1234567890abcdef
	if !strings.Contains(userID, "|") {
		return "", &validationErrors.ValidationError{
			Field:   "auth0_user_id",
			Message: "Invalid Auth0 user ID format",
		}
	}

	sanitized, err := uv.sanitizer.SanitizeString(userID)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "auth0_user_id",
			Message: fmt.Sprintf("invalid Auth0 user ID content: %v", err),
		}
	}

	return sanitized, nil
}

func (uv *UserValidator) validateEmail(email string) (string, error) {
	if email == "" {
		return "", &validationErrors.ValidationError{
			Field:   "email",
			Message: "Email cannot be empty",
		}
	}

	if len(email) > MaxEmailLength {
		return "", &validationErrors.ValidationError{
			Field:   "email",
			Message: fmt.Sprintf("Email must be less than %d characters", MaxEmailLength),
		}
	}

	// Basic email format validation
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return "", &validationErrors.ValidationError{
			Field:   "email",
			Message: "Invalid email format",
		}
	}

	sanitized, err := uv.sanitizer.SanitizeString(email)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "email",
			Message: fmt.Sprintf("invalid email content: %v", err),
		}
	}

	return sanitized, nil
}

func (uv *UserValidator) validateFirstName(firstName string) (string, error) {
	if firstName == "" {
		return "", &validationErrors.ValidationError{
			Field:   "first_name",
			Message: "First name cannot be empty",
		}
	}

	if len(firstName) < MinNameLength {
		return "", &validationErrors.ValidationError{
			Field:   "first_name",
			Message: fmt.Sprintf("First name must be at least %d character", MinNameLength),
		}
	}

	if len(firstName) > MaxNameLength {
		return "", &validationErrors.ValidationError{
			Field:   "first_name",
			Message: fmt.Sprintf("First name must be less than %d characters", MaxNameLength),
		}
	}

	sanitized, err := uv.sanitizer.SanitizeString(firstName)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "first_name",
			Message: fmt.Sprintf("invalid first name content: %v", err),
		}
	}

	return sanitized, nil
}

func (uv *UserValidator) validateLastName(lastName string) (string, error) {
	if lastName == "" {
		return "", &validationErrors.ValidationError{
			Field:   "last_name",
			Message: "Last name cannot be empty",
		}
	}

	if len(lastName) < MinNameLength {
		return "", &validationErrors.ValidationError{
			Field:   "last_name",
			Message: fmt.Sprintf("Last name must be at least %d character", MinNameLength),
		}
	}

	if len(lastName) > MaxNameLength {
		return "", &validationErrors.ValidationError{
			Field:   "last_name",
			Message: fmt.Sprintf("Last name must be less than %d characters", MaxNameLength),
		}
	}

	sanitized, err := uv.sanitizer.SanitizeString(lastName)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "last_name",
			Message: fmt.Sprintf("invalid last name content: %v", err),
		}
	}

	return sanitized, nil
}