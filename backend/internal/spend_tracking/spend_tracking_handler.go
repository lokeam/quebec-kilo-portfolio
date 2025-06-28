package spend_tracking

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
)

type SpendTrackingHandler struct {
	appContext *appcontext.AppContext
	spendTrackingService services.SpendTrackingService
}

type SpendTrackingRequestBody struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

func NewSpendTrackingHandler(
	appCtx *appcontext.AppContext,
	spendTrackingService services.SpendTrackingService,
) *SpendTrackingHandler {
	return &SpendTrackingHandler{
		appContext: appCtx,
		spendTrackingService: spendTrackingService,
	}
}

func RegisterSpendTrackingRoutes(
	r chi.Router,
	appCtx *appcontext.AppContext,
	spendTrackingService services.SpendTrackingService,
) {
	handler := NewSpendTrackingHandler(appCtx, spendTrackingService)
	// Base routes
	r.Post("/", handler.CreateOneTimePurchase)
	r.Delete("/", handler.DeleteSpendTrackingItems)

	// Nested routes with ID
	r.Route("/{id}", func(r chi.Router) {
		r.Put("/", handler.UpdateOneTimePurchase)
	})

	// BFF route
	r.Get("/bff", handler.GetAllSpendTrackingItemsBFF)
}

func (h *SpendTrackingHandler) GetAllSpendTrackingItemsBFF(w http.ResponseWriter, r *http.Request) {
	requestID := httputils.GetRequestID(r)

	userID := httputils.GetUserID(r)
	if userID == "" {
		h.appContext.Logger.Error("userID not found in request context", map[string]any{
			"requestID": requestID,
		})
		h.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	h.appContext.Logger.Info("Getting all spend tracking items", map[string]any{
		"requestID": requestID,
	})

	spendTrackingItems, err := h.spendTrackingService.GetSpendTrackingBFFResponse(r.Context(), userID)
	if err != nil {
		h.handleError(w, requestID, err, http.StatusInternalServerError)
		return
	}

	// Use standard response format
	// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "spend_tracking" key, DO NOT use a struct{}
	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"spend_tracking": spendTrackingItems,
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		h.appContext.Logger,
		http.StatusOK,
		response,
	)
}

func (h *SpendTrackingHandler) CreateOneTimePurchase(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracking
	requestID := httputils.GetRequestID(r)

	// Grab user ID from request context
	userID := httputils.GetUserID(r)
	if userID == "" {
		h.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	// Log request
	h.appContext.Logger.Info("Creating one-time purchase", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	// Parse request
	var spendTrackingRequest types.SpendTrackingRequest
	if err := json.NewDecoder(r.Body).Decode(&spendTrackingRequest); err != nil {
		h.handleError(w, requestID, err, http.StatusBadRequest)
		return
	}

	// Call service method
	createdOneTimePurchase, err := h.spendTrackingService.CreateOneTimePurchase(r.Context(), userID, spendTrackingRequest)
	if err != nil {
		h.handleError(w, requestID, err, http.StatusInternalServerError)
		return
	}

	// Convert to frontend format
	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"spend_tracking": map[string]any{
			"item": createdOneTimePurchase,
		},
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		h.appContext.Logger,
		http.StatusCreated,
		response,
	)
}

func (h *SpendTrackingHandler) UpdateOneTimePurchase(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracing
	requestID := httputils.GetRequestID(r)

	userID := httputils.GetUserID(r)
	if userID == "" {
		h.appContext.Logger.Error("userID not found in request context", map[string]any{
			"requestID": requestID,
		})
		h.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	oneTimePurchaseID := chi.URLParam(r, "id")
	h.appContext.Logger.Info("Updating one-time purchase", map[string]any{
		"requestID": requestID,
		"userID":    userID,
		"id":        oneTimePurchaseID,
	})

	if oneTimePurchaseID == "" {
		h.handleError(w, requestID, errors.New("one-time purchase ID is required"), http.StatusBadRequest)
		return
	}

	var req types.SpendTrackingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.appContext.Logger.Error("Failed to decode request body", map[string]any{
			"requestID": requestID,
			"error":     err,
		})
		h.handleError(w, requestID, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	// Set the ID in the request
	req.ID = oneTimePurchaseID

	// Call service method
	if err := h.spendTrackingService.UpdateOneTimePurchase(r.Context(), userID, req); err != nil {
		h.appContext.Logger.Error("Failed to update one-time purchase", map[string]any{
			"requestID": requestID,
			"error":     err,
		})
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrSpendTrackingItemNotFound) {
			statusCode = http.StatusNotFound
		}
		h.handleError(w, requestID, err, statusCode)
		return
	}

	// Get updated one time purchase and return
	updatedOneTimePurchase, err := h.spendTrackingService.GetSingleSpendTrackingItem(r.Context(), userID, oneTimePurchaseID)
	if err != nil {
		h.appContext.Logger.Error("Failed to get updated one-time purchase", map[string]any{
			"requestID": requestID,
			"error":     err,
		})
	}

	// Use standard response format
	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"spend_tracking": map[string]any{
			"item": updatedOneTimePurchase,
		},
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		h.appContext.Logger,
		http.StatusOK,
		response,
	)
}

func (h *SpendTrackingHandler) DeleteSpendTrackingItems(w http.ResponseWriter, r *http.Request) {
	requestID := httputils.GetRequestID(r)

	userID := httputils.GetUserID(r)
	if userID == "" {
		h.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
			"request_id": requestID,
		})
		h.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	// Get IDs from query parameters
	spendTrackingIDs := r.URL.Query().Get("ids")
	if spendTrackingIDs == "" {
		h.handleError(w, requestID, ErrEmptySpendTrackingIDs, http.StatusBadRequest)
		return
	}

	// If multiple IDs are provided, split into array
	spendTrackingIDsArr := strings.Split(spendTrackingIDs, ",")
	h.appContext.Logger.Info("Deleting spend tracking item(s)", map[string]any{
		"requestID":  requestID,
		"userID":     userID,
		"locationID": spendTrackingIDsArr,
	})

	// Call service method to delete locations
	deletedCount, err := h.spendTrackingService.DeleteSpendTrackingItems(r.Context(), userID, spendTrackingIDsArr)
	if err != nil {
		h.appContext.Logger.Error("Failed to delete spend tracking items", map[string]any{
			"error": err,
			"request_id": requestID,
			"location_ids": spendTrackingIDsArr,
		})

		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrSpendTrackingItemNotFound) {
			statusCode = http.StatusNotFound
		}
		h.handleError(w, requestID, err, statusCode)
		return
	}



	// Log success
	h.appContext.Logger.Info("Successfully deleted digital locations", map[string]any{
		"request_id": requestID,
		"user_id": userID,
		"deleted_count": deletedCount,
		"total_count": len(spendTrackingIDsArr),
	})

	// Use standard response format
	// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "digital" key, DO NOT use a struct{}
	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"spend_tracking": map[string]any{
			"success": true,
			"deleted_count": deletedCount,
			"spend_tracking_ids": spendTrackingIDsArr,
		},
	})



	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		h.appContext.Logger,
		http.StatusOK,
		response,
	)
}



// helper fn to standardize error handling
func (h *SpendTrackingHandler) handleError(
	w http.ResponseWriter,
	requestID string,
	err error,
	statusCode int,
) {
	httputils.RespondWithError(
		httputils.NewResponseWriterAdapter(w),
		h.appContext.Logger,
		requestID,
		err,
		statusCode,
	)
}

