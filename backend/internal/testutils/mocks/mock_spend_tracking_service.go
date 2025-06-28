package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type MockSpendTrackingService struct {
	GetSpendTrackingBFFResponseFunc func(ctx context.Context, userID string) (types.SpendTrackingBFFResponseFINAL, error)
	GetSingleSpendTrackingItemFunc func(ctx context.Context, userID string, itemID string) (models.SpendTrackingOneTimePurchaseDB, error)

	CreateOneTimePurchaseFunc func(ctx context.Context, userID string, request types.SpendTrackingRequest) (models.SpendTrackingOneTimePurchaseDB, error)
	UpdateOneTimePurchaseFunc func(ctx context.Context, userID string, request types.SpendTrackingRequest) error
	DeleteSpendTrackingItemsFunc func(ctx context.Context, userID string, itemIDs []string) (types.DeleteSpendTrackingResponse, error)
}


func (m *MockSpendTrackingService) GetSpendTrackingBFFResponse(
	ctx context.Context,
	userID string,
) (types.SpendTrackingBFFResponseFINAL, error) {
	if m.GetSpendTrackingBFFResponseFunc != nil {
		return m.GetSpendTrackingBFFResponseFunc(ctx, userID)
	}
	return types.SpendTrackingBFFResponseFINAL{}, nil
}

func (m *MockSpendTrackingService) GetSingleSpendTrackingItem(
	ctx context.Context,
	userID string,
	itemID string,
) (models.SpendTrackingOneTimePurchaseDB, error) {
	if m.GetSingleSpendTrackingItemFunc != nil {
		return m.GetSingleSpendTrackingItemFunc(ctx, userID, itemID)
	}
	return models.SpendTrackingOneTimePurchaseDB{}, nil
}

func (m *MockSpendTrackingService) CreateOneTimePurchase(
	ctx context.Context,
	userID string,
	request types.SpendTrackingRequest,
) (models.SpendTrackingOneTimePurchaseDB, error) {
	if m.CreateOneTimePurchaseFunc != nil {
		return m.CreateOneTimePurchaseFunc(ctx, userID, request)
	}
	return models.SpendTrackingOneTimePurchaseDB{}, nil
}

func (m *MockSpendTrackingService) UpdateOneTimePurchase(
	ctx context.Context,
	userID string,
	request types.SpendTrackingRequest,
) error {
	if m.UpdateOneTimePurchaseFunc != nil {
		return m.UpdateOneTimePurchaseFunc(ctx, userID, request)
	}
	return nil
}

func (m *MockSpendTrackingService) DeleteSpendTrackingItems(
	ctx context.Context,
	userID string,
	itemIDs []string,
) (types.DeleteSpendTrackingResponse, error) {
	if m.DeleteSpendTrackingItemsFunc != nil {
		return m.DeleteSpendTrackingItemsFunc(ctx, userID, itemIDs)
	}
	return types.DeleteSpendTrackingResponse{}, nil
}