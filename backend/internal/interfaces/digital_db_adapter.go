package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type DigitalDbAdapter interface {
	// Digital Location Operations
	GetUserDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error)
	FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error)
	AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error
	RemoveDigitalLocation(ctx context.Context, userID, locationID string) error

	// Game Management Operations
	AddGameToDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error
	RemoveGameFromDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error
	GetGamesByDigitalLocationID(ctx context.Context, userID string, locationID string) ([]models.Game, error)

	// Subscription Operations
	GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error)
	AddSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, subscription models.Subscription) error
	RemoveSubscription(ctx context.Context, locationID string) error
	ValidateSubscriptionExists(ctx context.Context, locationID string) (*models.Subscription, error)

	// Payment Operations
	GetPayments(ctx context.Context, locationID string) ([]models.Payment, error)
	AddPayment(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetPayment(ctx context.Context, paymentID int64) (*models.Payment, error)
}
