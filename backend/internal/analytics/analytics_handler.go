package analytics

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

// RegisterRoutes adds analytics routes to the router
func RegisterRoutes(r chi.Router, appCtx *appcontext.AppContext, service Service) {
	r.Get("/", GetAnalytics(appCtx, service))
}

// GetAnalytics handles requests for analytics data
func GetAnalytics(appCtx *appcontext.AppContext, service Service) http.HandlerFunc {
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
				errors.New("userID not found in request context"),
				http.StatusUnauthorized,
			)
			return
		}

		// Process domains from query parameters with a simpler approach
		domains := make([]string, 0)
		queryParams := r.URL.Query()["domains"]

		// If we have query parameters, process them
		if len(queryParams) > 0 {
			// Join all parameters with commas and split once
			allParams := strings.Join(queryParams, ",")
			parts := strings.Split(allParams, ",")

			// Process each part in a single loop
			for i := 0; i < len(parts); i++ {
				if trimmed := strings.TrimSpace(parts[i]); trimmed != "" {
					domains = append(domains, trimmed)
				}
			}
		}

		// If no domains were found or specified, use the default
		if len(domains) == 0 {
			domains = []string{DomainGeneral}
		}

		appCtx.Logger.Info("Getting analytics data", map[string]any{
			"requestID": requestID,
			"userID":    userID,
			"domains":   domains,
		})

		// Call service to get analytics data
		data, err := service.GetAnalytics(r.Context(), userID, domains)
		if err != nil {
			appCtx.Logger.Error("Failed to get analytics data", map[string]any{
				"requestID": requestID,
				"userID":    userID,
				"domains":   domains,
				"error":     err.Error(),
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

		// Format storage data if it's included in the response
		if storageData, ok := data["storage"].(*StorageStats); ok {
			// Log data before formatting
			fmt.Printf("\nData before formatting:\n")
			for _, loc := range storageData.DigitalLocations {
				fmt.Printf("Location: %s\n", loc.Name)
				fmt.Printf("  Billing Cycle: %v\n", loc.BillingCycle)
				fmt.Printf("  Cost Per Cycle: %v\n", loc.CostPerCycle)
				fmt.Printf("  Monthly Cost: %v\n", loc.MonthlyCost)
			}
			FormatStorageStats(storageData)
		}

		// Create the response structure
		response := httputils.NewAPIResponse(r, userID, map[string]any{
			"analytics": data,
		})

		// Respond with analytics data
		httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			appCtx.Logger,
			http.StatusOK,
			response,
		)
	}
}
