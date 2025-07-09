package audit

const (
	// ---------------- AUDIT LOG OPERATIONS ----------------
	CreateAuditLogQuery = `
		INSERT INTO audit_logs (user_id, action, timestamp, details, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	GetUserAuditLogQuery = `
		SELECT id, user_id, action, timestamp, details, created_at
		FROM audit_logs
		WHERE user_id = $1
		ORDER BY timestamp DESC
	`

	GetUserAuditLogByDateRangeQuery = `
		SELECT id, user_id, action, timestamp, details, created_at
		FROM audit_logs
		WHERE user_id = $1
		AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp DESC
		LIMIT $4
	`

	GetAllAuditLogsByDateRangeQuery = `
		SELECT id, user_id, action, timestamp, details, created_at
		FROM audit_logs
		WHERE timestamp BETWEEN $1 AND $2
		ORDER BY timestamp DESC
		LIMIT $3
	`

	GetAuditLogsByActionQuery = `
		SELECT id, user_id, action, timestamp, details, created_at
		FROM audit_logs
		WHERE action = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	GetUserAuditLogsByActionQuery = `
		SELECT id, user_id, action, timestamp, details, created_at
		FROM audit_logs
		WHERE user_id = $1 AND action = $2
		ORDER BY timestamp DESC
		LIMIT $3
	`

	// ---------------- STATISTICS QUERIES ----------------
	GetTotalAuditLogsCountQuery = `
		SELECT COUNT(*) FROM audit_logs
	`

	GetOldestAuditLogQuery = `
		SELECT MIN(timestamp) FROM audit_logs
	`

	GetAuditLogsCountByUserQuery = `
		SELECT user_id, COUNT(*) as log_count
		FROM audit_logs
		WHERE user_id = $1
		GROUP BY user_id
	`

	GetAuditLogsCountByActionQuery = `
		SELECT action, COUNT(*) as action_count
		FROM audit_logs
		WHERE timestamp >= $1
		GROUP BY action
		ORDER BY action_count DESC
	`

	GetAuditLogsCountByDateRangeQuery = `
		SELECT
			DATE(timestamp) as log_date,
			COUNT(*) as daily_count
		FROM audit_logs
		WHERE timestamp BETWEEN $1 AND $2
		GROUP BY DATE(timestamp)
		ORDER BY log_date DESC
	`

	// ---------------- CLEANUP QUERIES ----------------
	DeleteOldAuditLogsQuery = `
		DELETE FROM audit_logs WHERE timestamp < $1
	`

	DeleteAuditLogsByUserQuery = `
		DELETE FROM audit_logs WHERE user_id = $1
	`

	DeleteAuditLogsByActionQuery = `
		DELETE FROM audit_logs WHERE action = $1 AND timestamp < $2
	`

	// ---------------- EXPORT QUERIES ----------------
	GetUserAuditLogForExportQuery = `
		SELECT
			id,
			user_id,
			action,
			timestamp,
			details,
			created_at
		FROM audit_logs
		WHERE user_id = $1
		ORDER BY timestamp ASC
	`

	GetAllAuditLogsForExportQuery = `
		SELECT
			id,
			user_id,
			action,
			timestamp,
			details,
			created_at
		FROM audit_logs
		WHERE timestamp BETWEEN $1 AND $2
		ORDER BY timestamp ASC
	`

	// ---------------- COMPLIANCE QUERIES ----------------
	GetComplianceAuditLogsQuery = `
		SELECT
			id,
			user_id,
			action,
			timestamp,
			details,
			created_at
		FROM audit_logs
		WHERE action IN (
			'deletion_requested',
			'deletion_cancelled',
			'deletion_completed',
			'data_exported',
			'data_export_requested',
			'account_restored'
		)
		AND timestamp >= $1
		ORDER BY timestamp DESC
	`

	GetUserComplianceAuditLogsQuery = `
		SELECT
			id,
			user_id,
			action,
			timestamp,
			details,
			created_at
		FROM audit_logs
		WHERE user_id = $1
		AND action IN (
			'deletion_requested',
			'deletion_cancelled',
			'deletion_completed',
			'data_exported',
			'data_export_requested',
			'account_restored'
		)
		ORDER BY timestamp DESC
	`

	// ---------------- VALIDATION QUERIES ----------------
	CheckIfAuditLogExistsQuery = `
		SELECT EXISTS(SELECT 1 FROM audit_logs WHERE id = $1)
	`

	CheckIfUserHasAuditLogsQuery = `
		SELECT EXISTS(SELECT 1 FROM audit_logs WHERE user_id = $1)
	`

	GetAuditLogByIDQuery = `
		SELECT id, user_id, action, timestamp, details, created_at
		FROM audit_logs
		WHERE id = $1
	`
)