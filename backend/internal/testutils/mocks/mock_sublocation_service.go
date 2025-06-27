package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type MockSublocationService struct {
	GetSublocationsFunc     func(ctx context.Context, userID string) ([]models.Sublocation, error)
	GetSingleSublocationFunc func(ctx context.Context, userID, locationID string) (models.Sublocation, error)
	CreateSublocationFunc    func(ctx context.Context, userID string, req types.CreateSublocationRequest) (models.Sublocation, error)
	UpdateSublocationFunc    func(ctx context.Context, userID string, locationID string, req types.UpdateSublocationRequest) error
	DeleteSublocationFunc    func(ctx context.Context, userID string, sublocationIDs []string) (types.DeleteSublocationResponse, error)
	MoveGameFunc            func(ctx context.Context, userID string, req types.MoveGameRequest) error
	RemoveGameFunc          func(ctx context.Context, userID string, req types.RemoveGameRequest) error
}

func (m *MockSublocationService) GetSublocations(
	ctx context.Context,
	userID string,
) ([]models.Sublocation, error) {
	if m.GetSublocationsFunc != nil {
		return m.GetSublocationsFunc(ctx, userID)
	}
	return []models.Sublocation{}, nil
}

func (m *MockSublocationService) GetSingleSublocation(
	ctx context.Context,
	userID,
	locationID string,
) (models.Sublocation, error) {
	if m.GetSingleSublocationFunc != nil {
		return m.GetSingleSublocationFunc(ctx, userID, locationID)
	}
	return models.Sublocation{}, nil
}

func (m *MockSublocationService) CreateSublocation(
	ctx context.Context,
	userID string,
	req types.CreateSublocationRequest,
) (models.Sublocation, error) {
	if m.CreateSublocationFunc != nil {
		return m.CreateSublocationFunc(ctx, userID, req)
	}
	return models.Sublocation{}, nil
}

func (m *MockSublocationService) UpdateSublocation(
	ctx context.Context,
	userID string,
	locationID string,
	req types.UpdateSublocationRequest,
) error {
	if m.UpdateSublocationFunc != nil {
		return m.UpdateSublocationFunc(ctx, userID, locationID, req)
	}
	return nil
}

func (m *MockSublocationService) DeleteSublocation(
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

func (m *MockSublocationService) MoveGame(
	ctx context.Context,
	userID string,
	req types.MoveGameRequest,
) error {
	if m.MoveGameFunc != nil {
		return m.MoveGameFunc(ctx, userID, req)
	}
	return nil
}

func (m *MockSublocationService) RemoveGame(
	ctx context.Context,
	userID string,
	req types.RemoveGameRequest,
) error {
	if m.RemoveGameFunc != nil {
		return m.RemoveGameFunc(ctx, userID, req)
	}
	return nil
}
