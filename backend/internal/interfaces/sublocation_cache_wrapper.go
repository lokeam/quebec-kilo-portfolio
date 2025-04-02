package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type SublocationCacheWrapper interface {
	GetCachedSublocations(ctx context.Context, userID string) ([]models.Sublocation, error)
	SetCachedSublocations(ctx context.Context, userID string, sublocations []models.Sublocation) error
	GetSingleCachedSublocation(ctx context.Context, userID string, sublocationID string) (*models.Sublocation, bool, error)
	SetSingleCachedSublocation(ctx context.Context, userID string, sublocation models.Sublocation) error
	InvalidateUserCache(ctx context.Context, userID string) error
	InvalidateSublocationCache(ctx context.Context, userID string, locationID string) error
}
