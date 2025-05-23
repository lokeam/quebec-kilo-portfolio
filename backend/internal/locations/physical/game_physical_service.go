package physical

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/infrastructure/cache"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
)

type GamePhysicalService struct {
	dbAdapter      interfaces.PhysicalDbAdapter
	config         *config.Config
	cacheWrapper   interfaces.PhysicalCacheWrapper
	logger         interfaces.Logger
	sanitizer      interfaces.Sanitizer
	validator      interfaces.PhysicalValidator
}

type PhysicalService interface {
	GetUserPhysicalLocations(ctx context.Context, userID string) ([]models.PhysicalLocation, error)
	GetPhysicalLocation(ctx context.Context, userID, locationID string) (models.PhysicalLocation, error)
	AddPhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error)
	DeletePhysicalLocation(ctx context.Context, userID, locationID string) error
}

func NewGamePhysicalService(appContext *appcontext.AppContext) (*GamePhysicalService, error) {
	// Create + initialize db adapter
	dbAdapter, err := NewPhysicalDbAdapter(appContext)
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
	validator, err := NewPhysicalValidator(sanitizer)
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

	// Create physical media cache adapter
	physicalCacheAdapter, err := NewPhysicalCacheAdapter(cacheWrapper)
	if err != nil {
		return nil, err
	}

	// Sanity check that all dependencies are intialized
	appContext.Logger.Info("GamePhysicalService components initialized", map[string]any{
		"dbAdapter": dbAdapter != nil,
		"validator": validator != nil,
		"cacheWrapper": physicalCacheAdapter != nil,
	})

	return &GamePhysicalService{
		dbAdapter:       dbAdapter,
		validator:       validator,
		logger:          appContext.Logger,
		config:          appContext.Config,
		cacheWrapper:    physicalCacheAdapter,
		sanitizer:       sanitizer,
	}, nil
}

// GET
func (gps *GamePhysicalService) GetUserPhysicalLocations(
	ctx context.Context,
	userID string,
) ([]models.PhysicalLocation, error) {
	// Try to get from cache first
	cachedLocations, err := gps.cacheWrapper.GetCachedPhysicalLocations(ctx, userID)
	if err == nil && cachedLocations != nil {
		gps.logger.Debug("Cache hit for physical locations", map[string]any{"userID": userID})
		return cachedLocations, nil
	}

	// Cache miss or error, get from DB
	gps.logger.Debug("Cache miss for physical locations, fetching from DB", map[string]any{"userID": userID})
	locations, err := gps.dbAdapter.GetUserPhysicalLocations(ctx, userID)
	if err != nil {
		gps.logger.Error("Failed to fetch physical locations from DB", map[string]any{"error": err})
		return nil, err
	}

	// Cache the results
	if err := gps.cacheWrapper.SetCachedPhysicalLocations(ctx, userID, locations); err != nil {
		gps.logger.Error("Failed to cache physical locations", map[string]any{"error": err})
		// Continue with returning locations from DB
	}

	return locations, nil
}

func (gps *GamePhysicalService) GetPhysicalLocation(
	ctx context.Context,
	userID string,
	locationID string,
) (models.PhysicalLocation, error) {
	// Try to get from cache
	cachedLocation, found, err := gps.cacheWrapper.GetSingleCachedPhysicalLocation(ctx, userID, locationID)
	if err == nil && found {
		gps.logger.Debug("Cache hit for physical location", map[string]any{
			"userID": userID,
			"locationID": locationID,
		})
		return *cachedLocation, nil
	}

	// Cache miss or error, get from DB
	gps.logger.Debug("Cache miss for physical location, fetching from DB", map[string]any{
		"userID": userID,
		"locationID": locationID,
	})

	location, err := gps.dbAdapter.GetPhysicalLocation(ctx, userID, locationID)
	if err != nil {
		gps.logger.Error("Failed to fetch physical location from DB", map[string]any{"error": err})
		return models.PhysicalLocation{}, err
	}

	// Cache the location
	if err := gps.cacheWrapper.SetSingleCachedPhysicalLocation(ctx, userID, location); err != nil {
		gps.logger.Error("Failed to cache physical location", map[string]any{
			"error": err,
			"userID": userID,
			"locationID": location.ID,
		})
	}

	return location, nil
}

// POST
func (gps *GamePhysicalService) AddPhysicalLocation(
	ctx context.Context,
	userID string,
	location models.PhysicalLocation,
) (models.PhysicalLocation, error) {
	// Validate location
	validatedLocation, err := gps.validator.ValidatePhysicalLocation(location)
	if err != nil {
		gps.logger.Error("Location validation failed", map[string]any{"error": err})
		return models.PhysicalLocation{}, err
	}

	// Add to db
	createdLocation, err := gps.dbAdapter.AddPhysicalLocation(ctx, userID, validatedLocation)
	if err != nil {
		gps.logger.Error("Failed to add physical location to DB", map[string]any{"error": err})
		return models.PhysicalLocation{}, err
	}

	// Invalidate the user's locations cache
	if err := gps.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gps.logger.Error("Failed to invalidate user locations cache", map[string]any{
			"error": err,
			"userID": userID,
		})
		// Continue despite error, since the DB update was successful
	}

	return createdLocation, nil
}

// UPDATE
func (gps *GamePhysicalService) UpdatePhysicalLocation(ctx context.Context, userID string, location models.PhysicalLocation) (models.PhysicalLocation, error) {
	gps.logger.Debug("Updating physical location", map[string]any{
		"userID": userID,
		"location": location,
	})

	// Update in database
	updatedLocation, err := gps.dbAdapter.UpdatePhysicalLocation(ctx, userID, location)
	if err != nil {
		return models.PhysicalLocation{}, fmt.Errorf("failed to update physical location: %w", err)
	}

	// Invalidate both the user's locations cache and the specific location cache
	if err := gps.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gps.logger.Error("Failed to invalidate user locations cache", map[string]any{
			"error": err,
			"userID": userID,
		})
	}

	if err := gps.cacheWrapper.InvalidateLocationCache(ctx, userID, location.ID); err != nil {
		gps.logger.Error("Failed to invalidate location cache", map[string]any{
			"error": err,
			"userID": userID,
			"locationID": location.ID,
		})
	}

	return updatedLocation, nil
}

// DELETE
func (gps *GamePhysicalService) DeletePhysicalLocation(
	ctx context.Context,
	userID string,
	locationID string,
) error {
	// Remove from database
	if err := gps.dbAdapter.RemovePhysicalLocation(ctx, userID, locationID); err != nil {
		gps.logger.Error("Failed to delete physical location from DB", map[string]any{"error": err})
		return err
	}

	// Invalidate both the user's locations cache and the specific location cache
	if err := gps.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gps.logger.Error("Failed to invalidate user locations cache", map[string]any{
			"error": err,
			"userID": userID,
		})
	}

	if err := gps.cacheWrapper.InvalidateLocationCache(ctx, userID, locationID); err != nil {
		gps.logger.Error("Failed to invalidate location cache", map[string]any{
			"error": err,
			"userID": userID,
			"locationID": locationID,
		})
	}

	return nil
}