package types

// LibraryGameResponse represents a game in the library
type LibraryGameResponse struct {
	ID                 int64    `json:"id"`
	Name               string   `json:"name"`
	CoverURL           string   `json:"cover_url"`
	FirstReleaseDate   int64    `json:"first_release_date"`
	Rating             float64  `json:"rating"`
	ThemeNames         []string `json:"theme_names"`
	IsInLibrary        bool     `json:"is_in_library"`
	IsInWishlist       bool     `json:"is_in_wishlist"`
	GameType           struct {
		DisplayText      string   `json:"displayText"`
		NormalizedText   string   `json:"normalizedText"`
	}                           `json:"gameType"`
	Favorite           bool     `json:"favorite"`
	GamesByPlatformAndLocation []GameLocationDBResponse `json:"games_by_platform_and_location"`
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
	ID                int64  `json:"id"`
	PlatformID        int64  `json:"platform_id"`
	PlatformName      string `json:"platform_name"`
	LocationID        string `json:"location_id"`
	LocationName      string `json:"location_name"`
	LocationType      string `json:"location_type"`
	SublocationID     string `json:"sublocation_id"`
	SublocationName   string `json:"sublocation_name"`
	SublocationType   string `json:"sublocation_type"`
	SublocationBgColor string `json:"sublocation_bg_color"`
}

// Database result type for digital locations
type LibraryGameDigitalLocationDBResponse struct {
	ID           int64  `json:"id"`
	PlatformID   int64  `json:"platform_id"`
	PlatformName string `json:"platform_name"`
	LocationID   string `json:"location_id"`
	LocationName string `json:"location_name"`
	IsActive     bool   `json:"is_active"`
}

// LibraryGamesByPlatformAndLocationItemFINAL represents a platform and location item in the BFF response
type LibraryGamesByPlatformAndLocationItemFINAL struct {
	ID                 int64     `json:"id"`
	PlatformID         int64     `json:"platform_id"`
	PlatformName       string    `json:"platform_name"`
	IsPC               bool      `json:"is_pc"`
	IsMobile           bool      `json:"is_mobile"`
	DateAdded          int64     `json:"date_added"`
	ParentLocationID   string    `json:"parent_location_id"`
	ParentLocationName string    `json:"parent_location_name"`
	ParentLocationType string    `json:"parent_location_type"`
	ParentLocationBgColor string `json:"parent_location_bg_color"`
	SublocationID      string    `json:"sublocation_id"`
	SublocationName    string    `json:"sublocation_name"`
	SublocationType    string    `json:"sublocation_type"`
}

// LibraryGameItemBFFResponseFINAL represents a game item in the BFF response
type LibraryGameItemBFFResponseFINAL struct {
	ID                    int64                                    `json:"id"`
	Name                  string                                   `json:"name"`
	CoverURL              string                                   `json:"cover_url"`
	GameTypeDisplayText   string                                   `json:"game_type_display_text"`
	GameTypeNormalizedText string                                  `json:"game_type_normalized_text"`
	IsFavorite            bool                                     `json:"is_favorite"`
	GamesByPlatformAndLocation []LibraryGamesByPlatformAndLocationItemFINAL `json:"games_by_platform_and_location"`
}

// LibraryBFFResponseFINAL represents the final BFF response structure
type LibraryBFFResponseFINAL struct {
	LibraryItems  []LibraryGameItemBFFResponseFINAL `json:"library_items"`
	RecentlyAdded []LibraryGameItemBFFResponseFINAL `json:"recently_added"`
}

// GameLocationDBResponse represents a game's location in the database
type GameLocationDBResponse struct {
	// ID is the game's IGDB ID (int64)
	ID                 int64    `json:"id"`
	// PlatformID is the platform's IGDB ID (int64)
	PlatformID         int64    `json:"platform_id"`
	PlatformName       string   `json:"platform_name"`
	Type               string   `json:"type"`
	// LocationID is a UUID string for the location
	LocationID         string   `json:"location_id"`
	LocationName       string   `json:"location_name"`
	LocationType       string   `json:"location_type"`
	// SublocationID is a UUID string for the sublocation
	SublocationID      string   `json:"sublocation_id"`
	SublocationName    string   `json:"sublocation_name"`
	SublocationType    string   `json:"sublocation_type"`
	SublocationBgColor string   `json:"sublocation_bg_color"`
	IsActive           bool     `json:"is_active"`
}