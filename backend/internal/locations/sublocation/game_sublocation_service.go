package sublocation

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/infrastructure/cache"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/services"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
	"github.com/lokeam/qko-beta/internal/types"
)

type GameSublocationService struct {
	dbAdapter       interfaces.SublocationDbAdapter
	config          *config.Config
	cacheWrapper    interfaces.SublocationCacheWrapper
	logger          interfaces.Logger
	sanitizer       interfaces.Sanitizer
	validator       interfaces.SublocationValidator
	physicalService services.PhysicalService
}

type SublocationService interface {
	GetSublocations(ctx context.Context, userID string) ([]models.Sublocation, error)
	GetSingleSublocation(ctx context.Context, userID string, sublocationID string) (models.Sublocation, error)
	CreateSublocation(ctx context.Context, userID string, req types.CreateSublocationRequest) (models.Sublocation, error)
	DeleteSublocation(ctx context.Context, userID string, sublocationID string) error
	UpdateSublocation(ctx context.Context, userID string, locationID string, req types.UpdateSublocationRequest) error
}

func NewGameSublocationService(
	appContext *appcontext.AppContext,
	physicalService services.PhysicalService,
) (*GameSublocationService, error) {
	// Create initialize db adapter
	dbAdapter, err := NewSublocationDbAdapter(appContext)
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

	//Create cache wrapper to handle Redis caching
	sublocationCacheAdapter, err := NewSublocationCacheAdapter(cacheWrapper)
	if err != nil {
		return nil, err
	}

	// Create validator
	validator, err := NewSublocationValidator(
		sanitizer,
		sublocationCacheAdapter,
		appContext.Logger,
	)
	if err != nil {
		appContext.Logger.Error("Failed to create validator", map[string]any{"error": err})
		return nil, err
	}
	appContext.Logger.Info("validator created successfully", nil)

	// Create physical media cache adapter
	appContext.Logger.Info("GameSublocationService components intialized", map[string]any{
		"dbAdapter": dbAdapter != nil,
		"validator": validator != nil,
		"cacheWrapper": sublocationCacheAdapter != nil,
		"physicalService": physicalService != nil,
	})

	// Sanity check that all dependencies are intialized

	return &GameSublocationService{
		dbAdapter:       dbAdapter,
		validator:       validator,
		logger:          appContext.Logger,
		config:          appContext.Config,
		cacheWrapper:    sublocationCacheAdapter,
		sanitizer:       sanitizer,
		physicalService: physicalService,
	}, nil
}

// GET
func (gss *GameSublocationService) GetSingleSublocation(
	ctx context.Context,
	userID string,
	sublocationID string,
) (models.Sublocation, error) {
	// Use GetSingleCachedSublocation instead of GetCachedSublocations
	sublocation, found, err := gss.cacheWrapper.GetSingleCachedSublocation(ctx, userID, sublocationID)
	if err == nil && found {
		return *sublocation, nil
	}

	// Get from DB
	return gss.dbAdapter.GetSingleSublocation(ctx, userID, sublocationID)
}

func (gss *GameSublocationService) GetSublocations(
	ctx context.Context,
	userID string,
) ([]models.Sublocation, error) {
	// Attempt to cached locations
	cachedLocations, err := gss.cacheWrapper.GetCachedSublocations(ctx, userID)
	if err == nil {
		gss.logger.Debug("Cache hit for sublocations", map[string]any{"userID": userID})
		return cachedLocations, nil
	}

	// Cache miss or error, get from DB
	gss.logger.Debug("Cache miss for sublocations, fetching from DB", map[string]any{"userID": userID})
	sublocations, err := gss.dbAdapter.GetAllSublocations(ctx, userID)
	if err != nil {
		gss.logger.Error("Failed to fetch sublocations from DB", map[string]any{"error": err})
		return nil, err
	}

	// Cache the results for future requests
	if cacheErr := gss.cacheWrapper.SetCachedSublocations(
		ctx,
		userID,
		sublocations,
	); cacheErr != nil {
		gss.logger.Error("Failed to cache sublocations", map[string]any{"error": cacheErr})
		// Continue w/ return the locations from DB
	}

	return sublocations, nil
}

