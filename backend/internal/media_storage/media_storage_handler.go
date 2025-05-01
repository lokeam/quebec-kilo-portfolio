package media_storage

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/analytics"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

// RegisterRoutes adds media storage routes to the router
func RegisterRoutes(r chi.Router, appCtx *appcontext.AppContext, service MediaStorageService) {
	r.Get("/stats", GetStorageStats(appCtx, service))
}

// GetStorageStats handles requests for storage statistics
func GetStorageStats(appCtx *appcontext.AppContext, service MediaStorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Request ID for tracing
		requestID := httputils.GetRequestID(r)

		// Get user ID from the context
		userID := httputils.GetUserID(r)
		if userID == "" {
			appCtx.Logger.Error("userID not found in request context", map[string]any{
				"request_id": requestID,
			})
			httputils.RespondWithError(
				httputils.NewResponseWriterAdapter(w),
				appCtx.Logger,
				requestID,
				ErrInvalidUserID,
				http.StatusUnauthorized,
			)
			return
		}

		appCtx.Logger.Info("Getting storage stats", map[string]any{
			"requestID": requestID,
			"userID":    userID,
		})

		// Call service to get storage stats
		stats, err := service.GetStorageStats(r.Context(), userID)
		if err != nil {
			appCtx.Logger.Error("Failed to get storage stats", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"error":     err.Error(),
			})

			// Map service errors to appropriate HTTP status codes
			var statusCode int
			switch {
			case errors.Is(err, ErrInvalidUserID):
				statusCode = http.StatusUnauthorized
			case errors.Is(err, ErrStorageStatsNotFound):
				statusCode = http.StatusNotFound
			case errors.Is(err, ErrAnalyticsServiceUnavailable):
				statusCode = http.StatusServiceUnavailable
			default:
				statusCode = http.StatusInternalServerError
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

		// Create the response structure
		response := struct {
			Success bool                `json:"success"`
			UserID  string              `json:"user_id"`
			Stats   *analytics.StorageStats `json:"stats"`
		}{
			Success: true,
			UserID:  userID,
			Stats:   stats,
		}

		// Respond with storage stats
		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}