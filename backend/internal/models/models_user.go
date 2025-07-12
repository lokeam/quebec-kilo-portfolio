package models

import (
	"time"
)

type User struct {
	ID                   string     `json:"id" db:"id"`
	UserID               string     `json:"user_id" db:"user_id"`
	Email                string     `json:"email" db:"email"`
	FirstName            string     `json:"first_name" db:"first_name"`
	LastName             string     `json:"last_name" db:"last_name"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	// User deletion tracking fields
	DeletionRequestedAt *time.Time `json:"deletion_requested_at" db:"deletion_requested_at"`
	DeletionReason      *string    `json:"deletion_reason" db:"deletion_reason"`
	DeletedAt           *time.Time `json:"deleted_at" db:"deleted_at"`
}

// IsActive returns true if the user account is active (not deleted and no deletion requested)
func (u *User) IsActive() bool {
	return u.DeletedAt == nil && u.DeletionRequestedAt == nil
}

// IsDeleted returns true if the user account has been permanently deleted
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// IsDeletionRequested returns true if the user has requested account deletion
func (u *User) IsDeletionRequested() bool {
	return u.DeletionRequestedAt != nil && u.DeletedAt == nil
}

// IsInGracePeriod returns true if the user has requested deletion but is still in the grace period
func (u *User) IsInGracePeriod() bool {
	if !u.IsDeletionRequested() {
		return false
	}

	// Grace period is 30 days from deletion request
	gracePeriodEnd := u.GetDeletionGracePeriodEnd()
	if gracePeriodEnd == nil {
		return false
	}

	return time.Now().Before(*gracePeriodEnd)
}

// GetDeletionGracePeriodEnd returns the end date of the grace period (30 days from deletion request)
func (u *User) GetDeletionGracePeriodEnd() *time.Time {
	if u.DeletionRequestedAt == nil {
		return nil
	}

	gracePeriodEnd := u.DeletionRequestedAt.Add(30 * 24 * time.Hour)
	return &gracePeriodEnd
}