// POST
// CreateSublocation adds a new sublocation for a user.
//
// Cache Invalidation Strategy:
// - Invalidates the user cache to ensure cache consistency
// - Invalidates the parent physical location's cache to maintain consistency across related caches
// - Forces immediate refresh of physical location cache by fetching it
//
// IMPORTANT: The cache invalidation is INTENTIONAL and NECESSARY:
// 1. User cache invalidation ensures consistency for user's sublocations
// 2. Physical location cache invalidation ensures consistency for parent location
// 3. Force refresh ensures immediate consistency
//
// DO NOT attempt to "optimize" this by removing cache invalidation:
// - Removing user cache invalidation would break cache consistency
// - Removing physical location cache invalidation would miss the physical location cache update
// - Removing force refresh could lead to stale data
//
// This pattern is consistent with our production codebase and follows defensive programming
// principles to ensure cache consistency across all possible code paths.
func (gss *GameSublocationService) CreateSublocation(
	ctx context.Context,
	userID string,
	req types.CreateSublocationRequest,
) (models.Sublocation, error) {
	// Transform request to model
	sublocation := TransformCreateRequestToModel(req, userID)

	// Validate sublocation
	validatedLocation, err := gss.validator.ValidateSublocation(sublocation)
	if err != nil {
		gss.logger.Error("Location validation failed", map[string]any{"error": err})
		return models.Sublocation{}, err
	}

	// Add to db
	createdSublocation, err := gss.dbAdapter.CreateSublocation(ctx, userID, validatedLocation)
	if err != nil {
		gss.logger.Error("Failed to add sublocation to DB", map[string]any{"error": err})
		return models.Sublocation{}, err
	}

	// Invalidate caches
	if err := gss.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gss.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}

	// Also invalidate the parent physical location's cache
	physicalLocationID := createdSublocation.PhysicalLocationID
	if err := gss.cacheWrapper.InvalidateLocationCache(ctx, userID, physicalLocationID); err != nil {
		gss.logger.Error("Failed to invalidate parent physical location cache", map[string]any{"error": err})
	}

	// Force immediate refresh of physical location cache by fetching it
	gss.logger.Debug("Forcing refresh of physical location cache", map[string]any{
		"userID": userID,
		"physicalLocationID": physicalLocationID,
	})

	// Actively refresh the physical location data if service is available
	if gss.physicalService != nil {
		_, refreshErr := gss.physicalService.GetSinglePhysicalLocation(ctx, userID, physicalLocationID)
		if refreshErr != nil {
			gss.logger.Warn("Failed to refresh physical location cache", map[string]any{
				"error": refreshErr,
				"physicalLocationID": physicalLocationID,
			})
		} else {
			gss.logger.Debug("Successfully refreshed physical location cache", map[string]any{
				"physicalLocationID": physicalLocationID,
			})
		}
	} else {
		gss.logger.Warn("Physical service not available for cache refresh", nil)
	}

	return createdSublocation, nil
}

