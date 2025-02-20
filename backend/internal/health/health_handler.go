package health

import (
	"net/http"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/interfaces"
)

// NewHealthHandler returns an http.HandlerFunc which handles health check requests.
func NewHealthHandler(cfg *config.Config, logger interfaces.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the health-check request for debugging.
		logger.Info("health check requested", map[string]any{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
		})

		// Set the response headers and write a simple JSON response.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"available","env":"` + cfg.Env + `","version":"1.0.0"}`))
		if err != nil {
			logger.Error("health check write failed", map[string]any{"error": err.Error()})
		}
	}
}