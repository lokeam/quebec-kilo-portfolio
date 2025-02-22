package types

// Client holds the IGDB API client configuration.

// Data types returned by IGDB.
type Game struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	Summary            string  `json:"summary,omitempty"`
	Cover              int     `json:"cover,omitempty"`    // Original cover ID from IGDB
	CoverURL           string  `json:"cover_url,omitempty"` // Populated after cover lookup
	FirstReleaseDate   int     `json:"first_release_date,omitempty"`
	Rating             float64 `json:"rating,omitempty"`
}

type Cover struct {
	ID  int    `json:"id"`
	AlphaChannel    bool   `json:"alpha_channel"`
	Animated        bool   `json:"animated"`
	Checksum        string `json:"checksum"`
	Game            int    `json:"game"`
	GameLocalization int   `json:"game_localization"`
	Height          int    `json:"height"`
	ImageID         string `json:"image_id"`
	URL             string `json:"url"`
	Width           int    `json:"width"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Platform struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GameDetails represents a game with expanded related data.
type GameDetails struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Summary         string    `json:"summary"`
	FirstReleaseDate int64     `json:"first_release_date"`
	Rating          float64   `json:"rating"`
	Cover           Cover     `json:"cover"`
	CoverURL        string    `json:"cover_url,omitempty"` // Add this field
	Genres          []Genre   `json:"genres"`
	Platforms       []Platform `json:"platforms"`
}