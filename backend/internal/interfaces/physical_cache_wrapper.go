package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type PhysicalCacheWrapper interface {
	GetCachedPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
	SetCachedPhysicalLocations(ctx context.Context, userID string, locations []models.PhysicalLocation) error
	GetCachedPhysicalLocationsBFF(ctx context.Context, userID string) (types.LocationsBFFResponse, error)
	SetCachedPhysicalLocationsBFF(ctx context.Context, userID string, locations types.LocationsBFFResponse) error
	GetSingleCachedPhysicalLocation(ctx context.Context, userID string, locationID string) (*models.PhysicalLocation, bool, error)
	SetSingleCachedPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) error
	InvalidateUserCache(ctx context.Context, userID string) error
	InvalidateLocationCache(ctx context.Context, userID string, locationID string) error
	InvalidateCache(ctx context.Context, cacheKey string) error
}