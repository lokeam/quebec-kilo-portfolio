package search

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/igdb"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
	"github.com/sony/gobreaker"
)

/*
	A typical IGDB query needs to include the following fields:
		- id
		- name
		- summary
		- first_release_date
		- rating
		- cover
		- genres
		- themes
		- game_modes
		- platforms
*/
const (
	IGDBGameQueryTemplate = `fields id,name,summary,first_release_date,rating,cover,genres,themes,platforms,game_type.id,game_type.type; search "%s"; limit %d;`
)

// IGDBAdapter wraps the IGDB client and adds circuit breaker protection and logging.
// It acts as a bridge between our application's game search functionality and the IGDB API.
// The adapter pattern allows us to:
//   - Convert between IGDB's response format and our application's game model
//   - Add logging and error handling
//   - Implement caching if needed
//   - Make the IGDB API easier to use in our application
type IGDBAdapter struct {
	client  *igdb.IGDBClient             // Underlying IGDB client.
	breaker *gobreaker.CircuitBreaker    // Circuit breaker to protect IGDB calls.
	logger  interfaces.Logger            // Logger interface.
}

// NewIGDBAdapter creates a new IGDBAdapter instance.
// It sets up the connection to IGDB by:
//   1. Getting a Twitch token (required for IGDB authentication)
//   2. Creating an IGDB client with the token
//   3. Setting up logging
func NewIGDBAdapter(appContext *appcontext.AppContext) (*IGDBAdapter, error) {
	appContext.Logger.Debug("NewIGDBAdapter called", map[string]any{"appContext": appContext})

	// Retrieve Twitch token needed for IGDB authentication.
	twitchToken, err := appContext.TwitchTokenRetriever.GetToken(
			context.Background(),
			appContext.Config.IGDB.ClientID,
			appContext.Config.IGDB.ClientSecret,
			appContext.Config.IGDB.AuthURL,
			appContext.Logger,
	)
	if err != nil {
		appContext.Logger.Error("Failed to fetch Twitch token", map[string]any{
				"error": err,
				"clientID": appContext.Config.IGDB.ClientID,
				"authURL": appContext.Config.IGDB.AuthURL,
		})
		return nil, fmt.Errorf("failed to fetch Twitch token: %w", err)
	}
	appContext.Logger.Debug("Twitch token retrieved successfully", map[string]any{"token": twitchToken})

	// Log the headers for verification.
	appContext.Logger.Info("Headers for IGDB API request", map[string]any{
			"Client-ID": appContext.Config.IGDB.ClientID,
			"Authorization": fmt.Sprintf("Bearer %s", twitchToken),
	})

	// Create the IGDB client.
	client := igdb.NewIGDBClient(appContext, twitchToken)
	if client == nil {
			return nil, fmt.Errorf("failed to create IGDB client")
	}

	appContext.Logger.Info("Creating IGDBAdapter", map[string]any{
			"clientID": appContext.Config.IGDB.ClientID,
	})

	return &IGDBAdapter{
			client: client,
			logger: appContext.Logger,
	}, nil
}

// SearchGames searches for games in the IGDB database.
// It:
//   1. Creates a query using the QueryBuilder
//   2. Executes the query against IGDB
//   3. Converts the IGDB response into our application's game model
//   4. Handles any errors that occur during the process
//
// The query includes:
//   - Basic game information (id, name, summary)
//   - Release date and rating
//   - Cover image URL
//   - Platform, genre, and theme names
//   - Game type information
func (a *IGDBAdapter) SearchGames(ctx context.Context, query string, limit int) ([]*models.Game, error) {
	a.logger.Info("IGDB Adapter - UPDATED SearchGames called", map[string]any{
		"query": query,
		"limit": limit,
	})

	// Create query builder with logger
	queryBuilder := igdb.NewIGDBQueryBuilder(a.logger).
        Search(query).
        Fields(igdb.DefaultGameFields...).
        Where(igdb.GameTypeFilter).
        Limit(limit)

	// Execute the query
	responses, err := a.client.ExecuteQuery(ctx, queryBuilder)
	if err != nil {
		a.logger.Error("Failed to execute IGDB query", map[string]any{
			"error": err,
			"query": query,
		})
		return nil, err
	}

	// Convert responses to games
	games := a.convertResponsesToGames(responses)

	a.logger.Info("Successfully retrieved games from IGDB", map[string]any{
		"count": len(games),
	})

	return games, nil
}


