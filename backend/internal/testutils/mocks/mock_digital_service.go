package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/models"
)

type MockGameDigitalService struct {
	GetUserDigitalLocationsFunc func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetDigitalLocationFunc      func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error)
	AddDigitalLocationFunc      func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocationFunc   func(ctx context.Context, userID string, location models.DigitalLocation) error
	RemoveDigitalLocationFunc   func(ctx context.Context, userID, locationID string) error
}

func DefaultGameDigitalService() *MockGameDigitalService {
	return &MockGameDigitalService{
		GetUserDigitalLocationsFunc: func(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
			return []models.DigitalLocation{}, nil
		},
		GetDigitalLocationFunc: func(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
			return models.DigitalLocation{}, nil
		},
		AddDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
			// Set fields that would normally be set by the database
			if location.ID == "" {
				location.ID = uuid.New().String()
			}
			location.UserID = userID
			location.CreatedAt = time.Now()
			location.UpdatedAt = time.Now()

			// If subscription is provided, set its fields
			if location.Subscription != nil {
				location.Subscription.LocationID = location.ID
				location.Subscription.CreatedAt = time.Now()
				location.Subscription.UpdatedAt = time.Now()
			}

			return location, nil
		},
		UpdateDigitalLocationFunc: func(ctx context.Context, userID string, location models.DigitalLocation) error {
			return nil
		},
		RemoveDigitalLocationFunc: func(ctx context.Context, userID, locationID string) error {
			return nil
		},
	}
}

func (m *MockGameDigitalService) GetAllDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
	return m.GetUserDigitalLocationsFunc(ctx, userID)
}

func (m *MockGameDigitalService) GetSingleDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
	return m.GetDigitalLocationFunc(ctx, userID, locationID)
}

func (m *MockGameDigitalService) CreateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
	return m.AddDigitalLocationFunc(ctx, userID, location)
}

func (m *MockGameDigitalService) UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error {
	return m.UpdateDigitalLocationFunc(ctx, userID, location)
}

func (m *MockGameDigitalService) DeleteDigitalLocation(ctx context.Context, userID, locationID string) error {
	return m.RemoveDigitalLocationFunc(ctx, userID, locationID)
}