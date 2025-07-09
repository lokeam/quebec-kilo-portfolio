package audit

// Audit actions for user deletion and data management
const (
	// User deletion actions
	ActionDeletionRequested    = "deletion_requested"
	ActionDeletionCancelled    = "deletion_cancelled"
	ActionDeletionCompleted    = "deletion_completed"
	ActionAccountRestored      = "account_restored"

	// Data export actions
	ActionDataExported         = "data_exported"
	ActionDataExportRequested  = "data_export_requested"

	// Account management actions
	ActionAccountCreated       = "account_created"
	ActionAccountUpdated       = "account_updated"
	ActionAccountAccessed      = "account_accessed"

	// Data management actions
	ActionDataDeleted          = "data_deleted"
	ActionDataModified         = "data_modified"
	ActionDataAccessed         = "data_accessed"

	// Compliance actions
	ActionComplianceCheck      = "compliance_check"
	ActionRetentionPolicy      = "retention_policy_applied"
	ActionAuditLogCleaned      = "audit_log_cleaned"
)

// Audit categories for organizing audit logs
const (
	CategoryUserManagement     = "user_management"
	CategoryDataManagement     = "data_management"
	CategoryCompliance         = "compliance"
	CategorySecurity           = "security"
	CategorySystem             = "system"
)

// Audit severity levels
const (
	SeverityLow       = "low"
	SeverityMedium    = "medium"
	SeverityHigh      = "high"
	SeverityCritical  = "critical"
)

// Retention periods (in days)
const (
	RetentionDaysDefault = 2555 // 7 years
	RetentionDaysMin     = 1825 // 5 years (minimum for most regulations)
	RetentionDaysMax     = 3650 // 10 years (maximum for most regulations)
)

// File rotation settings
const (
	LogFileDateFormat = "2006-01-02"
	LogFilePrefix     = "audit_"
	LogFileExtension  = ".log"
)

// Database table names
const (
	TableAuditLogs = "audit_logs"
)

// Default log directory
const (
	DefaultLogDir = "logs/audit"
)