package models

import "time"

var PhysicalLocationType = struct {
	House         string
	Apartment     string
	Office        string
	Warehouse     string
	Vehicle       string
} {
	House:        "house",
	Apartment:    "apartment",
	Office:       "office",
	Warehouse:    "warehouse",
	Vehicle:      "vehicle",
}

// Check if string is a valid physical location type
func IsValidPhysicalLocationType(locationType string) bool {
	switch locationType {
		case PhysicalLocationType.House, PhysicalLocationType.Apartment,
		     PhysicalLocationType.Office, PhysicalLocationType.Warehouse,
				 PhysicalLocationType.Vehicle:
				 return true
	default:
		return false
	}
}

type PhysicalLocation struct {
	ID               string    `json:"id" db:"id"`
	UserID           string    `json:"user_id" db:"user_id"`
	Name             string    `json:"name" db:"name"`
	Label            string    `json:"label" db:"label"`
	LocationType     string    `json:"location_type" db:"location_type"`
	MapCoordinates   string    `json:"map_coordinates" db:"map_coordinates"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
  UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}
