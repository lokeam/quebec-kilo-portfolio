package audit

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
)

// JSONMap is a custom type for handling JSON fields in the database
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into JSONMap", value)
	}

	return json.Unmarshal(bytes, j)
}

// AuditLogWithJSON is a struct for scanning audit logs with proper JSON handling
type AuditLogWithJSON struct {
	ID        int       `db:"id"`
	UserID    string    `db:"user_id"`
	Action    string    `db:"action"`
	Timestamp time.Time `db:"timestamp"`
	Details   JSONMap   `db:"details"`
	CreatedAt time.Time `db:"created_at"`
}

// AuditDbAdapter implements the audit database operations
type AuditDbAdapter struct {
	appCtx *appcontext.AppContext
	db     *sqlx.DB
}

// NewAuditDbAdapter creates a new audit database adapter
func NewAuditDbAdapter(appCtx *appcontext.AppContext) (interfaces.AuditDbAdapter, error) {
	if appCtx == nil {
		return nil, fmt.Errorf("app context cannot be nil")
	}

	// Use shared DB pool
	db := appCtx.DB

	return &AuditDbAdapter{
		appCtx: appCtx,
		db:     db,
	}, nil
}

// Close closes the database connection
func (ada *AuditDbAdapter) Close() error {
	if ada.db != nil {
		return ada.db.Close()
	}
	return nil
}

// CreateAuditLog creates a new audit log entry
func (ada *AuditDbAdapter) CreateAuditLog(
	ctx context.Context,
	userID,
	action string,
	timestamp time.Time,
	details map[string]any) error {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal details: %w", err)
	}

	_, err = ada.db.ExecContext(
		ctx,
		CreateAuditLogQuery,
		userID,
		action,
		timestamp,
		detailsJSON,
		timestamp,
	)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetUserAuditLog retrieves audit logs for a specific user
func (ada *AuditDbAdapter) GetUserAuditLog(ctx context.Context, userID string) ([]models.AuditLog, error) {
	var logsWithJSON []AuditLogWithJSON
	err := ada.db.SelectContext(
		ctx,
		&logsWithJSON,
		GetUserAuditLogQuery,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user audit log: %w", err)
	}

	// Convert to models.AuditLog
	logs := make([]models.AuditLog, len(logsWithJSON))
	for i, logWithJSON := range logsWithJSON {
		logs[i] = models.AuditLog{
			ID:        logWithJSON.ID,
			UserID:    logWithJSON.UserID,
			Action:    logWithJSON.Action,
			Timestamp: logWithJSON.Timestamp,
			Details:   map[string]interface{}(logWithJSON.Details),
			CreatedAt: logWithJSON.CreatedAt,
		}
	}

	return logs, nil
}

// GetUserAuditLogByDateRange retrieves audit logs for a user within a date range
func (ada *AuditDbAdapter) GetUserAuditLogByDateRange(
	ctx context.Context,
	userID string,
	startDate,
	endDate time.Time,
	limit int,
) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := ada.db.SelectContext(
		ctx,
		&logs,
		GetUserAuditLogByDateRangeQuery,
		userID,
		startDate,
		endDate,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user audit log by date range: %w", err)
	}

	// Parse details JSON for each log entry
	for i := range logs {
		if logs[i].Details == nil {
			logs[i].Details = make(map[string]any)
		}
	}

	return logs, nil
}

// GetAllAuditLogsByDateRange retrieves all audit logs within a date range
func (ada *AuditDbAdapter) GetAllAuditLogsByDateRange(
	ctx context.Context,
	startDate,
	endDate time.Time,
	limit int,
) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := ada.db.SelectContext(
		ctx,
		&logs,
		GetAllAuditLogsByDateRangeQuery,
		startDate,
		endDate,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get all audit logs by date range: %w", err)
	}

	// Parse details JSON for each log entry
	for i := range logs {
		if logs[i].Details == nil {
			logs[i].Details = make(map[string]any)
		}
	}

	return logs, nil
}

