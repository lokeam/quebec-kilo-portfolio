package types

type Game struct {
    ID                int64    `json:"id"`
    Name              string   `json:"name"`
    Summary           string   `json:"summary,omitempty"`
    CoverID           int64    `json:"cover,omitempty"`
    CoverURL          string   `json:"cover_url,omitempty"`
    FirstReleaseDate  int64    `json:"first_release_date,omitempty"`
    Rating            float64  `json:"rating,omitempty"`
    Platforms         []int64  `json:"platforms,omitempty"`
    PlatformNames     []string `json:"platform_names,omitempty"`
    Genres            []int64  `json:"genres,omitempty"`
    GenreNames        []string `json:"genre_names,omitempty"`
    Themes            []int64  `json:"themes,omitempty"`
    ThemeNames        []string `json:"theme_names,omitempty"`
    IsInLibrary       bool     `json:"is_in_library"`
    IsInWishlist      bool     `json:"is_in_wishlist"`
}

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
}
