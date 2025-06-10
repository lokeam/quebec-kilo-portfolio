package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

type MockSublocationValidator struct {
	ValidateSublocationFunc func(sublocation models.Sublocation) (models.Sublocation, error)
	ValidateSublocationUpdateFunc func(update, existing models.Sublocation) (models.Sublocation, error)
	ValidateSublocationCreationFunc func(sublocation models.Sublocation) (models.Sublocation, error)
	ValidateGameOwnershipFunc func(userID string, userGameID string) error
	ValidateSublocationOwnershipFunc func(userID string, sublocationID string) error
	ValidateGameNotInSublocationFunc func(userGameID string, sublocationID string) error
	ValidateDeleteSublocationRequestFunc func(userID string, sublocationIDs []string) error
}

func (m *MockSublocationValidator) ValidateSublocation(sublocation models.Sublocation) (models.Sublocation, error) {
	if m.ValidateSublocationFunc != nil {
		return m.ValidateSublocationFunc(sublocation)
	}
	return sublocation, nil
}

func (m *MockSublocationValidator) ValidateSublocationUpdate(update, existing models.Sublocation) (models.Sublocation, error) {
	if m.ValidateSublocationUpdateFunc != nil {
		return m.ValidateSublocationUpdateFunc(update, existing)
	}
	return update, nil
}

func (m *MockSublocationValidator) ValidateSublocationCreation(sublocation models.Sublocation) (models.Sublocation, error) {
	if m.ValidateSublocationCreationFunc != nil {
		return m.ValidateSublocationCreationFunc(sublocation)
	}
	return sublocation, nil
}

func (m *MockSublocationValidator) ValidateGameOwnership(userID string, userGameID string) error {
	if m.ValidateGameOwnershipFunc != nil {
		return m.ValidateGameOwnershipFunc(userID, userGameID)
	}
	return nil
}

func (m *MockSublocationValidator) ValidateSublocationOwnership(userID string, sublocationID string) error {
	if m.ValidateSublocationOwnershipFunc != nil {
		return m.ValidateSublocationOwnershipFunc(userID, sublocationID)
	}
	return nil
}

func (m *MockSublocationValidator) ValidateGameNotInSublocation(userGameID string, sublocationID string) error {
	if m.ValidateGameNotInSublocationFunc != nil {
		return m.ValidateGameNotInSublocationFunc(userGameID, sublocationID)
	}
	return nil
}

func (m *MockSublocationValidator) ValidateDeleteSublocationRequest(userID string, sublocationIDs []string) error {
	if m.ValidateDeleteSublocationRequestFunc != nil {
		return m.ValidateDeleteSublocationRequestFunc(userID, sublocationIDs)
	}
	return nil
}

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
