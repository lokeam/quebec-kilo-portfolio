package sublocation

import (
	"context"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/infrastructure/cache"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/services"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
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
	GetSublocation(ctx context.Context, userID string, sublocationID string) (models.Sublocation, error)
	AddSublocation(ctx context.Context, userID string, sublocation models.Sublocation) (models.Sublocation, error)
	DeleteSublocation(ctx context.Context, userID string, sublocationID string) error
	UpdateSublocation(ctx context.Context, userID string, sublocation models.Sublocation) error
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

	// Create validator
	validator, err := NewSublocationValidator(sanitizer, dbAdapter)
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

	//Create cache wrapper to handle Redis caching
	sublocationCacheAdapter, err := NewSublocationCacheAdapter(cacheWrapper)
	if err != nil {
		return nil, err
	}

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
func (gss *GameSublocationService) GetSublocation(
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
	return gss.dbAdapter.GetSublocation(ctx, userID, sublocationID)
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
	sublocations, err := gss.dbAdapter.GetUserSublocations(ctx, userID)
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
// AddSublocation adds a new sublocation for a user.
//
// Cache Invalidation Strategy:
// - The private addSublocationImpl method invalidates the user cache to ensure cache consistency
//   even if called directly (defensive programming).
// - The public AddSublocation method invalidates both the user cache and the parent physical
//   location's cache to maintain consistency across related caches.
//
// IMPORTANT: The double cache invalidation is INTENTIONAL and NECESSARY:
// 1. addSublocationImpl handles basic user cache invalidation for consistency
// 2. AddSublocation handles both user cache and physical location cache invalidation
//
// DO NOT attempt to "optimize" this by removing either cache invalidation:
// - Removing addSublocationImpl's invalidation would break cache consistency if called directly
// - Removing AddSublocation's invalidation would miss the physical location cache update
//
// This pattern is consistent with our production codebase and follows defensive programming
// principles to ensure cache consistency across all possible code paths.
func (gss *GameSublocationService) AddSublocation(
	ctx context.Context,
	userID string,
	sublocation models.Sublocation,
) (models.Sublocation, error) {
	// Call existing implementation
	createdSublocation, err := gss.addSublocationImpl(ctx, userID, sublocation)
	if err != nil {
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
		_, refreshErr := gss.physicalService.GetPhysicalLocation(ctx, userID, physicalLocationID)
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

	return *createdSublocation, nil
}

func (gss *GameSublocationService) addSublocationImpl(
	ctx context.Context,
	userID string,
	location models.Sublocation,
) (*models.Sublocation, error) {
	// Validate sublocation
	validatedLocation, err := gss.validator.ValidateSublocation(location)
	if err != nil {
		gss.logger.Error("Location validation failed", map[string]any{"error": err})
		return nil, err
	}

	// Add to db
	createdSublocation, err := gss.dbAdapter.AddSublocation(ctx, userID, validatedLocation)
	if err != nil {
		gss.logger.Error("Failed to add sublocation to DB", map[string]any{"error": err})
		return nil, err
	}

	// Cache invalidation is handled by the public AddSublocation method
	return &createdSublocation, nil
}

// PUT
func (gss *GameSublocationService) UpdateSublocation(
	ctx context.Context,
	userID string,
	sublocation models.Sublocation,
) error {
	// First get the existing sublocation
	existingSublocation, err := gss.dbAdapter.GetSublocation(ctx, userID, sublocation.ID)
	if err != nil {
		gss.logger.Error("Failed to get existing sublocation", map[string]any{"error": err})
		return err
	}

	// Validate only the fields that are being updated
	validatedSublocation, err := gss.validator.ValidateSublocationUpdate(sublocation, existingSublocation)
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
		_, refreshErr := gss.physicalService.GetPhysicalLocation(ctx, userID, physicalLocationID)
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
func (gss *GameSublocationService) DeleteSublocation(
	ctx context.Context,
	userID string,
	sublocationID string,
) error {
	// First get the sublocation to get its physical location ID
	sublocation, err := gss.dbAdapter.GetSublocation(ctx, userID, sublocationID)
	if err != nil {
		gss.logger.Error("Failed to get sublocation before deletion", map[string]any{"error": err})
		return err
	}

	// Store the physical location ID before deletion
	physicalLocationID := sublocation.PhysicalLocationID

	// Remove from database
	if err := gss.dbAdapter.RemoveSublocation(
		ctx,
		userID,
		sublocationID,
	); err != nil {
		gss.logger.Error("Failed to delete sublocation from DB", map[string]any{"error": err})
		return err
	}

	// Invalidate caches
	if err := gss.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gss.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}
	if err := gss.cacheWrapper.InvalidateSublocationCache(ctx, userID, sublocationID); err != nil {
		gss.logger.Error("Failed to invalidate sublocation cache", map[string]any{"error": err})
	}

	// Also invalidate the parent physical location's cache
	if err := gss.cacheWrapper.InvalidateLocationCache(ctx, userID, physicalLocationID); err != nil {
		gss.logger.Error("Failed to invalidate parent physical location cache", map[string]any{"error": err})
	}

	// Force immediate refresh of physical location cache
	gss.logger.Debug("Forcing refresh of physical location cache after sublocation deletion", map[string]any{
		"userID": userID,
		"physicalLocationID": physicalLocationID,
	})

	// Actively refresh the physical location data if service is available
	if gss.physicalService != nil {
		_, refreshErr := gss.physicalService.GetPhysicalLocation(ctx, userID, physicalLocationID)
		if refreshErr != nil {
			gss.logger.Warn("Failed to refresh physical location cache after deletion", map[string]any{
				"error": refreshErr,
				"physicalLocationID": physicalLocationID,
			})
		} else {
			gss.logger.Debug("Successfully refreshed physical location cache after deletion", map[string]any{
				"physicalLocationID": physicalLocationID,
			})
		}
	} else {
		gss.logger.Warn("Physical service not available for cache refresh after deletion", nil)
	}

	return nil
}