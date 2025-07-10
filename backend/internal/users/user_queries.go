package users

const (
	/*  ---------- User Creation Queries ----------*/
	GetUserQuery = `
	SELECT id, email, first_name, last_name, created_at, updated_at,
				 deletion_requested_at, deletion_reason, deleted_at
	FROM users
	WHERE id = $1
`

CreateUserQuery = `
	INSERT INTO users (id, email, first_name, last_name, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, email, first_name, last_name, created_at, updated_at,
						deletion_requested_at, deletion_reason, deleted_at
`

UpdateUserProfileQuery = `
	UPDATE users
	SET first_name = $1, last_name = $2, updated_at = $3
	WHERE id = $4
	RETURNING id, email, first_name, last_name, created_at, updated_at,
						deletion_requested_at, deletion_reason, deleted_at
`

HasCompleteProfileQuery = `
	SELECT first_name, last_name
	FROM users
	WHERE id = $1
`

CheckUserExistsQuery = `
	SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)
`

GetSingleUserByEmailQuery = `
	SELECT id, email, first_name, last_name, created_at, updated_at,
				 deletion_requested_at, deletion_reason, deleted_at
	FROM users
	WHERE email = $1
`

	/*  ---------- User Deletion Queries ----------*/
)