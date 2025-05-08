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
}

// RegisterPhysicalRoutes registers all physical location routes
func RegisterPhysicalRoutes(r chi.Router, appCtx *appcontext.AppContext, service services.PhysicalService, analyticsService analytics.Service) {
	// Base routes
	r.Get("/", GetUserPhysicalLocations(appCtx, service))
	r.Post("/", AddPhysicalLocation(appCtx, service, analyticsService))

	// Nested routes with ID
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", GetPhysicalLocation(appCtx, service))
		r.Put("/", UpdatePhysicalLocation(appCtx, service, analyticsService))
		r.Delete("/", RemovePhysicalLocation(appCtx, service, analyticsService))
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

// GetUserPhysicalLocations handles GET requests for listing all physical locations
func GetUserPhysicalLocations(appCtx *appcontext.AppContext, service services.PhysicalService) http.HandlerFunc {
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

		locations, err := service.GetUserPhysicalLocations(r.Context(), userID)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		data := struct {
			UserID      string                  `json:"user_id"`
			Locations   []map[string]any        `json:"locations"`
		} {
			UserID:      userID,
			Locations:   formatters.FormatPhysicalLocationsToFrontend(locations),
		}

		response := httputils.NewAPIResponse(r, userID, data)

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// GetPhysicalLocation handles GET requests for a single physical location
func GetPhysicalLocation(appCtx *appcontext.AppContext, service services.PhysicalService) http.HandlerFunc {
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

		location, err := service.GetPhysicalLocation(r.Context(), userID, locationID)
		if err != nil {
			handleError(w, appCtx.Logger, requestID, err)
			return
		}

		data := struct {
			Location map[string]any     `json:"location"`
		}{
			Location: formatters.FormatPhysicalLocationToFrontend(&location),
		}

		response := httputils.NewAPIResponse(r, userID, data)

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// AddPhysicalLocation handles POST requests for creating a new physical location
func AddPhysicalLocation(
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
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		createdLocation, err := service.AddPhysicalLocation(r.Context(), userID, location)
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

		data := struct {
			Location map[string]any        `json:"location"`
		} {
			Location: formatters.FormatPhysicalLocationToFrontend(&createdLocation),
		}

		response := httputils.NewAPIResponse(r, userID, data)

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

		data := struct {
			Location map[string]any        `json:"location"`
		} {
			Location: formatters.FormatPhysicalLocationToFrontend(&updatedLocation),
		}

		response := httputils.NewAPIResponse(r, userID, data)

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// RemovePhysicalLocation handles DELETE requests for removing a physical location
func RemovePhysicalLocation(
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

		data := struct {
			ID      string   `json:"id"`
			Message string   `json:"message"`
		} {
			ID:      locationID,
			Message: "Physical location removed successfully",
		}

		response := httputils.NewAPIResponse(r, userID, data)

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}
