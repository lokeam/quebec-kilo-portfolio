package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type PhysicalDbAdapter interface {
	GetAllPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
	GetSinglePhysicalLocation(ctx context.Context, userID string, locationID string) (models.PhysicalLocation, error)
	CreatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	DeletePhysicalLocation(ctx context.Context, userID string, locationID string) error
}