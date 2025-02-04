package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/shared/logger"
)

type Server struct {
	config *config.Config
	logger *logger.Logger
	router chi.Router
}

func NewServer(
	config *config.Config,
	logger *logger.Logger,
) *Server {

	srv := &Server{
		config: config,
		logger: logger,
		router: chi.NewRouter(),
	}

	return srv
}

// Add http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Gracefully shut down server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server", map[string]any{
		"time": time.Now().Format(time.RFC3339),
	})

	// Clean up tasks
	return nil
}

// Force an immediate shutdown
func (s *Server) Close() error {
	s.logger.Info("Forcefully closing server", map[string]any{
		"time": time.Now().Format(time.RFC3339),
	})

	return nil
}
