package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/monitoring"
)

type PostgresClient struct {
	pool *pgxpool.Pool
	logger interfaces.Logger
	appContext *appcontext.AppContext
}

func NewPostgresClient(appContext *appcontext.AppContext) (*PostgresClient, error) {
	// TODO: Clean this up
	if appContext == nil {
		panic("appContext is nil")
	}

	// Validate configuration
	if appContext.Config.Postgres.ConnectionString == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	logger := appContext.Logger
	logger.Debug("PostgreSQL configuration", map[string]any{
    "connectionString": hideConnectionStringPw(appContext.Config.Postgres.ConnectionString),
    "maxConnections": appContext.Config.Postgres.MaxConnections,
    "maxLifetime": appContext.Config.Postgres.MaxLifetime,
    "maxIdleTime": appContext.Config.Postgres.MaxIdleTime,
})

	// Initialize client
	client := &PostgresClient{
		logger: logger,
		appContext: appContext,
	}

	// Use expontential backoff for connection retries
	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.MaxElapsedTime = 30 * time.Second // Max time to try connecting

	var pool *pgxpool.Pool
	operation := func() error {
		// Create a connection pool
		poolConfig, err := pgxpool.ParseConfig(appContext.Config.Postgres.ConnectionString)
		if err != nil {
			return fmt.Errorf("failed to parse connection string: %w", err)
		}

		// Configure the connection pool using config values or defaults
		poolConfig.MaxConns = int32(appContext.Config.Postgres.MaxConnections)
		poolConfig.MaxConnLifetime = appContext.Config.Postgres.MaxLifetime
		poolConfig.MaxConnIdleTime = appContext.Config.Postgres.MaxIdleTime

		// Create context with timeout for connection
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		logger.Debug("Attempting to connect to PostgreSQL", map[string]any{
			"attempt": "new",
			"maxConns": poolConfig.MaxConns,
			"maxLifetime": poolConfig.MaxConnLifetime,
			"maxIdleTime": poolConfig.MaxConnIdleTime,
		})

		// Create the pool
		newPool, err := pgxpool.New(ctx, poolConfig.ConnConfig.ConnString())
		if err != nil {
			logger.Error("Failed to create connection pool", map[string]any{
				"error": err.Error(),
			})
			// Close pool we just created to prevent memory leaks
			newPool.Close()
			return fmt.Errorf("unable to create connection pool: %w", err)
		}
		logger.Debug("Connection pool created, testing with ping", nil)

		// Test the connection - IMPORTANT: Do not close newPool on error
		if err := newPool.Ping(ctx); err != nil {
			logger.Error("Ping failed", map[string]any{
				"error": err.Error(),
			})
			// Close pool we just created to prevent memory leaks
			newPool.Close()
			return fmt.Errorf("unable to ping database: %w", err)
		}
		logger.Debug("Ping successful", nil)

		pool = newPool
		return nil
	}

	// Execute the opeartion with exponential backoff
	err := backoff.Retry(operation, exponentialBackoff)
	if err != nil {
		return nil, fmt.Errorf("failed to create Postgres client: %w", err)
	}

	client.pool = pool
	logger.Info("Successfully connected to PostgreSQL", nil)

	// Start collecting metrics from NewPostgresClient after successfully connecting
	client.StartMetricsCollection()

	return client, nil
}

func (pc *PostgresClient) UpdateMetrics() {
	if pc.pool != nil {
		stats := pc.pool.Stat()
		monitoring.DBConnectionsOpen.Set(float64(stats.AcquiredConns()))
		monitoring.DBConnectionsMax.Set(float64(stats.MaxConns()))
	}
}

// Go routine to periodically update metrics
func (pc *PostgresClient) StartMetricsCollection() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		// NOTE: use range when selecting on a single channel
		for range ticker.C {
			pc.UpdateMetrics()
		}
	}()
}

// Close the db connection pool
func (pc *PostgresClient) Close() {
	if pc.pool != nil {
		pc.pool.Close()
	}
}

// Returns the underlying connection pool
func (pc *PostgresClient) GetPool() *pgxpool.Pool {
	return pc.pool
}

// Helper fn to mask pw in connection string for logging
func hideConnectionStringPw(connectionString string) string {
	// NOTE - use a more sophisticated method in the future
	if connectionString == "" {
		return ""
	}

	return "postgres://[hidden]@[host]:[port]/[dbName]"
}
