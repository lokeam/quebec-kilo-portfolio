package library

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
)

type LibraryCacheAdapter struct {
	cacheWrapper interfaces.CacheWrapper
}

func NewLibraryCacheAdapter(
	cacheWrapper interfaces.CacheWrapper,
) (interfaces.LibraryCacheWrapper, error) {
	return &LibraryCacheAdapter{
		cacheWrapper: cacheWrapper,
	}, nil
}

func (lca *LibraryCacheAdapter) GetCachedLibraryItems(ctx context.Context, userID string) ([]models.Game, error) {
	cacheKey := fmt.Sprintf("library:%s", userID)

	var games []models.Game
	cacheHit, err := lca.cacheWrapper.GetCachedResults(ctx, cacheKey, &games)
	if err != nil {
		return nil, err
	}

	if cacheHit {
		return games, nil
	}

	return nil, nil
}

func (lca *LibraryCacheAdapter) SetCachedLibraryItems(
	ctx context.Context,
	userID string,
	games []models.Game,
) error {
	cacheKey := fmt.Sprintf("library:%s", userID)
	return lca.cacheWrapper.SetCachedResults(ctx, cacheKey, games)
}

func (lca *LibraryCacheAdapter) GetCachedGame(
	ctx context.Context,
	userID string,
	gameID int64,
) (*models.Game, bool, error) {
	cacheKey := fmt.Sprintf("library:%s:game:%d", userID, gameID)

	var game models.Game
	cacheHit, err := lca.cacheWrapper.GetCachedResults(ctx, cacheKey, &game)
	if err != nil {
		return nil, false, ErrDatabaseConnection
	}

	if cacheHit {
		return &game, true, nil
	}

	return nil, false, nil
}

func (lca *LibraryCacheAdapter) SetCachedGame(
	ctx context.Context,
	userID string,
	game models.Game,
	) error {
	cacheKey := fmt.Sprintf("library:%s:game:%d", userID, game.ID)
	return lca.cacheWrapper.SetCachedResults(ctx, cacheKey, game)
	}

// Invalidates all cache entries for a specific user
func (lca *LibraryCacheAdapter) InvalidateUserCache(
	ctx context.Context,
	userID string,
	) error {
		cacheKey := fmt.Sprintf("library:%s:game", userID)
		return lca.cacheWrapper.SetCachedResults(ctx, cacheKey, nil)
}

// Invalidates cache for a specific game
func (lca *LibraryCacheAdapter) InvalidateGameCache(
	ctx context.Context,
	userID string,
	gameID int64,
	) error {
		cacheKey := fmt.Sprintf("library:%s:game:%d", userID, gameID)
		return lca.cacheWrapper.SetCachedResults(ctx, cacheKey, nil)
	}
