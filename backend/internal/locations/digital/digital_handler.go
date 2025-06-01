package digital

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/analytics"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/locations/formatters"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

// DigitalLocationRequest represents the request payload for digital location operations
type DigitalLocationRequest struct {
	ID             string               `json:"id,omitempty" db:"id"`
	Name           string               `json:"name" db:"name"`
	IsSubscription bool                 `json:"is_subscription" db:"is_subscription"`
	IsActive       bool                 `json:"is_active" db:"is_active"`
	URL            string               `json:"url" db:"url"`
	PaymentMethod  string               `json:"payment_method" db:"payment_method"`
	CreatedAt      time.Time            `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at,omitempty" db:"updated_at"`
	Subscription   *models.Subscription `json:"subscription,omitempty"`
}

// RegisterDigitalRoutes registers all digital location routes
func RegisterDigitalRoutes(
	r chi.Router,
	appCtx *appcontext.AppContext,
	service services.DigitalService,
	analyticsService analytics.Service,
) {
	// Base routes
	r.Get("/", GetUserDigitalLocations(appCtx, service))
	r.Post("/", AddDigitalLocation(appCtx, service, analyticsService))
	r.Delete("/", RemoveDigitalLocation(appCtx, service, analyticsService))  // Handles both single and bulk deletion

	// Nested routes with ID
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", GetDigitalLocation(appCtx, service))
		r.Put("/", UpdateDigitalLocation(appCtx, service, analyticsService))
	})
}

// GetUserDigitalLocations handles GET requests for listing all digital locations
func GetUserDigitalLocations(appCtx *appcontext.AppContext, service services.DigitalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)
		userID := httputils.GetUserID(r)

		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("userID not found in request context"),
				http.StatusUnauthorized,
			)
			return
		}

		appCtx.Logger.Info("Listing digital locations", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		locations, err := service.GetUserDigitalLocations(r.Context(), userID)
		if err != nil {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				err,
				http.StatusInternalServerError,
			)
			return
		}

		// Debug: Check subscription data
		for i, loc := range locations {
			appCtx.Logger.Debug("Location from DB", map[string]any{
				"index": i,
				"id": loc.ID,
				"name": loc.Name,
				"has_subscription": loc.Subscription != nil,
			})
			if loc.Subscription != nil {
				appCtx.Logger.Debug("Subscription details", map[string]any{
					"sub_id": loc.Subscription.ID,
					"billing_cycle": loc.Subscription.BillingCycle,
					"cost": loc.Subscription.CostPerCycle,
				})
			}
		}

		// Convert backend model to frontend-compatible format
		frontendLocations := make([]map[string]any, len(locations))
		for i, loc := range locations {
			frontendLocations[i] = formatters.FormatDigitalLocationToFrontend(&loc)
		}

		// Log a sample of the transformed data
		if len(frontendLocations) > 0 {
			appCtx.Logger.Debug("Sample transformed location", map[string]any{
				"sample": frontendLocations[0],
				"has_billing": frontendLocations[0]["billing"] != nil,
				"has_subscription": frontendLocations[0]["subscription"] != nil,
			})
		}

		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"digital": frontendLocations,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// GetDigitalLocation handles GET requests for a single digital location
func GetDigitalLocation(appCtx *appcontext.AppContext, service services.DigitalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)

		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("userID not found in request context"),
				http.StatusUnauthorized,
			)
			return
		}

		locationID := chi.URLParam(r, "id")
		appCtx.Logger.Info("Getting digital location", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
			"locationID": locationID,
		})

		if locationID == "" {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("id is required"),
				http.StatusBadRequest,
			)
			return
		}

		var location models.DigitalLocation
		var err error

		if _, parseErr := uuid.Parse(locationID); parseErr == nil {
			location, err = service.GetDigitalLocation(r.Context(), userID, locationID)
		} else {
			location, err = service.FindDigitalLocationByName(r.Context(), userID, locationID)
		}

		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrDigitalLocationNotFound) {
				statusCode = http.StatusNotFound
			}
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				err,
				statusCode,
			)
			return
		}

		// Convert to frontend format
		frontendLocation := formatters.FormatDigitalLocationToFrontend(&location)

		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"digital": frontendLocation,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// AddDigitalLocation handles POST requests for creating a new digital location
