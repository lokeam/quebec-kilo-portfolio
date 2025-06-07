package physical

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/analytics"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/locations/formatters"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

// PhysicalLocationRequest represents the request payload for physical location operations
type PhysicalLocationRequest struct {
	Name           string  `json:"name"`
	Label          string  `json:"label"`
	LocationType   string  `json:"location_type"`
	MapCoordinates string  `json:"map_coordinates"`
	BgColor        string  `json:"bg_color"`
}

// RegisterPhysicalRoutes registers all physical location routes
func RegisterPhysicalRoutes(r chi.Router, appCtx *appcontext.AppContext, service services.PhysicalService, analyticsService analytics.Service) {
	// Base routes
	r.Get("/", GetAllPhysicalLocations(appCtx, service))
	r.Post("/", CreatePhysicalLocation(appCtx, service, analyticsService))
	r.Delete("/", DeletePhysicalLocation(appCtx, service, analyticsService))

	// Nested routes with ID
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", GetSinglePhysicalLocation(appCtx, service))
		r.Put("/", UpdatePhysicalLocation(appCtx, service, analyticsService))
	})

	// BFF route
	r.Get("/bff", GetAllPhysicalLocationsBFF(appCtx, service))
}

// handleError is a helper function to standardize error handling
func handleError(w http.ResponseWriter, logger interfaces.Logger, requestID string, err error) {
	statusCode := GetStatusCodeForError(err)
	httputils.RespondWithError(
		httputils.NewResponseWriterAdapter(w),
		logger,
		requestID,
		err,
		statusCode,
	)
}

