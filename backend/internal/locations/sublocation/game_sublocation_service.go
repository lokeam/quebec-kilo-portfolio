package sublocation

import (
	"context"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/infrastructure/cache"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
)

type GameSublocationService struct {
	dbAdapter      interfaces.SublocationDbAdapter
	config         *config.Config
	cacheWrapper   interfaces.SublocationCacheWrapper
	logger         interfaces.Logger
	sanitizer      interfaces.Sanitizer
	validator      interfaces.SublocationValidator
}

type SublocationService interface {
	GetSublocations(ctx context.Context, userID string) ([]models.Sublocation, error)
	AddSublocation(ctx context.Context, userID string, sublocation models.Sublocation) error
	DeleteSublocation(ctx context.Context, userID string, sublocationID string) error
	UpdateSublocation(ctx context.Context, userID string, sublocation models.Sublocation) error
}

func NewGameSublocationService(appContext *appcontext.AppContext) (*GameSublocationService, error) {
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
	validator, err := NewSublocationValidator(sanitizer)
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
	})

	// Sanity check that all dependencies are intialized

	return &GameSublocationService{
		dbAdapter:      dbAdapter,
		validator:      validator,
		logger:         appContext.Logger,
		config:         appContext.Config,
		cacheWrapper:   sublocationCacheAdapter,
		sanitizer:      sanitizer,
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
func (gss *GameSublocationService) AddSublocation(
	ctx context.Context,
	userID string,
	sublocation models.Sublocation,
) error {
	// Call existing implementation, discard first return value
	_, err := gss.addSublocationImpl(ctx, userID, sublocation)
	return err
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

	// Invalidate user cache
	if err := gss.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gss.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}

	return &createdSublocation, nil
}

// PUT
func (gss *GameSublocationService) UpdateSublocation(
	ctx context.Context,
	userID string,
	sublocation models.Sublocation,
) error {
	// Validate sublocation
	validatedSublocation, err := gss.validator.ValidateSublocation(sublocation)
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

	// Invalidate caches
	if err := gss.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
		gss.logger.Error("Failed to invalidate user cache", map[string]any{"error": err})
	}
	if err := gss.cacheWrapper.InvalidateSublocationCache(ctx, userID, validatedSublocation.ID); err != nil {
		gss.logger.Error("Failed to invalidate location cache", map[string]any{"error": err})
	}

	return nil
}

// DELETE
func (gss *GameSublocationService) DeleteSublocation(
	ctx context.Context,
	userID string,
	sublocationID string,
) error {
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

	return nil
}