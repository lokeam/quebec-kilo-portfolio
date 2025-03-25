package models

type Game struct {
	ID              int64     `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Summary         string    `json:"summary,omitempty" db:"summary"`
	CoverID         int64     `json:"cover,omitempty" db:"cover_id"`
	CoverURL        string    `json:"cover_url,omitempty" db:"cover_url"`
	FirstReleaseDate int64    `json:"first_release_date,omitempty" db:"first_release_date"`
	Rating          float64   `json:"rating,omitempty" db:"rating"`
	Platforms       []int64   `json:"platforms,omitempty" db:"platforms"`
	Genres          []int64   `json:"genres,omitempty" db:"genres"`
	Themes          []int64   `json:"themes,omitempty" db:"themes"`
	IsInLibrary     bool      `json:"is_in_library" db:"-"` // Use db:"-" for fields not in the database
	IsInWishlist    bool      `json:"is_in_wishlist" db:"-"`

	// NOTE: These fields won't be stored directly in the games table
	PlatformNames   []string  `json:"platform_names" db:"-"`
	GenreNames      []string  `json:"genre_names" db:"-"`
	ThemeNames      []string  `json:"theme_names" db:"-"`
}
