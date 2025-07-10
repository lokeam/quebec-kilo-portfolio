package mocks

import (
	"context"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

type MockAuditDbAdapter struct {
	// Audit Log Operations
	CreateAuditLogFunc                    func(ctx context.Context, userID, action string, timestamp time.Time, details map[string]interface{}) error
	GetUserAuditLogFunc                   func(ctx context.Context, userID string) ([]models.AuditLog, error)
	GetUserAuditLogByDateRangeFunc        func(ctx context.Context, userID string, startDate, endDate time.Time, limit int) ([]models.AuditLog, error)
	GetAllAuditLogsByDateRangeFunc        func(ctx context.Context, startDate, endDate time.Time, limit int) ([]models.AuditLog, error)
	GetAuditLogsByActionFunc              func(ctx context.Context, action string, limit int) ([]models.AuditLog, error)
	GetUserAuditLogsByActionFunc          func(ctx context.Context, userID, action string, limit int) ([]models.AuditLog, error)
	GetAuditLogByIDFunc                   func(ctx context.Context, logID int) (*models.AuditLog, error)

	// Statistics Operations
	GetTotalAuditLogsCountFunc            func(ctx context.Context) (int, error)
	GetOldestAuditLogFunc                 func(ctx context.Context) (*time.Time, error)
	GetAuditLogsCountByUserFunc           func(ctx context.Context, userID string) (int, error)
	GetAuditLogsCountByActionFunc         func(ctx context.Context, since time.Time) (map[string]int, error)
	GetAuditLogsCountByDateRangeFunc      func(ctx context.Context, startDate, endDate time.Time) (map[string]int, error)

	// Cleanup Operations
	DeleteOldAuditLogsFunc                func(ctx context.Context, cutoffTime time.Time) (int64, error)
	DeleteAuditLogsByUserFunc             func(ctx context.Context, userID string) (int64, error)
	DeleteAuditLogsByActionFunc           func(ctx context.Context, action string, cutoffTime time.Time) (int64, error)

	// Export Operations
	GetUserAuditLogForExportFunc          func(ctx context.Context, userID string) ([]models.AuditLog, error)
	GetAllAuditLogsForExportFunc          func(ctx context.Context, startDate, endDate time.Time) ([]models.AuditLog, error)

	// Compliance Operations
	GetComplianceAuditLogsFunc            func(ctx context.Context, since time.Time) ([]models.AuditLog, error)
	GetUserComplianceAuditLogsFunc        func(ctx context.Context, userID string) ([]models.AuditLog, error)

	// Validation Operations
	CheckIfAuditLogExistsFunc             func(ctx context.Context, logID int) (bool, error)
	CheckIfUserHasAuditLogsFunc           func(ctx context.Context, userID string) (bool, error)
}

// Audit Log Operations
func (m *MockAuditDbAdapter) CreateAuditLog(ctx context.Context, userID, action string, timestamp time.Time, details map[string]interface{}) error {
	if m.CreateAuditLogFunc != nil {
		return m.CreateAuditLogFunc(ctx, userID, action, timestamp, details)
	}
	return nil
}

func (m *MockAuditDbAdapter) GetUserAuditLog(ctx context.Context, userID string) ([]models.AuditLog, error) {
	if m.GetUserAuditLogFunc != nil {
		return m.GetUserAuditLogFunc(ctx, userID)
	}
	return []models.AuditLog{}, nil
}

func (m *MockAuditDbAdapter) GetUserAuditLogByDateRange(ctx context.Context, userID string, startDate, endDate time.Time, limit int) ([]models.AuditLog, error) {
	if m.GetUserAuditLogByDateRangeFunc != nil {
		return m.GetUserAuditLogByDateRangeFunc(ctx, userID, startDate, endDate, limit)
	}
	return []models.AuditLog{}, nil
}

func (m *MockAuditDbAdapter) GetAllAuditLogsByDateRange(ctx context.Context, startDate, endDate time.Time, limit int) ([]models.AuditLog, error) {
	if m.GetAllAuditLogsByDateRangeFunc != nil {
		return m.GetAllAuditLogsByDateRangeFunc(ctx, startDate, endDate, limit)
	}
	return []models.AuditLog{}, nil
}

func (m *MockAuditDbAdapter) GetAuditLogsByAction(ctx context.Context, action string, limit int) ([]models.AuditLog, error) {
	if m.GetAuditLogsByActionFunc != nil {
		return m.GetAuditLogsByActionFunc(ctx, action, limit)
	}
	return []models.AuditLog{}, nil
}

func (m *MockAuditDbAdapter) GetUserAuditLogsByAction(ctx context.Context, userID, action string, limit int) ([]models.AuditLog, error) {
	if m.GetUserAuditLogsByActionFunc != nil {
		return m.GetUserAuditLogsByActionFunc(ctx, userID, action, limit)
	}
	return []models.AuditLog{}, nil
}

func (m *MockAuditDbAdapter) GetAuditLogByID(ctx context.Context, logID int) (*models.AuditLog, error) {
	if m.GetAuditLogByIDFunc != nil {
		return m.GetAuditLogByIDFunc(ctx, logID)
	}
	return nil, nil
}

// Statistics Operations
func (m *MockAuditDbAdapter) GetTotalAuditLogsCount(ctx context.Context) (int, error) {
	if m.GetTotalAuditLogsCountFunc != nil {
		return m.GetTotalAuditLogsCountFunc(ctx)
	}
	return 0, nil
}

func (m *MockAuditDbAdapter) GetOldestAuditLog(ctx context.Context) (*time.Time, error) {
	if m.GetOldestAuditLogFunc != nil {
		return m.GetOldestAuditLogFunc(ctx)
	}
	return nil, nil
}

func (m *MockAuditDbAdapter) GetAuditLogsCountByUser(ctx context.Context, userID string) (int, error) {
	if m.GetAuditLogsCountByUserFunc != nil {
		return m.GetAuditLogsCountByUserFunc(ctx, userID)
	}
	return 0, nil
}

func (m *MockAuditDbAdapter) GetAuditLogsCountByAction(ctx context.Context, since time.Time) (map[string]int, error) {
	if m.GetAuditLogsCountByActionFunc != nil {
		return m.GetAuditLogsCountByActionFunc(ctx, since)
	}
	return map[string]int{}, nil
}

func (m *MockAuditDbAdapter) GetAuditLogsCountByDateRange(ctx context.Context, startDate, endDate time.Time) (map[string]int, error) {
	if m.GetAuditLogsCountByDateRangeFunc != nil {
		return m.GetAuditLogsCountByDateRangeFunc(ctx, startDate, endDate)
	}
	return map[string]int{}, nil
}

// Cleanup Operations
func (m *MockAuditDbAdapter) DeleteOldAuditLogs(ctx context.Context, cutoffTime time.Time) (int64, error) {
	if m.DeleteOldAuditLogsFunc != nil {
		return m.DeleteOldAuditLogsFunc(ctx, cutoffTime)
	}
	return 0, nil
}

func (m *MockAuditDbAdapter) DeleteAuditLogsByUser(ctx context.Context, userID string) (int64, error) {
	if m.DeleteAuditLogsByUserFunc != nil {
		return m.DeleteAuditLogsByUserFunc(ctx, userID)
	}
	return 0, nil
}

func (m *MockAuditDbAdapter) DeleteAuditLogsByAction(ctx context.Context, action string, cutoffTime time.Time) (int64, error) {
	if m.DeleteAuditLogsByActionFunc != nil {
		return m.DeleteAuditLogsByActionFunc(ctx, action, cutoffTime)
	}
	return 0, nil
}

// Export Operations
func (m *MockAuditDbAdapter) GetUserAuditLogForExport(ctx context.Context, userID string) ([]models.AuditLog, error) {
	if m.GetUserAuditLogForExportFunc != nil {
		return m.GetUserAuditLogForExportFunc(ctx, userID)
	}
	return []models.AuditLog{}, nil
}

func (m *MockAuditDbAdapter) GetAllAuditLogsForExport(ctx context.Context, startDate, endDate time.Time) ([]models.AuditLog, error) {
	if m.GetAllAuditLogsForExportFunc != nil {
		return m.GetAllAuditLogsForExportFunc(ctx, startDate, endDate)
	}
	return []models.AuditLog{}, nil
}

// Compliance Operations
func (m *MockAuditDbAdapter) GetComplianceAuditLogs(ctx context.Context, since time.Time) ([]models.AuditLog, error) {
	if m.GetComplianceAuditLogsFunc != nil {
		return m.GetComplianceAuditLogsFunc(ctx, since)
	}
	return []models.AuditLog{}, nil
}

func (m *MockAuditDbAdapter) GetUserComplianceAuditLogs(ctx context.Context, userID string) ([]models.AuditLog, error) {
	if m.GetUserComplianceAuditLogsFunc != nil {
		return m.GetUserComplianceAuditLogsFunc(ctx, userID)
	}
	return []models.AuditLog{}, nil
}

// Validation Operations
func (m *MockAuditDbAdapter) CheckIfAuditLogExists(ctx context.Context, logID int) (bool, error) {
	if m.CheckIfAuditLogExistsFunc != nil {
		return m.CheckIfAuditLogExistsFunc(ctx, logID)
	}
	return false, nil
}

func (m *MockAuditDbAdapter) CheckIfUserHasAuditLogs(ctx context.Context, userID string) (bool, error) {
	if m.CheckIfUserHasAuditLogsFunc != nil {
		return m.CheckIfUserHasAuditLogsFunc(ctx, userID)
	}
	return false, nil
}