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
	RemoveDigitalLocation(ctx context.Context, userID string, locationIDs []string) (int64, error)
	RemoveDigitalLocationsBulk(ctx context.Context, userID string, locationIDs []string) (int64, error)

	// Game Management Operations
	AddGameToDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error
	RemoveGameFromDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error
	GetGamesByDigitalLocationID(ctx context.Context, userID string, locationID string) ([]models.Game, error)

	// Subscription management
	GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error)
	AddSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, subscription models.Subscription) error
	RemoveSubscription(ctx context.Context, locationID string) error

	// Payment management
	GetPayments(ctx context.Context, locationID string) ([]models.Payment, error)
	AddPayment(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetPayment(ctx context.Context, paymentID int64) (*models.Payment, error)
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
	// Try to get from cache first
	cachedLocations, cacheErr := gds.cacheWrapper.GetCachedDigitalLocations(ctx, userID)
	if cacheErr == nil && cachedLocations != nil && len(cachedLocations) > 0 {
			gds.logger.Debug("Cache hit for digital locations", map[string]any{
					"userID": userID,
					"count": len(cachedLocations),
			})
			return cachedLocations, nil
	}

	gds.logger.Debug("Cache miss for digital locations, fetching from DB", map[string]any{
			"userID": userID,
			"error": cacheErr,
	})

	// Get fresh data from DB
	locations, err := gds.dbAdapter.GetUserDigitalLocations(ctx, userID)
	if err != nil {
			gds.logger.Error("Failed to fetch digital locations from DB", map[string]any{"error": err})
			return nil, err
	}

	// Cache the results
	if cacheErr := gds.cacheWrapper.SetCachedDigitalLocations(ctx, userID, locations); cacheErr != nil {
			gds.logger.Error("Failed to cache digital locations", map[string]any{"error": cacheErr})
	}

	return locations, nil
}

func (gds *GameDigitalService) GetDigitalLocation(ctx context.Context, userID, locationID string) (models.DigitalLocation, error) {
	// Try to get from cache
	cachedLocation, found, err := gds.cacheWrapper.GetSingleCachedDigitalLocation(ctx, userID, locationID)
	if err == nil && found {
		gds.logger.Debug("Cache hit for digital location", map[string]any{
			"userID": userID,
			"locationID": locationID,
		})
		return *cachedLocation, nil
	}

	// Cache miss or error, get from DB
	gds.logger.Debug("Cache miss for digital location, fetching from DB", map[string]any{
		"userID": userID,
		"locationID": locationID,
	})

	location, err := gds.dbAdapter.GetDigitalLocation(ctx, userID, locationID)
	if err != nil {
		gds.logger.Error("Failed to fetch digital location from DB", map[string]any{"error": err})
		return models.DigitalLocation{}, err
	}

	// Cache the location
	if err := gds.cacheWrapper.SetSingleCachedDigitalLocation(ctx, userID, location); err != nil {
		gds.logger.Error("Failed to cache digital location", map[string]any{
			"error": err,
			"userID": userID,
			"locationID": location.ID,
		})
	}

	return location, nil
}

func (gds *GameDigitalService) FindDigitalLocationByName(ctx context.Context, userID string, name string) (models.DigitalLocation, error) {
	// This method doesn't implement caching as it's a search operation
	// that could have varying results based on the search term
	return gds.dbAdapter.FindDigitalLocationByName(ctx, userID, name)
}

