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

	// Create all dependencies
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
		"firstName":   req.FirstName,
		"lastName":    req.LastName,
	})

	// Validate the request
	validatedReq, err := s.validator.ValidateCreateUserRequest(req)
	if err != nil {
		s.appCtx.Logger.Error("CreateUser validation failed", map[string]any{
			"error": err,
		})
		return models.User{}, fmt.Errorf("validation failed: %w", err)
	}

	// Create When Complete" approach to ENSURE WE HAVE COMPLETE DATA
	if validatedReq.FirstName == "" || validatedReq.LastName == "" {
		s.appCtx.Logger.Error("CreateUser failed - incomplete data", map[string]any{
			"auth0UserID": validatedReq.Auth0UserID,
			"firstName":    validatedReq.FirstName,
			"lastName":     validatedReq.LastName,
		})
		return models.User{}, fmt.Errorf("first name and last name are required for user creation")
	}

	// Check if user already exists
	userExists, err := s.dbAdapter.CheckUserExists(ctx, validatedReq.Auth0UserID)
	if err != nil {
		s.appCtx.Logger.Error("Failed to check if user exists", map[string]any{
			"auth0UserID": validatedReq.Auth0UserID,
			"error":        err,
		})
		return models.User{}, fmt.Errorf("failed to check user existence: %w", err)
	}

	// If user already exists, update them with complete data
	if userExists {
		s.appCtx.Logger.Info("User already exists, updating with complete data", map[string]any{
			"auth0UserID": validatedReq.Auth0UserID,
			"firstName":    validatedReq.FirstName,
			"lastName":     validatedReq.LastName,
		})

		updatedUser, err := s.dbAdapter.UpdateUserProfile(
			ctx,
			validatedReq.Auth0UserID,
			validatedReq.FirstName,
			validatedReq.LastName,
		)
		if err != nil {
			s.appCtx.Logger.Error("Failed to update existing user", map[string]any{
				"auth0UserID": validatedReq.Auth0UserID,
				"error":        err,
			})
			return models.User{}, fmt.Errorf("failed to update existing user: %w", err)
		}

		// Sync updated app metadata with Auth0
		auth0Metadata := map[string]any{
			"hasCompletedOnboarding": true,
		}

		// Log the metadata payload before sending to Auth0
		s.appCtx.Logger.Info("UserService: Syncing app metadata with Auth0 after user update", map[string]any{
			"userID": updatedUser.UserID,
			"metadata": auth0Metadata,
			"metadataKeys": getMapKeys(auth0Metadata),
			"hasCompletedOnboarding": true,
		})

		if err := s.auth0Adapter.PatchAppMetadata(
			ctx,
			updatedUser.UserID,
			auth0Metadata,
		); err != nil {
					s.appCtx.Logger.Error("Failed to sync app metadata with Auth0", map[string]any{
			"userID": updatedUser.UserID,
			"error":  err,
		})
			// NOTE: See below note about failing entire operation
		}

		s.appCtx.Logger.Info("Existing user updated successfully", map[string]any{
			"userID": updatedUser.UserID,
			"email":  updatedUser.Email,
		})

		return updatedUser, nil
	}

	// Create new user with complete data
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

	// Sync app metadata with Auth0
	auth0Metadata := map[string]any{
		"hasCompletedOnboarding": true,
	}

			// Log the metadata payload before sending to Auth0
		s.appCtx.Logger.Info("UserService: Syncing app metadata with Auth0 after user creation", map[string]any{
			"userID": createdUser.UserID,
			"metadata": auth0Metadata,
			"metadataKeys": getMapKeys(auth0Metadata),
			"hasCompletedOnboarding": true,
		})

	if err := s.auth0Adapter.PatchAppMetadata(
		ctx,
		createdUser.UserID,
		auth0Metadata,
	); err != nil {
		s.appCtx.Logger.Error("Failed to sync app metadata with Auth0", map[string]any{
			"userID": createdUser.UserID,
			"error":  err,
		})
		// NOTE: We don't fail the entire operation if Auth0 sync fails
		// The user is still created in our database, but we log the error
		// We may choose to fail here if Auth0 sync is critical
	}

	s.appCtx.Logger.Info("User created successfully", map[string]any{
		"userID": createdUser.UserID,
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
	updatedUser, err := s.dbAdapter.UpdateUserProfile(
		ctx,
		userID,
		validatedReq.FirstName,
		validatedReq.LastName,
	)
	if err != nil {
		s.appCtx.Logger.Error("Failed to update user profile", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return models.User{}, fmt.Errorf("failed to update user profile: %w", err)
	}

	// Sync app metadata with Auth0
	auth0Metadata := map[string]any{
		"hasCompletedOnboarding": true,
	}

			// Log the metadata payload before sending to Auth0
		s.appCtx.Logger.Info("UserService: Syncing app metadata with Auth0 after profile update", map[string]any{
			"userID": userID,
			"metadata": auth0Metadata,
			"metadataKeys": getMapKeys(auth0Metadata),
			"hasCompletedOnboarding": true,
		})

	if err := s.auth0Adapter.PatchAppMetadata(
		ctx,
		userID,
		auth0Metadata,
	); err != nil {
		s.appCtx.Logger.Error("Failed to sync app metadata with Auth0", map[string]any{
			"userID": userID,
			"error":  err,
		})
		// NOTE: See above note about failing entire operation in CreateUser method. To review later
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

// UpdateAppMetadata updates app metadata in Auth0
func (s *UserService) UpdateAppMetadata(
	ctx context.Context,
	userID string,
	metadata map[string]any,
) error {
	s.appCtx.Logger.Debug("UpdateAppMetadata called", map[string]any{
		"userID":   userID,
		"metadata": metadata,
		"metadataKeys": getMapKeys(metadata),
	})

	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	// Validate user exists first
	user, err := s.GetSingleUser(ctx, userID)
	if err != nil {
		s.appCtx.Logger.Error("Failed to get user for metadata update", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Update metadata through Auth0 adapter
	if err := s.auth0Adapter.PatchAppMetadata(
		ctx,
		userID,
		metadata,
	); err != nil {
		s.appCtx.Logger.Error("Failed to update app metadata with Auth0", map[string]any{
			"userID": userID,
			"error":  err,
		})
		return fmt.Errorf("failed to update app metadata: %w", err)
	}

	s.appCtx.Logger.Info("App metadata updated successfully", map[string]any{
		"userID":   userID,
		"email":    user.Email,
		"metadata": metadata,
		"metadataKeys": getMapKeys(metadata),
	})

	return nil
}


// CreateUserFromID creates a user with just the Auth0 user ID (for middleware)
func (s *UserService) CreateUserFromID(ctx context.Context, userID string) error {
	s.appCtx.Logger.Debug("CreateUserFromID called", map[string]any{"userID": userID})

	// Check if user already exists
	exists, err := s.dbAdapter.CheckUserExists(ctx, userID)
	if err != nil {
			s.appCtx.Logger.Error("Failed to check user existence", map[string]any{
					"userID": userID,
					"error":  err,
			})
			return fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
			s.appCtx.Logger.Debug("User already exists", map[string]any{"userID": userID})
			return nil
	}

	// Create minimal user record
	user := models.User{
			UserID: userID,
			Email:  fmt.Sprintf("%s@placeholder.com", userID), // Placeholder email
	}

	_, err = s.dbAdapter.CreateUser(ctx, user)
	if err != nil {
			s.appCtx.Logger.Error("Failed to create minimal user", map[string]any{
					"userID": userID,
					"error":  err,
			})
			return fmt.Errorf("failed to create user: %w", err)
	}

	s.appCtx.Logger.Debug("CreateUserFromID success", map[string]any{"userID": userID})
	return nil
}

// UserExists checks if a user exists in the database
func (s *UserService) UserExists(ctx context.Context, userID string) (bool, error) {
	s.appCtx.Logger.Debug("UserExists called", map[string]any{"userID": userID})

	exists, err := s.dbAdapter.CheckUserExists(ctx, userID)
	if err != nil {
			s.appCtx.Logger.Error("Failed to check user existence", map[string]any{
					"userID": userID,
					"error":  err,
			})
			return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	s.appCtx.Logger.Debug("UserExists result", map[string]any{
			"userID": userID,
			"exists": exists,
	})

	return exists, nil
}