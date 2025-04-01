package library

import (
	"context"
	"errors"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/infrastructure/cache"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
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
	GetLibraryItems(ctx context.Context, userID string) ([]models.Game, error)
	AddGameToLibrary(ctx context.Context, userID string, game models.Game) error
	DeleteGameFromLibrary(ctx context.Context, userID string, gameID int64) error
	GetGameByID(ctx context.Context, userID string, gameID int64) (models.Game, error)
	UpdateGameInLibrary(ctx context.Context, userID string, game models.Game) error
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

// GET
func (ls *GameLibraryService) GetLibraryItems(ctx context.Context, userID string) ([]models.Game, error) {
	// Validate the user ID
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return nil, err
	}

	// Attempt to get items from cache
	cachedGames, err := ls.cacheWrapper.GetCachedLibraryItems(ctx, userID)
	if err == nil && cachedGames != nil {
		ls.logger.Debug("Cache hit for library items", map[string]any{"userID": userID})
		return cachedGames, nil
	}

	// On cache miss, get from db
	ls.logger.Debug("Cache miss for user library, fetching from db", map[string]any{"userID": userID})
	games, err := ls.dbAdapter.GetUserLibraryItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache the results
	if err := ls.cacheWrapper.SetCachedLibraryItems(ctx, userID, games); err != nil {
		ls.logger.Error("Failed to cache library items", map[string]any{"error": err})
	}

	return games, nil
}

func (ls *GameLibraryService) GetUserGame(ctx context.Context, userID string, gameID int64) (models.Game, bool, error) {
	// Validate inputs
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return models.Game{}, false, err
	}
	if err := ls.validator.ValidateGameID(gameID); err != nil {
		return models.Game{}, false, err
	}

	// Try to get cache first
	cachedGame, found, err := ls.cacheWrapper.GetCachedGame(ctx, userID, gameID)
	if err == nil && found {
		ls.logger.Debug("Cache hit for user game", map[string]any{"userID": userID, "gameID": gameID})
		return *cachedGame, true, nil
	}

	// Cache miss, get from db
	ls.logger.Debug("Cache miss for user game, fetching from database", map[string]any{"userID": userID, "gameID": gameID})
	game, found, err := ls.dbAdapter.GetUserGame(ctx, userID, gameID)
	if err != nil {
		return models.Game{}, false, err
	}

	if found {
		// Cache the results
		if err := ls.cacheWrapper.SetCachedGame(ctx, userID, game); err != nil {
			ls.logger.Error("Failed to cache user game", map[string]any{"error": err})
		}
	}

	return game, found, nil
}

// POST
func (ls *GameLibraryService) AddGameToLibrary(ctx context.Context, userID string, game models.Game) error {
	// Validate inputs
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return err
	}
	if err := ls.validator.ValidateGameID(game.ID); err != nil {
		return err
	}

	// Add to db
	if err := ls.dbAdapter.AddGameToLibrary(ctx, userID, game.ID); err != nil {
		return err
	}

	// Invalidate cache
	if err := ls.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		ls.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}

	return nil
}

// DELETE
func (ls *GameLibraryService) DeleteGameFromLibrary(ctx context.Context, userID string, gameID int64) error {
	// Validate inputs
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return err
	}
	if err := ls.validator.ValidateGameID(gameID); err != nil {
		return err
	}

	// Remove from database
	if err := ls.dbAdapter.RemoveGameFromLibrary(ctx, userID, gameID); err != nil {
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

func (ls *GameLibraryService) GetGameByID(ctx context.Context, userID string, gameID int64) (models.Game, error) {
	// Single database call to get the game while verifying ownership
	game, exists, err := ls.dbAdapter.GetUserGame(ctx, userID, gameID)
	if err != nil {
		return models.Game{}, err
	}

	if !exists {
		return models.Game{}, ErrGameNotFound
	}

	return game, nil
}

func (ls *GameLibraryService) UpdateGameInLibrary(ctx context.Context, userID string, game models.Game) error {
	// NOTE: May need implementation if I want to support updating game metadata
	return errors.New("not implemented")
}
