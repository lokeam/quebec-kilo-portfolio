package interfaces

import (
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type UserValidator interface {
	ValidateUserProfile(user models.User) (models.User, error)
	ValidateCreateUserRequest(request types.CreateUserRequest) (types.CreateUserRequest, error)
	ValidateUpdateUserProfileRequest(request types.UpdateUserProfileRequest) (types.UpdateUserProfileRequest, error)
}