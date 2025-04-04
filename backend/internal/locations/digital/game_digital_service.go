package digital

import (
	"context"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/infrastructure/cache"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
)

type GameDigitalService struct {
	dbAdapter    interfaces.DigitalDbAdapter
	config       *config.Config
	cacheWrapper interfaces.DigitalCacheWrapper
	logger       interfaces.Logger
	sanitizer    interfaces.Sanitizer
	validator    interfaces.DigitalValidator
}
type DigitalService interface {
	GetDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error
	DeleteDigitalLocation(ctx context.Context, userID string, locationID string) error
	UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error
}

func NewGameDigitalService(appContext *appcontext.AppContext) (*GameDigitalService, error) {
	// Create + initialize db adapter
	dbAdapter, err := NewDigitalDbAdapter(appContext)
	if err != nil {
		appContext.Logger.Error("Failed to create dbAdapter", map[string]any{"error": err})
		return nil, err
	}
	appContext.Logger.Info("Digital dbAdapter created successfully", nil)

	// Create sanitizer to feed into validator
	sanitizer, err := security.NewSanitizer()
	if err != nil {
		appContext.Logger.Error("Failed to create sanitizer", map[string]any{"error": err})
		return nil, err
	}
	appContext.Logger.Info("Sanitizer created successfully", nil)

	// Create validator
	validator, err := NewDigitalValidator(sanitizer)
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

	// Create digital media cache adapter
	digitalCacheAdapter, err := NewDigitalCacheAdapter(cacheWrapper)
	if err != nil {
		return nil, err
	}

	// Sanity check that all deps are initialized
	appContext.Logger.Info("GameDigitalService dependencies intialized", map[string]any{
		"dbAdapter": dbAdapter,
		"validator": validator,
		"cacheAdapter": digitalCacheAdapter,
	})

	return &GameDigitalService{
		dbAdapter:     dbAdapter,
		validator:     validator,
		logger:        appContext.Logger,
		config:        appContext.Config,
		cacheWrapper:  digitalCacheAdapter,
		sanitizer:     sanitizer,
	}, nil
}

// GET
func (gds *GameDigitalService) GetDigitalLocations(
	ctx context.Context,
	userID string,
) ([]models.DigitalLocation, error) {
	// Attempt to get cached locations
	cachedLocations, err := gds.cacheWrapper.GetCachedDigitalLocations(ctx, userID)
	if err == nil && cachedLocations != nil {
		gds.logger.Debug("Cache hit for digital locations", map[string]any{"userID": userID})
		return cachedLocations, nil
	}

	// Cache miss or error, get from DB
	gds.logger.Debug("Cache miss for digital locations, fetching from DB", map[string]any{"userID": userID})
	locations, err := gds.dbAdapter.GetUserDigitalLocations(ctx, userID)
	if err != nil {
		gds.logger.Error("Failed to fetch digital locations from DB", map[string]any{"error": err})
		return nil, err
	}

	// Cache the results for future requests
	if cacheErr := gds.cacheWrapper.SetCachedDigitalLocations(
		ctx,
		userID,
		locations,
	); cacheErr != nil {
		gds.logger.Error("Failed to cache digital locations", map[string]any{"error": cacheErr})
		// Continue w/ returning the locations from DB
	}

	return locations, nil
}

func (gds *GameDigitalService) GetDigitalLocation(
	ctx context.Context,
	userID string,
	locationID string,
) (*models.DigitalLocation, error) {
	// Try to get from cache
	cachedLocation, found, err := gds.cacheWrapper.GetSingleCachedDigitalLocation(
		ctx,
		userID,
		locationID,
	)
	if err == nil && found {
		gds.logger.Debug("Cache hit for digital location", map[string]any{
			"userID":     userID,
			"locationID": locationID,
		})
		return cachedLocation, nil
	}

	// Cache miss or error, get from DB
	gds.logger.Debug("Cache miss for digital location, fetching from DB", map[string]any{
		"userID":     userID,
		"locationID": locationID,
	})

	location, err := gds.dbAdapter.GetDigitalLocation(
		ctx,
		userID, locationID,
	)
	if err != nil {
		gds.logger.Error("Failed to fetch digital location from DB", map[string]any{"error": err})
		return nil, err
	}

	// Cache the result for future requests
	if cacheErr := gds.cacheWrapper.SetSingleCachedDigitalLocation(
		ctx,
		userID,
		location,
	); cacheErr != nil {
		gds.logger.Error("Failed to cache digital location", map[string]any{"error": cacheErr})
		// Continue returning the location from DB
	}

	return &location, nil
}

// POST
func (gds *GameDigitalService) AddDigitalLocation(
	ctx context.Context,
	userID string,
	location models.DigitalLocation,
) error {
	// Call the existing implementation and discard the first return value
	_, err := gds.addDigitalLocationImpl(ctx, userID, location)
	return err
}

func (gds *GameDigitalService) addDigitalLocationImpl(
	ctx context.Context,
	userID string,
	location models.DigitalLocation,
) (*models.DigitalLocation, error) {
	// Validate location
	validatedLocation, err := gds.validator.ValidateDigitalLocation(location)
	if err != nil {
		gds.logger.Error("Location validation failed", map[string]any{"error": err})
		return nil, err
	}

	// Add to db
	createdLocation, err := gds.dbAdapter.AddDigitalLocation(ctx, userID, validatedLocation)
	if err != nil {
		gds.logger.Error("Failed to add digital location to DB", map[string]any{"error": err})
		return nil, err
	}

	// Invalidate user cache
	if err := gds.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gds.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}

	return &createdLocation, nil
}

// PUT
func (gds *GameDigitalService) UpdateDigitalLocation(
	ctx context.Context,
	userID string,
	location models.DigitalLocation,
) error {
	// Validate location
	validatedLocation, err := gds.validator.ValidateDigitalLocation(location)
	if err != nil {
		gds.logger.Error("Location validation failed", map[string]any{"error": err})
		return err
	}

	// Update in db
	if err := gds.dbAdapter.UpdateDigitalLocation(
		ctx,
		userID,
		validatedLocation,
	); err != nil {
		gds.logger.Error("Failed to update digital location in DB", map[string]any{"error": err})
		return err
	}

	// Invalidate caches
	if err := gds.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gds.logger.Error("failed to invalidate user cache", map[string]any{"error": err})
	}
	if err := gds.cacheWrapper.InvalidateDigitalLocationCache(ctx, userID, validatedLocation.ID); err != nil {
		gds.logger.Error("failed to invalidate location cache", map[string]any{"error": err})
	}

	return nil
}

// DELETE
func (gds *GameDigitalService) DeleteDigitalLocation(
	ctx context.Context,
	userID string,
	locationID string,
) error {
	// Remove from database
	if err := gds.dbAdapter.RemoveDigitalLocation(
		ctx,
		userID,
		locationID,
	); err != nil {
		gds.logger.Error("Failed to delete digital location from DB", map[string]any{"error": err})
		return err
	}

	// Invalidate caches
	if err := gds.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gds.logger.Error("failed to invalidate user cache", map[string]any{"error": err})
	}
	if err := gds.cacheWrapper.InvalidateDigitalLocationCache(ctx, userID, locationID); err != nil {
		gds.logger.Error("Failed to invalidate location cache", map[string]any{"error": err})
	}

	return nil
}