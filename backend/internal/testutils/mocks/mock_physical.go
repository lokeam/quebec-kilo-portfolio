package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type MockPhysicalValidator struct {
	ValidatePhysicalLocationFunc func(location models.PhysicalLocation) (models.PhysicalLocation, error)
}

func (m *MockPhysicalValidator) ValidatePhysicalLocation(location models.PhysicalLocation) (models.PhysicalLocation, error) {
	if m.ValidatePhysicalLocationFunc != nil {
		return m.ValidatePhysicalLocationFunc(location)
	}
	return location, nil
}

type MockPhysicalDbAdapter struct {
	GetPhysicalLocationsFunc    func(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
	GetPhysicalLocationFunc     func(ctx context.Context, userID, locationID string) (*models.PhysicalLocation, error)
	CreatePhysicalLocationFunc  func(ctx context.Context, userID string, location models.PhysicalLocation) error
	UpdatePhysicalLocationFunc  func(ctx context.Context, userID string, location models.PhysicalLocation) error
	DeletePhysicalLocationFunc  func(ctx context.Context, userID, locationID string) error
}

func (m *MockPhysicalDbAdapter) AddPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	err := m.CreatePhysicalLocationFunc(ctx, userID, location)
	return location, err
}

func (m *MockPhysicalDbAdapter) GetUserPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
	return m.GetPhysicalLocationsFunc(ctx, userID)
}

func (m *MockPhysicalDbAdapter) RemovePhysicalLocation(ctx context.Context, userID, locationID string) error {
	return m.DeletePhysicalLocationFunc(ctx, userID, locationID)
}

func (m *MockPhysicalDbAdapter) GetPhysicalLocation(ctx context.Context, userID, locationID string) (models.PhysicalLocation, error) {
	location, err := m.GetPhysicalLocationFunc(ctx, userID, locationID)
	if err != nil {
		return models.PhysicalLocation{}, err
	}
	return *location, nil
}

func (m *MockPhysicalDbAdapter) UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) error {
	return m.UpdatePhysicalLocationFunc(ctx, userID, location)
}

type MockPhysicalCacheWrapper struct {
	GetCachedPhysicalLocationsFunc      func(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
	SetCachedPhysicalLocationsFunc      func(ctx context.Context, userID string, locations []models.PhysicalLocation) error
	GetSingleCachedPhysicalLocationFunc func(ctx context.Context, userID, locationID string) (*models.PhysicalLocation, bool, error)
	SetSingleCachedPhysicalLocationFunc func(ctx context.Context, userID string, location models.PhysicalLocation) error
	InvalidateUserCacheFunc             func(ctx context.Context, userID string) error
	InvalidateLocationCacheFunc         func(ctx context.Context, userID, locationID string) error
}

func (m *MockPhysicalCacheWrapper) GetCachedPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
	return m.GetCachedPhysicalLocationsFunc(ctx, userID)
}

func (m *MockPhysicalCacheWrapper) SetCachedPhysicalLocations(ctx context.Context, userID string, locations []models.PhysicalLocation) error {
	return m.SetCachedPhysicalLocationsFunc(ctx, userID, locations)
}

func (m *MockPhysicalCacheWrapper) GetSingleCachedPhysicalLocation(ctx context.Context, userID, locationID string) (*models.PhysicalLocation, bool, error) {
	return m.GetSingleCachedPhysicalLocationFunc(ctx, userID, locationID)
}

func (m *MockPhysicalCacheWrapper) SetSingleCachedPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) error {
	return m.SetSingleCachedPhysicalLocationFunc(ctx, userID, location)
}

func (m *MockPhysicalCacheWrapper) InvalidateUserCache(ctx context.Context, userID string) error {
	return m.InvalidateUserCacheFunc(ctx, userID)
}

func (m *MockPhysicalCacheWrapper) InvalidateLocationCache(ctx context.Context, userID, locationID string) error {
	return m.InvalidateLocationCacheFunc(ctx, userID, locationID)
}