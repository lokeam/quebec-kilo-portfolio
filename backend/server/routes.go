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

	// Add trailing slash middleware
	mux.Use(middleware.StripSlashes)

	// Add mock authentication middleware
	mux.Use(authMiddleware.MockAuthMiddleware(appContext))

	// Initialize services
	libraryService, err := library.NewGameLibraryService(appContext)
	if err != nil {
		appContext.Logger.Error("Failed to initialize library service", map[string]any{
			"error": err,
		})
	}

	// Create library services map
	libraryServices := make(library.DomainLibraryServices)
	libraryServices["games"] = libraryService

	// Create library handler
	libraryHandler := library.NewLibraryHandler(appContext, libraryServices)


	wishlistService, err := wishlist.NewGameWishlistService(appContext)
	if err != nil {
		appContext.Logger.Error("Failed to initialize wishlist service", map[string]any{
			"error": err,
		})
	}

	// Initialize handlers using single App Context
	searchServiceFactory := search.NewSearchServiceFactory(appContext)
	healthHandler := health.NewHealthHandler(s.Config, s.Logger)

	// Initialize search services
	searchServices := make(search.DomainSearchServices)
	gameSearchService, err := searchServiceFactory.GetService("games")
	if err == nil {
		searchServices["games"] = gameSearchService
	}

	// Create search handler
	searchHandler := search.NewSearchHandler(
		appContext,
		searchServices,  // Pass the map instead of the factory
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
		// Health
		r.Get("/health", healthHandler)

		// Search
		r.Post("/search", searchHandler)

		// Library
		r.Get("/library", libraryHandler)
		r.Post("/library", libraryHandler)
		r.Route("/library/games", func(r chi.Router) {
			r.Delete("/{gameID}", libraryHandler)
		})

		appContext.Logger.Info("Routes registered", map[string]any{
			"health": "/api/v1/health",
			"search": "/api/v1/search",
			"library": "/api/v1/library",
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