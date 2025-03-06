package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

func EnrichRequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get or generated request ID
		requestID := r.Header.Get(httputils.XRequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
			r.Header.Set(httputils.XRequestIDHeader, requestID)
		}

		// Store request ID in context for easy access in handlers
		ctx := context.WithValue(r.Context(), httputils.RequestIDKey, requestID)

		// Carry on to next middleware/handler w/ enriched context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
