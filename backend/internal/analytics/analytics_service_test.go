package analytics

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/stretchr/testify/assert"
)

// MockAnalyticsRepository implements the Repository interface for testing
type MockAnalyticsRepository struct {
	GetGeneralStatsFunc    func(ctx context.Context, userID string) (*GeneralStats, error)
	GetFinancialStatsFunc  func(ctx context.Context, userID string) (*FinancialStats, error)
	GetStorageStatsFunc    func(ctx context.Context, userID string) (*StorageStats, error)
	GetInventoryStatsFunc  func(ctx context.Context, userID string) (*InventoryStats, error)
	GetWishlistStatsFunc   func(ctx context.Context, userID string) (*WishlistStats, error)
}

func (m *MockAnalyticsRepository) GetGeneralStats(ctx context.Context, userID string) (*GeneralStats, error) {
	return m.GetGeneralStatsFunc(ctx, userID)
}

func (m *MockAnalyticsRepository) GetFinancialStats(ctx context.Context, userID string) (*FinancialStats, error) {
	return m.GetFinancialStatsFunc(ctx, userID)
}

func (m *MockAnalyticsRepository) GetStorageStats(ctx context.Context, userID string) (*StorageStats, error) {
	return m.GetStorageStatsFunc(ctx, userID)
}

func (m *MockAnalyticsRepository) GetInventoryStats(ctx context.Context, userID string) (*InventoryStats, error) {
	return m.GetInventoryStatsFunc(ctx, userID)
}

func (m *MockAnalyticsRepository) GetWishlistStats(ctx context.Context, userID string) (*WishlistStats, error) {
	return m.GetWishlistStatsFunc(ctx, userID)
}

// MockCacheWrapper implements CacheWrapper interface for testing
type MockCacheWrapper struct {
	GetCachedResultsFunc func(ctx context.Context, key string, result any) (bool, error)
	SetCachedResultsFunc func(ctx context.Context, key string, data any) error
	DeleteCacheKeyFunc   func(ctx context.Context, key string) error
}

func (m *MockCacheWrapper) GetCachedResults(ctx context.Context, key string, result any) (bool, error) {
	return m.GetCachedResultsFunc(ctx, key, result)
}

func (m *MockCacheWrapper) SetCachedResults(ctx context.Context, key string, data any) error {
	return m.SetCachedResultsFunc(ctx, key, data)
}

func (m *MockCacheWrapper) DeleteCacheKey(ctx context.Context, key string) error {
	return m.DeleteCacheKeyFunc(ctx, key)
}

func newMockAnalyticsServiceWithDefaults(logger *testutils.TestLogger) Service {
	mockRepo := &MockAnalyticsRepository{
		GetGeneralStatsFunc: func(ctx context.Context, userID string) (*GeneralStats, error) {
			return &GeneralStats{
				TotalPhysicalLocations: 5,
				TotalDigitalLocations:  3,
				MonthlySubscriptionCost: 120.50,
				TotalGames:              42,
			}, nil
		},
		GetFinancialStatsFunc: func(ctx context.Context, userID string) (*FinancialStats, error) {
			return &FinancialStats{
				AnnualSubscriptionCost: 1440.00,
				RenewalsThisMonth:      2,
				TotalServices:          3,
				Services: []ServiceDetails{
					{Name: "Netflix", MonthlyFee: 14.99, BillingCycle: "monthly", NextPayment: "2023-05-15"},
					{Name: "Spotify", MonthlyFee: 9.99, BillingCycle: "monthly", NextPayment: "2023-05-20"},
				},
			}, nil
		},
		GetStorageStatsFunc: func(ctx context.Context, userID string) (*StorageStats, error) {
			return &StorageStats{
				TotalPhysicalLocations: 2,
				TotalDigitalLocations:  3,
				PhysicalLocations: []LocationSummary{
					{ID: "loc1", Name: "Living Room", ItemCount: 7, LocationType: "house"},
					{ID: "loc2", Name: "Game Shelf", ItemCount: 3, LocationType: "shelf"},
				},
				DigitalLocations: []LocationSummary{
					{ID: "loc3", Name: "Steam", ItemCount: 10, LocationType: "basic", IsSubscription: false, MonthlyCost: 0.0},
					{ID: "loc4", Name: "Xbox Game Pass", ItemCount: 5, LocationType: "subscription", IsSubscription: true, MonthlyCost: 14.99},
				},
			}, nil
		},
		GetInventoryStatsFunc: func(ctx context.Context, userID string) (*InventoryStats, error) {
			return &InventoryStats{
				TotalItemCount: 50,
				NewItemCount:   5,
				PlatformCounts: []PlatformItemCount{
					{Platform: "PlayStation", ItemCount: 20},
					{Platform: "Nintendo", ItemCount: 15},
					{Platform: "Xbox", ItemCount: 10},
					{Platform: "PC", ItemCount: 5},
				},
			}, nil
		},
		GetWishlistStatsFunc: func(ctx context.Context, userID string) (*WishlistStats, error) {
			return &WishlistStats{
				TotalWishlistItems: 8,
				ItemsOnSale: 2,
				StarredItem: "Game 1",
				StarredItemPrice: 59.99,
				CheapestSaleDiscount: 25.0,
			}, nil
		},
	}

	mockCache := &MockCacheWrapper{
		GetCachedResultsFunc: func(ctx context.Context, key string, result any) (bool, error) {
			return false, errors.New("cache miss")
		},
		SetCachedResultsFunc: func(ctx context.Context, key string, data any) error {
			return nil
		},
		DeleteCacheKeyFunc: func(ctx context.Context, key string) error {
			return nil
		},
	}

	// Create ttl for different domains
	ttl := map[string]time.Duration{
		DomainGeneral:   time.Duration(TTLGeneralMinutes) * time.Minute,
		DomainFinancial: time.Duration(TTLFinancialMinutes) * time.Minute,
		DomainStorage:   time.Duration(TTLStorageMinutes) * time.Minute,
		DomainInventory: time.Duration(TTLInventoryMinutes) * time.Minute,
		DomainWishlist:  time.Duration(TTLWishlistMinutes) * time.Minute,
	}

	return &service{
		repo:   mockRepo,
		cache:  mockCache,
		logger: logger,
		ttl:    ttl,
	}
}