// GetAuditLogsByAction retrieves audit logs by action type
func (ada *AuditDbAdapter) GetAuditLogsByAction(
	ctx context.Context,
	action string,
	limit int,
) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := ada.db.SelectContext(
		ctx,
		&logs,
		GetAuditLogsByActionQuery,
		action,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by action: %w", err)
	}

	// Parse details JSON for each log entry
	for i := range logs {
		if logs[i].Details == nil {
			logs[i].Details = make(map[string]any)
		}
	}

	return logs, nil
}

// GetUserAuditLogsByAction retrieves audit logs for a user by action type
func (ada *AuditDbAdapter) GetUserAuditLogsByAction(
	ctx context.Context,
	userID,
	action string,
	limit int,
) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := ada.db.SelectContext(
		ctx,
		&logs,
		GetUserAuditLogsByActionQuery,
		userID,
		action,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user audit logs by action: %w", err)
	}

	// Parse details JSON for each log entry
	for i := range logs {
		if logs[i].Details == nil {
			logs[i].Details = make(map[string]interface{})
		}
	}

	return logs, nil
}

// GetAuditLogByID retrieves a specific audit log by ID
func (ada *AuditDbAdapter) GetAuditLogByID(ctx context.Context, logID int) (*models.AuditLog, error) {
	var log models.AuditLog
	err := ada.db.GetContext(ctx, &log, GetAuditLogByIDQuery, logID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log by ID: %w", err)
	}

	if log.Details == nil {
		log.Details = make(map[string]interface{})
	}

	return &log, nil
}

// GetTotalAuditLogsCount returns the total number of audit logs
func (ada *AuditDbAdapter) GetTotalAuditLogsCount(ctx context.Context) (int, error) {
	var count int
	err := ada.db.GetContext(
		ctx,
		&count,
		GetTotalAuditLogsCountQuery,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get total audit logs count: %w", err)
	}
	return count, nil
}

// GetOldestAuditLog returns the timestamp of the oldest audit log
func (ada *AuditDbAdapter) GetOldestAuditLog(ctx context.Context) (*time.Time, error) {
	var oldestTime time.Time
	err := ada.db.GetContext(
		ctx,
		&oldestTime,
		GetOldestAuditLogQuery,
	)
	if err != nil {
		// No logs found
		return nil, nil
	}
	return &oldestTime, nil
}

// GetAuditLogsCountByUser returns the count of audit logs for a specific user
func (ada *AuditDbAdapter) GetAuditLogsCountByUser(ctx context.Context, userID string) (int, error) {
	var count int
	err := ada.db.GetContext(
		ctx,
		&count,
		GetAuditLogsCountByUserQuery,
		userID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get audit logs count by user: %w", err)
	}
	return count, nil
}

// GetAuditLogsCountByAction returns counts of audit logs grouped by action
func (ada *AuditDbAdapter) GetAuditLogsCountByAction(
	ctx context.Context,
	since time.Time,
) (map[string]int, error) {
	rows, err := ada.db.QueryContext(
		ctx,
		GetAuditLogsCountByActionQuery,
		since,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs count by action: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var action string
		var count int
		if err := rows.Scan(&action, &count); err != nil {
			return nil, fmt.Errorf("failed to scan audit logs count by action: %w", err)
		}
		result[action] = count
	}

	return result, nil
}

// GetAuditLogsCountByDateRange returns daily counts of audit logs within a date range
func (ada *AuditDbAdapter) GetAuditLogsCountByDateRange(
	ctx context.Context,
	startDate,
	endDate time.Time,
) (map[string]int, error) {
	rows, err := ada.db.QueryContext(
		ctx,
		GetAuditLogsCountByDateRangeQuery,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs count by date range: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var logDate string
		var count int
		if err := rows.Scan(&logDate, &count); err != nil {
			return nil, fmt.Errorf("failed to scan audit logs count by date range: %w", err)
		}
		result[logDate] = count
	}

	return result, nil
}