// UpdateToken updates the authentication token used by the IGDB client.
// This is needed because IGDB tokens expire and need to be refreshed.
func (a *IGDBAdapter) UpdateToken(token string) error {
	if a.client == nil {
		return fmt.Errorf("IGDB client is nil")
	}
	a.client.UpdateToken(token)
	return nil
}

// Helper fn - convertResponsesToGames converts IGDB API responses into our application's game model.
// It handles the conversion of:
//   - Basic game information
//   - Nested data like platforms, genres, and themes
//   - Game type information
func (a *IGDBAdapter) convertResponsesToGames(responses []*types.IGDBResponse) []*models.Game {
	games := make([]*models.Game, 0, len(responses))

	for i := 0; i < len(responses); i++ {
			resp := responses[i]

			// Guard against missing/malformed payload - skip if game_type, coverURL or platforms are missing
			if resp.GameType.Type == "" ||
					resp.Cover.URL == "" ||
					len(resp.Platforms) == 0 {
					continue
			}

			// Look up game type ID to get display text and normalized text
			gameTypeResponseField := getGameType(resp.GameType.ID)

			// Build full GameType value for DB
			gameTypeforDB := types.GameType{
					ID:             resp.GameType.ID,
					Type:           resp.GameType.Type,
					DisplayText:    gameTypeResponseField.DisplayText,
					NormalizedText: gameTypeResponseField.NormalizedText,
			}

			// Build GameTypeResponse for frontend JSON
			gameTypeforResponse := types.GameTypeResponse{
					DisplayText:    gameTypeforDB.DisplayText,
					NormalizedText: gameTypeforDB.NormalizedText,
			}

			// Pre-allocate slices with known capacity
			platformNames := make([]string, len(resp.Platforms))
			genreNames := make([]string, len(resp.Genres))
			themeNames := make([]string, len(resp.Themes))

			platforms := make([]models.PlatformInfo, len(resp.Platforms))
			// Convert platforms using index-based loop
			for j := 0; j < len(resp.Platforms); j++ {
					platforms[j] = models.PlatformInfo{
							ID:   resp.Platforms[j].ID,
							Name: resp.Platforms[j].Name,
					}
					platformNames[j] = resp.Platforms[j].Name
			}

			// Convert genres using index-based loop
			for j := 0; j < len(resp.Genres); j++ {
					genreNames[j] = resp.Genres[j].Name
			}

			// Convert themes using index-based loop
			for j := 0; j < len(resp.Themes); j++ {
					themeNames[j] = resp.Themes[j].Name
			}

			gameResponse := &models.Game{
					ID:               resp.ID,
					Name:             resp.Name,
					Summary:          resp.Summary,
					CoverURL:         resp.Cover.URL,
					FirstReleaseDate: resp.FirstReleaseDate,
					Rating:           resp.Rating,
					GameType:         gameTypeforDB,
					GameTypeResponse: gameTypeforResponse,
					Platforms:        platforms,
					PlatformNames:    platformNames,
					GenreNames:       genreNames,
					ThemeNames:       themeNames,
			}

			games = append(games, gameResponse)
	}

	return games
}

// Helper fn - getGameType looks up the game type in the types.GameTypes map
// Returns the corresponding types.GameType obj containing display text
// and normalized text for the frontend response.
// If the game type is not found, it returns a zero value.
func getGameType(id int64) types.GameType {
	if gameType, exists := types.GameTypes[id]; exists {
		return gameType
	}
	return types.GameType{} // Return zero value if not found
}
