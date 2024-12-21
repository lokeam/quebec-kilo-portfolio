package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/shared/logger"
	"github.com/lokeam/qko-beta/server"
)

func main() {
	// 1. Initialize root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	// Convert string env to logger.Environment type:
	    // Convert string env to logger.Environment type
			env := os.Getenv("ENV")
			var logEnv logger.Environment
			switch env {
			case "production":
					logEnv = logger.EnvProd
			case "development":
					logEnv = logger.EnvDev
			default:
					logEnv = logger.EnvDev // Default to development
			}


	// 3. Initialize logging (Zap + slog)
	log, err := logger.NewLogger(
		logger.WithEnv(logEnv),
		logger.WithAlertLevel(logger.LevelInfo),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize environment: %v\n", err)
		os.Exit(1)
	}
	defer log.Cleanup()

	// 4. Initialize app configuration
	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed to load configuration", map[string]any{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	// 5. Create server
	srv := server.NewServer(cfg, log)

	// 6. Configure HTTP server
	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:     srv,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 7. Set up graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// 8. Start server
	serverErrors := make(chan error, 1)
	go func() {
		log.Info("Starting server", map[string]any{
			"port": cfg.Server.Port,
			"env": cfg.Env,
			"time": time.Now().Format(time.RFC3339),
		})
		serverErrors <- httpServer.ListenAndServe()
	}()

	// 9. Wait for shutdown signal
	select {
	case err := <- serverErrors:
		log.Error("Server error", map[string]any{
			"error": err.Error(),
		})
	case sig := <- shutdown:
		log.Info("Shutdown signal received", map[string]any{
			"signal": sig,
			"time": time.Now().Format(time.RFC3339),
		})

		// 10. Graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("Server shutdown error", map[string]any{
				"error": err.Error(),
			})

			// Force shutdown if graceful shutdown fails
			srv.Close()
		}
	}
}
