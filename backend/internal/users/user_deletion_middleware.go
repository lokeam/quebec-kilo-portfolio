package users

import (
	"context"
	"net/http"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

// UserDeletionMiddleware checks if the user is deleted or has requested deletion
func UserDeletionMiddleware(appCtx *appcontext.AppContext, service *UserDeletionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userID := ctx.Value("userID")

			if userID == nil {
				// No user ID in context, continue (this should be handled by auth middleware)
				next.ServeHTTP(w, r)
				return
			}

			userIDStr := userID.(string)

			// Get user status
			user, err := service.GetUser(ctx, userIDStr)
			if err != nil {
				appCtx.Logger.Error("Failed to get user in deletion middleware", map[string]any{
					"user_id": userIDStr,
					"error":   err.Error(),
				})
				// Continue with request even if we can't get user status
				next.ServeHTTP(w, r)
				return
			}

			// Check if user is permanently deleted
			if user.IsDeleted() {
				appCtx.Logger.Warn("Access denied to deleted user", map[string]any{
					"user_id": userIDStr,
				})
				http.Error(w, "Account has been permanently deleted", http.StatusForbidden)
				return
			}

			// Check if user has requested deletion and is in grace period
			if user.IsDeletionRequested() {
				appCtx.Logger.Info("User with pending deletion accessing system", map[string]any{
					"user_id": userIDStr,
					"grace_period_end": user.GetDeletionGracePeriodEnd(),
				})
				// Allow access during grace period, but add warning header
				w.Header().Set("X-Account-Status", "pending-deletion")
			}

			// Add user status to context for handlers that need it
			ctx = context.WithValue(ctx, "userStatus", user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// RequireActiveUserMiddleware blocks access for users who have requested deletion
func RequireActiveUserMiddleware(appCtx *appcontext.AppContext, service *UserDeletionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userID := ctx.Value("userID")

			if userID == nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			userIDStr := userID.(string)

			// Get user status
			user, err := service.GetUser(ctx, userIDStr)
			if err != nil {
				appCtx.Logger.Error("Failed to get user in active user middleware", map[string]any{
					"user_id": userIDStr,
					"error":   err.Error(),
				})
				http.Error(w, "Failed to verify user status", http.StatusInternalServerError)
				return
			}

			// Block access for deleted users
			if user.IsDeleted() {
				appCtx.Logger.Warn("Access denied to deleted user", map[string]any{
					"user_id": userIDStr,
				})
				http.Error(w, "Account has been permanently deleted", http.StatusForbidden)
				return
			}

			// Block access for users who have requested deletion
			if user.IsDeletionRequested() {
				appCtx.Logger.Warn("Access denied to user with pending deletion", map[string]any{
					"user_id": userIDStr,
				})
				httputils.RespondWithJSON(
					w,
					appCtx.Logger,
					http.StatusForbidden,
					httputils.APIResponse{
						Success: false,
						UserID:  userIDStr,
						Data: map[string]any{
							"message": "Account deletion has been requested. Access is restricted.",
							"grace_period_end": user.GetDeletionGracePeriodEnd(),
						},
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}