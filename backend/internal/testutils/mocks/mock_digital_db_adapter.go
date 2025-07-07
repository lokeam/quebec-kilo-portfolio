package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)


type MockDigitalDbAdapter struct {
	GetDigitalLocationFunc func(ctx context.Context, userID, digitalLocationID string) (models.DigitalLocation, error)
	GetDigitalLocationsFunc func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	AddDigitalLocationFunc func(ctx context.Context, userID string, digitalLocation models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocationFunc func(ctx context.Context, userID string, digitalLocation models.DigitalLocation) error
	RemoveDigitalLocationFunc func(ctx context.Context, userID string, locationIDs []string) (int64, error)
	FindDigitalLocationByNameFunc func(ctx context.Context, userID string, name string) (models.DigitalLocation, error)
	GetAllDigitalLocationsBFFFunc func(ctx context.Context, userID string) (types.DigitalLocationsBFFResponse, error)

	// Game Management Operations
	AddGameToDigitalLocationFunc func(ctx context.Context, userID string, locationID string, gameID int64) error
	RemoveGameFromDigitalLocationFunc func(ctx context.Context, userID string, locationID string, gameID int64) error
	GetGamesByDigitalLocationIDFunc func(ctx context.Context, userID string, locationID string) ([]models.Game, error)

	// Subscription Operations
	GetSubscriptionFunc func(ctx context.Context, locationID string) (*models.Subscription, error)
	AddSubscriptionFunc func(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscriptionFunc func(ctx context.Context, subscription models.Subscription) error
	RemoveSubscriptionFunc func(ctx context.Context, locationID string) error
	ValidateSubscriptionExistsFunc func(ctx context.Context, locationID string) (*models.Subscription, error)

	// Payment Operations
	GetPaymentsFunc func(ctx context.Context, locationID string) ([]models.Payment, error)
	AddPaymentFunc func(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetPaymentFunc func(ctx context.Context, paymentID int64) (*models.Payment, error)
}

// GET
func (m *MockDigitalDbAdapter) GetSingleDigitalLocation(
	ctx context.Context,
	userID,
	digitalLocationID string,
) (models.DigitalLocation, error) {
	digitalLocation, err := m.GetDigitalLocationFunc(ctx, userID, digitalLocationID)
	return digitalLocation, err
}

func (m *MockDigitalDbAdapter) GetAllDigitalLocations(
	ctx context.Context,
	userID string,
) ([]models.DigitalLocation, error) {
	return m.GetDigitalLocationsFunc(ctx, userID)
}

func (m *MockDigitalDbAdapter) GetAllDigitalLocationsBFF(
	ctx context.Context,
	userID string,
) (types.DigitalLocationsBFFResponse, error) {
	if m.GetAllDigitalLocationsBFFFunc != nil {
		return m.GetAllDigitalLocationsBFFFunc(ctx, userID)
	}
	return types.DigitalLocationsBFFResponse{}, nil
}

// POST
func (m *MockDigitalDbAdapter) CreateDigitalLocation(
	ctx context.Context,
	userID string,
	digitalLocation models.DigitalLocation,
) (models.DigitalLocation, error) {
	return m.AddDigitalLocationFunc(ctx, userID, digitalLocation)
}

// PUT
func (m *MockDigitalDbAdapter) UpdateDigitalLocation(
	ctx context.Context,
	userID string,
	digitalLocation models.DigitalLocation,
) error {
	return m.UpdateDigitalLocationFunc(ctx, userID, digitalLocation)
}

// DELETE
func (m *MockDigitalDbAdapter) DeleteDigitalLocation(
	ctx context.Context,
	userID string,
	locationIDs []string,
) (int64, error) {
	if m.RemoveDigitalLocationFunc != nil {
		return m.RemoveDigitalLocationFunc(ctx, userID, locationIDs)
	}
	return 0, nil
}

// Find by name
func (m *MockDigitalDbAdapter) FindDigitalLocationByName(
	ctx context.Context,
	userID string,
	name string,
) (models.DigitalLocation, error) {
	return m.FindDigitalLocationByNameFunc(ctx, userID, name)
}




