package digital

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/locations/formatters"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
)

type DigitalHandler struct {
	appContext *appcontext.AppContext
	digitalService services.DigitalService
}

func NewDigitalHandler(
	appCtx *appcontext.AppContext,
	digitalService services.DigitalService,
) *DigitalHandler {
	return &DigitalHandler{
		appContext: appCtx,
		digitalService: digitalService,
	}
}


// RegisterDigitalRoutes registers all digital location routes
func RegisterDigitalRoutes(
	r chi.Router,
	appCtx *appcontext.AppContext,
	digitalService services.DigitalService,
) {
	handler := NewDigitalHandler(appCtx, digitalService)
	// Base routes
	r.Get("/", handler.GetAllDigitalLocations)
	r.Post("/", handler.CreateDigitalLocation)
	r.Delete("/", handler.DeleteDigitalLocation)  // Handles both single and bulk deletion

	// Nested routes with ID
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", handler.GetSingleDigitalLocation)
		r.Put("/", handler.UpdateDigitalLocation)
	})

	// BFF route
	r.Get("/bff", handler.GetAllDigitalLocationsBFF)
}

func (dh *DigitalHandler) handleError(
	w http.ResponseWriter,
	requestID string,
	err error,
	statusCode int,
) {
	httputils.RespondWithError(
		httputils.NewResponseWriterAdapter(w),
		dh.appContext.Logger,
		requestID,
		err,
		statusCode,
	)
}


// GetAllDigitalLocations handles GET requests for listing all digital locations
func (dh *DigitalHandler ) GetAllDigitalLocations(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracing
	requestID := httputils.GetRequestID(r)
	userID := httputils.GetUserID(r)

	if userID == "" {
		dh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
			"request_id": requestID,
		})
		dh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	dh.appContext.Logger.Info("Listing digital locations", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	locations, err := dh.digitalService.GetAllDigitalLocations(r.Context(), userID)
	if err != nil {
		dh.handleError(w, requestID, err, http.StatusInternalServerError)
		return
	}

	// Debug: Check subscription data
	for i, loc := range locations {
		dh.appContext.Logger.Debug("Location from DB", map[string]any{
			"index": i,
			"id": loc.ID,
			"name": loc.Name,
			"has_subscription": loc.Subscription != nil,
		})
		if loc.Subscription != nil {
			dh.appContext.Logger.Debug("Subscription details", map[string]any{
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
		dh.appContext.Logger.Debug("Sample transformed location", map[string]any{
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
		dh.appContext.Logger,
		http.StatusOK,
		response,
	)
}


// GetSingleDigitalLocation handles GET requests for a single digital location
func (dh *DigitalHandler) GetSingleDigitalLocation(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracing
	requestID := httputils.GetRequestID(r)

	userID := httputils.GetUserID(r)
	if userID == "" {
		dh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
			"request_id": requestID,
		})
		dh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized,
		)
		return
	}

	locationID := chi.URLParam(r, "id")

	dh.appContext.Logger.Info("Getting digital location", map[string]any{
		"requestID":  requestID,
		"userID":     userID,
		"locationID": locationID,
	})

	if locationID == "" {
		dh.handleError(w, requestID, errors.New("id is required"), http.StatusBadRequest)
		return
	}
	if _, err := uuid.Parse(locationID); err != nil {
		dh.handleError(w, requestID, errors.New("invalid location ID format"), http.StatusBadRequest)
		return
	}

	// NOTE: TODO - replace this with response type
	var location models.DigitalLocation
	var err error

	location, err = dh.digitalService.GetSingleDigitalLocation(r.Context(), userID, locationID)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrDigitalLocationNotFound) {
			statusCode = http.StatusNotFound
		}
		dh.handleError(w, requestID, err, statusCode)
		return
	}

	// Convert to frontend format
	frontendLocation := formatters.FormatDigitalLocationToFrontend(&location)

	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"digital": frontendLocation,
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		dh.appContext.Logger,
		http.StatusOK,
		response,
	)
}


// GetAllDigitalLocationsBFF handles GET requests for the /online-services page
func (dh *DigitalHandler) GetAllDigitalLocationsBFF(w http.ResponseWriter, r *http.Request) {
	// Get request ID
	requestID := httputils.GetRequestID(r)

	// Get userID
	userID := httputils.GetUserID(r)
	if userID == "" {
		dh.appContext.Logger.Error("userID not found in request context", map[string]any{
			"requestID": requestID,
		})
		dh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
	}

	// Log request
	dh.appContext.Logger.Info("Getting all digital location data", map[string]any{
		"requestID": requestID,
		"userID": userID,
	})


	// Call service
	digitalLocations, err := dh.digitalService.GetAllDigitalLocationsBFF(r.Context(), userID)
	if err != nil {
		dh.handleError(w, requestID, err, http.StatusInternalServerError)
		return
	}

	// Use standard response format
	// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "digital" key, DO NOT use a struct{}
	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"digital": digitalLocations,
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		dh.appContext.Logger,
		http.StatusOK,
		response,
	)
}


