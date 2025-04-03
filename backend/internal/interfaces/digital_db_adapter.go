package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type DigitalDbAdapter interface {
	GetUserPhysicalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error)
    GetPhysicalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error)
    AddPhysicalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error)
    UpdatePhysicalLocation(ctx context.Context, userID string, location models.DigitalLocation) error
    RemovePhysicalLocation(ctx context.Context, userID, locationID string) error
}
