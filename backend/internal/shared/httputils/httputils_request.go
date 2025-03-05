package httputils

import (
	"net/http"

	authMiddleware "github.com/lokeam/qko-beta/server/middleware"
)

const (
	ClientIDHeader = "Client-ID"
	AuthorizationHeader = "Authorization"
	BearerPrefix = "Bearer "
	ContentTypeHeader = "Content-Type"
	ContentTypeText = "text/plain"
	XRequestIDHeader = "X-Request-ID"
	RequestIDKey contextKey = "requestID"
)

type contextKey string

func GetRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(authMiddleware.UserIDKey).(string); ok {
		return id
	}

	return r.Header.Get(XRequestIDHeader)
}

func GetUserID(r *http.Request) string {
	if id, ok := r.Context().Value(authMiddleware.UserIDKey).(string); ok {
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
