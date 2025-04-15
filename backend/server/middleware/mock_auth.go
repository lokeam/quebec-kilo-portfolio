package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/appcontext"
)

type contextKey string

const UserIDKey contextKey = "userID"

// Mock user ID for testing
const MockUserID = "9a4aeee6-fb31-4839-a921-f61b0525046d"

func MockAuthMiddleware(appCtx *appcontext.AppContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Header.Get("X-Mock-UserID")
			if userID == "" {
				// Use consistent mock user ID for testing
				userID = MockUserID
			} else {
				// Validate that the provided ID is a valid UUID
				if _, err := uuid.Parse(userID); err != nil {
					appCtx.Logger.Error("Invalid UUID format in X-Mock-UserID header", map[string]any{
						"userID": userID,
						"error": err,
					})
					http.Error(w, "Invalid UUID format in X-Mock-UserID header", http.StatusBadRequest)
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
