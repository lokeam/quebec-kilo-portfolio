package types

type CreateLibraryGameRequest struct {
	GameID                       int64                           `json:"game_id"`
	GameName                     string                          `json:"game_name"`
	GameCoverURL                 string                          `json:"game_cover_url"`
	GameFirstReleaseDate         int64                           `json:"game_first_release_date"`
	GameType                     LibraryRequestGameType          `json:"game_type"`
	GameThemeNames               []string                        `json:"game_theme_names"`
	GameRating                   float64                         `json:"game_rating"`
	GamesByPlatformAndLocation   []LibraryRequestGameLocation    `json:"games_by_platform_and_location"`
}

type UpdateLibraryGameRequest struct {
	GamesByPlatformAndLocation    []LibraryRequestGameLocation    `json:"games_by_platform_and_location"`
}

// BatchDeleteLibraryGameRequest represents a request to delete specific platform versions of a game
type BatchDeleteLibraryGameRequest struct {
	GameID    int64                           `json:"game_id"`
	DeleteAll bool                            `json:"delete_all"`
	Versions  []BatchDeleteLibraryGameVersion `json:"versions"`
}

// BatchDeleteLibraryGameVersion represents a specific version to delete
type BatchDeleteLibraryGameVersion struct {
	Type        string `json:"type"`         // "physical" or "digital"
	LocationID  string `json:"location_id"`  // sublocation_id or digital_location_id
	PlatformID  int64  `json:"platform_id"`  // platform ID
}

type LibraryRequestGameType struct {
	DisplayText     string `json:"display_text"`
	NormalizedText  string `json:"normalized_text"`
}

type LibraryRequestGameLocation struct {
	PlatformID   int64        `json:"platform_id"`
	PlatformName string       `json:"platform_name"`
	Type         string       `json:"type"`
	Location     GameLocation `json:"location"`
}

type GameLocation struct {
	SublocationID     string  `json:"sublocation_id,omitempty"`
	DigitalLocationID string  `json:"digital_location_id,omitempty"`
}