// POST
// POST
func (gds *GameDigitalService) AddDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) (models.DigitalLocation, error) {
	// Add detailed logging
	gds.logger.Debug("Adding digital location with is_active", map[string]any{
			"userID": userID,
			"original_is_active": location.IsActive,
	})

	// Store original is_active value and subscription data
	originalIsActive := location.IsActive
	originalSubscription := location.Subscription

	// Validate the location
	validatedLocation, err := gds.validator.ValidateDigitalLocation(location)
	if err != nil {
			return models.DigitalLocation{}, fmt.Errorf("validation failed: %w", err)
	}

	// DEBUG: Log after validation
	gds.logger.Debug("After validation", map[string]any{
			"validated_is_active": validatedLocation.IsActive,
			"original_is_active": originalIsActive,
	})

	// CRITICAL FIX: Ensure is_active value is preserved
	validatedLocation.IsActive = originalIsActive

	// CRITICAL: Restore subscription data if it was provided
	validatedLocation.Subscription = originalSubscription

	createdLocation, err := gds.dbAdapter.AddDigitalLocation(ctx, userID, validatedLocation)
	if err != nil {
			gds.logger.Error("Failed to add digital location to DB", map[string]any{"error": err})
			return models.DigitalLocation{}, err
	}

	// CRITICAL: Save subscription data if provided
	if originalSubscription != nil {
			// Update the location ID in the subscription
			originalSubscription.LocationID = createdLocation.ID

			// Add the subscription to the database
			subscription, err := gds.dbAdapter.AddSubscription(ctx, *originalSubscription)
			if err != nil {
					gds.logger.Error("Failed to add subscription data", map[string]any{
							"error": err,
							"locationID": createdLocation.ID,
					})
					// Continue without subscription, just log the error
			} else {
					// Add the subscription to the created location
					createdLocation.Subscription = subscription
			}
	}

	// Invalidate the cache for this user
	if err := gds.cacheWrapper.InvalidateUserCache(ctx, userID); err != nil {
			gds.logger.Error("Failed to invalidate user cache after adding location", map[string]any{
					"error": err,
					"userID": userID,
			})
			// Continue despite error, since the DB update was successful
	}

	return createdLocation, nil
}

// PUT
func (gds *GameDigitalService) UpdateDigitalLocation(ctx context.Context, userID string, location models.DigitalLocation) error {
	gds.logger.Debug("Updating digital location", map[string]any{
		"userID": userID,
		"location": location,
		"is_active_value": location.IsActive,
	})

	// CRITICAL FIX: Invalidate the cache BEFORE doing anything else
	if err := gds.cacheWrapper.InvalidateDigitalLocationCache(ctx, userID, location.ID); err != nil {
		gds.logger.Error("Failed to invalidate location cache before update", map[string]any{
			"error": err,
			"userID": userID,
			"locationID": location.ID,
		})
		// Continue despite error, since we're going to update the DB anyway
	}

	// Fill in any missing fields from the DB if needed
	if location.UserID == "" {
		location.UserID = userID
	}

	// DEBUG: Log state before validation
	gds.logger.Debug("Before validation", map[string]any{
		"location_is_active": location.IsActive,
	})

	// Validate the location
	validatedLocation, err := gds.validator.ValidateDigitalLocation(location)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// DEBUG: Log state after validation
	gds.logger.Debug("After validation", map[string]any{
		"validatedLocation_is_active": validatedLocation.IsActive,
	})

	// CRITICAL FIX: Ensure is_active value from request is preserved
	validatedLocation.IsActive = location.IsActive

	// Update in database
	if err := gds.dbAdapter.UpdateDigitalLocation(ctx, userID, validatedLocation); err != nil {
		if err.Error() == "digital location not found" {
			return fmt.Errorf("digital location not found")
		}
		return fmt.Errorf("failed to update digital location: %w", err)
	}

	// Get the updated location to ensure we have the correct state after the update
	updatedLocation, err := gds.dbAdapter.GetDigitalLocation(ctx, userID, location.ID)
	if err != nil {
		gds.logger.Error("Failed to get updated location after update", map[string]any{"error": err})
		// Continue despite error, since the DB update was successful
	} else {
		// Update single location cache with the fresh data from DB
		if err := gds.cacheWrapper.SetSingleCachedDigitalLocation(ctx, userID, updatedLocation); err != nil {
			gds.logger.Error("Failed to update digital location in cache", map[string]any{
				"error": err,
				"userID": userID,
				"location": updatedLocation,
			})
		}
	}

	// Update the list cache as well to ensure consistency
	locations, err := gds.dbAdapter.GetUserDigitalLocations(ctx, userID)
	if err != nil {
		gds.logger.Error("Failed to get updated locations after update", map[string]any{"error": err})
		// As a fallback, invalidate the user's cache to force a refresh on next list request
		if invalidateErr := gds.cacheWrapper.InvalidateUserCache(ctx, userID); invalidateErr != nil {
			gds.logger.Error("Failed to invalidate user cache", map[string]any{"error": invalidateErr})
		}
	} else {
		// Update the cache with all locations
		if err := gds.cacheWrapper.SetCachedDigitalLocations(ctx, userID, locations); err != nil {
			gds.logger.Error("Failed to update locations cache after location update", map[string]any{"error": err})
		}
	}

	gds.logger.Debug("UpdateDigitalLocation success", map[string]any{
		"location": validatedLocation,
	})

	return nil
}

