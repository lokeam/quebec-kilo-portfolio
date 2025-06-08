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

type MockSublocationDbAdapter struct {
	GetSublocationFunc func(ctx context.Context, userID, sublocationID string) (models.Sublocation, error)
	GetSublocationsFunc func(ctx context.Context, userID string) ([]models.Sublocation, error)
	AddSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error)
	UpdateSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) error
	DeleteSublocationFunc func(ctx context.Context, userID string, sublocationID string) error
	CheckDuplicateSublocationFunc func(ctx context.Context, userID string, physicalLocationID string, name string) (bool, error)
	CheckGameOwnershipFunc func(ctx context.Context, userID string, userGameID string) (bool, error)
	CheckGameInSublocationFunc func(ctx context.Context, userGameID string, sublocationID string) (bool, error)
}

// GET
func (m *MockSublocationDbAdapter) GetAllSublocations(ctx context.Context, userID string) ([]models.Sublocation, error) {
	if m.GetSublocationsFunc != nil {
		return m.GetSublocationsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockSublocationDbAdapter) GetSingleSublocation(ctx context.Context, userID, sublocationID string) (models.Sublocation, error) {
	if m.GetSublocationFunc != nil {
		return m.GetSublocationFunc(ctx, userID, sublocationID)
	}
	return models.Sublocation{}, nil
}

// POST
func (m *MockSublocationDbAdapter) CreateSublocation(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error) {
	if m.AddSublocationFunc != nil {
		return m.AddSublocationFunc(ctx, userID, sublocation)
	}
	return sublocation, nil
}

// PUT
func (m *MockSublocationDbAdapter) UpdateSublocation(ctx context.Context, userID string, sublocation models.Sublocation) error {
	if m.UpdateSublocationFunc != nil {
		return m.UpdateSublocationFunc(ctx, userID, sublocation)
	}
	return nil
}

// DELETE
func (m *MockSublocationDbAdapter) DeleteSublocation(ctx context.Context, userID string, sublocationID string) error {
	if m.DeleteSublocationFunc != nil {
		return m.DeleteSublocationFunc(ctx, userID, sublocationID)
	}
	return nil
}

func (m *MockSublocationDbAdapter) CheckDuplicateSublocation(ctx context.Context, userID string, physicalLocationID string, name string) (bool, error) {
	if m.CheckDuplicateSublocationFunc != nil {
		return m.CheckDuplicateSublocationFunc(ctx, userID, physicalLocationID, name)
	}
	return false, nil
}

func (m *MockSublocationDbAdapter) CheckGameOwnership(ctx context.Context, userID string, userGameID string) (bool, error) {
	if m.CheckGameOwnershipFunc != nil {
		return m.CheckGameOwnershipFunc(ctx, userID, userGameID)
	}
	return true, nil
}

func (m *MockSublocationDbAdapter) CheckGameInSublocation(ctx context.Context, userGameID string, sublocationID string) (bool, error) {
	if m.CheckGameInSublocationFunc != nil {
		return m.CheckGameInSublocationFunc(ctx, userGameID, sublocationID)
	}
	return false, nil
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
	if m.InvalidateLocationCacheFunc != nil {
		return m.InvalidateLocationCacheFunc(ctx, userID, locationID)
	}
	return nil
}
