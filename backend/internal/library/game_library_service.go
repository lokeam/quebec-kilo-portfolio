package library

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/infrastructure/cache"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
	"github.com/lokeam/qko-beta/internal/types"
)

type GameLibraryService struct {
	dbAdapter      *LibraryDbAdapter
	config         *config.Config
	cacheWrapper   interfaces.LibraryCacheWrapper
	logger         interfaces.Logger
	sanitizer      interfaces.Sanitizer
	validator      interfaces.LibraryValidator
}

type LibraryService interface {
	CreateLibraryGame(ctx context.Context, userID string, game models.LibraryGame) error
	GetAllLibraryGames(
		ctx context.Context,
		userID string,
	) (
		[]types.LibraryGameDBResult,
		[]types.LibraryGamePhysicalLocationDBResponse,
		[]types.LibraryGameDigitalLocationDBResponse,
		error,
	)
	GetSingleLibraryGame(
		ctx context.Context,
		userID string,
		gameID int64,
	) (
		types.LibraryGameDBResult,
		[]types.LibraryGamePhysicalLocationDBResponse,
		[]types.LibraryGameDigitalLocationDBResponse,
		error,
	)
	UpdateLibraryGame(ctx context.Context, userID string, game models.LibraryGame) error
	DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error

	// GetAllLibraryItemsBFF returns a BFF response containing all library items and recently added items
	GetAllLibraryItemsBFF(ctx context.Context, userID string) (types.LibraryBFFResponse, error)

	// InvalidateUserCache invalidates all cache entries for a specific user
	InvalidateUserCache(ctx context.Context, userID string) error
}

// Constructor that properly initializes the adapter
func NewGameLibraryService(appContext *appcontext.AppContext) (*GameLibraryService, error) {
	// Create and initialize the database adapter
	dbAdapter, err := NewLibraryDbAdapter(appContext)
	if err != nil {
			appContext.Logger.Error("Failed to create dbAdapter", map[string]any{"error": err})
			return nil, err
	}
	appContext.Logger.Info("dbAdapter created successfully", nil)

	// Create sanitizer to feed into validator
	sanitizer, err := security.NewSanitizer()
	if err != nil {
			appContext.Logger.Error("Failed to create sanitizer", map[string]any{"error": err})
			return nil, err
	}
	appContext.Logger.Info("sanitizer created successfully", nil)

	// Create validator
	validator, err := NewLibraryValidator(sanitizer)
	if err != nil {
			appContext.Logger.Error("Failed to create validator", map[string]any{"error": err})
			return nil, err
	}
	appContext.Logger.Info("validator created successfully", nil)

	// Create cache wrapper to handle Redis caching
	cacheWrapper, err := cache.NewCacheWrapper(
		appContext.RedisClient,
		appContext.Config.Redis.RedisTTL,
		appContext.Config.Redis.RedisTimeout,
		appContext.Logger,
	)
	if err != nil {
		return nil, err
	}

	// Create library cache adapter
	libraryCacheAdapter, err := NewLibraryCacheAdapter(cacheWrapper)
	if err != nil {
		return nil, err
	}

	// Sanity check that all dependencies are initialized
	appContext.Logger.Info("GameLibraryService components initialized",
		map[string]any{
			"dbAdapter": dbAdapter != nil,
			"validator": validator != nil,
			"cacheWrapper": libraryCacheAdapter != nil,
		},
	)


	return &GameLibraryService{
		dbAdapter:    dbAdapter,
		validator:    validator,
		logger:       appContext.Logger,
		config:       appContext.Config,
		cacheWrapper: libraryCacheAdapter,
		sanitizer:    sanitizer,
	}, nil
}

// GetAllLibraryGames retrieves all games in a user's library
func (ls *GameLibraryService) GetAllLibraryGames(
	ctx context.Context,
	userID string,
) (
	[]types.LibraryGameDBResult,
	[]types.LibraryGamePhysicalLocationDBResponse,
	[]types.LibraryGameDigitalLocationDBResponse,
	error,
) {
	// Validate the user ID
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return nil, nil, nil, err
	}

	// Attempt to get items from cache
	games, physicalLocations, digitalLocations, err := ls.cacheWrapper.GetCachedLibraryItems(ctx, userID)
	if err == nil && games != nil {
		ls.logger.Debug("Cache hit for library items", map[string]any{"userID": userID})
		return games, physicalLocations, digitalLocations, nil
	}

	// On cache miss, get from db
	ls.logger.Debug("Cache miss for user library, fetching from db", map[string]any{"userID": userID})
	games, physicalLocations, digitalLocations, err = ls.dbAdapter.GetUserLibraryItems(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Cache the results
	if err := ls.cacheWrapper.SetCachedLibraryItems(
		ctx,
		userID,
		games,
		physicalLocations,
		digitalLocations,
	); err != nil {
		ls.logger.Error("Failed to cache library items", map[string]any{"error": err})
	}

	return games, physicalLocations, digitalLocations, nil
}

