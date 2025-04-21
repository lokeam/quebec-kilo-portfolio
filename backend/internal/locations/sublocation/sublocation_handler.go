package sublocation

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

type SublocationRequest struct {
	Name               string `json:"name"`
	LocationType       string `json:"location_type"`
	BgColor           string `json:"bg_color"`
	StoredItems       int    `json:"stored_items"`
	PhysicalLocationID string `json:"physical_location_id"`
}

// RegisterSublocationRoutes registers all sublocation routes
func RegisterSublocationRoutes(r chi.Router, appCtx *appcontext.AppContext, service services.SublocationService) {
	// Base routes
	r.Get("/", GetSublocations(appCtx, service))
	r.Post("/", AddSublocation(appCtx, service))

	// Nested routes with ID
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", GetSublocation(appCtx, service))
		r.Put("/", UpdateSublocation(appCtx, service))
		r.Delete("/", DeleteSublocation(appCtx, service))
	})

	// Add an explicit endpoint to invalidate physical location caches if needed
	r.Post("/refresh-physical-cache/{physicalLocationID}", func(w http.ResponseWriter, r *http.Request) {
		requestID := httputils.GetRequestID(r)
		userID := httputils.GetUserID(r)
		physicalLocationID := chi.URLParam(r, "physicalLocationID")

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

		if physicalLocationID == "" {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("physicalLocationID is required"),
				http.StatusBadRequest,
			)
			return
		}

		// Force cache refresh by invalidating the cache for this physical location
		// This endpoint can be called from the client side after sublocation operations
		// that might need immediate refresh
		cacheAdapter, ok := service.(*GameSublocationService)
		if ok && cacheAdapter.cacheWrapper != nil {
			err := cacheAdapter.cacheWrapper.InvalidateLocationCache(r.Context(), userID, physicalLocationID)
			if err != nil {
				appCtx.Logger.Error("Failed to invalidate physical location cache", map[string]any{
					"error": err,
					"physicalLocationID": physicalLocationID,
				})
				httputils.RespondWithError(
					httputils.NewResponseWriterAdapter(w),
					appCtx.Logger,
					requestID,
					err,
					http.StatusInternalServerError,
				)
				return
			}

			httputils.RespondWithJSON(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				http.StatusOK,
				map[string]interface{}{"success": true, "message": "Physical location cache refreshed"},
			)
		} else {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("unable to access cache adapter"),
				http.StatusInternalServerError,
			)
		}
	})
}

// GetSublocations handles GET requests for listing all sublocations
func GetSublocations(appCtx *appcontext.AppContext, service services.SublocationService) http.HandlerFunc {
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

		appCtx.Logger.Info("Listing sublocations", map[string]any{
			"requestID": requestID,
			"userID": userID,
		})

		locations, err := service.GetSublocations(r.Context(), userID)
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
			Success        bool                    `json:"success"`
			Sublocations   []models.Sublocation    `json:"sublocations"`
		} {
			Success: true,
			Sublocations: locations,
		}

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// GetSublocation handles GET requests for a single sublocation
func GetSublocation(appCtx *appcontext.AppContext, service services.SublocationService) http.HandlerFunc {
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
		appCtx.Logger.Info("Getting sublocation", map[string]any{
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

		sublocation, err := service.GetSublocation(r.Context(), userID, locationID)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrSublocationNotFound) {
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
			Success    bool                 `json:"success"`
			Location   models.Sublocation   `json:"location"`
		} {
			Success:  true,
			Location: sublocation,
		}

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// AddSublocation handles POST requests for creating a new sublocation
func AddSublocation(appCtx *appcontext.AppContext, service services.SublocationService) http.HandlerFunc {
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

		appCtx.Logger.Info("Creating sublocation", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
		})

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("error reading request body"),
				http.StatusBadRequest,
			)
			return
		}
		appCtx.Logger.Debug("Request body", map[string]any{
			"body": string(body),
		})

		// Reset the request body
		r.Body = io.NopCloser(bytes.NewReader(body))

		var locationRequest SublocationRequest
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

		appCtx.Logger.Debug("Decoded request", map[string]any{
			"request": locationRequest,
		})

		// Validate required fields
		if locationRequest.PhysicalLocationID == "" {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("physical_location_id is required"),
				http.StatusBadRequest,
			)
			return
		}

		// Validate UUID format
		if _, err := uuid.Parse(locationRequest.PhysicalLocationID); err != nil {
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				errors.New("physical_location_id must be a valid UUID"),
				http.StatusBadRequest,
			)
			return
		}

		// Create a new UUID for the location
		locationID := uuid.New().String()
		now := time.Now()

		sublocation := models.Sublocation{
			ID:                 locationID,
			Name:               locationRequest.Name,
			LocationType:       locationRequest.LocationType,
			BgColor:            locationRequest.BgColor,
			StoredItems:        locationRequest.StoredItems,
			UserID:             userID,
			PhysicalLocationID: locationRequest.PhysicalLocationID,
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		createdSublocation, err := service.AddSublocation(r.Context(), userID, sublocation)
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
			Success     bool                  `json:"success"`
			Location    models.Sublocation    `json:"location"`
		} {
			Success:    true,
			Location:   createdSublocation,
		}

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusCreated,
			response,
		)
	}
}

// UpdateSublocation handles PUT requests for updating a sublocation
func UpdateSublocation(appCtx *appcontext.AppContext, service services.SublocationService) http.HandlerFunc {
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
		appCtx.Logger.Info("Updating sublocation", map[string]any{
			"requestID":    requestID,
			"userID":       userID,
			"locationID":   locationID,
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

		var locationRequest SublocationRequest
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

		location := models.Sublocation{
			ID:             locationID,
			Name:           locationRequest.Name,
			LocationType:   locationRequest.LocationType,
			BgColor:        locationRequest.BgColor,
			StoredItems:    locationRequest.StoredItems,
			UserID:         userID,
			PhysicalLocationID: locationRequest.PhysicalLocationID,
			UpdatedAt:      time.Now(),
		}

		err := service.UpdateSublocation(r.Context(), userID, location)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrSublocationNotFound) {
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
			Success     bool                `json:"success"`
			Location    models.Sublocation  `json:"location"`
		} {
			Success:    true,
			Location:   location,
		}

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// DeleteSublocation handles DELETE requests for removing a sublocation
func DeleteSublocation(appCtx *appcontext.AppContext, service services.SublocationService) http.HandlerFunc {
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
		appCtx.Logger.Info("Deleting sublocation", map[string]any{
			"requestID":    requestID,
			"userID":       userID,
			"locationID":   locationID,
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

		err := service.DeleteSublocation(r.Context(), userID, locationID)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrSublocationNotFound) {
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
			Success    bool      `json:"success"`
			ID         string    `json:"id"`
		} {
			Success:   true,
			ID:        locationID,
		}

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}
