package models

import "time"

type DigitalLocation struct {
	ID         string      `json:"id" db:"id"`
	UserID     string      `json:"user_id" db:"user_id"`
	Name       string      `json:"name" db:"name"`
	IsActive   bool        `json:"is_active" db:"is_active"`
	URL        string      `json:"url" db:"url"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" db:"updated_at"`
}
