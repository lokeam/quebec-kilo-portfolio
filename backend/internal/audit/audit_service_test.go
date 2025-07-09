package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuditService(t *testing.T) {
	// Skip if integration tests are disabled
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create a proper mock config with Postgres configuration
	mockConfig := &config.Config{
		Postgres: &config.PostgresConfig{
			ConnectionString: "postgres://test:test@localhost:5432/testdb?sslmode=disable",
		},
	}

	appCtx := &appcontext.AppContext{
		Logger: testutils.NewTestLogger(),
		Config: mockConfig,
	}

	// Create mock database adapter instead of real connection
	mockDbAdapter := &mocks.MockAuditDbAdapter{}

	service := &AuditService{
		appCtx:    appCtx,
		dbAdapter: mockDbAdapter,
		logDir:    t.TempDir(),
		retention: 7 * 365 * 24 * time.Hour,
	}

	assert.NotNil(t, service)
	assert.NotNil(t, service.dbAdapter)
	assert.Equal(t, appCtx, service.appCtx)
	assert.Equal(t, 7*365*24*time.Hour, service.retention)

	// Test Close
	err := service.Close()
	assert.NoError(t, err)
}

func TestAuditService_LogAction(t *testing.T) {
	// Create mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Create temporary log directory
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "audit")

	service := &AuditService{
		appCtx: &appcontext.AppContext{
			Logger: testutils.NewTestLogger(),
		},
		dbAdapter: &AuditDbAdapter{
			appCtx: &appcontext.AppContext{
				Logger: testutils.NewTestLogger(),
			},
			db: sqlxDB,
		},
		logDir:    logDir,
		retention: 7 * 365 * 24 * time.Hour,
	}

	// Create log directory
	err = os.MkdirAll(logDir, 0755)
	require.NoError(t, err)

	userID := "test-user-123"
	action := ActionDeletionRequested
	details := map[string]interface{}{
		"reason": "user requested deletion",
		"grace_period_days": 30,
	}

	// Expect database insert
	mock.ExpectExec(`INSERT INTO audit_logs`).
		WithArgs(userID, action, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = service.LogAction(context.Background(), userID, action, details)
	require.NoError(t, err)

	// Verify database expectations
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check that log file was created
	logFiles, err := os.ReadDir(logDir)
	require.NoError(t, err)
	assert.Greater(t, len(logFiles), 0)
}

