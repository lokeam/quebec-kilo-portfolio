package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type MockSublocationValidator struct {
	ValidateSublocationFunc func(sublocation models.Sublocation) (models.Sublocation, error)
}

func (m *MockSublocationValidator) ValidateSublocation(sublocation models.Sublocation) (models.Sublocation, error) {
	if m.ValidateSublocationFunc != nil {
		return m.ValidateSublocationFunc(sublocation)
	}
	return sublocation, nil
}

type MockSublocationDbAdapter struct {
	GetSublocationFunc func(ctx context.Context, userID, sublocationID string) (models.Sublocation, error)
	GetSublocationsFunc func(ctx context.Context, userID string) ([]models.Sublocation, error)
	AddSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error)
	UpdateSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) error
	DeleteSublocationFunc func(ctx context.Context, userID string, sublocationID string) error
}


// GET
func (m *MockSublocationDbAdapter) GetSublocation(
	ctx context.Context,
	userID,
	sublocationID string,
) (models.Sublocation, error) {
	sublocation, err := m.GetSublocationFunc(ctx, userID, sublocationID)
	return sublocation, err
}

func (m *MockSublocationDbAdapter) GetUserSublocations(
	ctx context.Context,
	userID string,
) ([]models.Sublocation, error) {
	return m.GetSublocationsFunc(ctx, userID)
}

// POST
func (m *MockSublocationDbAdapter) AddSublocation(
	ctx context.Context,
	userID string,
	sublocation models.Sublocation,
) (models.Sublocation, error) {
	return m.AddSublocationFunc(ctx, userID, sublocation)
}

// PUT
func (m *MockSublocationDbAdapter) UpdateSublocation(
	ctx context.Context,
	userID string,
	sublocation models.Sublocation,
) error {
	return m.UpdateSublocationFunc(ctx, userID, sublocation)
}

// DELETE
func (m *MockSublocationDbAdapter) RemoveSublocation(
	ctx context.Context,
	userID string,
	sublocationID string,
) error {
	return m.DeleteSublocationFunc(ctx, userID, sublocationID)
}

type MockSublocationCacheWrapper struct {
	GetCachedSublocationsFunc           func(ctx context.Context, userID string) ([]models.Sublocation, error)
	SetCachedSublocationsFunc           func(ctx context.Context, userID string, locations []models.Sublocation) error
	GetSingleCachedSublocationFunc      func(ctx context.Context, userID, sublocationID string) (*models.Sublocation, bool, error)
	SetSingleCachedSublocationFunc      func(ctx context.Context, userID string, location models.Sublocation) error
	InvalidateUserCacheFunc             func(ctx context.Context, userID string) error
	InvalidateSublocationCacheFunc      func(ctx context.Context, userID, sublocationID string) error
}

// GET
func (m *MockSublocationCacheWrapper) GetCachedSublocations(
	ctx context.Context,
	userID string,
) ([]models.Sublocation, error) {
	return m.GetCachedSublocationsFunc(ctx, userID)
}

func (m *MockSublocationCacheWrapper) GetSingleCachedSublocation(
	ctx context.Context,
	userID,
	sublocationID string,
) (*models.Sublocation, bool, error) {
	return m.GetSingleCachedSublocationFunc(ctx, userID, sublocationID)
}

// SET
func (m *MockSublocationCacheWrapper) SetCachedSublocations(
	ctx context.Context,
	userID string,
	sublocations []models.Sublocation,
) error {
	return m.SetCachedSublocationsFunc(ctx, userID, sublocations)
}

func (m *MockSublocationCacheWrapper) SetSingleCachedSublocation(
	ctx context.Context,
	userID string,
	sublocation models.Sublocation,
) error {
	return m.SetSingleCachedSublocationFunc(ctx, userID, sublocation)
}

// CLEAR
func (m *MockSublocationCacheWrapper) InvalidateUserCache(
	ctx context.Context,
	userID string,
) error {
	return m.InvalidateUserCacheFunc(ctx, userID)
}

func (m *MockSublocationCacheWrapper) InvalidateSublocationCache(
	ctx context.Context,
	userID,
	sublocationID string,
) error {
	return m.InvalidateSublocationCacheFunc(ctx, userID, sublocationID)
}
