package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type DigitalDbAdapter interface {
	GetUserDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error)
	FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error)
	AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error
	RemoveDigitalLocation(ctx context.Context, userID, locationID string) error
}
