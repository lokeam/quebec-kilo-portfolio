package httputils

import (
	"net/http"

	"github.com/lokeam/qko-beta/internal/shared/constants"
)

// These are HTTP header constants used throughout the application.
const (
	ClientIDHeader = "Client-ID"
	AuthorizationHeader = "Authorization"
	BearerPrefix = "Bearer "
	ContentTypeHeader = "Content-Type"
	ContentTypeText = "text/plain"
	XRequestIDHeader = "X-Request-ID"
	XContentTypeOptions = "X-Content-Type-Options"
	ContentTypeJSON = "application/json"
	ContentTypeHTML = "text/html"
)

func GetRequestID(r *http.Request) string {
	// Try to get from context first
	if id, ok := r.Context().Value(constants.RequestIDKey).(string); ok {
		return id
	}

	// Fall back to header
	return r.Header.Get(XRequestIDHeader)
}

func GetUserID(r *http.Request) string {
	if id, ok := r.Context().Value(constants.UserIDKey).(string); ok {
		return id
	}

	return ""
}

func GetDomainFromRequest(r *http.Request, defaultDomain string) string {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		return defaultDomain
	}

	return domain
}
