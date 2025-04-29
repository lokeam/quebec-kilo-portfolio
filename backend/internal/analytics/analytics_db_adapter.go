package analytics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/postgres"
)

type AnalyticsDbAdapter struct {
	client   *postgres.PostgresClient
	db       *sqlx.DB
	logger   interfaces.Logger
}

func NewAnalyticsDbAdapter(appContext *appcontext.AppContext) (*AnalyticsDbAdapter, error) {
	appContext.Logger.Debug("Creating AnalyticsDbAdapter", map[string]any{"appContext": appContext})

	// Create a PostgresClient
	client, err := postgres.NewPostgresClient(appContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create Postgres client %w", err)
	}

	// Create sqlx from px pool
	db, err := sqlx.Connect("pgx", appContext.Config.Postgres.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlx connection: %w", err)
	}

	// Register custom types for PostgreSQL arrays so sqlx can handle string array types
	db.MapperFunc(strings.ToLower)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &AnalyticsDbAdapter{
		client: client,
		db:     db,
		logger: appContext.Logger,
	}, nil
}

// Implement the Repository interface by wrapping the existing implementation
func (adapter *AnalyticsDbAdapter) GetGeneralStats(ctx context.Context, userID string) (*GeneralStats, error) {
	repo := &repository{db: adapter.db}
	return repo.GetGeneralStats(ctx, userID)
}

func (adapter *AnalyticsDbAdapter) GetFinancialStats(ctx context.Context, userID string) (*FinancialStats, error) {
	repo := &repository{db: adapter.db}
	return repo.GetFinancialStats(ctx, userID)
}

func (adapter *AnalyticsDbAdapter) GetStorageStats(ctx context.Context, userID string) (*StorageStats, error) {
	repo := &repository{db: adapter.db}
	return repo.GetStorageStats(ctx, userID)
}

func (adapter *AnalyticsDbAdapter) GetInventoryStats(ctx context.Context, userID string) (*InventoryStats, error) {
	repo := &repository{db: adapter.db}
	return repo.GetInventoryStats(ctx, userID)
}

func (adapter *AnalyticsDbAdapter) GetWishlistStats(ctx context.Context, userID string) (*WishlistStats, error) {
	repo := &repository{db: adapter.db}
	return repo.GetWishlistStats(ctx, userID)
}