func AddDigitalLocation(
	appCtx *appcontext.AppContext,
	service services.DigitalService,
	analyticsService analytics.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)

		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("userID not found in request context"),
				http.StatusUnauthorized,
			)
			return
		}

		appCtx.Logger.Info("Creating digital location", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		var locationRequest DigitalLocationRequest
		if err := json.NewDecoder(r.Body).Decode(&locationRequest); err != nil {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("invalid request body"),
				http.StatusBadRequest,
			)
			return
		}

		now := time.Now()
		digitalLocation := models.DigitalLocation{
			ID:          uuid.New().String(),
			UserID:      userID,
			Name:        locationRequest.Name,
			IsSubscription: locationRequest.IsSubscription,
			IsActive:    locationRequest.IsActive,
			URL:         locationRequest.URL,
			PaymentMethod: locationRequest.PaymentMethod,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Handle subscription if provided
		if locationRequest.Subscription != nil {
			subscription := *locationRequest.Subscription
			subscription.LocationID = digitalLocation.ID
			subscription.CreatedAt = now
			subscription.UpdatedAt = now
			digitalLocation.Subscription = &subscription
		}

		createdLocation, err := service.AddDigitalLocation(r.Context(), userID, digitalLocation)
		if err != nil {
			statusCode := http.StatusInternalServerError
			appCtx.Logger.Error("Failed to create digital location", map[string]any{
				"error":      err,
				"request_id": requestID,
			})
			// Handle specific error types
			if errors.Is(err, ErrDigitalLocationExists) {
				statusCode = http.StatusConflict
			}
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				err,
				statusCode,
			)
			return
		}

		// Convert to frontend format
		adaptedResponse := formatters.FormatDigitalLocationToFrontend(&createdLocation)

		// After successful creation, invalidate analytics cache
		if err := analyticsService.InvalidateDomains(r.Context(), userID, []string{
			analytics.DomainGeneral,
			analytics.DomainFinancial,
			analytics.DomainStorage,
		}); err != nil {
			appCtx.Logger.Error("Failed to invalidate analytics cache after adding location", map[string]any{
				"error": err,
				"userID": userID,
			})
			// Continue despite error, since the location was created successfully
		}

		// IMPORTANT: All responses must use the NewAPIResponse function AND be wrapped in a map
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"digital": adaptedResponse,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusCreated,
			response,
		)
	}
}

// UpdateDigitalLocation handles PUT requests for updating a digital location
func UpdateDigitalLocation(
	appCtx *appcontext.AppContext,
	service services.DigitalService,
	analyticsService analytics.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)

		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("userID not found in request context"),
				http.StatusUnauthorized,
			)
			return
		}

		locationID := chi.URLParam(r, "id")
		appCtx.Logger.Info("Updating digital location", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
			"locationID": locationID,
		})

		if locationID == "" {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("location ID is required"),
				http.StatusBadRequest,
			)
			return
		}

		// Unmarshal the request body
		var req DigitalLocationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			appCtx.Logger.Error("Failed to decode request body", map[string]any{"error": err})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("invalid request body"),
				http.StatusBadRequest,
			)
			return
		}

		// Convert request to model
		location := models.DigitalLocation{
			ID:          locationID,
			UserID:      userID,
			Name:        req.Name,
			IsSubscription: req.IsSubscription,
			IsActive:    req.IsActive,
			URL:         req.URL,
			PaymentMethod: req.PaymentMethod,
			Subscription: req.Subscription,
		}

		// Get existing location
		existingLocation, err := service.GetDigitalLocation(r.Context(), userID, locationID)
		if err != nil {
			appCtx.Logger.Error("Failed to get existing location", map[string]any{"error": err})
			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrDigitalLocationNotFound) {
				statusCode = http.StatusNotFound
			}
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				err,
				statusCode,
			)
			return
		}

		// Update fields
		location.Name = req.Name
		location.IsSubscription = req.IsSubscription
		location.IsActive = req.IsActive
		location.URL = req.URL
		location.PaymentMethod = req.PaymentMethod
		location.UpdatedAt = time.Now()

		// Ensure the ID is set from the URL
		location.ID = locationID

		// Handle subscription
		if req.Subscription != nil {
			subscription := *req.Subscription
			subscription.LocationID = locationID
			subscription.UpdatedAt = time.Now()

			if existingLocation.Subscription == nil {
				// Add new subscription
				subscription.CreatedAt = time.Now()
				newSubscription, err := service.AddSubscription(r.Context(), subscription)
				if err != nil {
					appCtx.Logger.Error("Failed to add subscription during location update", map[string]any{"error": err})
					httputils.RespondWithError(
						httputils.NewResponseWriterAdapter(w),
						appCtx.Logger,
						requestID,
						errors.New("failed to add subscription"),
						http.StatusInternalServerError,
					)
					return
				}
				location.Subscription = newSubscription
			} else {
				// Update existing subscription
				subscription.ID = existingLocation.Subscription.ID
				subscription.CreatedAt = existingLocation.Subscription.CreatedAt
				if err := service.UpdateSubscription(r.Context(), subscription); err != nil {
					appCtx.Logger.Error("Failed to update subscription", map[string]any{"error": err})
					httputils.RespondWithError(
						httputils.NewResponseWriterAdapter(w),
						appCtx.Logger,
						requestID,
						errors.New("failed to update subscription"),
						http.StatusInternalServerError,
					)
					return
				}
				location.Subscription = &subscription
			}
		} else if existingLocation.Subscription != nil && !req.IsSubscription {
			// Remove subscription if service type changed and no subscription was provided
			if err := service.RemoveSubscription(r.Context(), locationID); err != nil {
				appCtx.Logger.Error("Failed to remove subscription", map[string]any{"error": err})
				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					errors.New("failed to remove subscription"),
					http.StatusInternalServerError,
				)
				return
			}
			location.Subscription = nil
		}

		// Call service method to update location in database
		if err := service.UpdateDigitalLocation(r.Context(), userID, location); err != nil {
			appCtx.Logger.Error("Failed to update digital location", map[string]any{"error": err})
			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrDigitalLocationNotFound) {
				statusCode = http.StatusNotFound
			}
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				err,
				statusCode,
			)
			return
		}

		// Get updated location to return
		updatedLocation, err := service.GetDigitalLocation(r.Context(), userID, locationID)
		if err != nil {
			appCtx.Logger.Error("Failed to get updated location", map[string]any{"error": err})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("location was updated but could not be retrieved"),
				http.StatusInternalServerError,
			)
			return
		}

		// Convert to frontend format
		frontendLocation := formatters.FormatDigitalLocationToFrontend(&updatedLocation)

		// After successful update, invalidate analytics cache
		if err := analyticsService.InvalidateDomains(r.Context(), userID, []string{
			analytics.DomainGeneral,
			analytics.DomainFinancial,
			analytics.DomainStorage,
		}); err != nil {
			appCtx.Logger.Error("Failed to invalidate analytics cache after updating location", map[string]any{
				"error": err,
				"userID": userID,
			})
			// Continue despite error, since the location was updated successfully
		}

		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"digital": frontendLocation,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// RemoveDigitalLocation handles DELETE requests for removing digital locations
