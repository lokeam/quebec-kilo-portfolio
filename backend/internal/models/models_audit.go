package models

import (
	"time"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        int                    `json:"id" db:"id"`
	UserID    string                 `json:"user_id" db:"user_id"`
	Action    string                 `json:"action" db:"action"`
	Timestamp time.Time              `json:"timestamp" db:"timestamp"`
	Details   map[string]interface{} `json:"details" db:"details"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}