// DeleteOldAuditLogs removes audit logs older than the specified cutoff time
func (ada *AuditDbAdapter) DeleteOldAuditLogs(
	ctx context.Context,
	cutoffTime time.Time,
) (int64, error) {
	result, err := ada.db.ExecContext(
		ctx,
		DeleteOldAuditLogsQuery,
		cutoffTime,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old audit logs: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// DeleteAuditLogsByUser removes all audit logs for a specific user
func (ada *AuditDbAdapter) DeleteAuditLogsByUser(ctx context.Context, userID string) (int64, error) {
	result, err := ada.db.ExecContext(ctx, DeleteAuditLogsByUserQuery, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete audit logs by user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// DeleteAuditLogsByAction removes audit logs by action type older than cutoff time
func (ada *AuditDbAdapter) DeleteAuditLogsByAction(
	ctx context.Context,
	action string,
	cutoffTime time.Time,
) (int64, error) {
	result, err := ada.db.ExecContext(
		ctx,
		DeleteAuditLogsByActionQuery,
		action,
		cutoffTime,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to delete audit logs by action: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetUserAuditLogForExport retrieves audit logs for export (chronological order)
func (ada *AuditDbAdapter) GetUserAuditLogForExport(ctx context.Context, userID string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := ada.db.SelectContext(
		ctx,
		&logs,
		GetUserAuditLogForExportQuery,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user audit log for export: %w", err)
	}

	// Parse details JSON for each log entry
	for i := range logs {
		if logs[i].Details == nil {
			logs[i].Details = make(map[string]interface{})
		}
	}

	return logs, nil
}

// GetAllAuditLogsForExport retrieves all audit logs for export within a date range
func (ada *AuditDbAdapter) GetAllAuditLogsForExport(
	ctx context.Context,
	startDate,
	endDate time.Time,
) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := ada.db.SelectContext(
		ctx,
		&logs,
		GetAllAuditLogsForExportQuery,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get all audit logs for export: %w", err)
	}

	// Parse details JSON for each log entry
	for i := range logs {
		if logs[i].Details == nil {
			logs[i].Details = make(map[string]any)
		}
	}

	return logs, nil
}

// GetComplianceAuditLogs retrieves compliance-related audit logs
func (ada *AuditDbAdapter) GetComplianceAuditLogs(ctx context.Context, since time.Time) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := ada.db.SelectContext(
		ctx,
		&logs,
		GetComplianceAuditLogsQuery,
		since,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get compliance audit logs: %w", err)
	}

	// Parse details JSON for each log entry
	for i := range logs {
		if logs[i].Details == nil {
			logs[i].Details = make(map[string]any)
		}
	}

	return logs, nil
}

// GetUserComplianceAuditLogs retrieves compliance-related audit logs for a specific user
func (ada *AuditDbAdapter) GetUserComplianceAuditLogs(ctx context.Context, userID string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := ada.db.SelectContext(
		ctx,
		&logs,
		GetUserComplianceAuditLogsQuery,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user compliance audit logs: %w", err)
	}

	// Parse details JSON for each log entry
	for i := range logs {
		if logs[i].Details == nil {
			logs[i].Details = make(map[string]any)
		}
	}

	return logs, nil
}

// CheckIfAuditLogExists checks if an audit log with the given ID exists
func (ada *AuditDbAdapter) CheckIfAuditLogExists(ctx context.Context, logID int) (bool, error) {
	var exists bool
	err := ada.db.GetContext(
		ctx,
		&exists,
		CheckIfAuditLogExistsQuery,
		logID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to check if audit log exists: %w", err)
	}
	return exists, nil
}

// CheckIfUserHasAuditLogs checks if a user has any audit logs
func (ada *AuditDbAdapter) CheckIfUserHasAuditLogs(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := ada.db.GetContext(
		ctx,
		&exists,
		CheckIfUserHasAuditLogsQuery,
		userID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to check if user has audit logs: %w", err)
	}
	return exists, nil
}