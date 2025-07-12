package users

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
	"github.com/lokeam/qko-beta/internal/types"
)

type UserService struct {
	appCtx      *appcontext.AppContext
	dbAdapter   interfaces.UserDbAdapter
	sanitizer   interfaces.Sanitizer
	validator   interfaces.UserValidator
	auth0Adapter *Auth0Adapter
}

func NewUserService(appCtx *appcontext.AppContext) (*UserService, error) {
	appCtx.Logger.Debug("UserService: constructor entered", nil)
	appCtx.Logger.Debug("UserService: Starting creation", nil)

	dbAdapter, err := NewUserDbAdapter(appCtx)
	if err != nil {
		appCtx.Logger.Error("UserService: Failed to create dbAdapter", map[string]any{"error": err})
		return nil, err
	}
	appCtx.Logger.Debug("UserService: dbAdapter created", nil)

	sanitizer, err := security.NewSanitizer()
	if err != nil {
		appCtx.Logger.Error("UserService: Failed to create sanitizer", map[string]any{"error": err})
		return nil, err
	}
	appCtx.Logger.Debug("UserService: sanitizer created", nil)

	validator, err := NewUserValidator(sanitizer)
	if err != nil {
		appCtx.Logger.Error("UserService: Failed to create validator", map[string]any{"error": err})
		return nil, err
	}
	appCtx.Logger.Debug("UserService: validator created", nil)

	auth0Adapter, err := NewAuth0Adapter(appCtx)
	if err != nil {
		appCtx.Logger.Error("UserService: Failed to create auth0Adapter", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to create auth0 adapter: %w", err)
	}
	appCtx.Logger.Debug("UserService: auth0Adapter created", nil)

	appCtx.Logger.Debug("UserService: Successfully created", nil)
	return &UserService{
		appCtx:       appCtx,
		dbAdapter:    dbAdapter,
		sanitizer:    sanitizer,
		validator:    validator,
		auth0Adapter: auth0Adapter,
	}, nil
}

func (s *UserService) GetSingleUser(ctx context.Context, userID string) (models.User, error) {
	s.appCtx.Logger.Debug("GetSingleUser called", map[string]any{"userID": userID})

	if userID == "" {
		return models.User{}, fmt.Errorf("user ID cannot be empty")
	}

	user, err := s.dbAdapter.GetSingleUser(ctx, userID)
	if err != nil {
		s.appCtx.Logger.Error("Failed to get user", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	s.appCtx.Logger.Debug("GetSingleUser success", map[string]any{
		"userID": userID,
		"email":  user.Email,
	})

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, req types.CreateUserRequest) (models.User, error) {
	s.appCtx.Logger.Debug("CreateUser called", map[string]any{
		"auth0UserID": req.Auth0UserID,
		"email":       req.Email,
	})

	// Validate the request
	validatedReq, err := s.validator.ValidateCreateUserRequest(req)
	if err != nil {
		s.appCtx.Logger.Error("CreateUser validation failed", map[string]any{
			"error": err,
		})
		return models.User{}, fmt.Errorf("validation failed: %w", err)
	}

	// Create user model
	user := models.User{
		UserID:    validatedReq.Auth0UserID,
		Email:     validatedReq.Email,
		FirstName: validatedReq.FirstName,
		LastName:  validatedReq.LastName,
	}

	// Create user in database
	createdUser, err := s.dbAdapter.CreateUser(ctx, user)
	if err != nil {
		s.appCtx.Logger.Error("Failed to create user in database", map[string]any{
			"auth0UserID": validatedReq.Auth0UserID,
			"error":        err,
		})
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	// Sync user metadata with Auth0
	auth0Metadata := map[string]any{
		"firstName": createdUser.FirstName,
		"lastName":  createdUser.LastName,
		"email":     createdUser.Email,
	}

	if err := s.auth0Adapter.PatchUserMetadata(ctx, createdUser.UserID, auth0Metadata); err != nil {
		s.appCtx.Logger.Error("Failed to sync user metadata with Auth0", map[string]any{
			"userID": createdUser.UserID,
			"error":  err,
		})
		// Note: We don't fail the entire operation if Auth0 sync fails
		// The user is still created in our database, but we log the error
		// You could choose to fail here if Auth0 sync is critical
	}

	s.appCtx.Logger.Info("User created successfully", map[string]any{
		"userID": createdUser.ID,
		"email":  createdUser.Email,
	})

	return createdUser, nil
}

func (s *UserService) UpdateUserProfile(
	ctx context.Context,
	userID string,
	req types.UpdateUserProfileRequest,
) (models.User, error) {
	s.appCtx.Logger.Debug("UpdateUserProfile called", map[string]any{
		"userID":    userID,
		"firstName": req.FirstName,
		"lastName":  req.LastName,
	})

	if userID == "" {
		return models.User{}, fmt.Errorf("user ID cannot be empty")
	}

	// Validate the request
	validatedReq, err := s.validator.ValidateUpdateUserProfileRequest(req)
	if err != nil {
		s.appCtx.Logger.Error("UpdateUserProfile validation failed", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return models.User{}, fmt.Errorf("validation failed: %w", err)
	}

	// Update user profile in database
	updatedUser, err := s.dbAdapter.UpdateUserProfile(ctx, userID, validatedReq.FirstName, validatedReq.LastName)
	if err != nil {
		s.appCtx.Logger.Error("Failed to update user profile", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return models.User{}, fmt.Errorf("failed to update user profile: %w", err)
	}

	// Sync user metadata with Auth0
	auth0Metadata := map[string]any{
		"firstName": updatedUser.FirstName,
		"lastName":  updatedUser.LastName,
		"email":     updatedUser.Email,
	}

	if err := s.auth0Adapter.PatchUserMetadata(ctx, userID, auth0Metadata); err != nil {
		s.appCtx.Logger.Error("Failed to sync user metadata with Auth0", map[string]any{
			"userID": userID,
			"error":  err,
		})
		// Note: We don't fail the entire operation if Auth0 sync fails
		// The user profile is still updated in our database, but we log the error
		// You could choose to fail here if Auth0 sync is critical
	}

	s.appCtx.Logger.Info("User profile updated successfully", map[string]any{
		"userID":    userID,
		"firstName": updatedUser.FirstName,
		"lastName":  updatedUser.LastName,
	})

	return updatedUser, nil
}



// Helper method to check if user has complete profile
func (s *UserService) HasCompleteProfile(ctx context.Context, userID string) (bool, error) {
	s.appCtx.Logger.Debug("HasCompleteProfile called", map[string]any{"userID": userID})

	if userID == "" {
		return false, fmt.Errorf("user ID cannot be empty")
	}

	hasComplete, err := s.dbAdapter.HasCompleteProfile(ctx, userID)
	if err != nil {
		s.appCtx.Logger.Error("Failed to check user profile completeness", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return false, fmt.Errorf("failed to check profile completeness: %w", err)
	}

	s.appCtx.Logger.Debug("HasCompleteProfile result", map[string]any{
		"userID":         userID,
		"hasComplete":    hasComplete,
	})

	return hasComplete, nil
}