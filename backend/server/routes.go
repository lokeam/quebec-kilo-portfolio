package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	auth0middleware "github.com/lokeam/qko-beta/internal/shared/middleware"
)

func (s *Server) Routes() http.Handler {
	// Create Router
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)  // Trace requests
	mux.Use(middleware.RealIP)     // Get actual client IP
	mux.Use(middleware.Logger)     // Request logging
	mux.Use(middleware.Recoverer)  // Panic recovery
	mux.Use(middleware.Timeout(60 * time.Second))

	// Initialize CORS
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:      []string{"https://*", "http://*"},
		AllowedMethods:      []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:      []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:      []string{"Link"},
		AllowCredentials:    true,
		MaxAge:              300,
	}))

	// Initialize Routes
	mux.Route("/api/v1", func(r chi.Router) {

		// Public routes (no auth required)
		r.Group(func(r chi.Router) {
			// Health Check - Documented in health.yaml
			r.Get("/health", s.handleHealthCheck)
		})

		// Protected routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(auth0middleware.EnsureValidToken())
			  // Games endpoints will go here
        // r.Get("/games", s.handleListGames)
        // r.Post("/games", s.handleCreateGame)
		})

	})

	return mux
}

/* Handler fns (move to appropriate handler files later) */
// Utility - health check
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("health check requested", map[string]any{
		"method": r.Method,
		"path": r.URL.Path,
		"remote_addr": r.RemoteAddr,
	})

	health := map[string]string{
		"status": "available",
		"env":    s.config.Env,
		"version": "1.0.0",
	}

	s.writeJSON(w, http.StatusOK, health, nil)
}

