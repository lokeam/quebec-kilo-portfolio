package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/lokeam/qko-beta/internal/models"
)

type MockPhysicalDbAdapter struct {
	mock.Mock
}

func (m *MockPhysicalDbAdapter) GetSinglePhysicalLocation(ctx context.Context, userID string, locationID string) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, locationID)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) GetAllPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) CreatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, location)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, location)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) DeletePhysicalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error) {
	args := m.Called(ctx, userID, locationIDs)
	return args.Get(0).(int64), args.Error(1)
}

func DefaultPhysicalDbAdapter() *MockPhysicalDbAdapter {
	m := &MockPhysicalDbAdapter{}
	m.On("GetSinglePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(models.PhysicalLocation{}, nil)
	m.On("GetAllPhysicalLocations", mock.Anything, mock.Anything).
		Return([]models.PhysicalLocation{}, nil)
	m.On("CreatePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(models.PhysicalLocation{}, nil)
	m.On("UpdatePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(models.PhysicalLocation{}, nil)
	m.On("DeletePhysicalLocation", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	return m
}