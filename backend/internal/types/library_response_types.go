package types

// LibraryGameResponse represents a game in the library
type LibraryGameResponse struct {
	ID                 int64    `json:"id"`
	Name               string   `json:"name"`
	CoverURL           string   `json:"coverUrl"`
	FirstReleaseDate   int64    `json:"firstReleaseDate"`
	Rating             float64  `json:"rating"`
	ThemeNames         []string `json:"themeNames"`
	IsInLibrary        bool     `json:"isInLibrary"`
	IsInWishlist       bool     `json:"isInWishlist"`
	GameType           struct {
		DisplayText      string   `json:"displayText"`
		NormalizedText   string   `json:"normalizedText"`
	}                           `json:"gameType"`
	Favorite           bool     `json:"favorite"`
	GamesByPlatformAndLocation []GameLocationDBResult `json:"gamesByPlatformAndLocation"`
}

// LibraryResponse represents the response for library operations
type LibraryResponse struct {
	Success   bool                    `json:"success"`
	Games     []LibraryGameResponse   `json:"games"`
}

// SingleGameResponse represents the response for a single game operation
type SingleGameResponse struct {
	Success   bool              `json:"success"`
	Game      LibraryGameResponse `json:"game"`
}

// Database result types for library games
type LibraryGameDBResult struct {
	ID                    int64     `db:"id"`
	Name                  string    `db:"name"`
	CoverURL              string    `db:"cover_url"`
	FirstReleaseDate      int64     `db:"first_release_date"`
	Rating                float64   `db:"rating"`
	ThemeNames            []string  `db:"theme_names"`
	Favorite              bool      `db:"favorite"`
	IsInWishlist          bool      `db:"is_in_wishlist"`
	GameTypeDisplay       string    `db:"game_type_display"`
	GameTypeNormalized    string    `db:"game_type_normalized"`
	PlatformID            int64     `db:"platform_id"`
	PlatformName          string    `db:"platform_name"`
}

// Database result type for physical locations
type LibraryGamePhysicalLocationDBResponse struct {
	GameID                int64   `db:"game_id"`
	PlatformID            int64   `db:"platform_id"`
	PlatformName          string  `db:"platform_name"`
	LocationID            string  `db:"location_id"`
	LocationName          string  `db:"location_name"`
	LocationType          string  `db:"location_type"`
	SublocationID         *string `db:"sublocation_id"`
	SublocationName       *string `db:"sublocation_name"`
	SublocationType       *string `db:"sublocation_type"`
	SublocationBgColor    *string `db:"sublocation_bg_color"`
}

// Database result type for digital locations
type LibraryGameDigitalLocationDBResponse struct {
	GameID                int64   `db:"game_id"`
	PlatformID            int64   `db:"platform_id"`
	PlatformName          string  `db:"platform_name"`
	LocationID            string  `db:"location_id"`
	LocationName          string  `db:"location_name"`
	IsActive              *bool   `db:"is_active"`
}