// POST
func (ls *GameLibraryService) CreateLibraryGame(
	ctx context.Context,
	userID string,
	game models.LibraryGame,
) error {
	// Validate inputs
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return err
	}
	if err := ls.validator.ValidateGameID(game.GameID); err != nil {
		return err
	}

	// Add to db
	if err := ls.dbAdapter.CreateLibraryGame(ctx, userID, game); err != nil {
		return err
	}

	// Invalidate cache
	if err := ls.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		ls.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}

	return nil
}

// DELETE
func (ls *GameLibraryService) DeleteLibraryGame(
	ctx context.Context,
	userID string,
	gameID int64,
) error {
	// Validate inputs
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return err
	}
	if err := ls.validator.ValidateGameID(gameID); err != nil {
		return err
	}

	// Remove from database
	if err := ls.dbAdapter.DeleteLibraryGame(ctx, userID, gameID); err != nil {
		return err
	}

	// Invalidate cache
	if err := ls.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		ls.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}
	if err := ls.cacheWrapper.InvalidateGameCache(ctx, userID, gameID); err != nil {
		ls.logger.Error("Failed to invalidate game cache", map[string]any{"error": err})
	}

	return nil
}

func (ls *GameLibraryService) GetSingleLibraryGame(
	ctx context.Context,
	userID string,
	gameID int64,
) (
	types.LibraryGameDBResult,
	[]types.LibraryGamePhysicalLocationDBResponse,
	[]types.LibraryGameDigitalLocationDBResponse,
	error,
) {
	// Validate inputs
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return types.LibraryGameDBResult{}, nil, nil, err
	}
	if err := ls.validator.ValidateGameID(gameID); err != nil {
		return types.LibraryGameDBResult{}, nil, nil, err
	}

	// Try to get cache first
	game, physicalLocations, digitalLocations, found, err := ls.cacheWrapper.GetCachedGame(ctx, userID, gameID)
	if err == nil && found {
			ls.logger.Debug("Cache hit for user game", map[string]any{
					"userID": userID,
					"gameID": gameID,
			})
			return game, physicalLocations, digitalLocations, nil
	}

	// Cache miss, get from db
	ls.logger.Debug("Cache miss for user game, fetching from database", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})
	game, physicalLocations, digitalLocations, err = ls.dbAdapter.GetSingleLibraryGame(ctx, userID, gameID)
	if err != nil {
			return types.LibraryGameDBResult{}, nil, nil, err
  }

	// Cache the results
	if err := ls.cacheWrapper.SetCachedGame(
		ctx,
		userID,
		game,
		physicalLocations,
		digitalLocations,
	); err != nil {
			ls.logger.Error("Failed to cache user game", map[string]any{"error": err})
	}

	return game, physicalLocations, digitalLocations, nil
}

func (ls *GameLibraryService) UpdateLibraryGame(
	ctx context.Context,
	userID string,
	game models.LibraryGame,
) error {
	// Validate inputs
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return err
	}
	if err := ls.validator.ValidateGame(game); err != nil {
			return err
	}

	// Update in db
	if err := ls.dbAdapter.UpdateLibraryGame(ctx, userID, game); err != nil {
			return err
	}

	// Invalidate both user and game cache
	if err := ls.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
			ls.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}
	if err := ls.cacheWrapper.InvalidateGameCache(ctx, userID, game.GameID); err != nil {
			ls.logger.Error("Failed to invalidate game cache", map[string]any{"error": err})
	}

	return nil
}

// GetAllLibraryItemsBFF returns a BFF response containing all library items and recently added items
func (ls *GameLibraryService) GetAllLibraryItemsBFF(
	ctx context.Context,
	userID string,
) (types.LibraryBFFResponse, error) {
	ls.logger.Debug("GetAllLibraryItemsBFF called", map[string]any{
		"userID": userID,
	})

	// Try to get from cache first
	response, err := ls.cacheWrapper.GetCachedLibraryItemsBFF(ctx, userID)
	if err == nil {
		ls.logger.Debug("Cache hit for library items BFF", map[string]any{
			"userID": userID,
			"libraryItemsCount": len(response.LibraryItems),
			"recentlyAddedCount": len(response.RecentlyAdded),
		})
		return response, nil
	}

	// Cache miss or error, get from dbAdapter
	response, err = ls.dbAdapter.GetLibraryBFFResponse(ctx, userID)
	if err != nil {
		ls.logger.Error("Error getting library items BFF from dbAdapter", map[string]any{
			"error": err,
			"userID": userID,
		})
		return types.LibraryBFFResponse{}, fmt.Errorf("error getting library items BFF from dbAdapter: %w", err)
	}

	// Cache the response
	if err := ls.cacheWrapper.SetCachedLibraryItemsBFF(ctx, userID, response); err != nil {
		ls.logger.Error("Error caching library items BFF", map[string]any{
			"error": err,
			"userID": userID,
		})
		// Don't return error here, just log it
	}

	ls.logger.Debug("GetAllLibraryItemsBFF success", map[string]any{
		"userID": userID,
		"libraryItemsCount": len(response.LibraryItems),
		"recentlyAddedCount": len(response.RecentlyAdded),
	})

	return response, nil
}

// InvalidateUserCache invalidates all cache entries for a specific user
func (ls *GameLibraryService) InvalidateUserCache(ctx context.Context, userID string) error {
	return ls.cacheWrapper.InvalidateUserCache(ctx, userID)
}
