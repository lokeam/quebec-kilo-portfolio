package models

import (
	"github.com/lokeam/qko-beta/internal/types"
)

type Game struct {
	ID               int64     `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Label            string    `json:"label" db:"label"`
	Summary          string    `json:"summary" db:"summary"`
	CoverID          int64     `json:"cover,omitempty" db:"cover_id"`
	CoverURL         string    `json:"cover_url" db:"cover_url"`
	FirstReleaseDate int64     `json:"first_release_date" db:"first_release_date"`
	Rating           float64   `json:"rating" db:"rating"`

	Genres           []int64   `json:"genres,omitempty" db:"genres"`
	Themes           []int64   `json:"themes,omitempty" db:"themes"`
	GameType         types.GameType `json:"-" db:"game_type_id"`        // used when saving to the database
	GameTypeResponse types.GameTypeResponse `json:"game_type" db:"-"`   // sent as JSON to the frontend
	IsInLibrary      bool      `json:"is_in_library" db:"-"`            // Use db:"-" for fields not in the database
	IsInWishlist     bool      `json:"is_in_wishlist" db:"-"`

	// NOTE: These fields won't be stored directly in the games table
	//Platforms        []int64   `json:"platforms,omitempty" db:"platforms"`
	Platforms       []PlatformInfo `json:"platforms" db:"-"`
	PlatformNames   []string  `json:"platform_names" db:"-"`
	GenreNames      []string  `json:"genre_names" db:"-"`
	ThemeNames      []string  `json:"theme_names" db:"-"`

	// Additional fields for frontend compatibility
	Platform        string    `json:"platform" db:"-"`
	PlatformVersion string    `json:"platform_version" db:"-"`
	AcquiredDate    string    `json:"acquired_date" db:"-"`
	Condition       string    `json:"condition,omitempty" db:"-"`
	HasOriginalCase bool      `json:"has_original_case,omitempty" db:"-"`
	HasManual       bool      `json:"has_manual,omitempty" db:"-"`
}

type PlatformInfo struct {
	ID   int64     `json:"id"`
	Name string    `json:"name"`
}