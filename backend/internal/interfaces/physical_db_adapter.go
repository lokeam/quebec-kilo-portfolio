package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type PhysicalDbAdapter interface {
	GetAllPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
	GetAllPhysicalLocationsBFF(ctx context.Context, userID string) (types.LocationsBFFResponse, error)
	GetSinglePhysicalLocation(ctx context.Context, userID string, locationID string) (models.PhysicalLocation, error)
	CreatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	DeletePhysicalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error)
}