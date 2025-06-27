package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type MockSpendTrackingService struct {
	GetSpendTrackingBFFResponseFunc func(ctx context.Context, userID string) (types.SpendTrackingBFFResponseFINAL, error)
	CreateOneTimePurchaseFunc func(ctx context.Context, userID string, request types.SpendTrackingRequest) (models.SpendTrackingOneTimePurchaseDB, error)
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