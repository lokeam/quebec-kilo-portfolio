package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type UserDeletionDbAdapter interface {
	RequestDeletion(ctx context.Context, userID string, reason string) error
	CancelDeletionRequest(ctx context.Context, userID string) error
	GetUser(ctx context.Context, userID string) (models.User, error)
	GetUsersPendingDeletion(ctx context.Context) ([]string, error)
	PermanentlyDeleteUser(ctx context.Context, userID string) error
}