package cleanup

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatabaseCleanupService(t *testing.T) {
	// Skip if integration tests are disabled
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	appCtx := &appcontext.AppContext{
		Logger: testutils.NewTestLogger(),
	}
	service, err := NewDatabaseCleanupService(appCtx)

	require.NoError(t, err)
	assert.NotNil(t, service)
	assert.NotNil(t, service.db)
	assert.Equal(t, appCtx, service.appCtx)

	// Test Close
	err = service.Close()
	assert.NoError(t, err)
}

func TestGetExpiredDemoUsers(t *testing.T) {
	// Create mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	service := &DatabaseCleanupService{
		appCtx: &appcontext.AppContext{
			Logger: testutils.NewTestLogger(),
		},
		db:     sqlxDB,
	}

	// Test case: Found expired users
	expectedUsers := []string{"demo|123_abc", "demo|456_def"}
	rows := sqlmock.NewRows([]string{"user_id"}).
		AddRow("demo|123_abc").
		AddRow("demo|456_def")

	mock.ExpectQuery(`SELECT user_id FROM users WHERE user_id LIKE 'demo\|%' AND created_at < NOW\(\) - INTERVAL '24 hours'`).
		WillReturnRows(rows)

	users, err := service.getExpiredDemoUsers(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedUsers, users)

	// Test case: No expired users
	mock.ExpectQuery(`SELECT user_id FROM users WHERE user_id LIKE 'demo\|%' AND created_at < NOW\(\) - INTERVAL '24 hours'`).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}))

	users, err = service.getExpiredDemoUsers(context.Background())
	require.NoError(t, err)
	assert.Empty(t, users)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCleanupExpiredDemoUsers(t *testing.T) {
	// Create mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	service := &DatabaseCleanupService{
		appCtx: &appcontext.AppContext{
			Logger: testutils.NewTestLogger(),
		},
		db:     sqlxDB,
	}

	// Test case: No expired users
	mock.ExpectQuery(`SELECT user_id FROM users WHERE user_id LIKE 'demo\|%' AND created_at < NOW\(\) - INTERVAL '24 hours'`).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}))

	err = service.cleanupExpiredDemoUsers(context.Background())
	require.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCleanupUserData(t *testing.T) {
	// Create mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	service := &DatabaseCleanupService{
		appCtx: &appcontext.AppContext{
			Logger: testutils.NewTestLogger(),
		},
		db:     sqlxDB,
	}

	userID := "demo|123_abc"

	// Set up transaction expectations
	mock.ExpectBegin()

	// Expect DELETE statements in order
	mock.ExpectExec(`DELETE FROM user_games WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`DELETE FROM physical_locations WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`DELETE FROM digital_locations WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`DELETE FROM spending_data WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`DELETE FROM wishlist_items WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`DELETE FROM users WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err = service.cleanupUserData(context.Background(), userID)
	require.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCleanupExpiredData(t *testing.T) {
	// Create mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	service := &DatabaseCleanupService{
		appCtx: &appcontext.AppContext{
			Logger: testutils.NewTestLogger(),
		},
		db:     sqlxDB,
	}

	// Test case: No expired users
	mock.ExpectQuery(`SELECT user_id FROM users WHERE user_id LIKE 'demo\|%' AND created_at < NOW\(\) - INTERVAL '24 hours'`).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}))

	err = service.CleanupExpiredData(context.Background())
	require.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCleanupExpiredDataWithTimeout(t *testing.T) {
	// Create mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	service := &DatabaseCleanupService{
		appCtx: &appcontext.AppContext{
			Logger: testutils.NewTestLogger(),
		},
		db:     sqlxDB,
	}

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = service.CleanupExpiredData(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")

	// Verify no expectations were set (since context was cancelled)
	assert.NoError(t, mock.ExpectationsWereMet())
}