package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/types"
)

type LibraryCacheWrapper interface {
	GetCachedLibraryItems(
		ctx context.Context,
		userID string,
	) ([]types.LibraryGameDBResult, []types.LibraryGamePhysicalLocationDBResponse, []types.LibraryGameDigitalLocationDBResponse, error)

	SetCachedLibraryItems(
		ctx context.Context,
		userID string,
		games []types.LibraryGameDBResult,
		physicalLocations []types.LibraryGamePhysicalLocationDBResponse,
		digitalLocations []types.LibraryGameDigitalLocationDBResponse,
	) error

	GetCachedGame(
		ctx context.Context,
		userID string,
		gameID int64,
	) (types.LibraryGameItemBFFResponseFINAL, bool, error)

	SetCachedGame(
		ctx context.Context,
		userID string,
		game types.LibraryGameItemBFFResponseFINAL,
	) error

	InvalidateUserCache(ctx context.Context, userID string) error
	InvalidateGameCache(ctx context.Context, userID string, gameID int64) error

	// BFF methods
	GetCachedLibraryItemsBFF(ctx context.Context, userID string) (types.LibraryBFFResponseFINAL, error)
	SetCachedLibraryItemsBFF(ctx context.Context, userID string, response types.LibraryBFFResponseFINAL) error
}