// RemoveDigitalLocation removes one or more digital locations for a user.
// It handles both single and bulk deletion operations.
func (gds *GameDigitalService) RemoveDigitalLocation(
	ctx context.Context,
	userID string,
	locationIDs []string,
) (int64, error) {
	gds.logger.Debug("RemoveDigitalLocation called", map[string]any{
		"userID":      userID,
		"locationIDs": locationIDs,
		"isBulk":      len(locationIDs) > 1,
	})

	// Validate input and get sanitized, deduplicated IDs
	validatedIDs, err := gds.validator.ValidateRemoveDigitalLocation(userID, locationIDs)
	if err != nil {
		gds.logger.Error("Validation failed for RemoveDigitalLocation", map[string]any{
			"error": err,
			"userID": userID,
			"locationIDs": locationIDs,
		})
		return 0, fmt.Errorf("validation failed: %w", err)
	}

	// Remove locations from database using validated IDs
	count, err := gds.dbAdapter.RemoveDigitalLocation(ctx, userID, validatedIDs)
	if err != nil {
		gds.logger.Error("Failed to remove digital locations from database", map[string]any{
			"error": err,
			"userID": userID,
			"locationIDs": validatedIDs,
		})
		return 0, fmt.Errorf("failed to remove digital locations: %w", err)
	}

	// Invalidate cache for each location
	for _, locationID := range validatedIDs {
		if err := gds.cacheWrapper.InvalidateDigitalLocationCache(ctx, userID, locationID); err != nil {
			gds.logger.Error("Failed to invalidate cache for location", map[string]any{
				"error": err,
				"userID": userID,
				"locationID": locationID,
			})
			// Continue with other invalidations even if one fails
		}
	}

	gds.logger.Debug("RemoveDigitalLocation completed successfully", map[string]any{
		"userID": userID,
		"locationIDs": validatedIDs,
		"deletedCount": count,
	})

	return count, nil
}

// Subscription management
func (gds *GameDigitalService) GetSubscription(ctx context.Context, locationID string) (*models.Subscription, error) {
	// Try to get from cache first
	cachedSubscription, found, err := gds.cacheWrapper.GetCachedSubscription(ctx, locationID)
	if err == nil && found {
		gds.logger.Debug("Cache hit for subscription", map[string]any{
			"locationID": locationID,
		})
		return cachedSubscription, nil
	}

	// Cache miss or error, get from DB
	gds.logger.Debug("Cache miss for subscription, fetching from DB", map[string]any{
		"locationID": locationID,
	})

	subscription, err := gds.dbAdapter.GetSubscription(ctx, locationID)
	if err != nil {
		gds.logger.Error("Failed to fetch subscription from DB", map[string]any{"error": err})
		return nil, err
	}

	// Cache the subscription
	if err := gds.cacheWrapper.SetCachedSubscription(ctx, locationID, *subscription); err != nil {
		gds.logger.Error("Failed to cache subscription", map[string]any{
			"error": err,
			"locationID": locationID,
		})
	}

	return subscription, nil
}

func (gds *GameDigitalService) AddSubscription(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
	// Add to DB
	result, err := gds.dbAdapter.AddSubscription(ctx, subscription)
	if err != nil {
		gds.logger.Error("Failed to add subscription to DB", map[string]any{"error": err})
		return nil, err
	}

	// Cache the new subscription
	if err := gds.cacheWrapper.SetCachedSubscription(ctx, subscription.LocationID, *result); err != nil {
		gds.logger.Error("Failed to cache new subscription", map[string]any{
			"error": err,
			"locationID": subscription.LocationID,
		})
	}

	return result, nil
}

func (gds *GameDigitalService) UpdateSubscription(ctx context.Context, subscription models.Subscription) error {
	// Update in DB
	err := gds.dbAdapter.UpdateSubscription(ctx, subscription)
	if err != nil {
		gds.logger.Error("Failed to update subscription in DB", map[string]any{"error": err})
		return err
	}

	// Invalidate cache
	if err := gds.cacheWrapper.InvalidateSubscriptionCache(ctx, subscription.LocationID); err != nil {
		gds.logger.Error("Failed to invalidate subscription cache", map[string]any{
			"error": err,
			"locationID": subscription.LocationID,
		})
	}

	return nil
}

func (gds *GameDigitalService) RemoveSubscription(ctx context.Context, locationID string) error {
	// Remove from DB
	err := gds.dbAdapter.RemoveSubscription(ctx, locationID)
	if err != nil {
		gds.logger.Error("Failed to remove subscription from DB", map[string]any{"error": err})
		return err
	}

	// Invalidate cache
	if err := gds.cacheWrapper.InvalidateSubscriptionCache(ctx, locationID); err != nil {
		gds.logger.Error("Failed to invalidate subscription cache", map[string]any{
			"error": err,
			"locationID": locationID,
		})
	}

	return nil
}

