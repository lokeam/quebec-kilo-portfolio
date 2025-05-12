package types

type Cover struct {
    ID               int64  `json:"id"`
    AlphaChannel     bool   `json:"alpha_channel"`
    Animated         bool   `json:"animated"`
    Checksum         string `json:"checksum"`
    Game             int    `json:"game"`
    GameLocalization int    `json:"game_localization"`
    Height           int    `json:"height"`
    ImageID          string `json:"image_id"`
    URL              string `json:"url"`
    Width            int    `json:"width"`
}

type Genre struct {
    ID   int64    `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
}

type Platform struct {
    ID   int64    `json:"id"`
    Name string   `json:"name"`
}

type Theme struct {
    ID   int64    `json:"id"`
    Name string   `json:"name"`
}

type IGDBGameType struct {
    ID   int    `json:"id"`
    Type string `json:"type"`
}

type GameDetails struct {
    ID                int64      `json:"id"`
    Name              string     `json:"name"`
    Summary           string     `json:"summary,omitempty"`
    FirstReleaseDate  int64      `json:"first_release_date,omitempty"`
    Rating            float64    `json:"rating,omitempty"`
    Cover             Cover      `json:"cover"`
    CoverURL          string     `json:"cover_url"`
    Genres            []Genre    `json:"genres"`
    GenreNames        []string   `json:"genre_names"`
    Platforms         []Platform `json:"platforms,omitempty"`
    PlatformNames     []string   `json:"platform_names,omitempty"`
    Themes            []Theme    `json:"themes"`
    ThemeNames        []string   `json:"theme_names"`
    GameType          IGDBGameType `json:"game_type"`
}
