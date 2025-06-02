package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type DigitalDbAdapter interface {
	// Digital Location Operations
	GetAllDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetSingleDigitalLocation(ctx context.Context, userID, digitalLocationID string) (models.DigitalLocation, error)
	FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error)
	CreateDigitalLocation(ctx context.Context, userID string, digitalLocation models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocation(ctx context.Context, userID string, digitalLocation models.DigitalLocation) error
	DeleteDigitalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error)

	// Game Management Operations
	AddGameToDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error
	RemoveGameFromDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error
	GetGamesByDigitalLocationID(ctx context.Context, userID string, locationID string) ([]models.Game, error)

	// Subscription Operations
	GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error)
	CreateSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, subscription models.Subscription) error
	DeleteSubscription(ctx context.Context, locationID string) error
	ValidateSubscriptionExists(ctx context.Context, locationID string) (*models.Subscription, error)

	// Payment Operations
	GetAllPayments(ctx context.Context, locationID string) ([]models.Payment, error)
	CreatePayment(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetSinglePayment(ctx context.Context, paymentID int64) (*models.Payment, error)
}