func TestAuditService_GetUserAuditLog(t *testing.T) {
	// Create mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	service := &AuditService{
		appCtx: &appcontext.AppContext{
			Logger: testutils.NewTestLogger(),
		},
		dbAdapter: &AuditDbAdapter{
			appCtx: &appcontext.AppContext{
				Logger: testutils.NewTestLogger(),
			},
			db: sqlxDB,
		},
		logDir:    t.TempDir(),
		retention: 7 * 365 * 24 * time.Hour,
	}

	userID := "test-user-123"
	expectedLogs := []models.AuditLog{
		{
			ID:        1,
			UserID:    userID,
			Action:    ActionDeletionRequested,
			Timestamp: time.Now(),
			Details:   map[string]interface{}{"reason": "test"},
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    userID,
			Action:    ActionDeletionCancelled,
			Timestamp: time.Now(),
			Details:   map[string]interface{}{"reason": "changed mind"},
			CreatedAt: time.Now(),
		},
	}

	// Mock database query with proper JSON handling
	details1, _ := json.Marshal(expectedLogs[0].Details)
	details2, _ := json.Marshal(expectedLogs[1].Details)

	rows := sqlmock.NewRows([]string{"id", "user_id", "action", "timestamp", "details", "created_at"}).
		AddRow(expectedLogs[0].ID, expectedLogs[0].UserID, expectedLogs[0].Action, expectedLogs[0].Timestamp, details1, expectedLogs[0].CreatedAt).
		AddRow(expectedLogs[1].ID, expectedLogs[1].UserID, expectedLogs[1].Action, expectedLogs[1].Timestamp, details2, expectedLogs[1].CreatedAt)

	mock.ExpectQuery(`SELECT id, user_id, action, timestamp, details, created_at FROM audit_logs WHERE user_id = \$1 ORDER BY timestamp DESC`).
		WithArgs(userID).
		WillReturnRows(rows)

	logs, err := service.GetUserAuditLog(context.Background(), userID)
	require.NoError(t, err)
	assert.Len(t, logs, 2)
	assert.Equal(t, expectedLogs[0].Action, logs[0].Action)
	assert.Equal(t, expectedLogs[1].Action, logs[1].Action)

	// Verify database expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuditService_CleanupOldLogs(t *testing.T) {
	// Create mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Create temporary log directory with old files
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "audit")
	err = os.MkdirAll(logDir, 0755)
	require.NoError(t, err)

	// Create old log file (8 years ago)
	oldDate := time.Now().AddDate(-8, 0, 0)
	oldFilename := filepath.Join(logDir, fmt.Sprintf("audit_%s.log", oldDate.Format("2006-01-02")))
	err = os.WriteFile(oldFilename, []byte("old log content"), 0644)
	require.NoError(t, err)

	// Create recent log file (1 year ago)
	recentDate := time.Now().AddDate(-1, 0, 0)
	recentFilename := filepath.Join(logDir, fmt.Sprintf("audit_%s.log", recentDate.Format("2006-01-02")))
	err = os.WriteFile(recentFilename, []byte("recent log content"), 0644)
	require.NoError(t, err)

	service := &AuditService{
		appCtx: &appcontext.AppContext{
			Logger: testutils.NewTestLogger(),
		},
		dbAdapter: &AuditDbAdapter{
			appCtx: &appcontext.AppContext{
				Logger: testutils.NewTestLogger(),
			},
			db: sqlxDB,
		},
		logDir:    logDir,
		retention: 7 * 365 * 24 * time.Hour,
	}

	// Mock database cleanup
	mock.ExpectExec(`DELETE FROM audit_logs WHERE timestamp < \$1`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 5))

	err = service.CleanupOldLogs(context.Background())
	require.NoError(t, err)

	// Verify database expectations
	assert.NoError(t, mock.ExpectationsWereMet())

	// Check that old file was deleted but recent file remains
	_, err = os.Stat(oldFilename)
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(recentFilename)
	assert.NoError(t, err)
}

func TestAuditService_GetAuditStats(t *testing.T) {
	// Create mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	service := &AuditService{
		appCtx: &appcontext.AppContext{
			Logger: testutils.NewTestLogger(),
		},
		dbAdapter: &AuditDbAdapter{
			appCtx: &appcontext.AppContext{
				Logger: testutils.NewTestLogger(),
			},
			db: sqlxDB,
		},
		logDir:    t.TempDir(),
		retention: 7 * 365 * 24 * time.Hour,
	}

	// Mock statistics queries
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM audit_logs`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(150))

	oldestTime := time.Now().AddDate(-1, 0, 0)
	mock.ExpectQuery(`SELECT MIN\(timestamp\) FROM audit_logs`).
		WillReturnRows(sqlmock.NewRows([]string{"min"}).AddRow(oldestTime))

	stats, err := service.GetAuditStats(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 150, stats["total_logs"])
	assert.Equal(t, oldestTime, stats["oldest_log"])
	assert.Equal(t, 2555, stats["retention_days"])

	// Verify database expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuditLog_JSONMarshaling(t *testing.T) {
	log := models.AuditLog{
		ID:        1,
		UserID:    "test-user",
		Action:    ActionDeletionRequested,
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"reason": "test reason",
			"count":  42,
		},
		CreatedAt: time.Now(),
	}

	// Test JSON marshaling
	data, err := json.Marshal(log)
	require.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaled models.AuditLog
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, log.ID, unmarshaled.ID)
	assert.Equal(t, log.UserID, unmarshaled.UserID)
	assert.Equal(t, log.Action, unmarshaled.Action)
	assert.Equal(t, log.Details["reason"], unmarshaled.Details["reason"])
	// JSON numbers are unmarshaled as float64, so we need to compare the values
	assert.Equal(t, float64(42), unmarshaled.Details["count"])
}