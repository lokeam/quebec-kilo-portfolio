package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
)

// RequestDeletionRequest represents a request to delete a user account
type RequestDeletionRequest struct {
	Reason string `json:"reason"`
}

// RequestDeletionResponse represents the response to a deletion request
type RequestDeletionResponse struct {
	Message     string     `json:"message"`
	GracePeriod *time.Time `json:"grace_period_end,omitempty"`
}

type UserHandler struct {
	appContext *appcontext.AppContext
	userService *UserService
	userDeletionService *UserDeletionService
}

func NewUserHandler(
	appCtx *appcontext.AppContext,
	userService *UserService,
	userDeletionService *UserDeletionService,
) *UserHandler {
	return &UserHandler{
		appContext: appCtx,
		userService: userService,
		userDeletionService: userDeletionService,
	}
}

// RegisterUserRoutes registers all user-related routes (profile + deletion)
func RegisterUserRoutes(
	r chi.Router,
	appCtx *appcontext.AppContext,
	userService *UserService,
	userDeletionService *UserDeletionService,
) {
	handler := NewUserHandler(appCtx, userService, userDeletionService)

	// Profile routes
	r.Get("/profile", handler.GetUserProfile)
	r.Put("/profile", handler.UpdateUserProfile)
	r.Post("/", handler.CreateUser)
	r.Get("/profile/complete", handler.CheckProfileComplete)

	// Deletion routes
	r.Post("/deletion/request", handler.RequestDeletion)
	r.Post("/deletion/cancel", handler.CancelDeletionRequest)
	r.Get("/deletion/status", handler.GetDeletionStatus)
}

func (uh *UserHandler) handleError(
	w http.ResponseWriter,
	requestID string,
	err error,
	statusCode int,
) {
	httputils.RespondWithError(
		httputils.NewResponseWriterAdapter(w),
		uh.appContext.Logger,
		requestID,
		err,
		statusCode,
	)
}

// GetUserProfile handles GET requests for retrieving user profile
func (uh *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracing
	requestID := httputils.GetRequestID(r)
	userID := httputils.GetUserID(r)

	if userID == "" {
		uh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
			"request_id": requestID,
		})
		uh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	uh.appContext.Logger.Info("Getting user profile", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	user, err := uh.userService.GetSingleUser(r.Context(), userID)
	if err != nil {
		uh.appContext.Logger.Error("Failed to get user profile", map[string]any{
			"error":      err,
			"request_id": requestID,
			"userID":     userID,
		})
		statusCode := http.StatusInternalServerError
		if errors.Is(err, errors.New("user not found")) {
			statusCode = http.StatusNotFound
		}
		uh.handleError(w, requestID, err, statusCode)
		return
	}

	// Convert to frontend format (only return safe fields)
	frontendUser := map[string]any{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"user": frontendUser,
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		uh.appContext.Logger,
		http.StatusOK,
		response,
	)
}

// CreateUser handles POST requests for creating a new user
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracking
	requestID := httputils.GetRequestID(r)

	// Grab user ID from request context
	userID := httputils.GetUserID(r)
	if userID == "" {
		uh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
			"request_id": requestID,
		})
		uh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	uh.appContext.Logger.Info("Creating user", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	// Parse request
	var userRequest types.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		uh.appContext.Logger.Error("Failed to decode request body", map[string]any{"error": err})
		uh.handleError(w, requestID, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	// Set the Auth0 user ID from the context
	userRequest.Auth0UserID = userID

	// Call service method
	createdUser, err := uh.userService.CreateUser(r.Context(), userRequest)
	if err != nil {
		uh.appContext.Logger.Error("Failed to create user", map[string]any{
			"error":      err,
			"request_id": requestID,
			"userID":     userID,
		})
		statusCode := http.StatusInternalServerError
		// Handle specific error types if needed
		uh.handleError(w, requestID, err, statusCode)
		return
	}

	// Convert to frontend format
	frontendUser := map[string]any{
		"id":         createdUser.ID,
		"email":      createdUser.Email,
		"first_name": createdUser.FirstName,
		"last_name":  createdUser.LastName,
		"created_at": createdUser.CreatedAt,
		"updated_at": createdUser.UpdatedAt,
	}

	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"user": frontendUser,
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		uh.appContext.Logger,
		http.StatusCreated,
		response,
	)
}

// UpdateUserProfile handles PUT requests for updating user profile
func (uh *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracing
	requestID := httputils.GetRequestID(r)

	userID := httputils.GetUserID(r)
	if userID == "" {
		uh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
			"request_id": requestID,
		})
		uh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	uh.appContext.Logger.Info("Updating user profile", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	// Parse request
	var req types.UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uh.appContext.Logger.Error("Failed to decode request body", map[string]any{"error": err})
		uh.handleError(w, requestID, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	// Call service method
	updatedUser, err := uh.userService.UpdateUserProfile(r.Context(), userID, req)
	if err != nil {
		uh.appContext.Logger.Error("Failed to update user profile", map[string]any{
			"error":      err,
			"request_id": requestID,
			"userID":     userID,
		})
		statusCode := http.StatusInternalServerError
		if errors.Is(err, errors.New("user not found")) {
			statusCode = http.StatusNotFound
		}
		uh.handleError(w, requestID, err, statusCode)
		return
	}

	// Convert to frontend format
	frontendUser := map[string]any{
		"id":         updatedUser.ID,
		"email":      updatedUser.Email,
		"first_name": updatedUser.FirstName,
		"last_name":  updatedUser.LastName,
		"created_at": updatedUser.CreatedAt,
		"updated_at": updatedUser.UpdatedAt,
	}

	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"user": frontendUser,
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		uh.appContext.Logger,
		http.StatusOK,
		response,
	)
}