// It supports both single location deletion (via URL param) and bulk deletion (via request body)
func RemoveDigitalLocation(
	appCtx *appcontext.AppContext,
	service services.DigitalService,
	analyticsService analytics.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)

		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("userID not found in request context"),
				http.StatusUnauthorized,
			)
			return
		}

		// Get IDs from query parameters
		digitalLocationIDs := r.URL.Query().Get("ids")
		if digitalLocationIDs == "" {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				ErrEmptyLocationIDs,
				http.StatusBadRequest,
			)
			return
		}

		// if multiple IDs are provided, split into array
		digitalLocationIDsArr := strings.Split(digitalLocationIDs, ",")
		appCtx.Logger.Info("Deleting digital location(s)", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
			"locationID": digitalLocationIDsArr,
		})

		// Call service method to delete locations
		deletedCount, err := service.RemoveDigitalLocation(r.Context(), userID, digitalLocationIDsArr)
		if err != nil {
			appCtx.Logger.Error("Failed to delete digital locations", map[string]any{
				"error": err,
				"request_id": requestID,
				"location_ids": digitalLocationIDsArr,
			})

			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrDigitalLocationNotFound) {
				statusCode = http.StatusNotFound
			}
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				err,
				statusCode,
			)
			return
		}

		// After successful deletion, invalidate analytics cache
		if err := analyticsService.InvalidateDomains(r.Context(), userID, []string{
			analytics.DomainGeneral,
			analytics.DomainFinancial,
			analytics.DomainStorage,
		}); err != nil {
			appCtx.Logger.Error("Failed to invalidate analytics cache after deleting location", map[string]any{
				"error": err,
				"userID": userID,
			})
			// Continue despite error, since the location was deleted successfully
		}

		// Log success
		appCtx.Logger.Info("Successfully deleted digital locations", map[string]any{
			"request_id": requestID,
			"user_id": userID,
			"deleted_count": deletedCount,
			"total_count": len(digitalLocationIDsArr),
		})

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "digital" key, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"digital": map[string]any{
				"success": true,
				"deleted_count": deletedCount,
				"location_ids": digitalLocationIDsArr,
			},
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// For backward compatibility - wrapper for the digital service catalog
func GetDigitalServicesCatalog(appCtx *appcontext.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)
		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
					"request_id": requestID,
			})
			httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					errors.New("userID not found in request context"),
					http.StatusUnauthorized,
			)
			return
		}

		appCtx.Logger.Info("Getting digital services catalog", map[string]any{
			"requestID": requestID,
			"userID": userID,
		})

		responseAdapter := NewDigitalResponseAdapter()
		adaptedResponse := responseAdapter.AdaptToCatalogResponse(DigitalServicesCatalog)

		// IMPORTANT: All responses MUST be wrapped in map[string]any{}, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"catalog": adaptedResponse,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}
