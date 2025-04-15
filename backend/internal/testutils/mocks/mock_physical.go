package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/stretchr/testify/mock"
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
	mock.Mock
}

func (m *MockPhysicalDbAdapter) GetPhysicalLocation(ctx context.Context, userID string, locationID string) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, locationID)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) GetUserPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) AddPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, location)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, location)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) RemovePhysicalLocation(ctx context.Context, userID string, locationID string) error {
	args := m.Called(ctx, userID, locationID)
	return args.Error(0)
}

type MockPhysicalCacheWrapper struct {
	mock.Mock
}

// GET
func (m *MockPhysicalCacheWrapper) GetCachedPhysicalLocations(
	ctx context.Context,
	userID string,
) ([]models.PhysicalLocation, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalCacheWrapper) GetSingleCachedPhysicalLocation(
	ctx context.Context,
	userID,
	locationID string,
) (*models.PhysicalLocation, bool, error) {
	args := m.Called(ctx, userID, locationID)
	if args.Get(0) == nil {
		return nil, args.Bool(1), args.Error(2)
	}
	return args.Get(0).(*models.PhysicalLocation), args.Bool(1), args.Error(2)
}

// SET
func (m *MockPhysicalCacheWrapper) SetCachedPhysicalLocations(
	ctx context.Context,
	userID string,
	locations []models.PhysicalLocation,
) error {
	args := m.Called(ctx, userID, locations)
	return args.Error(0)
}

func (m *MockPhysicalCacheWrapper) SetSingleCachedPhysicalLocation(
	ctx context.Context,
	userID string,
	location models.PhysicalLocation,
) error {
	args := m.Called(ctx, userID, location)
	return args.Error(0)
}

// CLEAR
func (m *MockPhysicalCacheWrapper) InvalidateUserCache(
	ctx context.Context,
	userID string,
) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPhysicalCacheWrapper) InvalidateLocationCache(
	ctx context.Context,
	userID,
	locationID string,
) error {
	args := m.Called(ctx, userID, locationID)
	return args.Error(0)
}