package models

import "time"

var SublocationType = struct {
	Shelf    string
	Console  string
	Cabinet  string
	Closet   string
	Drawer   string
	Box      string
	Device   string
} {
	Shelf:     "shelf",
	Console:   "console",
	Cabinet:   "cabinet",
	Closet:    "closet",
	Drawer:    "drawer",
	Box:       "box",
	Device:    "device",
}

// Check if a string is a valid sub-location type
func IsValidSublocationType(locationType string) bool {
	switch locationType {
	case SublocationType.Shelf, SublocationType.Console,
	      SublocationType.Cabinet, SublocationType.Closet,
				SublocationType.Drawer, SublocationType.Box,
				SublocationType.Device:
				return true
	default:
		return false
	}
}

type Sublocation struct {
	ID                 string      `json:"id" db:"id"`
	UserID             string      `json:"user_id" db:"user_id"`
	Name               string      `json:"name" db:"name"`
	Label              string      `json:"label" db:"label"` // TODO: determine if this is worthwhile in UAT
	Description        string      `json:"description" db:"description"` // TODO: determine if this is worthwhile in UAT
	LocationType       string      `json:"location_type" db:"location_type"`
	Capacity           int         `json:"capacity" db:"capacity"`
	IsAccessible       bool        `json:"is_accessible" db:"is_accessible"`
	CreatedAt          time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at" db:"updated_at"`

	// Field not directly stored in DB
	// Populated by a separate query
	Items              []Game      `json:"items" db:"-"`
}
