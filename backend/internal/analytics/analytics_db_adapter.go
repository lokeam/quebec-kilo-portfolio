package analytics

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
)

type AnalyticsDbAdapter struct {
	db       *sqlx.DB
	logger   interfaces.Logger
}

func NewAnalyticsDbAdapter(appContext *appcontext.AppContext) (*AnalyticsDbAdapter, error) {
	appContext.Logger.Debug("Creating AnalyticsDbAdapter", map[string]any{"appContext": appContext})

	// Create sqlx from px pool
	db := appContext.DB

	return &AnalyticsDbAdapter{
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