// Payment management
func (gds *GameDigitalService) GetPayments(ctx context.Context, locationID string) ([]models.Payment, error) {
	// Try to get from cache first
	cachedPayments, err := gds.cacheWrapper.GetCachedPayments(ctx, locationID)
	if err == nil {
		gds.logger.Debug("Cache hit for payments", map[string]any{
			"locationID": locationID,
		})
		return cachedPayments, nil
	}

	// Cache miss or error, get from DB
	gds.logger.Debug("Cache miss for payments, fetching from DB", map[string]any{
		"locationID": locationID,
	})

	payments, err := gds.dbAdapter.GetPayments(ctx, locationID)
	if err != nil {
		gds.logger.Error("Failed to fetch payments from DB", map[string]any{"error": err})
		return nil, err
	}

	// Cache the payments
	if err := gds.cacheWrapper.SetCachedPayments(ctx, locationID, payments); err != nil {
		gds.logger.Error("Failed to cache payments", map[string]any{
			"error": err,
			"locationID": locationID,
		})
	}

	return payments, nil
}

func (gds *GameDigitalService) AddPayment(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	// Add to DB
	result, err := gds.dbAdapter.AddPayment(ctx, payment)
	if err != nil {
		gds.logger.Error("Failed to add payment to DB", map[string]any{"error": err})
		return nil, err
	}

	// Invalidate payments cache
	if err := gds.cacheWrapper.InvalidatePaymentsCache(ctx, payment.LocationID); err != nil {
		gds.logger.Error("Failed to invalidate payments cache", map[string]any{
			"error": err,
			"locationID": payment.LocationID,
		})
	}

	return result, nil
}

func (gds *GameDigitalService) GetPayment(ctx context.Context, paymentID int64) (*models.Payment, error) {
	// Get from DB (no caching for single payment)
	return gds.dbAdapter.GetPayment(ctx, paymentID)
}

// AddGameToDigitalLocation adds a game to a digital location
func (gds *GameDigitalService) AddGameToDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error {
	gds.logger.Debug("AddGameToDigitalLocation called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
		"gameID":     gameID,
	})

	// Add to db
	if err := gds.dbAdapter.AddGameToDigitalLocation(ctx, userID, locationID, gameID); err != nil {
		gds.logger.Error("Failed to add game to digital location", map[string]any{"error": err})
		return err
	}

	// Invalidate cache for this location
	if err := gds.cacheWrapper.SetSingleCachedDigitalLocation(ctx, userID, models.DigitalLocation{ID: locationID}); err != nil {
		gds.logger.Error("Failed to invalidate cache", map[string]any{"error": err})
		// Don't return error here, just log it
	}

	return nil
}

// RemoveGameFromDigitalLocation removes a game from a digital location
func (gds *GameDigitalService) RemoveGameFromDigitalLocation(ctx context.Context, userID string, locationID string, gameID int64) error {
	gds.logger.Debug("RemoveGameFromDigitalLocation called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
		"gameID":     gameID,
	})

	// Remove from db
	if err := gds.dbAdapter.RemoveGameFromDigitalLocation(ctx, userID, locationID, gameID); err != nil {
		gds.logger.Error("Failed to remove game from digital location", map[string]any{"error": err})
		return err
	}

	// Invalidate cache for this location
	if err := gds.cacheWrapper.SetSingleCachedDigitalLocation(ctx, userID, models.DigitalLocation{ID: locationID}); err != nil {
		gds.logger.Error("Failed to invalidate cache", map[string]any{"error": err})
		// Don't return error here, just log it
	}

	return nil
}

// GetGamesByDigitalLocationID gets all games in a digital location
func (gds *GameDigitalService) GetGamesByDigitalLocationID(ctx context.Context, userID string, locationID string) ([]models.Game, error) {
	gds.logger.Debug("GetGamesByDigitalLocationID called", map[string]any{
		"userID":     userID,
		"locationID": locationID,
	})

	// Get from db
	games, err := gds.dbAdapter.GetGamesByDigitalLocationID(ctx, userID, locationID)
	if err != nil {
		gds.logger.Error("Failed to get games for digital location", map[string]any{"error": err})
		return nil, err
	}

	return games, nil
}

// isRetryableError determines if an error is retryable
