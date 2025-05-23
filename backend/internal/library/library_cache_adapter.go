package library

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/types"
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

type cachedLibraryItems struct {
	Games             []types.LibraryGameDBResult
	PhysicalLocations []types.LibraryGamePhysicalLocationDBResponse
	DigitalLocations  []types.LibraryGameDigitalLocationDBResponse
}

func (lca *LibraryCacheAdapter) GetCachedLibraryItems(
	ctx context.Context,
	userID string,
) ([]types.LibraryGameDBResult, []types.LibraryGamePhysicalLocationDBResponse, []types.LibraryGameDigitalLocationDBResponse, error) {
	cacheKey := fmt.Sprintf("library:%s", userID)

	var items cachedLibraryItems
	cacheHit, err := lca.cacheWrapper.GetCachedResults(ctx, cacheKey, &items)
	if err != nil {
		return nil, nil, nil, err
	}

	if cacheHit {
		return items.Games, items.PhysicalLocations, items.DigitalLocations, nil
	}

	return nil, nil, nil, nil
}

func (lca *LibraryCacheAdapter) SetCachedLibraryItems(
	ctx context.Context,
	userID string,
	games []types.LibraryGameDBResult,
	physicalLocations []types.LibraryGamePhysicalLocationDBResponse,
	digitalLocations []types.LibraryGameDigitalLocationDBResponse,
) error {
	cacheKey := fmt.Sprintf("library:%s", userID)

	items := cachedLibraryItems{
		Games:             games,
		PhysicalLocations: physicalLocations,
		DigitalLocations:  digitalLocations,
	}

	return lca.cacheWrapper.SetCachedResults(ctx, cacheKey, items)
}

type cachedGame struct {
	Game             types.LibraryGameDBResult
	PhysicalLocations []types.LibraryGamePhysicalLocationDBResponse
	DigitalLocations  []types.LibraryGameDigitalLocationDBResponse
}

func (lca *LibraryCacheAdapter) GetCachedGame(
	ctx context.Context,
	userID string,
	gameID int64,
) (types.LibraryGameDBResult, []types.LibraryGamePhysicalLocationDBResponse, []types.LibraryGameDigitalLocationDBResponse, bool, error) {
	cacheKey := fmt.Sprintf("library:%s:game:%d", userID, gameID)

	var item cachedGame
	cacheHit, err := lca.cacheWrapper.GetCachedResults(ctx, cacheKey, &item)
	if err != nil {
		return types.LibraryGameDBResult{}, nil, nil, false, err
	}

	if cacheHit {
		return item.Game, item.PhysicalLocations, item.DigitalLocations, true, nil
	}

	return types.LibraryGameDBResult{}, nil, nil, false, nil
}

func (lca *LibraryCacheAdapter) SetCachedGame(
	ctx context.Context,
	userID string,
	game types.LibraryGameDBResult,
	physicalLocations []types.LibraryGamePhysicalLocationDBResponse,
	digitalLocations []types.LibraryGameDigitalLocationDBResponse,
) error {
	cacheKey := fmt.Sprintf("library:%s:game:%d", userID, game.ID)

	item := cachedGame{
		Game:             game,
		PhysicalLocations: physicalLocations,
		DigitalLocations:  digitalLocations,
	}

	return lca.cacheWrapper.SetCachedResults(ctx, cacheKey, item)
}

// Invalidates all cache entries for a specific user
func (lca *LibraryCacheAdapter) InvalidateUserCache(
	ctx context.Context,
	userID string,
) error {
	cacheKey := fmt.Sprintf("library:%s", userID)
	return lca.cacheWrapper.DeleteCacheKey(ctx, cacheKey)
}

// Invalidates cache for a specific game
func (lca *LibraryCacheAdapter) InvalidateGameCache(
	ctx context.Context,
	userID string,
	gameID int64,
) error {
	cacheKey := fmt.Sprintf("library:%s:game:%d", userID, gameID)
	return lca.cacheWrapper.DeleteCacheKey(ctx, cacheKey)
}
