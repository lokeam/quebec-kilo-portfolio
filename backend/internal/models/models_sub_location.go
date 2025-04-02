package models

import "time"

var SublocationBgColor = struct {
	Red          string
	Green        string
	Blue         string
	Orange       string
	Gold         string
	Purple       string
	Brown        string
	Gray         string
}{
	Red:          "red",
	Green:        "green",
	Blue:         "blue",
	Orange:       "orange",
	Gold:         "gold",
	Purple:       "purple",
	Brown:        "brown",
	Gray:         "gray",
}
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

func IsValidSublocationBgColor(color string) bool {
	switch color {
	case SublocationBgColor.Red, SublocationBgColor.Green,
	     SublocationBgColor.Blue, SublocationBgColor.Orange,
			 SublocationBgColor.Gold, SublocationBgColor.Purple,
			 SublocationBgColor.Brown, SublocationBgColor.Gray:
		return true
	}
	return false
}

type Sublocation struct {
	ID                 string      `json:"id" db:"id"`
	UserID             string      `json:"user_id" db:"user_id"`
	Name               string      `json:"name" db:"name"`
	LocationType       string      `json:"location_type" db:"location_type"`
	BgColor            string      `json:"bg_color" db:"bg_color"`
	Capacity           int         `json:"capacity" db:"capacity"`
	CreatedAt          time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at" db:"updated_at"`

	// Field not directly stored in DB
	// Populated by a separate query
	Items              []Game      `json:"items" db:"-"`
}
