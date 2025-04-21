package physical

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/appcontext"
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
func RegisterPhysicalRoutes(r chi.Router, appCtx *appcontext.AppContext, service services.PhysicalService) {
	// Base routes
	r.Get("/", GetUserPhysicalLocations(appCtx, service))
	r.Post("/", AddPhysicalLocation(appCtx, service))

	// Nested routes with ID
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", GetPhysicalLocation(appCtx, service))
		r.Put("/", UpdatePhysicalLocation(appCtx, service))
		r.Delete("/", RemovePhysicalLocation(appCtx, service))
	})
}

// GetUserPhysicalLocations handles GET requests for listing all physical locations
func GetUserPhysicalLocations(appCtx *appcontext.AppContext, service services.PhysicalService) http.HandlerFunc {
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

		appCtx.Logger.Info("Listing physical locations", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		locations, err := service.GetUserPhysicalLocations(r.Context(), userID)
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

		response := struct {
			Success   bool                     `json:"success"`
			UserID    string                   `json:"user_id"`
			Locations []models.PhysicalLocation `json:"locations"`
		}{
			Success:   true,
			UserID:    userID,
			Locations: locations,
		}

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
		appCtx.Logger.Info("Getting physical location", map[string]any{
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

		location, err := service.GetPhysicalLocation(r.Context(), userID, locationID)
		if err != nil {
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

		response := struct {
			Success  bool                    `json:"success"`
			Location models.PhysicalLocation `json:"location"`
		}{
			Success:  true,
			Location: location,
		}

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// AddPhysicalLocation handles POST requests for creating a new physical location
func AddPhysicalLocation(appCtx *appcontext.AppContext, service services.PhysicalService) http.HandlerFunc {
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

		appCtx.Logger.Info("Creating physical location", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		var locationRequest PhysicalLocationRequest
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
			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrValidationFailed) {
				statusCode = http.StatusBadRequest
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

		response := struct {
			Success  bool                    `json:"success"`
			Location models.PhysicalLocation `json:"location"`
		}{
			Success:  true,
			Location: createdLocation,
		}

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusCreated,
			response,
		)
	}
}

// UpdatePhysicalLocation handles PUT requests for updating a physical location
func UpdatePhysicalLocation(appCtx *appcontext.AppContext, service services.PhysicalService) http.HandlerFunc {
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
		appCtx.Logger.Info("Updating physical location", map[string]any{
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

		var locationRequest PhysicalLocationRequest
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
			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrLocationNotFound) {
				statusCode = http.StatusNotFound
			} else if errors.Is(err, ErrValidationFailed) {
				statusCode = http.StatusBadRequest
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

		response := struct {
			Success  bool                    `json:"success"`
			Location models.PhysicalLocation `json:"location"`
		}{
			Success:  true,
			Location: updatedLocation,
		}

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// RemovePhysicalLocation handles DELETE requests for removing a physical location
func RemovePhysicalLocation(appCtx *appcontext.AppContext, service services.PhysicalService) http.HandlerFunc {
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
		appCtx.Logger.Info("Deleting physical location", map[string]any{
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

		err := service.DeletePhysicalLocation(r.Context(), userID, locationID)
		if err != nil {
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

		response := struct {
			Success bool   `json:"success"`
			ID      string `json:"id"`
		}{
			Success: true,
			ID:      locationID,
		}

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}
