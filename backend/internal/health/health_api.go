package health

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
	"github.com/lokeam/qko-beta/internal/shared/logger"
)

// HealthHandler manages HTTP requests that check if our server is running.
// It when called, tells clients:
// - If the server is running ("status": "available")
// - What environment we're in (development, production, etc)
// - What version of the software we're running
type HealthHandler struct {
	config *config.Config
	logger *logger.Logger
}

type HealthChecker interface {
	CheckHealth(w http.ResponseWriter, r *http.Request) error
}

// Creates a new handler for health check requests.
// It needs:
// - cfg: Settings for the application
// - logger: For recording what happens
//
// Returns an http.Handler that serves health check requests.
func NewHealthHandler(cfg *config.Config, logger *logger.Logger) http.Handler {
	handler := &HealthHandler{
		config: cfg,
		logger: logger,
}

r := chi.NewRouter()
// Wrap the handler to handle the error return
r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := handler.CheckHealth(w, r); err != nil {
				handler.logger.Error("health check failed", map[string]any{"error": err})
		}
})

	return r
}

// CheckHealth responds to web requests asking if server is running.
// Stops working if the client disconnects (to save server resources).
//
// Returns an error if:
// - Client disconnects
// - We can't send the response back
func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	select {
	case <- ctx.Done():
		return ctx.Err()
	default:
		h.logger.Info("health check requested", map[string]any{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
		})

		health := map[string]string{
			"status":  "available",
			"env":     h.config.Env,
			"version": "1.0.0",
		}

		return httputils.RespondWithJSON(w, h.logger, http.StatusOK, health)
	}
}