func (dh *DigitalHandler) CreateDigitalLocation(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracking
	requestID := httputils.GetRequestID(r)

	// Grab user ID from request context
	userID := httputils.GetUserID(r)
	if userID == "" {
		dh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
				"request_id": requestID,
		})
		dh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}
	dh.appContext.Logger.Info("Creating digital location", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	// Parse request
	var locationRequest types.DigitalLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&locationRequest); err != nil {
			dh.handleError(w, requestID, errors.New("invalid request body"), http.StatusBadRequest)
			return
	}

	// Call service method - business logic moved to service
	createdLocation, err := dh.digitalService.CreateDigitalLocation(r.Context(), userID, locationRequest)
	if err != nil {
			statusCode := http.StatusInternalServerError
			dh.appContext.Logger.Error("Failed to create digital location", map[string]any{
					"error":      err,
					"request_id": requestID,
			})
			// Handle specific error types
			if errors.Is(err, ErrDigitalLocationExists) {
					statusCode = http.StatusConflict
			}
			dh.handleError(w, requestID, err, statusCode)
			return
	}

	// Convert to frontend format
	adaptedResponse := formatters.FormatDigitalLocationToFrontend(&createdLocation)

	// IMPORTANT: All responses must use the NewAPIResponse function AND be wrapped in a map
	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"digital": adaptedResponse,
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		dh.appContext.Logger,
		http.StatusCreated,
		response,
	)
}


func (dh *DigitalHandler) UpdateDigitalLocation(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracing
	requestID := httputils.GetRequestID(r)

	userID := httputils.GetUserID(r)
	if userID == "" {
			dh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
					"request_id": requestID,
			})
			dh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
			return
	}

	locationID := chi.URLParam(r, "id")
	dh.appContext.Logger.Info("Updating digital location", map[string]any{
			"requestID":  requestID,
			"userID":     userID,
			"locationID": locationID,
	})

	if locationID == "" {
			dh.handleError(w, requestID, errors.New("location ID is required"), http.StatusBadRequest)
			return
	}

	// Parse request
	var req types.DigitalLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			dh.appContext.Logger.Error("Failed to decode request body", map[string]any{"error": err})
			dh.handleError(w, requestID, errors.New("invalid request body"), http.StatusBadRequest)
			return
	}

	// Set the locationID in the request
	req.ID = locationID

	// Call service method - business logic moved to service
	if err := dh.digitalService.UpdateDigitalLocation(r.Context(), userID, req); err != nil {
			dh.appContext.Logger.Error("Failed to update digital location", map[string]any{"error": err})
			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrDigitalLocationNotFound) {
					statusCode = http.StatusNotFound
			}
			dh.handleError(w, requestID, err, statusCode)
			return
	}

	// Get updated location to return
	updatedLocation, err := dh.digitalService.GetSingleDigitalLocation(r.Context(), userID, locationID)
	if err != nil {
			dh.appContext.Logger.Error("Failed to get updated location", map[string]any{"error": err})
			dh.handleError(w, requestID, errors.New("location was updated but could not be retrieved"), http.StatusInternalServerError)
			return
	}

	// Convert to frontend format
	frontendLocation := formatters.FormatDigitalLocationToFrontend(&updatedLocation)

	// Use standard response format
	response := httputils.NewAPIResponse(r, userID, map[string]any{
			"digital": frontendLocation,
	})

	httputils.RespondWithJSON(
			httputils.NewResponseWriterAdapter(w),
			dh.appContext.Logger,
			http.StatusOK,
			response,
	)
}



// DeleteDigitalLocation handles DELETE requests for removing digital locations
// It supports both single location deletion (via URL param) and bulk deletion (via request body)
func (dh *DigitalHandler) DeleteDigitalLocation(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracing
	requestID := httputils.GetRequestID(r)

	userID := httputils.GetUserID(r)
	if userID == "" {
		dh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
			"request_id": requestID,
		})
		dh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	// Get IDs from query parameters
	digitalLocationIDs := r.URL.Query().Get("ids")
	if digitalLocationIDs == "" {
		dh.handleError(w, requestID, ErrEmptyLocationIDs, http.StatusBadRequest)
		return
	}

	// if multiple IDs are provided, split into array
	digitalLocationIDsArr := strings.Split(digitalLocationIDs, ",")
	dh.appContext.Logger.Info("Deleting digital location(s)", map[string]any{
		"requestID":  requestID,
		"userID":     userID,
		"locationID": digitalLocationIDsArr,
	})

	// Call service method to delete locations
	deletedCount, err := dh.digitalService.DeleteDigitalLocation(r.Context(), userID, digitalLocationIDsArr)
	if err != nil {
		dh.appContext.Logger.Error("Failed to delete digital locations", map[string]any{
			"error": err,
			"request_id": requestID,
			"location_ids": digitalLocationIDsArr,
		})

		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrDigitalLocationNotFound) {
			statusCode = http.StatusNotFound
		}
		dh.handleError(w, requestID, err, statusCode)
		return
	}

	// Log success
	dh.appContext.Logger.Info("Successfully deleted digital locations", map[string]any{
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
		dh.appContext.Logger,
		http.StatusOK,
		response,
	)
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
