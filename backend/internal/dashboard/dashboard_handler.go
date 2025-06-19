package dashboard

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

type DashboardHandler struct {
	appContext *appcontext.AppContext
	dashboardService services.DashboardService
}

// NOTE: not sure if this is needed
type DashboardRequestBody struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

func NewDashboardHandler(
	appCtx *appcontext.AppContext,
	dashboardService services.DashboardService,
) *DashboardHandler {
	return &DashboardHandler{
		appContext: appCtx,
		dashboardService: dashboardService,
	}
}

func RegisterDashboardRoutes(
	r chi.Router,
	appCtx *appcontext.AppContext,
	dashboardService services.DashboardService,
) {
	handler := NewDashboardHandler(appCtx, dashboardService)
	// BFF route
	r.Get("/bff", handler.GetAllDashboardItemsBFF)
}

func (h *DashboardHandler) GetAllDashboardItemsBFF(w http.ResponseWriter, r *http.Request) {
	// Get request ID
	requestID := httputils.GetRequestID(r)

	// Get userID
	userID := httputils.GetUserID(r)
	if userID == "" {
		h.appContext.Logger.Error("userID not found in request context", map[string]any{
			"requestID": requestID,
		})
		h.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
	}

	// Log request
	h.appContext.Logger.Info("Getting all dashboard data", map[string]any{
		"requestID": requestID,
	})

	// Call service
	dashboardItems, err := h.dashboardService.GetDashboardBFFResponse(r.Context(), userID)
	if err != nil {
		h.handleError(w, requestID, err, http.StatusInternalServerError)
		return
	}

	// Build response from service output

	// Use standard response format
	// IMPORTANT: All responses MUST be wrapped in map[string]any{} along with a "dashboard" key, DO NOT use a struct{}
	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"dashboard": dashboardItems,
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		h.appContext.Logger,
		http.StatusOK,
		response,
	)
}

func (h *DashboardHandler) handleError(
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