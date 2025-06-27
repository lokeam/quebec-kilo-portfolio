package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

// MockDigitalService implements services.DigitalService
type MockDigitalService struct {
	// Read ops for digital locations
	GetAllDigitalLocationsBFFFunc func(ctx context.Context, userID string) (types.DigitalLocationsBFFResponse, error)
	GetUserDigitalLocationsFunc    func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetDigitalLocationFunc         func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error)
	FindDigitalLocationByNameFunc  func(ctx context.Context, userID string, name string) (models.DigitalLocation, error)

	// Write ops for digital locations
	AddDigitalLocationFunc         func(ctx context.Context, userID string, location types.DigitalLocationRequest) (models.DigitalLocation, error)
	UpdateDigitalLocationFunc      func(ctx context.Context, userID string, location types.DigitalLocationRequest) error
	RemoveDigitalLocationFunc      func(ctx context.Context, userID string, locationIDs []string) (int64, error)

	// Games
	AddGameToDigitalLocationFunc   func(ctx context.Context, userID string, locationID string, gameID int64) error
	RemoveGameFromDigitalLocationFunc func(ctx context.Context, userID string, locationID string, gameID int64) error
	GetGamesByDigitalLocationIDFunc func(ctx context.Context, userID string, locationID string) ([]models.Game, error)

	// Subscriptions
	GetSubscriptionFunc           func(ctx context.Context, locationID string) (*models.Subscription, error)
	AddSubscriptionFunc           func(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscriptionFunc        func(ctx context.Context, subscription models.Subscription) error
	RemoveSubscriptionFunc        func(ctx context.Context, locationID string) error

	// Payments
	GetPaymentsFunc              func(ctx context.Context, locationID string) ([]models.Payment, error)
	AddPaymentFunc               func(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetPaymentFunc               func(ctx context.Context, paymentID int64) (*models.Payment, error)
}

// DefaultGameDigitalService creates a MockDigitalService with sensible defaults for testing
func DefaultGameDigitalService() *MockDigitalService {
	return &MockDigitalService{
		GetAllDigitalLocationsBFFFunc: func(
			ctx context.Context,
			userID string,
		) (types.DigitalLocationsBFFResponse, error) {
			return types.DigitalLocationsBFFResponse{
				DigitalLocations: []types.SingleDigitalLocationBFFResponse{
					{ID: "1", Name: "Location 1", URL: "http://example.com/1", IsSubscription: false, IsActive: true},
					{ID: "2", Name: "Location 2", URL: "http://example.com/2", IsSubscription: true, IsActive: true},
				},
			}, nil
		},
		GetUserDigitalLocationsFunc: func(
			ctx context.Context,
			userID string,
		) ([]models.DigitalLocation, error) {
			return []models.DigitalLocation{
				{ID: "1", Name: "Location 1", URL: "http://example.com/1", IsSubscription: false, IsActive: true},
				{ID: "2", Name: "Location 2", URL: "http://example.com/2", IsSubscription: true, IsActive: true},
			}, nil
		},
		GetDigitalLocationFunc: func(
			ctx context.Context,
			userID,
			locationID string,
		) (models.DigitalLocation, error) {
			return models.DigitalLocation{
				ID: locationID,
				Name: "Test Location",
				URL: "http://example.com",
				IsSubscription: false,
				IsActive: true,
			},
			nil
		},
		FindDigitalLocationByNameFunc: func(
			ctx context.Context,
			userID string,
			name string,
		) (models.DigitalLocation, error) {
			return models.DigitalLocation{
				ID: "test-id",
				Name: name,
				URL: "http://example.com",
				IsSubscription: false,
				IsActive: true,
			}, nil
		},
		AddDigitalLocationFunc: func(
			ctx context.Context,
			userID string, location types.DigitalLocationRequest,
		) (models.DigitalLocation, error) {
			// Convert request to model and set fields that would normally be set by the database
			digitalLocation := models.DigitalLocation{
				ID:             uuid.New().String(),
				Name:           location.Name,
				URL:            location.URL,
				IsSubscription: location.IsSubscription,
				IsActive:       true,
				UserID:         userID,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			// If subscription is provided, set its fields
			if location.Subscription != nil {
				digitalLocation.Subscription = &models.Subscription{
					LocationID: digitalLocation.ID,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
			}

			return digitalLocation, nil
		},
		UpdateDigitalLocationFunc: func(
			ctx context.Context,
			userID string,
			location types.DigitalLocationRequest,
		) error {
			return nil
		},
		RemoveDigitalLocationFunc: func(
			ctx context.Context,
			userID string,
			locationIDs []string,
		) (int64, error) {
			return int64(len(locationIDs)), nil
		},
		AddGameToDigitalLocationFunc: func(
			ctx context.Context,
			userID string,
			locationID string,
			gameID int64,
		) error {
			return nil
		},
		RemoveGameFromDigitalLocationFunc: func(
			ctx context.Context,
			userID string, locationID string, gameID int64) error {
			return nil
		},
		GetGamesByDigitalLocationIDFunc: func(
			ctx context.Context,
			userID string,
			locationID string,
		) ([]models.Game, error) {
			return []models.Game{{ID: 1, Name: "Test Game"}}, nil
		},
		GetSubscriptionFunc: func(
			ctx context.Context,
			locationID string,
		) (*models.Subscription, error) {
			return &models.Subscription{ID: 1, LocationID: locationID}, nil
		},
		AddSubscriptionFunc: func(
			ctx context.Context,
			subscription models.Subscription,
		) (*models.Subscription, error) {
			return &models.Subscription{ID: 1, LocationID: subscription.LocationID}, nil
		},
		UpdateSubscriptionFunc: func(
			ctx context.Context,
			subscription models.Subscription,
		) error {
			return nil
		},
		RemoveSubscriptionFunc: func(
			ctx context.Context,
			locationID string,
		) error {
			return nil
		},
		GetPaymentsFunc: func(
			ctx context.Context,
			locationID string,
		) ([]models.Payment, error) {
			return []models.Payment{{ID: 1, LocationID: locationID}}, nil
		},
		AddPaymentFunc: func(
			ctx context.Context,
			payment models.Payment,
		) (*models.Payment, error) {
			return &models.Payment{ID: 1, LocationID: payment.LocationID}, nil
		},
		GetPaymentFunc: func(
			ctx context.Context,
			paymentID int64,
		) (*models.Payment, error) {
			return &models.Payment{ID: paymentID}, nil
		},
	}
}

// Interface implementation methods
func (m *MockDigitalService) GetAllDigitalLocationsBFF(
	ctx context.Context,
	userID string,
) (types.DigitalLocationsBFFResponse, error) {
	if m.GetAllDigitalLocationsBFFFunc != nil {
		return m.GetAllDigitalLocationsBFFFunc(ctx, userID)
	}
	return types.DigitalLocationsBFFResponse{}, nil
}

func (m *MockDigitalService) GetAllDigitalLocations(
	ctx context.Context,
	userID string,
) ([]models.DigitalLocation, error) {
	if m.GetUserDigitalLocationsFunc != nil {
		return m.GetUserDigitalLocationsFunc(ctx, userID)
	}
	return []models.DigitalLocation{}, nil
}

func (m *MockDigitalService) GetSingleDigitalLocation(
	ctx context.Context,
	userID,
	locationID string,
) (models.DigitalLocation, error) {
	if m.GetDigitalLocationFunc != nil {
		return m.GetDigitalLocationFunc(ctx, userID, locationID)
	}
	return models.DigitalLocation{}, nil
}

func (m *MockDigitalService) FindDigitalLocationByName(
	ctx context.Context,
	userID string,
	name string,
) (models.DigitalLocation, error) {
	if m.FindDigitalLocationByNameFunc != nil {
		return m.FindDigitalLocationByNameFunc(ctx, userID, name)
	}
	return models.DigitalLocation{}, nil
}

func (m *MockDigitalService) CreateDigitalLocation(
	ctx context.Context,
	userID string,
	location types.DigitalLocationRequest,
) (models.DigitalLocation, error) {
	if m.AddDigitalLocationFunc != nil {
		return m.AddDigitalLocationFunc(ctx, userID, location)
	}
	return models.DigitalLocation{}, nil
}

func (m *MockDigitalService) UpdateDigitalLocation(
	ctx context.Context,
	userID string,
	location types.DigitalLocationRequest,
) error {
	if m.UpdateDigitalLocationFunc != nil {
		return m.UpdateDigitalLocationFunc(ctx, userID, location)
	}
	return nil
}

func (m *MockDigitalService) DeleteDigitalLocation(
	ctx context.Context,
	userID string,
	locationIDs []string,
) (int64, error) {
	if m.RemoveDigitalLocationFunc != nil {
		return m.RemoveDigitalLocationFunc(ctx, userID, locationIDs)
	}
	return 0, nil
}

func (m *MockDigitalService) AddGameToDigitalLocation(
	ctx context.Context,
	userID string,
	locationID string,
	gameID int64,
) error {
	if m.AddGameToDigitalLocationFunc != nil {
		return m.AddGameToDigitalLocationFunc(ctx, userID, locationID, gameID)
	}
	return nil
}

func (m *MockDigitalService) RemoveGameFromDigitalLocation(
	ctx context.Context,
	userID string,
	locationID string,
	gameID int64,
) error {
	if m.RemoveGameFromDigitalLocationFunc != nil {
		return m.RemoveGameFromDigitalLocationFunc(ctx, userID, locationID, gameID)
	}
	return nil
}

func (m *MockDigitalService) GetGamesByDigitalLocationID(
	ctx context.Context,
	userID string,
	locationID string,
) ([]models.Game, error) {
	if m.GetGamesByDigitalLocationIDFunc != nil {
		return m.GetGamesByDigitalLocationIDFunc(ctx, userID, locationID)
	}
	return []models.Game{}, nil
}

func (m *MockDigitalService) GetSubscription(
	ctx context.Context,
	locationID string,
) (*models.Subscription, error) {
	if m.GetSubscriptionFunc != nil {
		return m.GetSubscriptionFunc(ctx, locationID)
	}
	return &models.Subscription{}, nil
}

func (m *MockDigitalService) CreateSubscription(
	ctx context.Context,
	subscription models.Subscription,
) (*models.Subscription, error) {
	if m.AddSubscriptionFunc != nil {
		return m.AddSubscriptionFunc(ctx, subscription)
	}
	return &models.Subscription{}, nil
}

func (m *MockDigitalService) UpdateSubscription(
	ctx context.Context,
	subscription models.Subscription,
) error {
	if m.UpdateSubscriptionFunc != nil {
		return m.UpdateSubscriptionFunc(ctx, subscription)
	}
	return nil
}

func (m *MockDigitalService) DeleteSubscription(
	ctx context.Context,
	locationID string,
) error {
	if m.RemoveSubscriptionFunc != nil {
		return m.RemoveSubscriptionFunc(ctx, locationID)
	}
	return nil
}

func (m *MockDigitalService) GetAllPayments(
	ctx context.Context,
	locationID string,
) ([]models.Payment, error) {
	if m.GetPaymentsFunc != nil {
		return m.GetPaymentsFunc(ctx, locationID)
	}
	return []models.Payment{}, nil
}

func (m *MockDigitalService) CreatePayment(
	ctx context.Context,
	payment models.Payment,
) (*models.Payment, error) {
	if m.AddPaymentFunc != nil {
		return m.AddPaymentFunc(ctx, payment)
	}
	return &models.Payment{}, nil
}

func (m *MockDigitalService) GetSinglePayment(
	ctx context.Context,
	paymentID int64,
) (*models.Payment, error) {
	if m.GetPaymentFunc != nil {
		return m.GetPaymentFunc(ctx, paymentID)
	}
	return &models.Payment{}, nil
}