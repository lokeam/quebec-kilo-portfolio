package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
)

type Server struct {
	Config      *config.Config
	AppContext  *appcontext.AppContext
	Logger      interfaces.Logger
	Router      chi.Router
}

func NewServer(
	config *config.Config,
	logger interfaces.Logger,
	appContext *appcontext.AppContext,
) *Server {

	s := &Server{
		Config:      config,
		Logger:      logger,
		AppContext:  appContext,
		Router:      chi.NewRouter(),
	}

	s.Router = s.SetupRoutes(appContext)
	return s
}

// Add http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// Gracefully shut down server
func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info("Shutting down server", map[string]any{
		"time": time.Now().Format(time.RFC3339),
	})

	// Clean up tasks
	return nil
}

// Force an immediate shutdown
func (s *Server) Close() error {
	s.Logger.Info("Forcefully closing server", map[string]any{
		"time": time.Now().Format(time.RFC3339),
	})

	return nil
}
