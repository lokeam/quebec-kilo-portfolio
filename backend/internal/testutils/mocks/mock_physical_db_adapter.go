package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/stretchr/testify/mock"
)


type MockPhysicalDbAdapter struct {
	mock.Mock
}

func (m *MockPhysicalDbAdapter) GetAllPhysicalLocations(
	ctx context.Context,
	userID string,
) ([]models.PhysicalLocation, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) GetAllPhysicalLocationsBFF(
	ctx context.Context,
	userID string,
) (types.LocationsBFFResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(types.LocationsBFFResponse), args.Error(1)
}

func (m *MockPhysicalDbAdapter) GetSinglePhysicalLocation(
	ctx context.Context,
	userID string,
	locationID string,
) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, locationID)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) CreatePhysicalLocation(
	ctx context.Context,
	userID string,
	location models.PhysicalLocation,
) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, location)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) UpdatePhysicalLocation(
	ctx context.Context,
	userID string,
	location models.PhysicalLocation,
) (models.PhysicalLocation, error) {
	args := m.Called(ctx, userID, location)
	return args.Get(0).(models.PhysicalLocation), args.Error(1)
}

func (m *MockPhysicalDbAdapter) DeletePhysicalLocation(
	ctx context.Context,
	userID string,
	locationIDs []string,
) (int64, error) {
	args := m.Called(ctx, userID, locationIDs)
	return args.Get(0).(int64), args.Error(1)
}