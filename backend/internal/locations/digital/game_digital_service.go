package digital

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

type GameDigitalService struct {
	dbAdapter    interfaces.DigitalDbAdapter
	config       *config.Config
	cacheWrapper interfaces.DigitalCacheWrapper
	logger       interfaces.Logger
	sanitizer    interfaces.Sanitizer
	validator    interfaces.DigitalValidator
}

type DigitalService interface {
	GetUserDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error)
	GetDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error)
	FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error)
	AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error)
	UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error
	RemoveDigitalLocation(ctx context.Context, userID, locationID string) error
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
func (gds *GameDigitalService) GetUserDigitalLocations(ctx context.Context, userID string) ([]models.DigitalLocation, error) {
	gds.logger.Debug("Getting user digital locations", map[string]any{"userID": userID})
	return gds.dbAdapter.GetUserDigitalLocations(ctx, userID)
}

func (gds *GameDigitalService) GetDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
	return gds.dbAdapter.GetDigitalLocation(ctx, userID, locationID)
}

func (gds *GameDigitalService) FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
	return gds.dbAdapter.FindDigitalLocationByName(ctx, userID, name)
}

// POST
func (gds *GameDigitalService) AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
	// Validate the location
	validatedLocation, err := gds.validator.ValidateDigitalLocation(location)
	if err != nil {
		return models.DigitalLocation{}, fmt.Errorf("validation failed: %w", err)
	}

	// Add the location
	return gds.dbAdapter.AddDigitalLocation(ctx, userID, validatedLocation)
}

// PUT
func (gds *GameDigitalService) UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error {
	// Validate the location
	validatedLocation, err := gds.validator.ValidateDigitalLocation(location)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update the location
	return gds.dbAdapter.UpdateDigitalLocation(ctx, userID, validatedLocation)
}

// DELETE
func (gds *GameDigitalService) RemoveDigitalLocation(ctx context.Context, userID, locationID string) error {
	gds.logger.Debug("Removing digital location", map[string]any{"userID": userID, "locationID": locationID})
	return gds.dbAdapter.RemoveDigitalLocation(ctx, userID, locationID)
}