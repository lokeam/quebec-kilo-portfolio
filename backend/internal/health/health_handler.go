package health

import (
	"net/http"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
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

		// Determine status based on health status
		status := http.StatusOK
		if cfg.HealthStatus == "unavailable" {
			status = http.StatusInternalServerError
		}

		response := map[string]string{
			"status":  cfg.HealthStatus,
			"env":     cfg.Env,
		}

		if err := httputils.RespondWithJSON(w, logger, status, response); err != nil {
			logger.Error("health check write failed", map[string]any{"error": err.Error()})
		}
	}
}