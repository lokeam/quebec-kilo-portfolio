package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/types"
)

type MockDashboardService struct {
	GetDashboardBFFResponseFunc func(ctx context.Context, userID string) (types.DashboardBFFResponse, error)
}

func (m *MockDashboardService) GetDashboardBFFResponse(
	ctx context.Context,
	userID string,
) (types.DashboardBFFResponse, error) {
	if m.GetDashboardBFFResponseFunc != nil {
		return m.GetDashboardBFFResponseFunc(ctx, userID)
	}
	return types.DashboardBFFResponse{}, nil
}