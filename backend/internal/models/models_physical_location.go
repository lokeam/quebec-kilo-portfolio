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

var PhysicalLocationBgColor = struct {
	Red          string
	Green        string
	Blue         string
	Orange       string
	Gold         string
	Purple       string
	Brown        string
	Gray         string
	Pink         string
}{
	Red:          "red",
	Green:        "green",
	Blue:         "blue",
	Orange:       "orange",
	Gold:         "gold",
	Purple:       "purple",
	Brown:        "brown",
	Gray:         "gray",
	Pink:         "pink",
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

func IsValidPhysicalLocationBgColor(color string) bool {
	switch color {
	case PhysicalLocationBgColor.Red, PhysicalLocationBgColor.Green,
	     PhysicalLocationBgColor.Blue, PhysicalLocationBgColor.Orange,
			 PhysicalLocationBgColor.Gold, PhysicalLocationBgColor.Purple,
			 PhysicalLocationBgColor.Brown, PhysicalLocationBgColor.Gray,
			 PhysicalLocationBgColor.Pink:
		return true
	}
	return false
}

type PhysicalLocation struct {
	ID               string                     `json:"id" db:"id"`
	UserID           string                     `json:"user_id" db:"user_id"`
	Name             string                     `json:"name" db:"name"`
	Label            string                     `json:"label" db:"label"`
	BgColor          string                     `json:"bg_color" db:"bg_color"`
	LocationType     string                     `json:"location_type" db:"location_type"`
	MapCoordinates   PhysicalMapCoordinates     `json:"map_coordinates" db:"map_coordinates"`
	CreatedAt        time.Time                  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at" db:"updated_at"`
	SubLocations     *[]Sublocation             `json:"sub_locations" db:"sub_locations,json"`
}

type PhysicalMapCoordinates struct {
	Coords          string     `json:"coords"`
	GoogleMapsLink  string     `json:"google_maps_link"`
}