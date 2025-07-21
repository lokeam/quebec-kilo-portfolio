package types

import "time"

// ---------- Physical Locations ----------

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

type LocationsBFFStoredGameResponse struct {
    ID               int64  `json:"id"`
    Name             string `json:"name"`
    Platform         string `json:"platform"`
    IsUniqueCopy     bool   `json:"is_unique_copy"`
    HasDigitalCopy   bool   `json:"has_digital_copy"`
}

type LocationsBFFSublocationResponse struct {
    // Sublocation fields
    SublocationID         string                           `json:"sublocation_id"`
    SublocationName       string                           `json:"sublocation_name"`
    SublocationType       string                           `json:"sublocation_type"`
    StoredItems           int                              `json:"stored_items"`
    StoredGames           []LocationsBFFStoredGameResponse `json:"stored_games"`


    // Parent physical location fields
    ParentLocationID      string                   `json:"parent_location_id"`
    ParentLocationName    string                   `json:"parent_location_name"`
    ParentLocationType    string                   `json:"parent_location_type"`
    ParentLocationBgColor string                   `json:"parent_location_bg_color"`
    MapCoordinates        MapCoordinatesResponse   `json:"map_coordinates"`
    CreatedAt             time.Time                `json:"created_at"`
    UpdatedAt             time.Time                `json:"updated_at"`
}

// DeletedGameDetails represents a game that was deleted along with a sublocation
type DeletedGameDetails struct {
    UserGameID   int    `json:"user_game_id"`
    GameID       int64  `json:"game_id"`
    GameName     string `json:"game_name"`
    PlatformName string `json:"platform_name"`
}

// DeleteSublocationResponse represents the response from a sublocation deletion operation
type DeleteSublocationResponse struct {
    Success         bool                    `json:"success"`
    DeletedCount    int                     `json:"deleted_count"`
    SublocationIDs  []string                `json:"sublocation_ids"`
    DeletedGames    []DeletedGameDetails    `json:"deleted_games"`
    Error           string                  `json:"error,omitempty"`
}

// ---------- Digital Locations ----------
type DigitalLocationsBFFResponse struct {
    DigitalLocations []SingleDigitalLocationBFFResponse `json:"digital_locations"`
}

// Digital Location BFF Response
// Field removed: LocationType string `json:"location_type"` // MARKED FOR DELETION
type SingleDigitalLocationBFFResponse struct {
    ID                 string                              `json:"id"`
	Name               string                              `json:"name"`
    URL                string                              `json:"url"`
    IsSubscription     bool                                `json:"is_subscription"`
	IsActive           bool                                `json:"is_active"`
	BillingCycle       string                              `json:"billing_cycle"`
	CostPerCycle       float64                             `json:"cost_per_cycle"`
    MonthlyCost        float64                             `json:"monthly_cost"`
	NextPaymentDate    *time.Time                          `json:"next_payment_date"`
    PaymentMethod      string                              `json:"payment_method"`
	ItemCount          int                                 `json:"item_count"`
	StoredGames        []DigitalLocationGameResponse       `json:"stored_games"`
	CreatedAt          time.Time                           `json:"created_at"`
	UpdatedAt          time.Time                           `json:"updated_at"`
}

type DigitalLocationGameResponse struct {
    ID               int64  `json:"id"`
    Name             string `json:"name"`
    Platform         string `json:"platform"`
    IsUniqueCopy     bool   `json:"is_unique_copy"`
    HasPhysicalCopy  bool   `json:"has_physical_copy"`
}