// backend/internal/types/physical_response_types.go
package types

import "time"

// PhysicalLocationResponse represents a physical location for the frontend
type PhysicalLocationResponse struct {
    ID             string                    `json:"id"`
    Name           string                    `json:"name"`
    Label          string                    `json:"label"`
    LocationType   string                    `json:"location_type"`
    MapCoordinates MapCoordinatesResponse    `json:"map_coordinates"`
    BgColor        string                    `json:"bg_color"`
    CreatedAt      time.Time                 `json:"created_at"`
    UpdatedAt      time.Time                 `json:"updated_at"`
    Sublocations   []SublocationResponse     `json:"sublocations"`
}

// MapCoordinatesResponse represents map coordinates for a location
type MapCoordinatesResponse struct {
    Coords         string `json:"coords"`
    GoogleMapsLink string `json:"google_maps_link"`
}

// ItemResponse represents an item stored in a sublocation
type ItemResponse struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    Platform      string    `json:"platform"`
    AcquiredDate  time.Time `json:"acquired_date"`
}

// Used for unmarshalling API responses FROM DB requests; not FOR DB operations
type SublocationResponse struct {
	ID                              string           `json:"id"`
	UserID                          string           `json:"user_id"`
	PhysicalLocationID              string           `json:"physical_location_id"`
	Name                            string           `json:"name"`
	LocationType                    string           `json:"location_type"`
	StoredItems                     int              `json:"stored_items"`
	CreatedAt                       time.Time        `json:"created_at"`
	UpdatedAt                       time.Time        `json:"updated_at"`
	ParentPhysicalLocationName      string           `json:"parent_physical_location_name"`
	ParentPhysicalLocationBgColor   string           `json:"parent_physical_location_bg_color"`
	ParentPhysicalLocationType      string           `json:"parent_physical_location_type"`
}

// Location page BFF Specific response structures
type LocationsBFFResponse struct {
    PhysicalLocations []LocationsBFFPhysicalLocationResponse  `json:"physical_locations"`
    Sublocations      []LocationsBFFSublocationResponse       `json:"sublocations"`
}


type LocationsBFFPhysicalLocationResponse struct {
    PhysicalLocationID     string                  `json:"physical_location_id"`
    Name                   string                  `json:"name"`
    PhysicalLocationType   string                  `json:"physical_location_type"`
    MapCoordinates         MapCoordinatesResponse  `json:"map_coordinates"`
    BgColor                string                  `json:"bg_color"`
    CreatedAt              time.Time               `json:"created_at"`
    UpdatedAt              time.Time               `json:"updated_at"`
}

type LocationsBFFSublocationResponse struct {
    // Sublocation fields
    SublocationID         string    `json:"sublocation_id"`
    SublocationName       string    `json:"sublocation_name"`
    SublocationType       string    `json:"sublocation_type"`
    StoredItems           int       `json:"stored_items"`

    // Parent physical location fields
    ParentLocationID      string                   `json:"parent_location_id"`
    ParentLocationName    string                   `json:"parent_location_name"`
    ParentLocationType    string                   `json:"parent_location_type"`
    ParentLocationBgColor string                   `json:"parent_location_bg_color"`
    MapCoordinates        MapCoordinatesResponse   `json:"map_coordinates"`
    CreatedAt             time.Time                `json:"created_at"`
    UpdatedAt             time.Time                `json:"updated_at"`
}