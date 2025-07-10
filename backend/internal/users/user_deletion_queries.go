package users

const (
	RequestDeletionQuery = `
		UPDATE users
		SET deletion_requested_at = NOW(),
				deletion_reason = $2,
				updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	CancelDeletionRequestQuery = `
		UPDATE users
		SET deletion_requested_at = NULL,
				deletion_reason = NULL,
				updated_at = NOW()
		WHERE id = $1 AND deletion_requested_at IS NOT NULL AND deleted_at IS NULL
	`

	GetUsersPendingDeletionQuery = `
		SELECT id
		FROM users
		WHERE deletion_requested_at IS NOT NULL
			AND deletion_requested_at < NOW() - INTERVAL '30 days'
			AND deleted_at IS NULL
	`

	PermanentlyDeleteUserQuery = `
		UPDATE users
		SET deleted_at = NOW(),
				updated_at = NOW()
		WHERE id = $1 AND deletion_requested_at IS NOT NULL
	`

)