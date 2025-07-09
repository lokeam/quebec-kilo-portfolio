package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/lokeam/qko-beta/internal/appcontext"
)

// RateLimitMiddleware creates a rate limiting middleware using Redis
// limit: maximum number of requests allowed
// window: time window for the limit (e.g., 1 hour)
func RateLimitMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get app context from request context
			appCtx := r.Context().Value("appContext")
			if appCtx == nil {
				// If no app context, continue without rate limiting
				next.ServeHTTP(w, r)
				return
			}

			appContext := appCtx.(*appcontext.AppContext)
			clientIP := getClientIP(r)

			// Check if rate limited
			if isRateLimited(r.Context(), appContext, clientIP, limit, window) {
				http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
				return
			}

			// Increment request count
			incrementRequestCount(r.Context(), appContext, clientIP, window)

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddlewareWithContext creates a rate limiting middleware using Redis with app context
// limit: maximum number of requests allowed
// window: time window for the limit (e.g., 1 hour)
func RateLimitMiddlewareWithContext(appCtx *appcontext.AppContext, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)

			// Check if rate limited
			if isRateLimited(r.Context(), appCtx, clientIP, limit, window) {
				http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
				return
			}

			// Increment request count
			incrementRequestCount(r.Context(), appCtx, clientIP, window)

			next.ServeHTTP(w, r)
		})
	}
}

// isRateLimited checks if the IP is currently rate limited
func isRateLimited(ctx context.Context, appCtx *appcontext.AppContext, ip string, limit int, window time.Duration) bool {
	key := fmt.Sprintf("rate_limit:general:%s", ip)

	// Get current count from Redis
	countStr, err := appCtx.RedisClient.Get(ctx, key)
	if err != nil || countStr == "" {
		return false // No previous requests, not rate limited
	}

	// Parse count
	currentCount, err := strconv.Atoi(countStr)
	if err != nil {
		return false // Invalid count, treat as no limit
	}

	return currentCount >= limit
}

// incrementRequestCount increments the request count for an IP
func incrementRequestCount(ctx context.Context, appCtx *appcontext.AppContext, ip string, window time.Duration) {
	key := fmt.Sprintf("rate_limit:general:%s", ip)

	// Get current count
	currentStr, err := appCtx.RedisClient.Get(ctx, key)
	currentCount := 0
	if err == nil && currentStr != "" {
		currentCount, _ = strconv.Atoi(currentStr)
	}

	// Increment count
	newCount := currentCount + 1

	// Set new count with expiration
	err = appCtx.RedisClient.Set(ctx, key, strconv.Itoa(newCount), window)
	if err != nil {
		// Log error but don't fail the request
		appCtx.Logger.Error("Failed to increment rate limit counter", map[string]any{
			"error": err,
			"ip": ip,
		})
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check for forwarded headers first (for proxy/load balancer scenarios)
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// Fall back to remote address
	return r.RemoteAddr
}