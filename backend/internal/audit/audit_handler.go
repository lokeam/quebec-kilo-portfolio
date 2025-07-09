package audit

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

type AuditHandler struct {
	appCtx *appcontext.AppContext
	service *AuditService
}

func NewAuditHandler(appCtx *appcontext.AppContext, service *AuditService) *AuditHandler {
	return &AuditHandler{
		appCtx:  appCtx,
		service: service,
	}
}


// RegisterAuditRoutes registers audit-related routes
func RegisterAuditRoutes(r chi.Router, appCtx *appcontext.AppContext, service *AuditService) {
	handler := NewAuditHandler(appCtx, service)

	// Audit routes
	r.Route("/audit", func(r chi.Router) {
		// Get current user's audit log
		r.Get("/logs", handler.GetUserAuditLog)

		// Admin routes (require admin middleware)
		r.Route("/admin", func(r chi.Router) {
			// Get audit statistics
			r.Get("/stats", handler.GetAuditStats)

			// Cleanup old logs
			r.Post("/cleanup", handler.CleanupOldLogs)

			// Get audit logs for specific user
			r.Get("/logs/user/{userID}", handler.GetAuditLogsByUser)

			// Get audit logs by date range
			r.Get("/logs/range", handler.GetAuditLogsByDateRange)
		})
	})
}


// GetUserAuditLog returns audit logs for the authenticated user
func (h *AuditHandler) GetUserAuditLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userID").(string)

	logs, err := h.service.GetUserAuditLog(ctx, userID)
	if err != nil {
		h.appCtx.Logger.Error("Failed to get user audit log", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		http.Error(w, "Failed to get audit log", http.StatusInternalServerError)
		return
	}

	httputils.RespondWithJSON(
		w,
		h.appCtx.Logger,
		http.StatusOK,
		httputils.APIResponse{
			Success: true,
			UserID:  userID,
			Data:    logs,
		})
}

// GetAuditStats returns audit statistics (admin only)
func (h *AuditHandler) GetAuditStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userID").(string)

	stats, err := h.service.GetAuditStats(ctx)
	if err != nil {
		h.appCtx.Logger.Error("Failed to get audit stats", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		http.Error(w, "Failed to get audit stats", http.StatusInternalServerError)
		return
	}

	httputils.RespondWithJSON(
		w,
		h.appCtx.Logger,
		http.StatusOK,
		httputils.APIResponse{
			Success: true,
			UserID:  userID,
			Data:    stats,
		})
}

// CleanupOldLogs removes audit logs older than retention period (admin only)
func (h *AuditHandler) CleanupOldLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userID").(string)

	if err := h.service.CleanupOldLogs(ctx); err != nil {
		h.appCtx.Logger.Error("Failed to cleanup old logs", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		http.Error(w, "Failed to cleanup old logs", http.StatusInternalServerError)
		return
	}

	httputils.RespondWithJSON(
		w,
		h.appCtx.Logger,
		http.StatusOK,
		httputils.APIResponse{
			Success: true,
			UserID:  userID,
			Data: map[string]string{
				"message": "Old audit logs cleaned up successfully",
			},
		})
}

// GetAuditLogsByUser returns audit logs for a specific user (admin only)
func (h *AuditHandler) GetAuditLogsByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminUserID := ctx.Value("userID").(string)

	// Get user ID from URL parameter
	targetUserID := chi.URLParam(r, "userID")
	if targetUserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	logs, err := h.service.GetUserAuditLog(ctx, targetUserID)
	if err != nil {
		h.appCtx.Logger.Error("Failed to get audit logs for user", map[string]any{
			"admin_user_id": adminUserID,
			"target_user_id": targetUserID,
			"error":         err.Error(),
		})
		http.Error(w, "Failed to get audit logs", http.StatusInternalServerError)
		return
	}

	httputils.RespondWithJSON(
		w,
		h.appCtx.Logger,
		http.StatusOK,
		httputils.APIResponse{
			Success: true,
			UserID:  adminUserID,
			Data: map[string]interface{}{
				"user_id": targetUserID,
				"logs":    logs,
			},
		})
}

// GetAuditLogsByDateRange returns audit logs within a date range (admin only)
func (h *AuditHandler) GetAuditLogsByDateRange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userID").(string)

	// Parse query parameters
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	limitStr := r.URL.Query().Get("limit")

	if startDate == "" || endDate == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	limit := 100 // default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// This would require adding a new method to AuditService
	// For now just return a placeholder response
	httputils.RespondWithJSON(
		w,
		h.appCtx.Logger,
		http.StatusOK,
		httputils.APIResponse{
			Success: true,
			UserID:  userID,
			Data: map[string]any{
				"message": "Date range query not yet implemented",
				"start_date": startDate,
				"end_date":   endDate,
				"limit":      limit,
			},
		})
}
