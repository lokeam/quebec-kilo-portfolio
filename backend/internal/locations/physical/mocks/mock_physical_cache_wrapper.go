package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/lokeam/qko-beta/internal/models"
)

type MockPhysicalCacheWrapper struct {
	mock.Mock
}

func (m *MockPhysicalCacheWrapper) GetSingleCachedPhysicalLocation(ctx context.Context, userID string, locationID string) (*models.PhysicalLocation, bool, error) {
	args := m.Called(ctx, userID, locationID)
	return args.Get(0).(*models.PhysicalLocation), args.Bool(1), args.Error(2)
}

func (m *MockPhysicalCacheWrapper) SetSingleCachedPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) error {
	args := m.Called(ctx, userID, location)
	return args.Error(0)
}

func (m *MockPhysicalCacheWrapper) GetCachedPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, bool, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.PhysicalLocation), args.Bool(1), args.Error(2)
}

func (m *MockPhysicalCacheWrapper) SetCachedPhysicalLocations(ctx context.Context, userID string, locations []models.PhysicalLocation) error {
	args := m.Called(ctx, userID, locations)
	return args.Error(0)
}

func DefaultPhysicalCacheWrapper() *MockPhysicalCacheWrapper {
	m := &MockPhysicalCacheWrapper{}
	m.On("GetSingleCachedPhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, false, nil)
	m.On("SetSingleCachedPhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	m.On("GetCachedPhysicalLocations", mock.Anything, mock.Anything).
		Return([]models.PhysicalLocation{}, false, nil)
	m.On("SetCachedPhysicalLocations", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	return m
}