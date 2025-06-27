package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type MockSublocationCacheWrapper struct {
	GetCachedSublocationsFunc           func(ctx context.Context, userID string) ([]models.Sublocation, error)
	SetCachedSublocationsFunc           func(ctx context.Context, userID string, locations []models.Sublocation) error
	GetSingleCachedSublocationFunc      func(ctx context.Context, userID, sublocationID string) (*models.Sublocation, bool, error)
	SetSingleCachedSublocationFunc      func(ctx context.Context, userID string, location models.Sublocation) error
	InvalidateUserCacheFunc             func(ctx context.Context, userID string) error
	InvalidateSublocationCacheFunc      func(ctx context.Context, userID, sublocationID string) error
	InvalidateLocationCacheFunc         func(ctx context.Context, userID, locationID string) error
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

func (m *MockSublocationCacheWrapper) InvalidateLocationCache(
	ctx context.Context,
	userID,
	locationID string,
) error {
	return m.InvalidateLocationCacheFunc(ctx, userID, locationID)
}
