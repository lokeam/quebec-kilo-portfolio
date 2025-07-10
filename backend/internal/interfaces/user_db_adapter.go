package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type UserDbAdapter interface {
	// User Profile Ops
	GetSingleUser(ctx context.Context, userID string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	UpdateUserProfile(ctx context.Context, userID string, firstName, lastName string) (models.User, error)
	HasCompleteProfile(ctx context.Context, userID string) (bool, error)

	// User Management Operations
	GetSingleUserByEmail(ctx context.Context, email string) (models.User, error)
}