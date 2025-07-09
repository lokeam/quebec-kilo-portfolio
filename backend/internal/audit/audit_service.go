package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
)

// AuditService handles audit logging for compliance
type AuditService struct {
	appCtx     *appcontext.AppContext
	dbAdapter  interfaces.AuditDbAdapter
	logDir     string
	retention  time.Duration // 7 years
}

// NewAuditService creates a new audit service
func NewAuditService(appCtx *appcontext.AppContext) (*AuditService, error) {
	if appCtx == nil {
		return nil, fmt.Errorf("app context cannot be nil")
	}

	// Create database adapter
	dbAdapter, err := NewAuditDbAdapter(appCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to create audit database adapter: %w", err)
	}

	// Create audit log directory
	logDir := "logs/audit"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create audit log directory: %w", err)
	}

	return &AuditService{
		appCtx:    appCtx,
		dbAdapter: dbAdapter,
		logDir:    logDir,
		retention: 7 * 365 * 24 * time.Hour, // 7 years
	}, nil
}

// Close closes the database connection
func (as *AuditService) Close() error {
	// The interface doesn't have a Close method, so we can't call it
	// In a real implementation, we might want to add Close() to the interface
	// For now, just return nil since the mock doesn't need cleanup
	return nil
}

// LogAction logs an audit action to both database and file
func (as *AuditService) LogAction(
	ctx context.Context,
	userID,
	action string,
	details map[string]any,
) error {
	timestamp := time.Now()

	// Log to database
	if err := as.logToDatabase(
		ctx,
		userID,
		action,
		timestamp,
		details,
	); err != nil {
		as.appCtx.Logger.Error("Failed to log to database", map[string]any{
			"user_id": userID,
			"action":  action,
			"error":   err.Error(),
		})
		// Continue with file logging even if database fails
	}

	// Log to file
	if err := as.logToFile(
		userID,
		action,
		timestamp,
		details,
	); err != nil {
		as.appCtx.Logger.Error("Failed to log to file", map[string]any{
			"user_id": userID,
			"action":  action,
			"error":   err.Error(),
		})
		// Don't return error, as database logging succeeded
	}

	as.appCtx.Logger.Info("Audit action logged", map[string]any{
		"user_id": userID,
		"action":  action,
		"timestamp": timestamp,
	})

	return nil
}

// logToDatabase logs an audit action to the database
func (as *AuditService) logToDatabase(
	ctx context.Context,
	userID,
	action string,
	timestamp time.Time,
	details map[string]any,
) error {
	return as.dbAdapter.CreateAuditLog(
		ctx,
		userID,
		action,
		timestamp,
		details,
	)
}

// logToFile logs an audit action to a file with automatic rotation
func (as *AuditService) logToFile(
	userID, action string,
	timestamp time.Time,
	details map[string]any,
) error {
	// Create filename with date for rotation
	filename := filepath.Join(as.logDir, fmt.Sprintf("audit_%s.log", timestamp.Format("2006-01-02")))

	// Prepare log entry
	logEntry := map[string]interface{}{
		"user_id":   userID,
		"action":    action,
		"timestamp": timestamp.Format(time.RFC3339),
		"details":   details,
	}

	logJSON, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// Append to file
	file, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Write log entry with newline
	if _, err := file.Write(append(logJSON, '\n')); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}

// GetUserAuditLog retrieves audit logs for a specific user
func (as *AuditService) GetUserAuditLog(ctx context.Context, userID string) ([]models.AuditLog, error) {
	return as.dbAdapter.GetUserAuditLog(ctx, userID)
}

// CleanupOldLogs removes audit logs older than the retention period
func (as *AuditService) CleanupOldLogs(ctx context.Context) error {
	cutoffTime := time.Now().Add(-as.retention)

	// Clean up database logs
	rowsAffected, err := as.dbAdapter.DeleteOldAuditLogs(ctx, cutoffTime)
	if err != nil {
		return fmt.Errorf("failed to cleanup database logs: %w", err)
	}

	as.appCtx.Logger.Info("Cleaned up old database audit logs", map[string]any{
		"rows_deleted": rowsAffected,
		"cutoff_time":  cutoffTime,
	})

	// Clean up file logs
	if err := as.cleanupOldFileLogs(cutoffTime); err != nil {
		as.appCtx.Logger.Error("Failed to cleanup old file logs", map[string]any{
			"error": err.Error(),
		})
		// Don't return error, as database cleanup succeeded
	}

	return nil
}

// cleanupOldFileLogs removes audit log files older than the retention period
func (as *AuditService) cleanupOldFileLogs(cutoffTime time.Time) error {
	entries, err := os.ReadDir(as.logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	deletedCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Parse date from filename (audit_2006-01-02.log)
		filename := entry.Name()
		if len(filename) < 15 || filename[:6] != "audit_" || filename[len(filename)-4:] != ".log" {
			continue
		}

		dateStr := filename[6 : len(filename)-4] // Extract date part
		fileDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			as.appCtx.Logger.Warn("Failed to parse date from log filename", map[string]any{
				"filename": filename,
				"error":    err.Error(),
			})
			continue
		}

		// Delete file if older than retention period
		if fileDate.Before(cutoffTime) {
			filePath := filepath.Join(as.logDir, filename)
			if err := os.Remove(filePath); err != nil {
				as.appCtx.Logger.Error("Failed to delete old log file", map[string]any{
					"filename": filename,
					"error":    err.Error(),
				})
				continue
			}
			deletedCount++
		}
	}

	as.appCtx.Logger.Info("Cleaned up old file audit logs", map[string]any{
		"files_deleted": deletedCount,
		"cutoff_time":   cutoffTime,
	})

	return nil
}

// GetAuditStats returns audit statistics
func (as *AuditService) GetAuditStats(ctx context.Context) (map[string]interface{}, error) {
	totalLogs, err := as.dbAdapter.GetTotalAuditLogsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total logs count: %w", err)
	}

	oldestLog, err := as.dbAdapter.GetOldestAuditLog(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get oldest log: %w", err)
	}

	stats := map[string]interface{}{
		"total_logs":    totalLogs,
		"retention_days": int(as.retention.Hours() / 24),
	}

	if oldestLog != nil {
		stats["oldest_log"] = *oldestLog
	} else {
		stats["oldest_log"] = time.Time{}
	}

	return stats, nil
}