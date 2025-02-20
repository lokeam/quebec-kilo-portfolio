package server

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/health"
	"github.com/lokeam/qko-beta/internal/search"
)

func (s *Server) SetupRoutes(appContext *appcontext.AppContext) chi.Router {
	// Create Router
	mux := chi.NewRouter()

	// Setup middleware
	s.setupMiddleware(mux)
	s.setupCORS(mux)

	// Initialize handlers using single App Context
	searchHandler := search.NewSearchHandler(appContext)
	healthHandler := health.NewHealthHandler(s.Config, s.Logger)

	// Initialize Routes
	mux.Route("/api/v1", func(r chi.Router) {

		// Protected routes
		//r.Use(auth0middleware.EnsureValidToken())
		// Mounted protected features below

		// Core + utility routes
		r.Get("/health", healthHandler)

		// Feature routes (NOTE: add to protected routes post testing)
		r.Handle("/search", searchHandler)

		// Feature routes (NOTE: add to protected routes post testing)
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
		AllowedOrigins:      s.AppContext.Config.CORS.AllowedOrigins,
		AllowedMethods:      s.AppContext.Config.CORS.AllowedMethods,
		AllowedHeaders:      s.AppContext.Config.CORS.AllowedHeaders,
		ExposedHeaders:      s.AppContext.Config.CORS.ExposedHeaders,
		AllowCredentials:    s.AppContext.Config.CORS.AllowCredentials,
		MaxAge:              s.AppContext.Config.CORS.MaxAge,
	}))
}
