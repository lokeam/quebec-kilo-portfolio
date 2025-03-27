package middleware

import (
	"net/http"

	"github.com/getsentry/sentry-go"
)

// Captures errors then sends them to Sentry
func SentryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currHub := sentry.CurrentHub().Clone()
		currHub.Scope().SetRequest(r)

		// Add user information if available
		if userID := getUserIDFromRequest(r); userID != "" {
			currHub.Scope().SetUser(sentry.User{
				ID: userID,
			})
		}

		// Add the hub to the request context
		ctx := sentry.SetHubOnContext(r.Context(), currHub)

		// Call next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper function to get user ID
// Note: replace this with a real implementation
func getUserIDFromRequest(r *http.Request) string {
	return r.Header.Get("X-User-ID")
}
