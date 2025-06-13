package types

import "time"

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
	ID                    int64   `db:"id"`
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
	ID                    int64   `db:"game_id"`
	PlatformID            int64   `db:"platform_id"`
	PlatformName          string  `db:"platform_name"`
	LocationID            string  `db:"location_id"`
	LocationName          string  `db:"location_name"`
	IsActive              *bool   `db:"is_active"`
}

// New types

// LibraryGameItemResponse represents a game item in the library response
type LibraryGameItemResponse struct {
	ID                    int64    `json:"id"`
	Name                  string   `json:"name"`
	CoverURL             string   `json:"coverUrl"`
	FirstReleaseDate     int64    `json:"firstReleaseDate"`
	Rating               float64  `json:"rating"`
	IsInLibrary          bool     `json:"isInLibrary"`
	IsInWishlist         bool     `json:"isInWishlist"`
	IsUniqueCopy         bool     `json:"isUniqueCopy"`
	GameType             struct {
		DisplayText    string `json:"displayText"`
		NormalizedText string `json:"normalizedText"`
	} `json:"gameType"`
	Favorite             bool     `json:"favorite"`
	ThemeNames           []string `json:"themeNames,omitempty"`
	GamesByPlatformAndLocation []GameLocationDBResult `json:"gamesByPlatformAndLocation"`
}

// LibraryBFFResponse represents the response for the BFF endpoint
type LibraryBFFResponse struct {
	LibraryItems  []LibraryGameItemResponse `json:"libraryItems"`
	RecentlyAdded []LibraryGameItemResponse `json:"recentlyAdded"`
}


type GameLocationDBResult struct {
	ID                   int64   `db:"game_id"`
	PlatformID           int64   `db:"platform_id"`
	PlatformName         string  `db:"platform_name"`
	Type                 string  `db:"type"`
	LocationID           string  `db:"location_id"`
	LocationName         string  `db:"location_name"`
	LocationType         string  `db:"location_type"`
	SublocationID        *string `db:"sublocation_id"`
	SublocationName      *string `db:"sublocation_name"`
	SublocationType      *string `db:"sublocation_type"`
	SublocationBgColor   *string `db:"sublocation_bg_color"`
	IsActive             *bool   `db:"is_active"`
}

// --- YET ANOTHER REFACTOR FOR LIBRARY BFF RESPONSE
type LibraryGameItemGameTypeResponseFINAL struct {
	DisplayText    string `json:"display_text"`
	NormalizedText string `json:"normalizedText"`
}

type LibraryGameItemBFFResponseFINAL struct {
	ID       int64  `json:"id" db:"ID"`
	Name     string `json:"name" db:"Name"`
	CoverURL string `json:"cover_url" db:"CoverURL"`
	GameTypeDisplayText    string `json:"displayText" db:"GameTypeDisplayText"`
	GameTypeNormalizedText string `json:"normalizedText" db:"GameTypeNormalizedText"`
	IsFavorite                 bool      `json:"is_favorite" db:"IsFavorite"`
	CreatedAt                  time.Time `json:"created_at" db:"CreatedAt"`
	GamesByPlatformAndLocation []LibraryGamesByPlatformAndLocationItemFINAL `json:"games_by_platform_and_location"`
}

type LibraryGamesByPlatformAndLocationItemFINAL struct {
	ID                 int64     `json:"id" db:"game_id"`           // IGDB ID
	PlatformID         int64     `json:"platform_id" db:"platform_id"`
	PlatformName       string    `json:"platform_name" db:"platform_name"`
	IsPC               bool      `json:"is_pc" db:"is_pc"`      // must be computed
	IsMobile           bool      `json:"is_mobile" db:"is_mobile"`  // must be computed
	DateAdded          int64     `json:"date_added" db:"date_added"`

	// Parent location fields
	ParentLocationID       string    `json:"parent_location_id" db:"parent_location_id"`
	ParentLocationName     string    `json:"parent_location_name" db:"parent_location_name"`
	ParentLocationType     string    `json:"parent_location_type" db:"parent_location_type"`
	ParentLocationBgColor  string    `json:"parent_location_bg_color" db:"parent_location_bg_color"`

	// Sublocation fields
	SublocationID         string     `json:"sublocation_id" db:"sublocation_id"`
	SublocationName       string     `json:"sublocation_name" db:"sublocation_name"`
	SublocationType       string     `json:"sublocation_type" db:"sublocation_type"`
}

type LibraryBFFResponseFINAL struct {
	LibraryItems  []LibraryGameItemBFFResponseFINAL `json:"libraryItems"`
	RecentlyAdded []LibraryGameItemBFFResponseFINAL `json:"recentlyAdded"` // Only add recent games within the last 6 months
}