func TestAnalyticsService(t *testing.T) {
	// Set up test logger
	logger := testutils.NewTestLogger()

	// Test GetAnalytics with specific domains
	t.Run("GetAnalytics - Success with specific domains", func(t *testing.T) {
		// Setup
		service := newMockAnalyticsServiceWithDefaults(logger)

		// Execute
		domains := []string{DomainGeneral, DomainFinancial}
		result, err := service.GetAnalytics(context.Background(), "test-user", domains)

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.NotNil(t, result[DomainGeneral])
		assert.NotNil(t, result[DomainFinancial])
	})

	// Test GetAnalytics with no domains (should default to general)
	t.Run("GetAnalytics - Success with no domains", func(t *testing.T) {
		// Setup
		service := newMockAnalyticsServiceWithDefaults(logger)

		// Execute
		result, err := service.GetAnalytics(context.Background(), "test-user", nil)

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.NotNil(t, result[DomainGeneral])
	})

	// Test GetAnalytics with repository error
	t.Run("GetAnalytics - Repository error", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)
		expectedErr := errors.New("repository error")

		// Override the default mock
		mockRepo := &MockAnalyticsRepository{}
		serviceImpl := svc.(*service)
		*mockRepo = *serviceImpl.repo.(*MockAnalyticsRepository)
		mockRepo.GetGeneralStatsFunc = func(ctx context.Context, userID string) (*GeneralStats, error) {
			return nil, expectedErr
		}
		serviceImpl.repo = mockRepo

		// Execute
		result, err := svc.GetAnalytics(context.Background(), "test-user", []string{DomainGeneral})

		// Verify
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), expectedErr.Error())
	})

	// Test GetAnalytics with cache hit
	t.Run("GetAnalytics - Cache hit", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)
		cachedStats := &GeneralStats{
			TotalPhysicalLocations: 10,
			TotalDigitalLocations:  5,
			MonthlySubscriptionCost: 150.75,
			TotalGames:              60,
		}

		// Override the default mock cache
		mockCache := &MockCacheWrapper{}
		serviceImpl := svc.(*service)
		*mockCache = *serviceImpl.cache.(*MockCacheWrapper)
		mockCache.GetCachedResultsFunc = func(ctx context.Context, key string, result any) (bool, error) {
			// Type assert to get the pointer to the result
			if ptr, ok := result.(*any); ok {
				// Assign the cached value
				*ptr = cachedStats
				return true, nil
			}
			return false, errors.New("cache miss")
		}
		serviceImpl.cache = mockCache

		// Execute
		result, err := svc.GetAnalytics(context.Background(), "test-user", []string{DomainGeneral})

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, cachedStats, result[DomainGeneral])
	})

	// Test individual domain getters
	t.Run("GetGeneralStats - Success", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)

		// Execute
		stats, err := svc.GetGeneralStats(context.Background(), "test-user")

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 5, stats.TotalPhysicalLocations)
		assert.Equal(t, 3, stats.TotalDigitalLocations)
		assert.Equal(t, 120.50, stats.MonthlySubscriptionCost)
		assert.Equal(t, 42, stats.TotalGames)
	})

	t.Run("GetFinancialStats - Success", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)

		// Execute
		stats, err := svc.GetFinancialStats(context.Background(), "test-user")

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 1440.00, stats.AnnualSubscriptionCost)
		assert.Equal(t, 2, stats.RenewalsThisMonth)
		assert.Equal(t, 3, stats.TotalServices)
		assert.Len(t, stats.Services, 2)
	})

	t.Run("GetStorageStats - Success", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)

		// Execute
		stats, err := svc.GetStorageStats(context.Background(), "test-user")

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 2, stats.TotalPhysicalLocations)
		assert.Equal(t, 3, stats.TotalDigitalLocations)
		assert.Len(t, stats.PhysicalLocations, 2)
		assert.Len(t, stats.DigitalLocations, 2)
	})

	t.Run("GetInventoryStats - Success", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)

		// Execute
		stats, err := svc.GetInventoryStats(context.Background(), "test-user")

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 50, stats.TotalItemCount)
		assert.Equal(t, 5, stats.NewItemCount)
		assert.Len(t, stats.PlatformCounts, 4)
	})

	t.Run("GetWishlistStats - Success", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)

		// Execute
		stats, err := svc.GetWishlistStats(context.Background(), "test-user")

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 8, stats.TotalWishlistItems)
		assert.Equal(t, 2, stats.ItemsOnSale)
		assert.Equal(t, "Game 1", stats.StarredItem)
	})

	// Test cache invalidation
	t.Run("InvalidateDomain - Success", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)

		// Execute
		err := svc.InvalidateDomain(context.Background(), "test-user", DomainGeneral)

		// Verify
		assert.NoError(t, err)
	})

	t.Run("InvalidateDomain - Error", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)
		expectedErr := errors.New("cache error")

		// Override the default mock cache
		mockCache := &MockCacheWrapper{}
		serviceImpl := svc.(*service)
		*mockCache = *serviceImpl.cache.(*MockCacheWrapper)
		mockCache.DeleteCacheKeyFunc = func(ctx context.Context, key string) error {
			return expectedErr
		}
		serviceImpl.cache = mockCache

		// Execute
		err := svc.InvalidateDomain(context.Background(), "test-user", DomainGeneral)

		// Verify
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
	})

	t.Run("InvalidateDomains - Success", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)

		// Execute
		err := svc.InvalidateDomains(context.Background(), "test-user", []string{DomainGeneral, DomainFinancial})

		// Verify
		assert.NoError(t, err)
	})

	t.Run("InvalidateDomains - Partial Error", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)
		expectedErr := errors.New("cache error")

		// Override the default mock cache
		mockCache := &MockCacheWrapper{}
		serviceImpl := svc.(*service)
		*mockCache = *serviceImpl.cache.(*MockCacheWrapper)
		mockCache.DeleteCacheKeyFunc = func(ctx context.Context, key string) error {
			// Use the correct key format to match what the service uses
			// CacheKeyFormat is "user:%s:analytics:%s" where first %s is userID and second is domain
			financialKey := fmt.Sprintf("user:%s:analytics:%s", "test-user", DomainFinancial)
			if key == financialKey {
				return expectedErr
			}
			return nil
		}
		serviceImpl.cache = mockCache

		// Execute
		err := svc.InvalidateDomains(context.Background(), "test-user", []string{DomainGeneral, DomainFinancial})

		// Verify
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
	})

	t.Run("setInCache - Error handling", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)
		expectedErr := errors.New("cache error")

		// Override the default mock cache
		mockCache := &MockCacheWrapper{}
		serviceImpl := svc.(*service)
		*mockCache = *serviceImpl.cache.(*MockCacheWrapper)
		mockCache.SetCachedResultsFunc = func(ctx context.Context, key string, data any) error {
			return expectedErr
		}
		serviceImpl.cache = mockCache

		// This should not fail the overall operation
		stats, err := svc.GetGeneralStats(context.Background(), "test-user")

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, stats)
	})

	t.Run("GetAnalytics - Empty user ID", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)
		mockRepo := &MockAnalyticsRepository{}
		serviceImpl := svc.(*service)
		*mockRepo = *serviceImpl.repo.(*MockAnalyticsRepository)
		mockRepo.GetGeneralStatsFunc = func(ctx context.Context, userID string) (*GeneralStats, error) {
			if userID == "" {
				return nil, errors.New("user ID cannot be empty")
			}
			return &GeneralStats{}, nil
		}
		serviceImpl.repo = mockRepo

		// Execute
		result, err := svc.GetAnalytics(context.Background(), "", []string{DomainGeneral})

		// Verify
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user ID cannot be empty")
	})
}

func TestAnalyticsServiceEdgeCases(t *testing.T) {
	// Set up test logger
	logger := testutils.NewTestLogger()

	// Test handling of unknown domain
	t.Run("GetAnalytics - Unknown domain", func(t *testing.T) {
		// Setup
		svc := newMockAnalyticsServiceWithDefaults(logger)

		// Execute
		domains := []string{"unknown_domain", DomainGeneral}
		result, err := svc.GetAnalytics(context.Background(), "test-user", domains)

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.NotNil(t, result[DomainGeneral])
		assert.Nil(t, result["unknown_domain"])
	})
}