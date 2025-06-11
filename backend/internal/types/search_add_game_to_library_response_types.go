package types

type AddGameFormStorageLocationsResponse struct {
	PhysicalLocations []AddGameFormPhysicalLocationsResponse `json:"physical_locations"`
	DigitalLocations  []AddGameFormDigitalLocationsResponse  `json:"digital_locations"`
}

type AddGameFormPhysicalLocationsResponse struct {
	ParentLocationID      string `json:"parentLocationId" db:"parent_location_id"`
	ParentLocationName    string `json:"parentLocationName" db:"parent_location_name"`
	ParentLocationType    string `json:"parentLocationType" db:"parent_location_type"`
	ParentLocationBgColor string `json:"parentLocationBgColor" db:"parent_location_bg_color"`
	SublocationID         string `json:"sublocationId" db:"sublocation_id"`
	SublocationName       string `json:"sublocationName" db:"sublocation_name"`
	SublocationType       string `json:"sublocationType" db:"sublocation_type"`
}

type AddGameFormDigitalLocationsResponse struct {
	DigitalLocationID    string `json:"digitalLocationId" db:"digital_location_id"`
	DigitalLocationName  string `json:"digitalLocationName" db:"digital_location_name"`
	IsSubscription       bool   `json:"isSubscription" db:"is_subscription"`
	IsActive             bool   `json:"isActive" db:"is_active"`
}