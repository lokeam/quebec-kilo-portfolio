package models

import (
	"database/sql"
	"time"
)

// This struct replaces LibraryGame
type GameToSave struct {
	GameID                int64
	GameName              string
	GameCoverURL          string
	GameFirstReleaseDate  int64
	GameType              GameToSaveIGDBType
	GameThemeNames        []string
	PlatformLocations     []GameToSaveLocation
	GameRating            float64
}

type GameToSaveLocation struct {
	PlatformID   int64
	PlatformName string
	Type         string  // "physical" or "digital"
	Location     GameToSaveLocationDetails
}

type GameToSaveIGDBType struct {
	DisplayText     string `json:"display_text"`
	NormalizedText  string `json:"normalized_text"`
}

type GameToSaveLocationDetails struct {
	SublocationID     string  `json:"sublocation_id,omitempty"`
	DigitalLocationID string  `json:"digital_location_id,omitempty"`
}




// LibraryGame model, matches the CreateLibraryGameRequest struct
//go:deprecated Use LibraryGameCreateRequest or LibraryGameUpdateRequest instead
// type LibraryGame struct {
// 	GameID              int64
// 	GameName            string
// 	GameCoverURL        string
// 	GameFirstReleaseDate int64
// 	GameType            LibraryGameType
// 	GameThemeNames      []string
// 	PlatformLocations   []CreateLibraryGameLocation
// 	GameRating          float64
// }

// // Deprecated: Use GameToSaveIGDBType instead
// type LibraryGameType struct {
// 	DisplayText     string `json:"display_text"`
// 	NormalizedText  string `json:"normalized_text"`
// }


// type CreateLibraryGameLocation struct {
// 	PlatformID   int64                `json:"platform_id"`
// 	PlatformName string               `json:"platform_name"`
// 	Type         string               `json:"type"`
// 	Location     LibraryGameLocation  `json:"location"`
// }

// type LibraryGameLocation struct {
// 	SublocationID     string  `json:"sublocation_id,omitempty"`
// 	DigitalLocationID string  `json:"digital_location_id,omitempty"`
// }

// LibraryLocationDB represents the database model for library BFF operations
type GameLocationDatabaseEntry struct {
	GameID                 int64          `db:"game_id"`
	PlatformID             int64          `db:"platform_id"`
	PlatformName           string         `db:"platform_name"`
	Category               string         `db:"category"`
	CreatedAt              time.Time      `db:"created_at"`
	ParentLocationID       sql.NullString `db:"parent_location_id"`
	ParentLocationName     sql.NullString `db:"parent_location_name"`
	ParentLocationType     sql.NullString `db:"parent_location_type"`
	ParentLocationBgColor  sql.NullString `db:"parent_location_bg_color"`
	SublocationID          sql.NullString `db:"sublocation_id"`
	SublocationName        sql.NullString `db:"sublocation_name"`
	SublocationType        sql.NullString `db:"sublocation_type"`
}

// Deprecated: Use GameLocationDatabaseEntry instead
// type LibraryLocationDB struct {
// 	GameID              int64          `db:"game_id"`
// 	PlatformID          int64          `db:"platform_id"`
// 	PlatformName        string         `db:"platform_name"`
// 	Category            string         `db:"category"`
// 	CreatedAt           time.Time      `db:"created_at"`
// 	ParentLocationID    sql.NullString `db:"parent_location_id"`
// 	ParentLocationName  sql.NullString `db:"parent_location_name"`
// 	ParentLocationType  sql.NullString `db:"parent_location_type"`
// 	ParentLocationBgColor sql.NullString `db:"parent_location_bg_color"`
// 	SublocationID       sql.NullString `db:"sublocation_id"`
// 	SublocationName     sql.NullString `db:"sublocation_name"`
// 	SublocationType     sql.NullString `db:"sublocation_type"`
//}

// LibraryGameDB represents the database model for library games
type LibraryGameDB struct {
	ID                  int64     `db:"id"`
	Name                string    `db:"name"`
	CoverURL            string    `db:"cover_url"`
	GameTypeDisplayText string    `db:"game_type_display_text"`
	GameTypeNormalizedText string `db:"game_type_normalized_text"`
	IsFavorite          bool      `db:"is_favorite"`
	CreatedAt           time.Time `db:"created_at"`
}
