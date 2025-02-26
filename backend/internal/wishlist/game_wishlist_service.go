package wishlist

import (
	"context"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/types"
)

type GameWishlistService struct {
	logger interfaces.Logger
}

type WishlistService interface {
	GetWishlistItems(ctx context.Context, userID string) ([]types.Game, error)
}

func NewGameWishlistService(appContext *appcontext.AppContext) *GameWishlistService {
	return &GameWishlistService{
		logger: appContext.Logger,
	}
}

func (w *GameWishlistService) GetWishlistItems(ctx context.Context, userID string) ([]types.Game, error) {
	// TODO: Implement, return an empty list for testing
	return []types.Game{}, nil
}
