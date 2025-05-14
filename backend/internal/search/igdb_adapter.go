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
	games := convertResponsesToGames(responses)

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
func convertResponsesToGames(responses []*types.IGDBResponse) []*models.Game {
	games := make([]*models.Game, 0, len(responses))
	for _, resp := range responses {
		gameResponse := &models.Game{
			ID:               resp.ID,
			Name:             resp.Name,
			Summary:          resp.Summary,
			FirstReleaseDate: resp.FirstReleaseDate,
			Rating:           resp.Rating,
			CoverURL:         resp.Cover.URL,
			PlatformNames:    make([]string, len(resp.Platforms)),
			GenreNames:       make([]string, len(resp.Genres)),
			ThemeNames:       make([]string, len(resp.Themes)),
			GameType:         convertGameType(resp.GameType),
		}

		// Convert platform names
		for i, pl := range resp.Platforms {
			gameResponse.PlatformNames[i] = pl.Name
		}

		// Convert genre names
		for i, ge := range resp.Genres {
			gameResponse.GenreNames[i] = ge.Name
		}

		// Convert theme names
		for i, th := range resp.Themes {
			gameResponse.ThemeNames[i] = th.Name
		}

		games = append(games, gameResponse)
	}
	return games
}

// Helper fn - convertGameType converts an IGDB game type to our application's game type.
// This is needed because IGDB's response format doesn't exactly match our model.
func convertGameType(igdbType types.IGDBResponseGameType) types.GameType {
	return types.GameType{
			ID:   int(igdbType.ID),
			Type: igdbType.Type,
	}
}