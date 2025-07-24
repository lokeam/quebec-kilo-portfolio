package library

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type GameLibraryService struct {
	dbAdapter interfaces.LibraryDbAdapter
	cacheWrapper interfaces.LibraryCacheWrapper
	dashboardCacheWrapper interfaces.DashboardCacheWrapper
	validator interfaces.LibraryValidator
	logger interfaces.Logger
}

type LibraryService interface {
	CreateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error

	UpdateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error
	DeleteLibraryGame(ctx context.Context, userID string, gameID int64) error

	// IsGameInLibraryBFF checks if a game is in a user's library, first checks cache then db as fallback
	IsGameInLibraryBFF(ctx context.Context, userID string, gameID int64) (bool, error)

	// GetAllLibraryItemsBFF returns a BFF response containing all library items and recently added items
	GetAllLibraryItemsBFF(ctx context.Context, userID string) (types.LibraryBFFResponseFINAL, error)

	// InvalidateUserCache invalidates all cache entries for a specific user
	InvalidateUserCache(ctx context.Context, userID string) error

	// REFACTORED RESPONSE
	GetLibraryRefactoredBFFResponse(ctx context.Context, userID string) (types.LibraryBFFRefactoredResponse, error)
}

func NewGameLibraryService(
	appContext *appcontext.AppContext,
	dbAdapter interfaces.LibraryDbAdapter,
	cacheWrapper interfaces.LibraryCacheWrapper,
	dashboardCacheWrapper interfaces.DashboardCacheWrapper,
) (*GameLibraryService, error) {
	if dbAdapter == nil {
		return nil, fmt.Errorf("dbAdapter is required")
	}
	if cacheWrapper == nil {
		return nil, fmt.Errorf("cacheWrapper is required")
	}
	if dashboardCacheWrapper == nil {
		return nil, fmt.Errorf("dashboardCacheWrapper is required")
	}

	return &GameLibraryService{
		dbAdapter: dbAdapter,
		cacheWrapper: cacheWrapper,
		dashboardCacheWrapper: dashboardCacheWrapper,
		validator: NewLibraryValidator(),
		logger: appContext.Logger,
	}, nil
}

// POST
func (ls *GameLibraryService) CreateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error {
	ls.logger.Info("GameLibraryService - CreateLibraryGame called", map[string]any{
		"userID": userID,
		"gameID": game.GameID,
	})

	// Validate inputs
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	if err := ls.validator.ValidateLibraryGame(game); err != nil {
		return fmt.Errorf("invalid game: %w", err)
	}

	// Add to db
	if err := ls.dbAdapter.CreateLibraryGame(ctx, userID, game); err != nil {
		return err
	}

	// Invalidate library cache
	if err := ls.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		ls.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}

	// Invalidate dashboard cache to refresh statistics
	if err := ls.dashboardCacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		ls.logger.Error("Failed to invalidate dashboard cache after adding game", map[string]any{
			"error": err,
			"userID": userID,
		})
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

	// Invalidate library cache
	if err := ls.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		ls.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}
	if err := ls.cacheWrapper.InvalidateGameCache(ctx, userID, gameID); err != nil {
		ls.logger.Error("Failed to invalidate game cache", map[string]any{"error": err})
	}

	// Invalidate dashboard cache to refresh statistics
	if err := ls.dashboardCacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		ls.logger.Error("Failed to invalidate dashboard cache after deleting game", map[string]any{
			"error": err,
			"userID": userID,
		})
	}

	return nil
}

// DeleteGameVersions handles batch deletion of specific platform versions of a game
func (ls *GameLibraryService) DeleteGameVersions(
	ctx context.Context,
	userID string,
	gameID int64,
	request types.BatchDeleteLibraryGameRequest,
) (types.BatchDeleteLibraryGameResponse, error) {
	ls.logger.Info("GameLibraryService - DeleteGameVersions called", map[string]any{
		"userID":    userID,
		"gameID":    gameID,
		"deleteAll": request.DeleteAll,
		"versions":  request.Versions,
	})

	// Validate inputs
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return types.BatchDeleteLibraryGameResponse{}, err
	}
	if err := ls.validator.ValidateGameID(gameID); err != nil {
		return types.BatchDeleteLibraryGameResponse{}, err
	}

	// Validate request
	if !request.DeleteAll && len(request.Versions) == 0 {
		return types.BatchDeleteLibraryGameResponse{}, fmt.Errorf("no versions specified for deletion")
	}

	// Remove from database
	response, err := ls.dbAdapter.DeleteGameVersions(ctx, userID, gameID, request)
	if err != nil {
		return types.BatchDeleteLibraryGameResponse{}, err
	}

	// Invalidate cache
	if err := ls.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		ls.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}
	if err := ls.cacheWrapper.InvalidateGameCache(ctx, userID, gameID); err != nil {
		ls.logger.Error("Failed to invalidate game cache", map[string]any{"error": err})
	}

	return response, nil
}

