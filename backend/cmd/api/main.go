package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/cmd/resourceinitializer"
	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/shared/logger"
	"github.com/lokeam/qko-beta/internal/shared/worker"
	"github.com/lokeam/qko-beta/server"
)

// Fail fast and set up local env variables
func init() {
	initEnv()
}

// Main - Now we actually do the thing
func main() {


	// 1. Initialize root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Sort out the runtime environment for logging configuration
	logEnv := os.Getenv("ENV")
	var loggerEnvironment logger.Environment
	switch logEnv {
	case "prod":
		loggerEnvironment = logger.EnvProd
	case "dev":
		loggerEnvironment = logger.EnvDev
	default:
		loggerEnvironment = logger.EnvDev // Default to "dev"
	}

	// 3. Initialize logging (Zap + slog)
	log, err := logger.NewLogger(
		logger.WithEnv(loggerEnvironment),
		logger.WithAlertLevel(logger.LevelInfo),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize environment: %v\n", err)
		os.Exit(1)
	}
	defer log.Cleanup()

	// 3.5 Initialize Sentry
	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_BACKEND_DSN"),
		Environment: os.Getenv("SENTRY_ENVIRONMENT"),
		// Set traces sample rate to capture transactions for performance monitoring
		TracesSampleRate: 0.2,
		// Use the same environment as your app
		Release: os.Getenv("APP_VERSION"), // NOTE:Consider adding a version env var?
	})
	if err != nil {
		log.Error("Failed to initialize Sentry", map[string]any{
			"error": err.Error(),
		})
		// Don't exit - continue without Sentry rather than crashing the app
	} else {
		log.Info("Sentry initialized successfully", nil)
		// Flush buffered events before the program terminates
		defer sentry.Flush(2 * time.Second)
	}

	// 4. Load app configuration
	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed to load configuration", map[string]any{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	// 5.Initialize all resources
	resources, err := resourceinitializer.NewResourceInitializer(ctx, cfg, log)
	if err != nil {
		log.Error("Failed to initialize resources", map[string]any{"error": err.Error()})
		os.Exit(1)
	}

	// 6. Build global app context to be passed into server
	appCtx := appcontext.NewAppContext(cfg, log, resources.MemCache, resources.RedisClient)

	// Initialize shared DB pool
	db, err := sqlx.Connect("pgx", appCtx.Config.Postgres.ConnectionString)
	if err != nil {
		log.Error("Failed to connect to Postgres", map[string]any{"error": err})
		os.Exit(1)
	}
	db.SetMaxOpenConns(appCtx.Config.Postgres.MaxConnections)
	db.SetMaxIdleConns(appCtx.Config.Postgres.MaxConnections)
	db.SetConnMaxIdleTime(appCtx.Config.Postgres.MaxIdleTime)
	db.SetConnMaxLifetime(appCtx.Config.Postgres.MaxLifetime)
	appCtx.DB = db
	defer db.Close()

	// 7. Create HTTP server
	srv := server.NewServer(cfg, log, appCtx)

	// 8. Start background workers
	worker.StartInitIGDBJob(
		ctx,
		cfg.IGDB.AccessTokenKey,
		&worker.CacheClients{
			RedisClient: resources.RedisClient,
			MemCache:    resources.MemCache,
		},
		cfg.IGDB.ClientID,
		cfg.IGDB.ClientSecret,
		cfg.IGDB.AuthURL,
		log,
	)

	// 9. Configure HTTP server timeouts
	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:     srv,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Info("Server listening", map[string]any{
		"port":  cfg.Server.Port,
		"env":   cfg.Env,
		"time":  time.Now().Format(time.RFC3339),
	})

	// 10. Set up graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// 11. Start server in separate goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Info("Starting server", map[string]any{
			"port": cfg.Server.Port,
			"env": cfg.Env,
			"time": time.Now().Format(time.RFC3339),
		})
		serverErrors <- httpServer.ListenAndServe()
	}()

	// 12. Wait for shutdown signal
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

		// 13. Graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
		// Clean up resources associated with the context
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
