package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
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

	if appContext.Config.Postgres.ConnectionString == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	// Create a connect pool
	poolConfig, err := pgxpool.ParseConfig(appContext.Config.Postgres.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Configure the connection pool
	poolConfig.MaxConns = 10
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// Create said pool
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, poolConfig.ConnConfig.ConnString())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &PostgresClient{
		pool: pool,
		logger: appContext.Logger,
		appContext: appContext,
	}, nil
}

// Close the db connection pool
func (pc *PostgresClient) Close() {
	if pc.pool != nil {
		pc.pool.Close()
	}
}

func (pc *PostgresClient) GetPool() *pgxpool.Pool {
	return pc.pool
}
