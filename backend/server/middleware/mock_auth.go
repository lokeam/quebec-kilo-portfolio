package middleware

import (
	"context"
	"net/http"

	"github.com/lokeam/qko-beta/internal/appcontext"
)

type contextKey string

const UserIDKey contextKey = "userID"

func MockAuthMiddleware(appCtx *appcontext.AppContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID := r.Header.Get("X-Mock-UserID")
			if userID == "" {
				userID = "mockUser123" // Default dummy userID
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
