package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lokeam/qko-beta/internal/health"
	auth0middleware "github.com/lokeam/qko-beta/internal/shared/middleware"
)

func (s *Server) Routes() http.Handler {
	// Create Router
	mux := chi.NewRouter()

	// Setup middleware
	s.setupMiddleware(mux)
	s.setupCORS(mux)

	// Initialize Routes
	mux.Route("/api/v1", func(r chi.Router) {

		// Core + utility routes
		r.Mount("/health", health.NewHealthHandler(s.config, s.logger))

		// Feature routes (NOTE: add to protected routes post testing)

		// Feature routes (NOTE: add to protected routes post testing)

		// Protected routes
		r.Use(auth0middleware.EnsureValidToken())
		// Mounted protected features below
	})

	return mux
}

func (s *Server) setupMiddleware(mux *chi.Mux) {
	mux.Use(middleware.RequestID)  // Trace requests
	mux.Use(middleware.RealIP)     // Get actual client IP
	mux.Use(middleware.Logger)     // Request logging
	mux.Use(middleware.Recoverer)  // Panic recovery
	mux.Use(middleware.Timeout(60 * time.Second))
}

func (s *Server) setupCORS(mux *chi.Mux) {
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:      s.config.CORS.AllowedOrigins,
		AllowedMethods:      s.config.CORS.AllowedMethods,
		AllowedHeaders:      s.config.CORS.AllowedHeaders,
		ExposedHeaders:      s.config.CORS.ExposedHeaders,
		AllowCredentials:    s.config.CORS.AllowCredentials,
		MaxAge:              s.config.CORS.MaxAge,
	}))
}
