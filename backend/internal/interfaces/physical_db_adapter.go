package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type PhysicalDbAdapter interface {
	  GetUserPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
    GetPhysicalLocation(ctx context.Context, userID, locationID string) (models.PhysicalLocation, error)
    AddPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
    UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) error
    RemovePhysicalLocation(ctx context.Context, userID, locationID string) error
}