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
	appCtx     *appcontext.AppContext
	dbAdapter  interfaces.UserDbAdapter
	sanitizer  interfaces.Sanitizer
	validator  interfaces.UserValidator
}

func NewUserService(appCtx *appcontext.AppContext) (*UserService, error) {
	dbAdapter, err := NewUserDbAdapter(appCtx)
	if err != nil {
		appCtx.Logger.Error("Failed to create dbAdapter", map[string]any{"error": err})
		return nil, err
	}
	appCtx.Logger.Info("User dbAdapter created successfully", nil)

	sanitizer, err := security.NewSanitizer()
	if err != nil {
		appCtx.Logger.Error("Failed to create sanitizer", map[string]any{"error": err})
		return nil, err
	}
	appCtx.Logger.Info("User sanitizer created successfully", nil)

	validator, err := NewUserValidator(sanitizer)
	if err != nil {
		appCtx.Logger.Error("Failed to create validator", map[string]any{"error": err})
		return nil, err
	}
	appCtx.Logger.Info("User validator created successfully", nil)

	return &UserService{
		appCtx:    appCtx,
		dbAdapter: dbAdapter,
		sanitizer: sanitizer,
		validator: validator,
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