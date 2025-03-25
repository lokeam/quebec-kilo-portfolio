package wishlist

import (
	"context"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
)

type GameWishlistService struct {
	logger interfaces.Logger
}

type WishlistService interface {
	GetWishlistItems(ctx context.Context, userID string) ([]models.Game, error)
}

func NewGameWishlistService(appContext *appcontext.AppContext) (*GameWishlistService, error) {
	return &GameWishlistService{
		logger: appContext.Logger,
	}, nil
}

func (w *GameWishlistService) GetWishlistItems(ctx context.Context, userID string) ([]models.Game, error) {
	// TODO: Implement, return an empty list for testing
	return []models.Game{}, nil
}
