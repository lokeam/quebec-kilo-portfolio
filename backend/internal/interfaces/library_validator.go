package interfaces

import "github.com/lokeam/qko-beta/internal/models"

type LibraryValidator interface {
	ValidateLibraryGame(game models.GameToSave) error
	ValidateUserID(userID string) error
	ValidateGameID(gameID int64) error
}