// PUT
func (gss *GameSublocationService) UpdateSublocation(
	ctx context.Context,
	userID string,
	locationID string,
	req types.UpdateSublocationRequest,
) error {
	// First get the existing sublocation
	existingSublocation, err := gss.dbAdapter.GetSingleSublocation(ctx, userID, locationID)
	if err != nil {
		gss.logger.Error("Failed to get existing sublocation", map[string]any{"error": err})
		return err
	}

	// Transform request to model
	updatedSublocation := TransformUpdateRequestToModel(req, existingSublocation)

	// Validate only the fields that are being updated
	validatedSublocation, err := gss.validator.ValidateSublocationUpdate(updatedSublocation, existingSublocation)
	if err != nil {
		gss.logger.Error("Location validation failed", map[string]any{"error": err})
		return err
	}

	// Update in db
	if err := gss.dbAdapter.UpdateSublocation(
		ctx,
		userID,
		validatedSublocation,
	); err != nil {
		gss.logger.Error("Failed to update sublocation in DB", map[string]any{"error": err})
		return err
	}

	// Store the physical location ID for cache invalidation and logging
	physicalLocationID := validatedSublocation.PhysicalLocationID

	// Invalidate caches
	if err := gss.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gss.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}
	if err := gss.cacheWrapper.InvalidateSublocationCache(ctx, userID, validatedSublocation.ID); err != nil {
		gss.logger.Error("Failed to invalidate sublocation cache", map[string]any{"error": err})
	}
	// Also invalidate the parent physical location's cache
	if err := gss.cacheWrapper.InvalidateLocationCache(ctx, userID, physicalLocationID); err != nil {
		gss.logger.Error("Failed to invalidate parent physical location cache", map[string]any{"error": err})
	}

	// Force immediate refresh of physical location cache
	gss.logger.Debug("Forcing refresh of physical location cache after sublocation update", map[string]any{
		"userID": userID,
		"physicalLocationID": physicalLocationID,
	})

	// Actively refresh the physical location data if service is available
	if gss.physicalService != nil {
		_, refreshErr := gss.physicalService.GetSinglePhysicalLocation(ctx, userID, physicalLocationID)
		if refreshErr != nil {
			gss.logger.Warn("Failed to refresh physical location cache after update", map[string]any{
				"error": refreshErr,
				"physicalLocationID": physicalLocationID,
			})
		} else {
			gss.logger.Debug("Successfully refreshed physical location cache after update", map[string]any{
				"physicalLocationID": physicalLocationID,
			})
		}
	} else {
		gss.logger.Warn("Physical service not available for cache refresh after update", nil)
	}

	return nil
}

// DELETE
// DeleteSublocation deletes one or more sublocations and handles orphaned games
func (gss *GameSublocationService) DeleteSublocation(
	ctx context.Context,
	userID string,
	sublocationIDs []string,
) (types.DeleteSublocationResponse, error) {
	gss.logger.Debug("DeleteSublocation called", map[string]any{
		"userID": userID,
		"sublocationIDs": sublocationIDs,
	})

	// Validate the request
	if err := gss.validator.ValidateDeleteSublocationRequest(userID, sublocationIDs); err != nil {
		gss.logger.Error("Delete request validation failed", map[string]any{
			"error": err,
			"userID": userID,
			"sublocationIDs": sublocationIDs,
		})
		return types.DeleteSublocationResponse{
			Success: false,
			Error: err.Error(),
		}, err
	}

	// Get all sublocations to get their physical location IDs for cache invalidation
	sublocations, err := gss.dbAdapter.GetAllSublocations(ctx, userID)
	if err != nil {
		gss.logger.Error("Failed to get sublocations for cache invalidation", map[string]any{
			"error": err,
			"userID": userID,
		})
		return types.DeleteSublocationResponse{
			Success: false,
			Error: fmt.Sprintf("failed to get sublocations: %v", err),
		}, err
	}

	// Create a map of sublocation IDs to their physical location IDs
	physicalLocationIDs := make(map[string]string)
	for _, sublocation := range sublocations {
		physicalLocationIDs[sublocation.ID] = sublocation.PhysicalLocationID
	}

	// Delete sublocations and handle orphaned games
	response, err := gss.dbAdapter.DeleteSublocation(ctx, userID, sublocationIDs)
	if err != nil {
		gss.logger.Error("Failed to delete sublocations", map[string]any{
			"error": err,
			"userID": userID,
			"sublocationIDs": sublocationIDs,
		})
		return response, err
	}

	// Invalidate caches
	if err := gss.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gss.logger.Error("Failed to invalidate user cache", map[string]any{
			"error": err,
			"userID": userID,
		})
	}

	// Invalidate each sublocation's cache
	for _, sublocationID := range sublocationIDs {
		if err := gss.cacheWrapper.InvalidateSublocationCache(ctx, userID, sublocationID); err != nil {
			gss.logger.Error("Failed to invalidate sublocation cache", map[string]any{
				"error": err,
				"sublocationID": sublocationID,
			})
		}

		// Invalidate the parent physical location's cache
		if physicalLocationID, exists := physicalLocationIDs[sublocationID]; exists {
			if err := gss.cacheWrapper.InvalidateLocationCache(ctx, userID, physicalLocationID); err != nil {
				gss.logger.Error("Failed to invalidate parent physical location cache", map[string]any{
					"error": err,
					"physicalLocationID": physicalLocationID,
				})
			}

			// Force immediate refresh of physical location cache
			if gss.physicalService != nil {
				_, refreshErr := gss.physicalService.GetSinglePhysicalLocation(ctx, userID, physicalLocationID)
				if refreshErr != nil {
					gss.logger.Warn("Failed to refresh physical location cache", map[string]any{
						"error": refreshErr,
						"physicalLocationID": physicalLocationID,
					})
				} else {
					gss.logger.Debug("Successfully refreshed physical location cache", map[string]any{
						"physicalLocationID": physicalLocationID,
					})
				}
			}
		}
	}

	gss.logger.Debug("DeleteSublocation completed successfully", map[string]any{
		"userID": userID,
		"sublocationIDs": sublocationIDs,
		"deletedCount": response.DeletedCount,
		"deletedGames": len(response.DeletedGames),
	})

	return response, nil
}

