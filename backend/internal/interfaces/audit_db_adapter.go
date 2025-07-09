package interfaces

import (
	"context"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

type AuditDbAdapter interface {
	// Audit Log Operations
	CreateAuditLog(ctx context.Context, userID, action string, timestamp time.Time, details map[string]interface{}) error
	GetUserAuditLog(ctx context.Context, userID string) ([]models.AuditLog, error)
	GetUserAuditLogByDateRange(ctx context.Context, userID string, startDate, endDate time.Time, limit int) ([]models.AuditLog, error)
	GetAllAuditLogsByDateRange(ctx context.Context, startDate, endDate time.Time, limit int) ([]models.AuditLog, error)
	GetAuditLogsByAction(ctx context.Context, action string, limit int) ([]models.AuditLog, error)
	GetUserAuditLogsByAction(ctx context.Context, userID, action string, limit int) ([]models.AuditLog, error)
	GetAuditLogByID(ctx context.Context, logID int) (*models.AuditLog, error)

	// Statistics Operations
	GetTotalAuditLogsCount(ctx context.Context) (int, error)
	GetOldestAuditLog(ctx context.Context) (*time.Time, error)
	GetAuditLogsCountByUser(ctx context.Context, userID string) (int, error)
	GetAuditLogsCountByAction(ctx context.Context, since time.Time) (map[string]int, error)
	GetAuditLogsCountByDateRange(ctx context.Context, startDate, endDate time.Time) (map[string]int, error)

	// Cleanup Operations
	DeleteOldAuditLogs(ctx context.Context, cutoffTime time.Time) (int64, error)
	DeleteAuditLogsByUser(ctx context.Context, userID string) (int64, error)
	DeleteAuditLogsByAction(ctx context.Context, action string, cutoffTime time.Time) (int64, error)

	// Export Operations
	GetUserAuditLogForExport(ctx context.Context, userID string) ([]models.AuditLog, error)
	GetAllAuditLogsForExport(ctx context.Context, startDate, endDate time.Time) ([]models.AuditLog, error)

	// Compliance Operations
	GetComplianceAuditLogs(ctx context.Context, since time.Time) ([]models.AuditLog, error)
	GetUserComplianceAuditLogs(ctx context.Context, userID string) ([]models.AuditLog, error)

	// Validation Operations
	CheckIfAuditLogExists(ctx context.Context, logID int) (bool, error)
	CheckIfUserHasAuditLogs(ctx context.Context, userID string) (bool, error)
}