package types

// Represents the IGDB API response structure as defined by the current search query
type IGDBResponse struct {
	ID                   int64                        `json:"id"`
    Name               string                       `json:"name"`
    Summary            string                       `json:"summary"`
    Cover              IGDBResponseGameCover        `json:"cover"`
    Platforms          []IGDBResponseGamePlatform   `json:"platforms"`
    Genres             []IGDBResponseGameGenre      `json:"genres"`
    Themes             []IGDBResponseGameTheme      `json:"themes"`
    GameType           IGDBResponseGameType         `json:"game_type"`
    FirstReleaseDate   int64                        `json:"first_release_date"`
    Rating             float64                      `json:"rating"`
}

// Cover represents a game cover image
type IGDBResponseGameCover struct {
	ID    int64  `json:"id"`
	URL   string `json:"url"`
}

// Platform represents a gaming platform
type IGDBResponseGamePlatform struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
}

// Genre represents a game genre (e.g. "Role-playing (RPG)", "Adventure")
type IGDBResponseGameGenre struct {
	ID   int64  `json:"id"`
  Name string `json:"name"`
}

// Theme represents a game theme (e.g. "Action", "Fantasy", "Open World")
type IGDBResponseGameTheme struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GameType represents the type of game (e.g. "Main Game", "DLC", "Port")
type IGDBResponseGameType struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}
