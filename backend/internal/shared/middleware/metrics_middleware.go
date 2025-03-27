package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/lokeam/qko-beta/internal/monitoring"
)

// Collects metrics for Prometheus
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create response writer to grab status code
		responseWriter := NewResponseWriter(w)

		// Call next handler
		next.ServeHTTP(responseWriter, r)

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		monitoring.HTTPRequestDuration.WithLabelValues(
			r.Method,
			r.URL.Path,
			strconv.Itoa(responseWriter.Status()),
		).Observe(duration)
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// New ResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Status() int {
	return rw.statusCode
}
