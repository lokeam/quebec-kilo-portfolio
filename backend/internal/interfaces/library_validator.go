package interfaces

import "github.com/lokeam/qko-beta/internal/models"

type LibraryValidator interface {
	ValidateGame(game models.LibraryGame) error
	ValidateUserID(userID string) error
	ValidateGameID(gameID int64) error
}