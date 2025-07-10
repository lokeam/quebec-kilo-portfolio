package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/types"
)

type UserHandler struct {
	appContext *appcontext.AppContext
	userService *UserService
}

func NewUserHandler(
	appCtx *appcontext.AppContext,
	userService *UserService,
) *UserHandler {
	return &UserHandler{
		appContext: appCtx,
		userService: userService,
	}
}

// RegisterUserProfileRoutes registers all user profile routes
func RegisterUserProfileRoutes(
	r chi.Router,
	appCtx *appcontext.AppContext,
	userService *UserService,
) {
	handler := NewUserHandler(appCtx, userService)

		// Base routes
	r.Get("/profile", handler.GetUserProfile)
	r.Put("/profile", handler.UpdateUserProfile)
	r.Post("/", handler.CreateUser)

	// Profile completion check
	r.Get("/profile/complete", handler.CheckProfileComplete)
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