// MoveGame moves a game from its current sublocation to a target sublocation
func (gss *GameSublocationService) MoveGame(
	ctx context.Context,
	userID string,
	req types.MoveGameRequest,
) error {
	gss.logger.Debug("MoveGame called", map[string]any{
		"userID": userID,
		"request": req,
	})

	// 1. Validate input format
	if err := gss.validator.ValidateGameOwnership(userID, req.UserGameID); err != nil {
		gss.logger.Error("Game ownership validation failed", map[string]any{
			"error": err,
			"userID": userID,
			"userGameID": req.UserGameID,
		})
		return err
	}

	if err := gss.validator.ValidateSublocationOwnership(userID, req.TargetSublocationID); err != nil {
		gss.logger.Error("Sublocation ownership validation failed", map[string]any{
			"error": err,
			"userID": userID,
			"sublocationID": req.TargetSublocationID,
		})
		return err
	}

	if err := gss.validator.ValidateGameNotInSublocation(req.UserGameID, req.TargetSublocationID); err != nil {
		gss.logger.Error("Game in sublocation validation failed", map[string]any{
			"error": err,
			"userGameID": req.UserGameID,
			"sublocationID": req.TargetSublocationID,
		})
		return err
	}

	// 2. Move game in database
	if err := gss.dbAdapter.MoveGameToSublocation(ctx, userID, req.UserGameID, req.TargetSublocationID); err != nil {
		gss.logger.Error("Failed to move game in database", map[string]any{
			"error": err,
			"userID": userID,
			"userGameID": req.UserGameID,
			"targetSublocationID": req.TargetSublocationID,
		})
		return err
	}

	// 3. Invalidate ALL related caches
	if err := gss.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gss.logger.Error("Failed to invalidate user cache", map[string]any{
			"error": err,
			"userID": userID,
		})
	}

	// Get all sublocations to invalidate their caches
	sublocations, err := gss.dbAdapter.GetAllSublocations(ctx, userID)
	if err != nil {
		gss.logger.Error("Failed to get sublocations for cache invalidation", map[string]any{
			"error": err,
			"userID": userID,
		})
	} else {
		for _, sublocation := range sublocations {
			// Invalidate sublocation cache
			if err := gss.cacheWrapper.InvalidateSublocationCache(ctx, userID, sublocation.ID); err != nil {
				gss.logger.Error("Failed to invalidate sublocation cache", map[string]any{
					"error": err,
					"sublocationID": sublocation.ID,
				})
			}

			// Invalidate physical location cache
			if err := gss.cacheWrapper.InvalidateLocationCache(ctx, userID, sublocation.PhysicalLocationID); err != nil {
				gss.logger.Error("Failed to invalidate physical location cache", map[string]any{
					"error": err,
					"physicalLocationID": sublocation.PhysicalLocationID,
				})
			}

			// Force refresh physical location cache
			if gss.physicalService != nil {
				_, refreshErr := gss.physicalService.GetSinglePhysicalLocation(ctx, userID, sublocation.PhysicalLocationID)
				if refreshErr != nil {
					gss.logger.Warn("Failed to refresh physical location cache", map[string]any{
						"error": refreshErr,
						"physicalLocationID": sublocation.PhysicalLocationID,
					})
				}
			}
		}
	}

	gss.logger.Debug("MoveGame completed successfully", map[string]any{
		"userID": userID,
		"userGameID": req.UserGameID,
		"targetSublocationID": req.TargetSublocationID,
	})

	return nil
}

