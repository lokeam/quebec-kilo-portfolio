package interfaces

import (
	"github.com/lokeam/qko-beta/internal/models"
)

type UserDeletionValidator interface {
	ValidateDeletionRequest(reason string) (string, error)
	ValidateUserID(userID string) (string, error)
	ValidateGracePeriod(user models.User) error
}