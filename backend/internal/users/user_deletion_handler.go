package users

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

type UserDeletionHandler struct {
	appCtx *appcontext.AppContext
	service *UserDeletionService
}

func NewUserDeletionHandler(appCtx *appcontext.AppContext, service *UserDeletionService) *UserDeletionHandler {
	return &UserDeletionHandler{
		appCtx:  appCtx,
		service: service,
	}
}

// RequestDeletionRequest represents a request to delete a user account
type RequestDeletionRequest struct {
	Reason string `json:"reason"`
}

// RequestDeletionResponse represents the response to a deletion request
type RequestDeletionResponse struct {
	Message     string     `json:"message"`
	GracePeriod *time.Time `json:"grace_period_end,omitempty"`
}

// RequestDeletion handles a user's request to delete their account
func (h *UserDeletionHandler) RequestDeletion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userID").(string)

	// Parse request
	var req RequestDeletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate reason
	if req.Reason == "" {
		http.Error(w, "Deletion reason is required", http.StatusBadRequest)
		return
	}

	// Request deletion
	if err := h.service.RequestDeletion(ctx, userID, req.Reason); err != nil {
		h.appCtx.Logger.Error("Failed to request deletion", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		http.Error(w, "Failed to request deletion", http.StatusInternalServerError)
		return
	}

	// Get user to calculate grace period
	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		h.appCtx.Logger.Error("Failed to get user after deletion request", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		// Don't fail the request, just don't include grace period
	}

	response := RequestDeletionResponse{
		Message: "Account deletion requested. Your account will be permanently deleted in 30 days unless you cancel the request.",
	}

	if user != nil {
		gracePeriod := user.GetDeletionGracePeriodEnd()
		if gracePeriod != nil {
			response.GracePeriod = gracePeriod
		}
	}

	httputils.RespondWithJSON(w, h.appCtx.Logger, http.StatusOK, httputils.APIResponse{
		Success: true,
		UserID:  userID,
		Data:    response,
	})
}

// CancelDeletionRequest handles canceling a pending deletion request
func (h *UserDeletionHandler) CancelDeletionRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userID").(string)

	// Cancel deletion request
	if err := h.service.CancelDeletionRequest(ctx, userID); err != nil {
		h.appCtx.Logger.Error("Failed to cancel deletion request", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		http.Error(w, "Failed to cancel deletion request", http.StatusInternalServerError)
		return
	}

	httputils.RespondWithJSON(w, h.appCtx.Logger, http.StatusOK, httputils.APIResponse{
		Success: true,
		UserID:  userID,
		Data: map[string]string{
			"message": "Account deletion request cancelled. Your account is now active.",
		},
	})
}

// GetDeletionStatus returns the current deletion status for a user
func (h *UserDeletionHandler) GetDeletionStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("userID").(string)

	// Get user
	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		h.appCtx.Logger.Error("Failed to get user", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	status := map[string]interface{}{
		"is_active":              user.IsActive(),
		"is_deleted":             user.IsDeleted(),
		"is_deletion_requested":  user.IsDeletionRequested(),
		"is_in_grace_period":     user.IsInGracePeriod(),
	}

	if user.IsDeletionRequested() {
		gracePeriod := user.GetDeletionGracePeriodEnd()
		if gracePeriod != nil {
			status["grace_period_end"] = gracePeriod
		}
		if user.DeletionReason != nil {
			status["deletion_reason"] = *user.DeletionReason
		}
	}

	httputils.RespondWithJSON(w, h.appCtx.Logger, http.StatusOK, httputils.APIResponse{
		Success: true,
		UserID:  userID,
		Data:    status,
	})
}

// RegisterUserRoutes registers user-related routes
func RegisterUserRoutes(r chi.Router, appCtx *appcontext.AppContext, service *UserDeletionService) {
	handler := NewUserDeletionHandler(appCtx, service)

	// User deletion routes
	r.Route("/users", func(r chi.Router) {
		// Request account deletion
		r.Post("/deletion/request", handler.RequestDeletion)

		// Cancel deletion request
		r.Post("/deletion/cancel", handler.CancelDeletionRequest)

		// Get deletion status
		r.Get("/deletion/status", handler.GetDeletionStatus)
	})
}