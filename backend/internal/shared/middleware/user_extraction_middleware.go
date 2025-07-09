package middleware

import (
	"context"
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/lokeam/qko-beta/internal/shared/constants"
)

// UserExtractionMiddleware extracts user ID from validated JWT token
func UserExtractionMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("UserExtractionMiddleware: Processing request for %s", r.URL.Path)

			// Get validated token from context
			token := r.Context().Value(jwtmiddleware.ContextKey{})
			if token == nil {
				log.Printf("No token found in context")
				http.Error(w, "No token found", http.StatusUnauthorized)
				return
			}

			claims, ok := token.(*validator.ValidatedClaims)
			if !ok {
				log.Printf("Failed to cast token to ValidatedClaims")
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			// Extract user ID from claims
			userID := claims.RegisteredClaims.Subject
			if userID == "" {
				log.Printf("No user ID found in token")
				http.Error(w, "No user ID in token", http.StatusUnauthorized)
				return
			}

			log.Printf("Extracted userID from JWT: %s", userID)

			// Add userID to request context
			ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
			r = r.WithContext(ctx)
			log.Printf("Successfully added userID to context: %s", userID)

			// Pass the updated request with context to the next handler
			next.ServeHTTP(w, r)
		})
	}
}