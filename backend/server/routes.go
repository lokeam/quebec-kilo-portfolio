package server

import (
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/health"
	"github.com/lokeam/qko-beta/internal/library"
	"github.com/lokeam/qko-beta/internal/search"
	"github.com/lokeam/qko-beta/internal/wishlist"
	authMiddleware "github.com/lokeam/qko-beta/server/middleware"
)

func (s *Server) SetupRoutes(appContext *appcontext.AppContext) chi.Router {
	// Create Router
	mux := chi.NewRouter()

	// Setup middleware
	s.setupMiddleware(mux)
	s.setupCORS(mux)

	// Add mock authentication middleware
	mux.Use(authMiddleware.MockAuthMiddleware(appContext))

	// Initialize services
	libraryService := library.NewGameLibraryService(appContext)
	wishlistService := wishlist.NewGameWishlistService(appContext)

	// Initialize handlers using single App Context
	searchServiceFactory := search.NewSearchServiceFactory(appContext)
	healthHandler := health.NewHealthHandler(s.Config, s.Logger)
	searchHandler := search.NewSearchHandler(
		appContext,
		searchServiceFactory,
		libraryService,
		wishlistService,
	)

	// Debug logging for searchHandler
	if searchHandler == nil {
		appContext.Logger.Error("searchHandler is nil", map[string]any{
			"appContext": appContext,
		})
	} else {
		appContext.Logger.Info("searchHandler initialized successfully", map[string]any{
			"appContext": appContext,
		})
	}

	if reflect.TypeOf(searchHandler) != reflect.TypeOf(http.HandlerFunc(nil)) {
		appContext.Logger.Error("searchHandler is not of type http.HandlerFunc", map[string]any{
			"appContext": appContext,
		})
	} else {
		appContext.Logger.Info("searchHandler is of type http.HandlerFunc", map[string]any{
			"appContext": appContext,
		})
	}

	// Initialize Routes
	mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler)
		r.Post("/search", searchHandler) // Use searchHandler directly
		r.Post("/search/", searchHandler) // Handle trailing slash
		appContext.Logger.Info("Routes registered", map[string]any{
			"health": "/api/v1/health",
			"search": "/api/v1/search",
		})
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
		AllowedOrigins:   s.AppContext.Config.CORS.AllowedOrigins,
		AllowedMethods:   s.AppContext.Config.CORS.AllowedMethods,
		AllowedHeaders:   s.AppContext.Config.CORS.AllowedHeaders,
		ExposedHeaders:   s.AppContext.Config.CORS.ExposedHeaders,
		AllowCredentials: s.AppContext.Config.CORS.AllowCredentials,
		MaxAge:           s.AppContext.Config.CORS.MaxAge,
	}))
}