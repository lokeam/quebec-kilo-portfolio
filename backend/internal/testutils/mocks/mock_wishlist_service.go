package mocks

import (
	"context"

	"github.com/lokeam/qko-beta/internal/models"
)

// MockWishlistService implements services.WishlistService
type MockWishlistService struct {
	GetWishlistItemsFunc func(ctx context.Context, userID string) ([]models.GameToSave, error)
}

func (m *MockWishlistService) GetWishlistItems(
	ctx context.Context,
	userID string,
) ([]models.GameToSave, error) {
	if m.GetWishlistItemsFunc != nil {
		return m.GetWishlistItemsFunc(ctx, userID)
	}
	return []models.GameToSave{}, nil
}