package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type DigitalCacheWrapper interface {
	GetCachedDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	SetCachedDigitalLocations(ctx context.Context, userID string, locations []models.DigitalLocation) error
	GetSingleCachedDigitalLocation(ctx context.Context, userID string, locationID string) (*models.DigitalLocation, bool, error)
	SetSingleCachedDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error
	InvalidateUserCache(ctx context.Context, userID string) error
	InvalidateDigitalLocationCache(ctx context.Context, userID string, locationID string) error

	// Subscription caching
	GetCachedSubscription(ctx context.Context, locationID string) (*models.Subscription, bool, error)
	SetCachedSubscription(ctx context.Context, locationID string, subscription models.Subscription) error
	InvalidateSubscriptionCache(ctx context.Context, locationID string) error

	// Payment caching
	GetCachedPayments(ctx context.Context, locationID string) ([]models.Payment, error)
	SetCachedPayments(ctx context.Context, locationID string, payments []models.Payment) error
	InvalidatePaymentsCache(ctx context.Context, locationID string) error

	// BFF Response
	GetCachedDigitalLocationsBFF(ctx context.Context, userID string) (types.DigitalLocationsBFFResponse, error)
	SetCachedDigitalLocationsBFF(ctx context.Context, userID string, response types.DigitalLocationsBFFResponse) error
}