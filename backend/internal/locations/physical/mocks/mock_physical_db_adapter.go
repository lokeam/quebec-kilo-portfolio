package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/lokeam/qko-beta/internal/models"
)

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

func DefaultPhysicalDbAdapter() *MockPhysicalDbAdapter {
	m := &MockPhysicalDbAdapter{}
	m.On("GetPhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(models.PhysicalLocation{}, nil)
	m.On("GetUserPhysicalLocations", mock.Anything, mock.Anything).
		Return([]models.PhysicalLocation{}, nil)
	m.On("AddPhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(models.PhysicalLocation{}, nil)
	m.On("UpdatePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(models.PhysicalLocation{}, nil)
	m.On("RemovePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	return m
}