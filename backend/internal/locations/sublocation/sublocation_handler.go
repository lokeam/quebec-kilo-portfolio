package sublocation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/analytics"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
)

type SublocationRequest struct {
	Name               string `json:"name"`
	LocationType       string `json:"location_type"`
	BgColor           string `json:"bg_color"`
	StoredItems       int    `json:"stored_items"`
	PhysicalLocationID string `json:"physical_location_id"`
}

// RegisterSublocationRoutes registers all sublocation routes
func RegisterSublocationRoutes(
	r chi.Router,
	appCtx *appcontext.AppContext,
	service services.SublocationService,
	analyticsService analytics.Service,
) {
	// Base routes
	r.Get("/", GetSublocations(appCtx, service))
	r.Post("/", CreateSublocation(appCtx, service, analyticsService))

	// Game management routes
	r.Post("/move-game", MoveGame(appCtx, service))
	r.Post("/remove-game", RemoveGame(appCtx, service))

	// Nested routes with ID
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", GetSingleSublocation(appCtx, service))
		r.Put("/", UpdateSublocation(appCtx, service, analyticsService))
		r.Delete("/", DeleteSublocation(appCtx, service, analyticsService))
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

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "sublocation" key, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"sublocation": locations,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// GetSingleSublocation handles GET requests for a single sublocation
func GetSingleSublocation(appCtx *appcontext.AppContext, service services.SublocationService) http.HandlerFunc {
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

		sublocation, err := service.GetSingleSublocation(r.Context(), userID, locationID)
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

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "sublocation" key, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"sublocation": sublocation,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// CreateSublocation handles POST requests for creating a new sublocation
func CreateSublocation(
	appCtx *appcontext.AppContext,
	service services.SublocationService,
	analyticsService analytics.Service,
) http.HandlerFunc {
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

		// Parse request body into CreateSublocationRequest
		var req types.CreateSublocationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			appCtx.Logger.Error("Failed to decode request body", map[string]any{
				"request_id": requestID,
				"error":     err,
			})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				fmt.Errorf("invalid request body: %w", err),
				http.StatusBadRequest,
			)
			return
		}

		// Create sublocation using service
		createdLocation, err := service.CreateSublocation(r.Context(), userID, req)
		if err != nil {
			appCtx.Logger.Error("Failed to create sublocation", map[string]any{
				"request_id": requestID,
				"error":     err,
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

		// Invalidate analytics cache
		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Warn("Failed to invalidate analytics cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "sublocation" key, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"sublocation": createdLocation,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusCreated,
			response,
		)
	}
}

// UpdateSublocation handles PUT requests for updating a sublocation
func UpdateSublocation(appCtx *appcontext.AppContext, service services.SublocationService, analyticsService analytics.Service) http.HandlerFunc {
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

		locationID := chi.URLParam(r, "id")
		appCtx.Logger.Info("Updating sublocation", map[string]any{
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

		var req types.UpdateSublocationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			appCtx.Logger.Error("Failed to decode request body", map[string]any{
				"request_id": requestID,
				"error":     err,
			})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				fmt.Errorf("invalid request body: %w", err),
				http.StatusBadRequest,
			)
			return
		}

		err := service.UpdateSublocation(r.Context(), userID, locationID, req)
		if err != nil {
			appCtx.Logger.Error("Failed to update sublocation", map[string]any{
				"request_id": requestID,
				"error":     err,
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

		// Invalidate analytics cache for inventory domain
		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Warn("Failed to invalidate analytics cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		// Get the updated sublocation to return in response
		updatedSublocation, err := service.GetSingleSublocation(r.Context(), userID, locationID)
		if err != nil {
			appCtx.Logger.Error("Failed to get updated sublocation", map[string]any{
				"request_id": requestID,
				"error":     err,
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

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "sublocation" key, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"sublocation": updatedSublocation,
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// DeleteSublocation handles DELETE requests for removing a sublocation
func DeleteSublocation(
	appCtx *appcontext.AppContext,
	service services.SublocationService,
	analyticsService analytics.Service,
) http.HandlerFunc {
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

		locationID := chi.URLParam(r, "id")
		appCtx.Logger.Info("Deleting sublocation", map[string]any{
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

		err := service.DeleteSublocation(r.Context(), userID, locationID)
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

		// Invalidate analytics cache for inventory domain
		if err := analyticsService.InvalidateDomain(r.Context(), userID, analytics.DomainInventory); err != nil {
			appCtx.Logger.Warn("Failed to invalidate analytics cache", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err,
			})
		}

		// Use standard response format
		// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "sublocation" key, DO NOT use a struct{}
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"sublocation": map[string]any{
				"id":      locationID,
				"message": "Sublocation deleted successfully",
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

// MoveGame handles POST requests for moving a game to a different sublocation
func MoveGame(appCtx *appcontext.AppContext, service services.SublocationService) http.HandlerFunc {
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

		// Parse request body
		var req types.MoveGameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			appCtx.Logger.Error("Failed to decode request body", map[string]any{
				"request_id": requestID,
				"error":     err,
			})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				fmt.Errorf("invalid request body: %w", err),
				http.StatusBadRequest,
			)
			return
		}

		// Move game using service
		if err := service.MoveGame(r.Context(), userID, req); err != nil {
			appCtx.Logger.Error("Failed to move game", map[string]any{
				"request_id": requestID,
				"error":     err,
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

		// Use standard response format
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"message": "Game moved successfully",
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}

// RemoveGame handles POST requests for removing a game from its current sublocation
func RemoveGame(appCtx *appcontext.AppContext, service services.SublocationService) http.HandlerFunc {
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

		// Parse request body
		var req types.RemoveGameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			appCtx.Logger.Error("Failed to decode request body", map[string]any{
				"request_id": requestID,
				"error":     err,
			})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				fmt.Errorf("invalid request body: %w", err),
				http.StatusBadRequest,
			)
			return
		}

		// Remove game using service
		if err := service.RemoveGame(r.Context(), userID, req); err != nil {
			appCtx.Logger.Error("Failed to remove game", map[string]any{
				"request_id": requestID,
				"error":     err,
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

		// Use standard response format
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"message": "Game removed successfully",
		})

		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}