// RemoveGame removes a game from its current sublocation
func (gss *GameSublocationService) RemoveGame(
	ctx context.Context,
	userID string,
	req types.RemoveGameRequest,
) error {
	gss.logger.Debug("RemoveGame called", map[string]any{
		"userID": userID,
		"request": req,
	})

	// 1. Validate input format
	if err := gss.validator.ValidateGameOwnership(userID, req.UserGameID); err != nil {
		gss.logger.Error("Game ownership validation failed", map[string]any{
			"error": err,
			"userID": userID,
			"userGameID": req.UserGameID,
		})
		return err
	}

	// 2. Remove game from sublocation in database
	if err := gss.dbAdapter.RemoveGameFromSublocation(ctx, userID, req.UserGameID); err != nil {
		gss.logger.Error("Failed to remove game from sublocation", map[string]any{
			"error": err,
			"userID": userID,
			"userGameID": req.UserGameID,
		})
		return err
	}

	// 3. Invalidate caches
	if err := gss.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gss.logger.Error("Failed to invalidate user cache", map[string]any{
			"error": err,
			"userID": userID,
		})
	}

	// Get the game's current sublocation to invalidate its cache
	exists, err := gss.dbAdapter.CheckGameInAnySublocation(ctx, req.UserGameID)
	if err != nil {
		gss.logger.Error("Failed to check game location for cache invalidation", map[string]any{
			"error": err,
			"userGameID": req.UserGameID,
		})
	} else if exists {
		// If game was in a sublocation, invalidate that sublocation's cache
		// Note: We don't know which sublocation it was in, so we need to invalidate all sublocations
		sublocations, err := gss.dbAdapter.GetAllSublocations(ctx, userID)
		if err != nil {
			gss.logger.Error("Failed to get sublocations for cache invalidation", map[string]any{
				"error": err,
				"userID": userID,
			})
		} else {
			for _, sublocation := range sublocations {
				if err := gss.cacheWrapper.InvalidateSublocationCache(ctx, userID, sublocation.ID); err != nil {
					gss.logger.Error("Failed to invalidate sublocation cache", map[string]any{
						"error": err,
						"sublocationID": sublocation.ID,
					})
				}

				// Invalidate physical location cache
				if err := gss.cacheWrapper.InvalidateLocationCache(ctx, userID, sublocation.PhysicalLocationID); err != nil {
					gss.logger.Error("Failed to invalidate physical location cache", map[string]any{
						"error": err,
						"physicalLocationID": sublocation.PhysicalLocationID,
					})
				}

				// Force refresh physical location cache
				if gss.physicalService != nil {
					_, refreshErr := gss.physicalService.GetSinglePhysicalLocation(ctx, userID, sublocation.PhysicalLocationID)
					if refreshErr != nil {
						gss.logger.Warn("Failed to refresh physical location cache", map[string]any{
							"error": refreshErr,
							"physicalLocationID": sublocation.PhysicalLocationID,
						})
					}
				}
			}
		}
	}

	gss.logger.Debug("RemoveGame completed successfully", map[string]any{
		"userID": userID,
		"userGameID": req.UserGameID,
	})

	return nil
}