package interfaces

import "github.com/lokeam/qko-beta/internal/models"

type SublocationValidator interface {
	ValidateSublocation(sublocation models.Sublocation) (models.Sublocation, error)
	ValidateSublocationUpdate(update, existing models.Sublocation) (models.Sublocation, error)
	ValidateSublocationCreation(sublocation models.Sublocation) (models.Sublocation, error)

	// New methods for game operations
	ValidateGameOwnership(userID string, userGameID string) error
	ValidateSublocationOwnership(userID string, sublocationID string) error
	ValidateGameNotInSublocation(userGameID string, sublocationID string) error
}