func (ls *GameLibraryService) GetSingleLibraryGame(
	ctx context.Context,
	userID string,
	gameID int64,
) (types.LibraryGameItemBFFResponseFINAL, error) {
	// Validate inputs
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return types.LibraryGameItemBFFResponseFINAL{}, err
	}
	if err := ls.validator.ValidateGameID(gameID); err != nil {
		return types.LibraryGameItemBFFResponseFINAL{}, err
	}

	// Try to get cache first
	game, found, err := ls.cacheWrapper.GetCachedGame(ctx, userID, gameID)
	if err == nil && found {
		ls.logger.Debug("Cache hit for user game", map[string]any{
			"userID": userID,
			"gameID": gameID,
		})
		return game, nil
	}

	// Cache miss, get from db
	ls.logger.Debug("Cache miss for user game, fetching from database", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})
	game, err = ls.dbAdapter.GetSingleLibraryGame(ctx, userID, gameID)
	if err != nil {
		return types.LibraryGameItemBFFResponseFINAL{}, err
	}

	// Cache the results
	if err := ls.cacheWrapper.SetCachedGame(ctx, userID, game); err != nil {
		ls.logger.Error("Failed to cache user game", map[string]any{"error": err})
	}

	return game, nil
}

func (ls *GameLibraryService) UpdateLibraryGame(ctx context.Context, userID string, game models.GameToSave) error {
	ls.logger.Info("GameLibraryService - UpdateLibraryGame called", map[string]any{
		"userID": userID,
		"gameID": game.GameID,
	})

	// Validate inputs
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	if err := ls.validator.ValidateLibraryGame(game); err != nil {
		return fmt.Errorf("invalid game: %w", err)
	}

	// Check if game exists in library
	exists, err := ls.dbAdapter.IsGameInLibrary(ctx, userID, game.GameID)
	if err != nil {
		return fmt.Errorf("error checking if game exists in library: %w", err)
	}
	if !exists {
		return ErrGameNotFound
	}

	// Update game in database
	if err := ls.dbAdapter.UpdateLibraryGame(ctx, game); err != nil {
		return fmt.Errorf("error updating game in library: %w", err)
	}

	// Invalidate library cache
	if err := ls.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		ls.logger.Error("Failed to invalidate user cache", map[string]any{
			"error": err,
			"userID": userID,
		})
	}
	if err := ls.cacheWrapper.InvalidateGameCache(ctx, userID, game.GameID); err != nil {
		ls.logger.Error("Failed to invalidate game cache", map[string]any{
			"error": err,
			"userID": userID,
			"gameID": game.GameID,
		})
	}

	// Invalidate dashboard cache to refresh statistics
	if err := ls.dashboardCacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		ls.logger.Error("Failed to invalidate dashboard cache after updating game", map[string]any{
			"error": err,
			"userID": userID,
		})
	}

	return nil
}

func (ls *GameLibraryService) IsGameInLibraryBFF(
	ctx context.Context,
	userID string,
	gameID int64,
) (bool, error) {
	// Validate userID
	if err := ls.validator.ValidateUserID(userID); err != nil {
		return false, err
	}
	if err := ls.validator.ValidateGameID(gameID); err != nil {
		return false, err
	}

	// Try to grab data from cache first
	cachedResponse, err := ls.cacheWrapper.GetCachedLibraryItemsBFF(ctx, userID)
	if err == nil && len(cachedResponse.LibraryItems) > 0 {
		ls.logger.Debug("Cache hit for library items, checking if game is in user library", map[string]any{
			"userID": userID,
			"gameID": gameID,
		})

		// Check if game is in cached items
		for i := 0; i < len(cachedResponse.LibraryItems); i++ {
			item := cachedResponse.LibraryItems[i]
			if item.ID == gameID {
				return true, nil
			}
		}
		return false, nil
	}

	// If cache miss or error, fallback to db
	ls.logger.Debug("Cache miss for is game in library search, falling back to database", map[string]any{
		"userID": userID,
		"gameID": gameID,
	})

	return ls.dbAdapter.IsGameInLibrary(ctx, userID, gameID)
}

// GetAllLibraryItemsBFF retrieves all library items for a user in BFF format
func (ls *GameLibraryService) GetAllLibraryItemsBFF(
	ctx context.Context,
	userID string,
) (types.LibraryBFFResponseFINAL, error) {
	if userID == "" {
		return types.LibraryBFFResponseFINAL{}, errors.New("user ID is required")
	}

	// Try to get from cache first
	cachedResponse, err := ls.cacheWrapper.GetCachedLibraryItemsBFF(ctx, userID)
	if err == nil && len(cachedResponse.LibraryItems) > 0 {
		return cachedResponse, nil
	}

	// If not in cache, get from database
	response, err := ls.dbAdapter.GetLibraryBFFResponse(ctx, userID)
	if err != nil {
		return types.LibraryBFFResponseFINAL{}, err
	}

	// Cache the response
	if err := ls.cacheWrapper.SetCachedLibraryItemsBFF(ctx, userID, response); err != nil {
		// Log the error but don't fail the request
		log.Printf("Failed to cache library BFF response: %v", err)
	}

	return response, nil
}


// REFACTORED RESPONSE
func (ls *GameLibraryService) GetLibraryRefactoredBFFResponse(
	ctx context.Context,
	userID string,
) (types.LibraryBFFRefactoredResponse, error) {
	if userID == "" {
			return types.LibraryBFFRefactoredResponse{}, errors.New("user ID is required")
	}

	// Try to get from cache first (if cache supports new structure)
	// For now, always get from database
	response, err := ls.dbAdapter.GetLibraryRefactoredBFFResponse(ctx, userID)
	if err != nil {
			return types.LibraryBFFRefactoredResponse{}, err
	}

	return response, nil
}





// InvalidateUserCache invalidates all cache entries for a specific user
func (ls *GameLibraryService) InvalidateUserCache(ctx context.Context, userID string) error {
	return ls.cacheWrapper.InvalidateUserCache(ctx, userID)
}