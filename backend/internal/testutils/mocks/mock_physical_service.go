package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type MockPhysicalService struct {
	GetAllPhysicalLocationsBFFFunc  func(ctx context.Context, userID string) (types.LocationsBFFResponse, error)
	GetAllPhysicalLocationsFunc     func(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
	GetSinglePhysicalLocationFunc   func(ctx context.Context, userID, locationID string) (models.PhysicalLocation, error)
	CreatePhysicalLocationFunc      func(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	UpdatePhysicalLocationFunc      func(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	DeletePhysicalLocationFunc      func(ctx context.Context, userID string, locationIDs []string) (int64, error)
	InvalidateCacheFunc             func(ctx context.Context, cacheKey string) error
}

func (m *MockPhysicalService) GetAllPhysicalLocationsBFF(
	ctx context.Context,
	userID string,
) (types.LocationsBFFResponse, error) {
	if m.GetAllPhysicalLocationsBFFFunc != nil {
		return m.GetAllPhysicalLocationsBFFFunc(ctx, userID)
	}
	return types.LocationsBFFResponse{}, nil
}

func (m *MockPhysicalService) GetAllPhysicalLocations(
	ctx context.Context,
	userID string,
) ([]models.PhysicalLocation, error) {
	if m.GetAllPhysicalLocationsFunc != nil {
		return m.GetAllPhysicalLocationsFunc(ctx, userID)
	}
	return []models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) GetSinglePhysicalLocation(
	ctx context.Context,
	userID,
	locationID string,
) (models.PhysicalLocation, error) {
	if m.GetSinglePhysicalLocationFunc != nil {
		return m.GetSinglePhysicalLocationFunc(ctx, userID, locationID)
	}
	return models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) CreatePhysicalLocation(
	ctx context.Context,
	userID string,
	location models.PhysicalLocation,
) (models.PhysicalLocation, error) {
	if m.CreatePhysicalLocationFunc != nil {
		return m.CreatePhysicalLocationFunc(ctx, userID, location)
	}
	return models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) UpdatePhysicalLocation(
	ctx context.Context,
	userID string,
	location models.PhysicalLocation,
) (models.PhysicalLocation, error) {
	if m.UpdatePhysicalLocationFunc != nil {
		return m.UpdatePhysicalLocationFunc(ctx, userID, location)
	}
	return models.PhysicalLocation{}, nil
}

func (m *MockPhysicalService) DeletePhysicalLocation(
	ctx context.Context,
	userID string,
	locationIDs []string,
) (int64, error) {
	if m.DeletePhysicalLocationFunc != nil {
		return m.DeletePhysicalLocationFunc(ctx, userID, locationIDs)
	}
	return 0, nil
}

func (m *MockPhysicalService) InvalidateCache(
	ctx context.Context,
	cacheKey string,
) error {
	if m.InvalidateCacheFunc != nil {
		return m.InvalidateCacheFunc(ctx, cacheKey)
	}
	return nil
}
