package igdb

// Field names for the IGDB API Query
const (
	// Basic fields
	FieldName               = "name"
	FieldSummary            = "summary"
	FieldFirstReleaseDate   = "first_release_date"
	FieldRating             = "rating"

	// Related entity fields
	FieldCoverURL        = "cover.url"
	FieldPlatformID      = "platforms.id"
	FieldPlatformName    = "platforms.name"
	FieldGenreName       = "genres.name"
	FieldThemeName       = "themes.name"
	FieldGameTypeType    = "game_type.type"
)

var (
	DefaultGameFields = []string{
		FieldName,
		FieldSummary,
		FieldFirstReleaseDate,
		FieldRating,
		FieldCoverURL,
		FieldPlatformID,
		FieldPlatformName,
		FieldGenreName,
		FieldThemeName,
		FieldGameTypeType,
	}
)

// Query filters
const (
	// GameTypeFilter represents the valid game types for our application
	// This is a business rule that defines which game types we want to include
	// See https://api-docs.igdb.com/#game-enums for more details
	GameTypeFilter = "game_type = (0,1,2,3,4,5,8,9)"
)

// Root URL for the IGDB API
const (
	BASE_IGDB_API_URL = "https://api.igdb.com/v4"
)
