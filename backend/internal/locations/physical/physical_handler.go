package physical

import (
	"encoding/json"
	"errors"
	"net/http"
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
	Name           string `json:"name"`
	Label          string `json:"label"`
	LocationType   string `json:"location_type"`
	MapCoordinates string `json:"map_coordinates"`
	BgColor        string `json:"bg_color"`
}

// RegisterPhysicalRoutes registers all physical location routes
func RegisterPhysicalRoutes(r chi.Router, appCtx *appcontext.AppContext, service services.PhysicalService, analyticsService analytics.Service) {
	// Base routes
	r.Get("/", GetAllPhysicalLocations(appCtx, service))
	r.Post("/", CreatePhysicalLocation(appCtx, service, analyticsService))

	// Nested routes with ID
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", GetSinglePhysicalLocation(appCtx, service))
		r.Put("/", UpdatePhysicalLocation(appCtx, service, analyticsService))
		r.Delete("/", DeletePhysicalLocation(appCtx, service, analyticsService))
	})
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

		location := models.PhysicalLocation{
			ID:             locationID,
			UserID:         userID,
			Name:           locationRequest.Name,
			Label:          locationRequest.Label,
			LocationType:   locationRequest.LocationType,
			MapCoordinates: locationRequest.MapCoordinates,
			BgColor:        locationRequest.BgColor,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		createdLocation, err := service.CreatePhysicalLocation(r.Context(), userID, location)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Invalidate analytics cache for inventory domain
		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Warn("Failed to invalidate analytics cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "physical" key, DO NOT use a struct{}
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

		location := models.PhysicalLocation{
			ID:             locationID,
			UserID:         userID,
			Name:           locationRequest.Name,
			Label:          locationRequest.Label,
			LocationType:   locationRequest.LocationType,
			MapCoordinates: locationRequest.MapCoordinates,
			UpdatedAt:      time.Now(),
		}

		updatedLocation, err := service.UpdatePhysicalLocation(r.Context(), userID, location)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Invalidate analytics cache for inventory domain
		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Warn("Failed to invalidate analytics cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "physical" key, DO NOT use a struct{}
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

// DeletePhysicalLocation handles DELETE requests for removing a physical location
func DeletePhysicalLocation(
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
		appCtx.Logger.Info("Removing physical location", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
			"locationID": locationID,
		})

		if locationID == "" {
			handleError(w, appCtx.Logger, requestID, errors.New("location ID is required"))
			return
		}

		err := service.DeletePhysicalLocation(r.Context(), userID, locationID)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		// Invalidate analytics cache for inventory domain
		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Warn("Failed to invalidate analytics cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "physical" key, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"physical": map[string]any{
				"id":      locationID,
				"message": "Physical location removed successfully",
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