// CheckProfileComplete handles GET requests for checking if user profile is complete
func (uh *UserHandler) CheckProfileComplete(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracing
	requestID := httputils.GetRequestID(r)
	userID := httputils.GetUserID(r)

	if userID == "" {
		uh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
			"request_id": requestID,
		})
		uh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	uh.appContext.Logger.Info("Checking user profile completeness", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	hasComplete, err := uh.userService.HasCompleteProfile(r.Context(), userID)
	if err != nil {
		uh.appContext.Logger.Error("Failed to check profile completeness", map[string]any{
			"error":      err,
			"request_id": requestID,
			"userID":     userID,
		})
		uh.handleError(w, requestID, err, http.StatusInternalServerError)
		return
	}

	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"user": map[string]any{
			"has_complete_profile": hasComplete,
			"user_id":             userID,
		},
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		uh.appContext.Logger,
		http.StatusOK,
		response,
	)
}

// RequestDeletion handles POST requests for requesting user account deletion
func (uh *UserHandler) RequestDeletion(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracing
	requestID := httputils.GetRequestID(r)
	userID := httputils.GetUserID(r)

	if userID == "" {
		uh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
			"request_id": requestID,
		})
		uh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	uh.appContext.Logger.Info("Requesting user account deletion", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	var req RequestDeletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uh.appContext.Logger.Error("Failed to decode deletion request body", map[string]any{"error": err})
		uh.handleError(w, requestID, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	// Validate reason
	if req.Reason == "" {
		uh.handleError(w, requestID, errors.New("deletion reason is required"), http.StatusBadRequest)
		return
	}

	err := uh.userDeletionService.RequestDeletion(r.Context(), userID, req.Reason)
	if err != nil {
		uh.appContext.Logger.Error("Failed to request user account deletion", map[string]any{
			"error":      err,
			"request_id": requestID,
			"userID":     userID,
		})
		statusCode := http.StatusInternalServerError
		if errors.Is(err, errors.New("user not found")) {
			statusCode = http.StatusNotFound
		}
		uh.handleError(w, requestID, err, statusCode)
		return
	}

	// Get user to calculate grace period
	user, err := uh.userDeletionService.GetUser(r.Context(), userID)
	if err != nil {
		uh.appContext.Logger.Error("Failed to get user after deletion request", map[string]any{
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

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		uh.appContext.Logger,
		http.StatusOK,
		response,
	)
}

// CancelDeletionRequest handles POST requests for canceling a user account deletion request
func (uh *UserHandler) CancelDeletionRequest(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracing
	requestID := httputils.GetRequestID(r)
	userID := httputils.GetUserID(r)

	if userID == "" {
		uh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
			"request_id": requestID,
		})
		uh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	uh.appContext.Logger.Info("Cancelling user account deletion request", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	err := uh.userDeletionService.CancelDeletionRequest(r.Context(), userID)
	if err != nil {
		uh.appContext.Logger.Error("Failed to cancel user account deletion request", map[string]any{
			"error":      err,
			"request_id": requestID,
			"userID":     userID,
		})
		statusCode := http.StatusInternalServerError
		if errors.Is(err, errors.New("user not found")) {
			statusCode = http.StatusNotFound
		}
		uh.handleError(w, requestID, err, statusCode)
		return
	}

	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"message": "Deletion request cancelled",
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		uh.appContext.Logger,
		http.StatusOK,
		response,
	)
}

// GetDeletionStatus handles GET requests for retrieving the status of a user account deletion request
func (uh *UserHandler) GetDeletionStatus(w http.ResponseWriter, r *http.Request) {
	// Get Request ID for tracing
	requestID := httputils.GetRequestID(r)
	userID := httputils.GetUserID(r)

	if userID == "" {
		uh.appContext.Logger.Error("userID NOT FOUND in request context", map[string]any{
			"request_id": requestID,
		})
		uh.handleError(w, requestID, errors.New("userID not found in request context"), http.StatusUnauthorized)
		return
	}

	uh.appContext.Logger.Info("Getting user account deletion status", map[string]any{
		"requestID": requestID,
		"userID":    userID,
	})

	user, err := uh.userDeletionService.GetUser(r.Context(), userID)
	if err != nil {
		uh.appContext.Logger.Error("Failed to get user account deletion status", map[string]any{
			"error":      err,
			"request_id": requestID,
			"userID":     userID,
		})
		statusCode := http.StatusInternalServerError
		if errors.Is(err, errors.New("user not found")) {
			statusCode = http.StatusNotFound
		}
		uh.handleError(w, requestID, err, statusCode)
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

	response := httputils.NewAPIResponse(r, userID, map[string]any{
		"status": status,
	})

	httputils.RespondWithJSON(
		httputils.NewResponseWriterAdapter(w),
		uh.appContext.Logger,
		http.StatusOK,
		response,
	)
}
