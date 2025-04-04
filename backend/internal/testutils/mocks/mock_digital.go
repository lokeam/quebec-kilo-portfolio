package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type MockDigitalValidator struct {
	ValidateDigitalLocationFunc func(digitalLocation models.DigitalLocation) (models.DigitalLocation, error)
}

func (m *MockDigitalValidator) ValidateDigitalLocation(
	digitalLocation models.DigitalLocation,
) (models.DigitalLocation, error) {
	if m.ValidateDigitalLocationFunc != nil {
		return m.ValidateDigitalLocationFunc(digitalLocation)
	}
	return digitalLocation, nil
}

type MockDigitalDbAdapter struct {
	GetDigitalLocationFunc func(ctx context.Context, userID, digitalLocationID string) (models.DigitalLocation, error)
	GetDigitalLocationsFunc func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	AddDigitalLocationFunc func(ctx context.Context, userID string, digitalLocation models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocationFunc func(ctx context.Context, userID string, digitalLocation models.DigitalLocation) error
	DeleteDigitalLocationFunc func(ctx context.Context, userID, digitalLocationID string) error
}

// GET
func (m *MockDigitalDbAdapter) GetDigitalLocation(
	ctx context.Context,
	userID,
	digitalLocationID string,
) (models.DigitalLocation, error) {
	digitalLocation, err := m.GetDigitalLocationFunc(ctx, userID, digitalLocationID)
	return digitalLocation, err
}

func (m *MockDigitalDbAdapter) GetUserDigitalLocations(
	ctx context.Context,
	userID string,
) ([]models.DigitalLocation, error) {
	return m.GetDigitalLocationsFunc(ctx, userID)
}

// POST
func (m *MockDigitalDbAdapter) AddDigitalLocation(
	ctx context.Context,
	userID string,
	digitalLocation models.DigitalLocation,
) (models.DigitalLocation, error) {
	return m.AddDigitalLocationFunc(ctx, userID, digitalLocation)
}

// PUT
func (m *MockDigitalDbAdapter) UpdateDigitalLocation(
	ctx context.Context,
	userID string,
	digitalLocation models.DigitalLocation,
) error {
	return m.UpdateDigitalLocationFunc(ctx, userID, digitalLocation)
}

// DELETE
func (m *MockDigitalDbAdapter) RemoveDigitalLocation(
	ctx context.Context,
	userID string,
	digitalLocationID string,
) error {
	return m.DeleteDigitalLocationFunc(ctx, userID, digitalLocationID)
}

// ---------

type MockDigitalCacheWrapper struct {
	GetCachedDigitalLocationsFunc           func(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	SetCachedDigitalLocationsFunc           func(ctx context.Context, userID string, locations []models.DigitalLocation) error
	GetSingleCachedDigitalLocationFunc      func(ctx context.Context, userID, digitalLocationID string) (*models.DigitalLocation, bool, error)
	SetSingleCachedDigitalLocationFunc      func(ctx context.Context, userID string, location models.DigitalLocation) error
	InvalidateUserCacheFunc             func(ctx context.Context, userID string) error
	InvalidateDigitalLocationCacheFunc      func(ctx context.Context, userID, digitalLocationID string) error
}

// GET
func (m *MockDigitalCacheWrapper) GetCachedDigitalLocations(
	ctx context.Context,
	userID string,
) ([]models.DigitalLocation, error) {
	return m.GetCachedDigitalLocationsFunc(ctx, userID)
}

func (m *MockDigitalCacheWrapper) GetSingleCachedDigitalLocation(
	ctx context.Context,
	userID,
	digitalLocationID string,
) (*models.DigitalLocation, bool, error) {
	return m.GetSingleCachedDigitalLocationFunc(ctx, userID, digitalLocationID)
}

// SET
func (m *MockDigitalCacheWrapper) SetCachedDigitalLocations(
	ctx context.Context,
	userID string,
	digitalLocations []models.DigitalLocation,
) error {
	return m.SetCachedDigitalLocationsFunc(ctx, userID, digitalLocations)
}

func (m *MockDigitalCacheWrapper) SetSingleCachedDigitalLocation(
	ctx context.Context,
	userID string,
	digitalLocation models.DigitalLocation,
) error {
	return m.SetSingleCachedDigitalLocationFunc(ctx, userID, digitalLocation)
}

// CLEAR
func (m *MockDigitalCacheWrapper) InvalidateUserCache(
	ctx context.Context,
	userID string,
) error {
	return m.InvalidateUserCacheFunc(ctx, userID)
}

func (m *MockDigitalCacheWrapper) InvalidateDigitalLocationCache(
	ctx context.Context,
	userID,
	digitalLocationID string,
) error {
	return m.InvalidateDigitalLocationCacheFunc(ctx, userID, digitalLocationID)
}
