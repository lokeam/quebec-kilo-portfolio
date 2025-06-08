package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type SublocationDbAdapter interface {
	GetAllSublocations(ctx context.Context, userID string) ([]models.Sublocation, error)
	GetSingleSublocation(ctx context.Context, userID, locationID string) (models.Sublocation, error)
	CreateSublocation(ctx context.Context, userID string, location models.Sublocation) (models.Sublocation, error)
	UpdateSublocation(ctx context.Context, userID string, location models.Sublocation) error
	DeleteSublocation(ctx context.Context, userID, locationID string) error
	CheckDuplicateSublocation(ctx context.Context, userID string, physicalLocationID string, name string) (bool, error)

	// Game management methods
	MoveGameToSublocation(ctx context.Context, userID string, userGameID string, targetSublocationID string) error
	RemoveGameFromSublocation(ctx context.Context, userID string, userGameID string) error
	CheckGameInAnySublocation(ctx context.Context, userGameID string) (bool, error)
	CheckGameInSublocation(ctx context.Context, userGameID string, sublocationID string) (bool, error)
	CheckGameOwnership(ctx context.Context, userID string, userGameID string) (bool, error)
}
