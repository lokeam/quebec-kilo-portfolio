package analytics

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/testutils"
)

func TestAnalyticsRepository(t *testing.T) {
	// Set up base app context for testing
	testLogger := testutils.NewTestLogger()
	_ = testLogger // Not currently used in repository, but can be used for logging in the future

	// Create mock DB + adapter
	setupMockDB := func() (*repository, sqlmock.Sqlmock, error) {
		// Create mock sqldatabase
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			return nil, nil, err
		}

		// Create a sqlx wrapper around mock data
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		// Create the repository with the mock DB and logger
		repo := &repository{
			db:     sqlxDB,
		}

		return repo, mock, nil
	}

	/*
		GIVEN a request to get general stats for a user
		WHEN all required data exists in the database
		THEN the repository returns the complete stats
	*/
	t.Run("GetGeneralStats - Successfully retrieves complete stats", func(t *testing.T) {
		// Setup
		repo, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer repo.db.Close()

		userID := "test-user-id"

		// Mock physical locations count
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM physical_locations WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		// Mock digital locations count
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM digital_locations WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		// Mock monthly subscription cost
		mock.ExpectQuery("SELECT COALESCE\\(SUM\\(CASE").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"cost"}).AddRow(120.50))

		// Mock total games count
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM user_games WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(42))

		// Execute
		stats, err := repo.GetGeneralStats(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if stats == nil {
			t.Fatalf("Expected stats, got nil")
		}
		if stats.TotalPhysicalLocations != 5 {
			t.Errorf("Expected TotalPhysicalLocations=5, got %d", stats.TotalPhysicalLocations)
		}
		if stats.TotalDigitalLocations != 3 {
			t.Errorf("Expected TotalDigitalLocations=3, got %d", stats.TotalDigitalLocations)
		}
		if stats.MonthlySubscriptionCost != 120.50 {
			t.Errorf("Expected MonthlySubscriptionCost=120.50, got %f", stats.MonthlySubscriptionCost)
		}
		if stats.TotalGames != 42 {
			t.Errorf("Expected TotalGames=42, got %d", stats.TotalGames)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get general stats for a user
		WHEN physical location query fails
		THEN the repository returns an error
	*/
	t.Run("GetGeneralStats - Handles physical locations query error", func(t *testing.T) {
		// Setup
		repo, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer repo.db.Close()

		userID := "test-user-id"
		dbError := errors.New("database error")

		// Mock physical locations count with error
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM physical_locations WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnError(dbError)

		// Execute
		stats, err := repo.GetGeneralStats(context.Background(), userID)

		// Verify
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if stats != nil {
			t.Errorf("Expected nil stats, got %+v", stats)
		}
		if !errors.Is(err, dbError) {
			t.Errorf("Expected error to contain %v, but got %v", dbError, err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get general stats for a user
		WHEN games table doesn't exist yet
		THEN the repository returns stats with games count as 0
	*/
	t.Run("GetGeneralStats - Handles missing games table", func(t *testing.T) {
		// Setup
		repo, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer repo.db.Close()

		userID := "test-user-id"

		// Mock physical locations count
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM physical_locations WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		// Mock digital locations count
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM digital_locations WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		// Mock monthly subscription cost
		mock.ExpectQuery("SELECT COALESCE\\(SUM\\(CASE").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"cost"}).AddRow(120.50))

		// Mock games table doesn't exist
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM user_games WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		// Execute
		stats, err := repo.GetGeneralStats(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if stats == nil {
			t.Fatalf("Expected stats, got nil")
		}
		if stats.TotalGames != 0 {
			t.Errorf("Expected TotalGames=0, got %d", stats.TotalGames)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get financial stats for a user
		WHEN all required data exists in the database
		THEN the repository returns the complete financial stats
	*/
	t.Run("GetFinancialStats - Successfully retrieves complete stats", func(t *testing.T) {
		// Setup
		repo, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer repo.db.Close()

		userID := "test-user-id"
		currentMonth := time.Now().Format("2006-01")

		// Mock annual subscription cost
		mock.ExpectQuery("SELECT COALESCE\\(SUM\\(CASE").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"cost"}).AddRow(1440.0))

		// Mock renewals this month
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM digital_location_subscriptions").
			WithArgs(userID, currentMonth).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// Mock total services count
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM digital_locations").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		// Mock service details
		serviceRows := sqlmock.NewRows([]string{"name", "monthly_fee", "billing_cycle", "next_payment"}).
			AddRow("Netflix", 14.99, "monthly", "2023-05-15").
			AddRow("Spotify", 9.99, "monthly", "2023-05-20").
			AddRow("Xbox Game Pass", 14.99, "quarterly", "2023-06-01")

		mock.ExpectQuery("SELECT l.name, CASE WHEN").
			WithArgs(userID).
			WillReturnRows(serviceRows)

		// Execute
		stats, err := repo.GetFinancialStats(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if stats == nil {
			t.Fatalf("Expected stats, got nil")
		}
		if stats.AnnualSubscriptionCost != 1440.0 {
			t.Errorf("Expected AnnualSubscriptionCost=1440.0, got %f", stats.AnnualSubscriptionCost)
		}
		if stats.RenewalsThisMonth != 2 {
			t.Errorf("Expected RenewalsThisMonth=2, got %d", stats.RenewalsThisMonth)
		}
		if stats.TotalServices != 3 {
			t.Errorf("Expected TotalServices=3, got %d", stats.TotalServices)
		}
		if len(stats.Services) != 3 {
			t.Errorf("Expected 3 services, got %d", len(stats.Services))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get storage stats for a user
		WHEN all required data exists in the database
		THEN the repository returns the complete storage stats
	*/
	t.Run("GetStorageStats - Successfully retrieves complete stats", func(t *testing.T) {
		// Setup
		repo, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer repo.db.Close()

		userID := "test-user-id"

		// Mock location counts
		mock.ExpectQuery("SELECT COUNT\\(CASE WHEN location_type = 'physical' THEN 1 END\\)").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"physical_count", "digital_count"}).AddRow(2, 3))

		// Mock digital locations
		digitalRows := sqlmock.NewRows([]string{"id", "name", "item_count", "location_type", "is_subscription", "monthly_cost"}).
			AddRow("loc1", "Steam", 10, "basic", false, 0.0).
			AddRow("loc2", "Xbox Game Pass", 5, "subscription", true, 14.99)

		mock.ExpectQuery("SELECT l.id, l.name, COUNT\\(dgl.id\\) as item_count").
			WithArgs(userID).
			WillReturnRows(digitalRows)

		// Mock physical locations
		physicalRows := sqlmock.NewRows([]string{"id", "name", "item_count", "location_type"}).
			AddRow("loc3", "Living Room", 7, "house").
			AddRow("loc4", "Game Shelf", 3, "shelf")

		mock.ExpectQuery("SELECT l.id, l.name, COUNT\\(pgl.id\\) as item_count").
			WithArgs(userID).
			WillReturnRows(physicalRows)

		// Execute
		stats, err := repo.GetStorageStats(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if stats == nil {
			t.Fatalf("Expected stats, got nil")
		}
		if stats.TotalPhysicalLocations != 2 {
			t.Errorf("Expected TotalPhysicalLocations=2, got %d", stats.TotalPhysicalLocations)
		}
		if stats.TotalDigitalLocations != 3 {
			t.Errorf("Expected TotalDigitalLocations=3, got %d", stats.TotalDigitalLocations)
		}
		if len(stats.DigitalLocations) != 2 {
			t.Errorf("Expected 2 digital locations, got %d", len(stats.DigitalLocations))
		}
		if len(stats.PhysicalLocations) != 2 {
			t.Errorf("Expected 2 physical locations, got %d", len(stats.PhysicalLocations))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get wishlist stats for a user
		WHEN the wishlist table does not exist
		THEN the repository returns empty wishlist stats
	*/
	t.Run("GetWishlistStats - Handles missing wishlist table", func(t *testing.T) {
		// Setup
		repo, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer repo.db.Close()

		userID := "test-user-id"

		// Mock wishlist table check
		mock.ExpectQuery("SELECT EXISTS").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Execute
		stats, err := repo.GetWishlistStats(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if stats == nil {
			t.Fatalf("Expected empty stats, got nil")
		}
		if stats.TotalWishlistItems != 0 {
			t.Errorf("Expected TotalWishlistItems=0, got %d", stats.TotalWishlistItems)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})

	/*
		GIVEN a request to get inventory stats for a user
		WHEN games tables exist with data
		THEN the repository returns complete inventory stats
	*/
	t.Run("GetInventoryStats - Successfully retrieves complete stats", func(t *testing.T) {
		// Setup
		repo, mock, err := setupMockDB()
		if err != nil {
			t.Fatalf("Failed to setup mock DB: %v", err)
		}
		defer repo.db.Close()

		userID := "test-user-id"
		currentMonth := time.Now().Format("2006-01")

		// Mock total item count
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM user_games WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(50))

		// Mock new items this month
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM user_games").
			WithArgs(userID, currentMonth).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		// Mock platform counts
		platformRows := sqlmock.NewRows([]string{"platform", "item_count"}).
			AddRow("PlayStation", 20).
			AddRow("Nintendo", 15).
			AddRow("Xbox", 10).
			AddRow("PC", 5)

		mock.ExpectQuery("SELECT p.name as platform, COUNT\\(ug.id\\) as item_count").
			WithArgs(userID).
			WillReturnRows(platformRows)

		// Execute
		stats, err := repo.GetInventoryStats(context.Background(), userID)

		// Verify
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if stats == nil {
			t.Fatalf("Expected stats, got nil")
		}
		if stats.TotalItemCount != 50 {
			t.Errorf("Expected TotalItemCount=50, got %d", stats.TotalItemCount)
		}
		if stats.NewItemCount != 5 {
			t.Errorf("Expected NewItemCount=5, got %d", stats.NewItemCount)
		}
		if len(stats.PlatformCounts) != 4 {
			t.Errorf("Expected 4 platforms, got %d", len(stats.PlatformCounts))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unmet expectations: %v", err)
		}
	})
}