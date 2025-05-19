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
	) (types.LibraryGameDBResult, []types.LibraryGamePhysicalLocationDBResponse, []types.LibraryGameDigitalLocationDBResponse, bool, error)

	SetCachedGame(
		ctx context.Context,
		userID string,
		game types.LibraryGameDBResult,
		physicalLocations []types.LibraryGamePhysicalLocationDBResponse,
		digitalLocations []types.LibraryGameDigitalLocationDBResponse,
	) error

	InvalidateUserCache(ctx context.Context, userID string) error
	InvalidateGameCache(ctx context.Context, userID string, gameID int64) error
}