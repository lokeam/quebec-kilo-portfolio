package models

// LibraryGame model, matches the CreateLibraryGameRequest struct
type LibraryGame struct {
	GameID              int64
	GameName            string
	GameCoverURL        string
	GameFirstReleaseDate int64
	GameType            LibraryGameType
	GameThemeNames      []string
	PlatformLocations   []CreateLibraryGameLocation
	GameRating          float64
}

type LibraryGameType struct {
	DisplayText     string `json:"display_text"`
	NormalizedText  string `json:"normalized_text"`
}

type CreateLibraryGameLocation struct {
	PlatformID   int64                `json:"platform_id"`
	PlatformName string               `json:"platform_name"`
	Type         string               `json:"type"`
	Location     LibraryGameLocation  `json:"location"`
}

type LibraryGameLocation struct {
	SublocationID     string  `json:"sublocation_id,omitempty"`
	DigitalLocationID string  `json:"digital_location_id,omitempty"`
}
