package spend_tracking

import (
	"encoding/json"
	"errors"
	"net/http"

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
	// r.Post("/", handler.CreateOneTimePurchase)

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
		"spend_tracking": createdOneTimePurchase,
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		h.appContext.Logger,
		http.StatusCreated,
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

