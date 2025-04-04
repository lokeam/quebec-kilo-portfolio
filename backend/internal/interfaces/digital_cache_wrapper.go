package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type DigitalCacheWrapper interface {
	GetCachedDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	SetCachedDigitalLocations(ctx context.Context, userID string, locations []models.DigitalLocation) error
	GetSingleCachedDigitalLocation(ctx context.Context, userID string, locationID string) (*models.DigitalLocation, bool, error)
	SetSingleCachedDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error
	InvalidateUserCache(ctx context.Context, userID string) error
	InvalidateDigitalLocationCache(ctx context.Context, userID string, locationID string) error
}