package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/stretchr/testify/mock"
)


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

func (m *MockPhysicalCacheWrapper) GetCachedPhysicalLocationsBFF(
	ctx context.Context,
	userID string,
) (types.LocationsBFFResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(types.LocationsBFFResponse), args.Error(1)
}

func (m *MockPhysicalCacheWrapper) GetSingleCachedPhysicalLocation(
	ctx context.Context,
	userID string,
	locationID string,
) (*models.PhysicalLocation, bool, error) {
	args := m.Called(ctx, userID, locationID)
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

func (m *MockPhysicalCacheWrapper) SetCachedPhysicalLocationsBFF(
	ctx context.Context,
	userID string,
	response types.LocationsBFFResponse,
) error {
	args := m.Called(ctx, userID, response)
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
	userID string,
	locationID string,
) error {
	args := m.Called(ctx, userID, locationID)
	return args.Error(0)
}