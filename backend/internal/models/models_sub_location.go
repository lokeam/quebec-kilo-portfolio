package models

import "time"

var SubLocationType = struct {
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
	case SubLocationType.Shelf, SubLocationType.Console,
	      SubLocationType.Cabinet, SubLocationType.Closet,
				SubLocationType.Drawer, SubLocationType.Box,
				SubLocationType.Device:
				return true
	default:
		return false
	}
}

type SubLocation struct {
	ID                 string      `json:"id" db:"id"`
	UserID             string      `json:"user_id" db:"user_id"`
	Name               string      `json:"name" db:"name"`
	Label              string      `json:"label" db:"label"` // TODO: determine if this is worthwhile in UAT
	Description        string      `json:"description" db:"description"` // TODO: determine if this is worthwhile in UAT
	LocationType       string      `json:"location_type" db:"location_type"`
	Items              []Game      `json:"items" db:"items"`
	Capacity           int         `json:"capacity" db:"capacity"`
	IsAccessible       bool        `json:"is_accessible" db:"is_accessible"`
	CreatedAt          time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at" db:"updated_at"`
}
