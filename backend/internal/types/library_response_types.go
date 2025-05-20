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
	Favorite           bool                      `json:"favorite"`
	PhysicalLocations []PhysicalLocationResponse `json:"physicalLocations"`
	DigitalLocations  []DigitalLocationResponse  `json:"digitalLocations"`
}

// PhysicalLocationResponse represents a physical location for a game
type PhysicalLocationResponse struct {
	Name        string              `json:"name"`
	Type        string              `json:"type"`
	Sublocation SublocationResponse `json:"sublocation"`
	Platform    string              `json:"platform"`
}

// DigitalLocationResponse represents a digital location for a game
type DigitalLocationResponse struct {
	Name            string `json:"name"`
	NormalizedName  string `json:"normalizedName"`
	Platform        string `json:"platform"`
}

// SublocationResponse represents a sublocation within a physical location
type SublocationResponse struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	BgColor string `json:"bgColor"`
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
}

// Database result type for physical locations
type LibraryGamePhysicalLocationDBResponse struct {
	GameID                int64  `db:"game_id"`
	LocationName          string `db:"location_name"`
	LocationType          string `db:"location_type"`
	SublocationName       string `db:"sublocation_name"`
	SublocationType       string `db:"sublocation_type"`
	SublocationBgColor    string `db:"sublocation_bg_color"`
	PlatformName          string `db:"platform_name"`
}

// Database result type for digital locations
type LibraryGameDigitalLocationDBResponse struct {
	GameID          int64  `db:"game_id"`
	LocationName    string `db:"location_name"`
	NormalizedName  string `db:"normalized_name"`
	PlatformName    string `db:"platform_name"`
}

