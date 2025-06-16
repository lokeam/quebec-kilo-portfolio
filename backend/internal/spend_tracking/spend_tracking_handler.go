package spend_tracking

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
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
	// BFF route
	r.Get("/bff", handler.GetAllSpendTrackingItemsBFF)
}

func (h *SpendTrackingHandler) GetAllSpendTrackingItemsBFF(w http.ResponseWriter, r *http.Request) {
	requestID := httputils.GetRequestID(r)

	userID := httputils.GetUserID(r)
	if userID == "" {
		h.appContext.Logger.Error("userID not found in request context", map[string]any{
			"request_id": requestID,
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
	// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "physical" key, DO NOT use a struct{}
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

