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
	PhysicalLocationID string      `json:"physical_location_id" db:"physical_location_id"`
	Name               string      `json:"name" db:"name"`
	LocationType       string      `json:"location_type" db:"location_type"`
	StoredItems        int         `json:"stored_items" db:"stored_items"`
	CreatedAt          time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at" db:"updated_at"`

	// Field not directly stored in DB
	// Populated by a separate query
	Items              []Game      `json:"items" db:"-"`
}

// GameLocationMove represents a request to move a game from one sublocation to another
type GameLocationMove struct {
	UserGameID         string `json:"user_game_id" db:"user_game_id"`
	TargetSublocationID string `json:"target_sublocation_id" db:"target_sublocation_id"`
}

// GameLocationRemove represents a request to remove a game from a sublocation
type GameLocationRemove struct {
	UserGameID string `json:"user_game_id" db:"user_game_id"`
}
