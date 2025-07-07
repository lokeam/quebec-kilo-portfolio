package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)


type MockSublocationDbAdapter struct {
	GetSublocationFunc func(ctx context.Context, userID string, sublocationID string) (models.Sublocation, error)
	GetUserSublocationsFunc func(ctx context.Context, userID string) ([]models.Sublocation, error)
	AddSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error)
	UpdateSublocationFunc func(ctx context.Context, userID string, sublocation models.Sublocation) error
	DeleteSublocationFunc func(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error)
	CheckGameInAnySublocationFunc func(ctx context.Context, userGameID string) (bool, error)
	CheckGameInSublocationFunc func(ctx context.Context, userGameID string, sublocationID string) (bool, error)
	CheckGameOwnershipFunc func(ctx context.Context, userID string, userGameID string) (bool, error)
	MoveGameToSublocationFunc func(ctx context.Context, userID string, userGameID string, targetSublocationID string) error
	RemoveGameFromSublocationFunc func(ctx context.Context, userID string, userGameID string) error
}

func (m *MockSublocationDbAdapter) CheckGameInAnySublocation(
	ctx context.Context,
	userGameID string,
) (bool, error) {
	if m.CheckGameInAnySublocationFunc != nil {
		return m.CheckGameInAnySublocationFunc(ctx, userGameID)
	}
	return false, nil
}

func (m *MockSublocationDbAdapter) CheckGameInSublocation(
	ctx context.Context,
	userGameID string,
	sublocationID string,
) (bool, error) {
	if m.CheckGameInSublocationFunc != nil {
		return m.CheckGameInSublocationFunc(ctx, userGameID, sublocationID)
	}
	return false, nil
}

func (m *MockSublocationDbAdapter) CheckGameOwnership(
	ctx context.Context,
	userID string,
	userGameID string,
) (bool, error) {
	if m.CheckGameOwnershipFunc != nil {
		return m.CheckGameOwnershipFunc(ctx, userID, userGameID)
	}
	return true, nil
}

func (m *MockSublocationDbAdapter) MoveGameToSublocation(
	ctx context.Context,
	userID string,
	userGameID string,
	targetSublocationID string,
) error {
	if m.MoveGameToSublocationFunc != nil {
		return m.MoveGameToSublocationFunc(ctx, userID, userGameID, targetSublocationID)
	}
	return nil
}

func (m *MockSublocationDbAdapter) RemoveGameFromSublocation(
	ctx context.Context,
	userID string,
	userGameID string,
) error {
	if m.RemoveGameFromSublocationFunc != nil {
		return m.RemoveGameFromSublocationFunc(ctx, userID, userGameID)
	}
	return nil
}

func (m *MockSublocationDbAdapter) GetSublocation(
	ctx context.Context,
	userID string,
	sublocationID string,
) (models.Sublocation, error) {
	if m.GetSublocationFunc != nil {
		return m.GetSublocationFunc(ctx, userID, sublocationID)
	}
	return models.Sublocation{}, nil
}

func (m *MockSublocationDbAdapter) GetAllSublocations(
	ctx context.Context,
	userID string,
) ([]models.Sublocation, error) {
	if m.GetUserSublocationsFunc != nil {
		return m.GetUserSublocationsFunc(ctx, userID)
	}
	return []models.Sublocation{}, nil
}

func (m *MockSublocationDbAdapter) CreateSublocation(
	ctx context.Context,
	userID string,
	sublocation models.Sublocation,
) (models.Sublocation, error) {
	if m.AddSublocationFunc != nil {
		return m.AddSublocationFunc(ctx, userID, sublocation)
	}
	return models.Sublocation{}, nil
}

func (m *MockSublocationDbAdapter) UpdateSublocation(
	ctx context.Context,
	userID string,
	sublocation models.Sublocation,
) error {
	if m.UpdateSublocationFunc != nil {
		return m.UpdateSublocationFunc(ctx, userID, sublocation)
	}
	return nil
}

func (m *MockSublocationDbAdapter) DeleteSublocation(
	ctx context.Context,
	userID string,
	sublocationIDs []string,
) (types.DeleteSublocationResponse, error) {
	if m.DeleteSublocationFunc != nil {
		return m.DeleteSublocationFunc(ctx, userID, sublocationIDs)
	}
	return types.DeleteSublocationResponse{
		Success: true,
		DeletedCount: len(sublocationIDs),
		SublocationIDs: sublocationIDs,
	}, nil
}

func (m *MockSublocationDbAdapter) CheckDuplicateSublocation(
	ctx context.Context,
	userID string,
	physicalLocationID string,
	name string,
) (bool, error) {
	// If you want to customize, add a CheckDuplicateSublocationFunc field and use it here
	return false, nil
}

func (m *MockSublocationDbAdapter) GetSingleSublocation(
	ctx context.Context,
	userID string,
	sublocationID string,
) (models.Sublocation, error) {
	// If you want to customize, add a GetSingleSublocationFunc field and use it here
	return models.Sublocation{}, nil
}
