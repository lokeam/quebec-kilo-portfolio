package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type LibraryDbAdapter interface {
	GetSingleLibraryGame(ctx context.Context, userID string, gameID int64) (types.LibraryGameItemBFFResponseFINAL, error)
	GetUserLibraryItems(ctx context.Context, userID string) ([]models.GameToSave, error)
	UpdateLibraryGame(ctx context.Context, game models.GameToSave) error
	CreateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error
	DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error
	IsGameInLibrary(ctx context.Context, userID string, gameID int64) (bool, error)

	// REFACTORED RESPONSE WITH HELPER METHODS
	GetLibraryRefactoredBFFResponse(ctx context.Context,userID string) (types.LibraryBFFRefactoredResponse, error)
	GetPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocationDB, error)
	GetDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocationDB, error)
	GetGamesMetadata(ctx context.Context, userID string) ([]models.LibraryGameRefactoredDB, error)

	// MARKED FOR DELETION - LEGACY RESPONSE
	GetLibraryBFFResponse(ctx context.Context, userID string) (types.LibraryBFFResponseFINAL, error)
}