// GetAllPhysicalLocations handles GET requests for listing all physical locations
func GetAllPhysicalLocations(appCtx *appcontext.AppContext, service services.PhysicalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)

		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			handleError(w, appCtx.Logger, requestID, errors.New("userID not found in request context"))
			return
		}

		appCtx.Logger.Info("Listing physical locations", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		locations, err := service.GetAllPhysicalLocations(r.Context(), userID)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		physicalLocations := formatters.FormatPhysicalLocationsToFrontend(locations)

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "physical" key, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"physical": physicalLocations,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// GetSinglePhysicalLocation handles GET requests for a single physical location
func GetSinglePhysicalLocation(appCtx *appcontext.AppContext, service services.PhysicalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)

		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			handleError(w, appCtx.Logger, requestID, errors.New("userID not found in request context"))
			return
		}

		locationID := chi.URLParam(r, "id")
		appCtx.Logger.Info("Getting physical location", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
			"locationID": locationID,
		})

		if locationID == "" {
			handleError(w, appCtx.Logger, requestID, errors.New("location ID is required"))
			return
		}

		location, err := service.GetSinglePhysicalLocation(r.Context(), userID, locationID)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "physical" key, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"physical": formatters.FormatPhysicalLocationToFrontend(&location),
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// GetAllPhysicalLocationsBFF handles GET requests for the /physical-locations page
func GetAllPhysicalLocationsBFF(appCtx *appcontext.AppContext, service services.PhysicalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)

		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			handleError(w, appCtx.Logger, requestID, errors.New("userID not found in request context"))
			return
		}

		appCtx.Logger.Info("Listing physical locations", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		locations, err := service.GetAllPhysicalLocationsBFF(r.Context(), userID)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "physical" key, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"physical": locations,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// CreatePhysicalLocation handles POST requests for creating a new physical location
func CreatePhysicalLocation(
	appCtx *appcontext.AppContext,
	service services.PhysicalService,
	analyticsService analytics.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)

		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			handleError(w, appCtx.Logger, requestID, errors.New("userID not found in request context"))
			return
		}

		appCtx.Logger.Info("Creating physical location", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		var locationRequest PhysicalLocationRequest
		if err := json.NewDecoder(r.Body).Decode(&locationRequest); err != nil {
			handleError(w, appCtx.Logger, requestID, errors.New("invalid request body"))
			return
		}

		// Create a new UUID for the location
		locationID := uuid.New().String()
		now := time.Now()

		// Convert string coordinates to PhysicalMapCoordinates struct
		mapCoords := models.PhysicalMapCoordinates{
			Coords: locationRequest.MapCoordinates,
		}

		location := models.PhysicalLocation{
			ID:             locationID,
			UserID:         userID,
			Name:           locationRequest.Name,
			Label:          locationRequest.Label,
			LocationType:   locationRequest.LocationType,
			MapCoordinates: mapCoords,
			BgColor:        locationRequest.BgColor,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		createdLocation, err := service.CreatePhysicalLocation(r.Context(), userID, location)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Invalidate both BFF and analytics caches
		bffCacheKey := fmt.Sprintf("physical:bff:%s", userID)
		if err := service.InvalidateCache(r.Context(), bffCacheKey); err != nil {
			appCtx.Logger.Warn("Failed to invalidate BFF cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Warn("Failed to invalidate analytics cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		// Use standard response format
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"physical": formatters.FormatPhysicalLocationToFrontend(&createdLocation),
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusCreated,
			response,
		)
	}
}

// UpdatePhysicalLocation handles PUT requests for updating a physical location
func UpdatePhysicalLocation(
	appCtx *appcontext.AppContext,
	service services.PhysicalService,
	analyticsService analytics.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)

		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
			})
			handleError(w, appCtx.Logger, requestID, errors.New("userID not found in request context"))
			return
		}

		locationID := chi.URLParam(r, "id")
		appCtx.Logger.Info("Updating physical location", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
			"locationID": locationID,
		})

		if locationID == "" {
			handleError(w, appCtx.Logger, requestID, errors.New("location ID is required"))
			return
		}

		var locationRequest PhysicalLocationRequest
		if err := json.NewDecoder(r.Body).Decode(&locationRequest); err != nil {
			handleError(w, appCtx.Logger, requestID, errors.New("invalid request body"))
			return
		}

		// Convert string coordinates to PhysicalMapCoordinates struct
		mapCoords := models.PhysicalMapCoordinates{
			Coords: locationRequest.MapCoordinates,
		}

		location := models.PhysicalLocation{
			ID:             locationID,
			UserID:         userID,
			Name:           locationRequest.Name,
			Label:          locationRequest.Label,
			LocationType:   locationRequest.LocationType,
			MapCoordinates: mapCoords,
			BgColor:        locationRequest.BgColor,
			UpdatedAt:      time.Now(),
		}

		updatedLocation, err := service.UpdatePhysicalLocation(r.Context(), userID, location)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Invalidate both BFF and analytics caches
		bffCacheKey := fmt.Sprintf("physical:bff:%s", userID)
		if err := service.InvalidateCache(r.Context(), bffCacheKey); err != nil {
			appCtx.Logger.Warn("Failed to invalidate BFF cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Warn("Failed to invalidate analytics cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		// Use standard response format
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"physical": formatters.FormatPhysicalLocationToFrontend(&updatedLocation),
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// DeletePhysicalLocation handles DELETE requests for removing physical locations
// It supports both single location deletion (via URL param) and bulk deletion (via request body)
func DeletePhysicalLocation(
	appCtx *appcontext.AppContext,
	service services.PhysicalService,
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
		physicalLocationIDs := r.URL.Query().Get("ids")
		if physicalLocationIDs == "" {
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
		physicalLocationIDsArr := strings.Split(physicalLocationIDs, ",")
		appCtx.Logger.Info("Deleting physical location(s)", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
			"locationID": physicalLocationIDsArr,
		})

		// Call service method to delete locations
		deletedCount, err := service.DeletePhysicalLocation(r.Context(), userID, physicalLocationIDsArr)
		if err != nil {
			appCtx.Logger.Error("Failed to delete physical locations", map[string]any{
				"error": err,
				"request_id": requestID,
				"location_ids": physicalLocationIDsArr,
			})

			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrLocationNotFound) {
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

		// Invalidate both BFF and analytics caches
		bffCacheKey := fmt.Sprintf("physical:bff:%s", userID)
		if err := service.InvalidateCache(r.Context(), bffCacheKey); err != nil {
			appCtx.Logger.Error("Failed to invalidate BFF cache after deleting location", map[string]any{
				"error": err,
				"userID": userID,
			})
		}

		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Error("Failed to invalidate analytics cache after deleting location", map[string]any{
				"error": err,
				"userID": userID,
			})
		}

		// Log success
		appCtx.Logger.Info("Successfully deleted physical locations", map[string]any{
			"request_id": requestID,
			"user_id": userID,
			"deleted_count": deletedCount,
			"total_count": len(physicalLocationIDsArr),
		})

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "physical" key, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"physical": map[string]any{
				"success": true,
				"deleted_count": deletedCount,
				"location_ids": physicalLocationIDsArr,
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
