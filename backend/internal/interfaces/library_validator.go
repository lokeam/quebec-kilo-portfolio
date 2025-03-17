package interfaces

import "github.com/lokeam/qko-beta/internal/types"

type LibraryValidator interface {
	ValidateGame(game types.Game) error
	ValidateUserID(userID string) error
}