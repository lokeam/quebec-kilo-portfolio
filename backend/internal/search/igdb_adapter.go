package search

import (
	"context"
	"fmt"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/igdb"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
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
	IGDBGameQueryTemplate = `fields id,name,summary,first_release_date,rating,cover,genres,themes,game_modes,platforms; search "%s"; limit %d;`
)

// IGDBAdapter wraps the IGDB client and adds circuit breaker protection and logging.
type IGDBAdapter struct {
	client  *igdb.IGDBClient             // Underlying IGDB client.
	breaker *gobreaker.CircuitBreaker // Circuit breaker to protect IGDB calls.
	logger  interfaces.Logger         // Logger interface.
}

// NewIGDBAdapter creates a new IGDBAdapter instance.
// It retrieves the Twitch token required by IGDB, creates an IGDB client,
// and sets up the circuit breaker with appropriate settings.
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

// SearchGames wraps the IGDB client's SearchGames method using the circuit breaker.
// It logs the request and response, and returns an error if the circuit is open or the call fails.
func (a *IGDBAdapter) SearchGames(ctx context.Context, query string, limit int) ([]*models.Game, error) {
	a.logger.Info("IGDB Adapter - SearchGames called", map[string]any{
		"query": query,
		"limit": limit,
	})

	// Log the query being sent to IGDB
	igdbQuery := fmt.Sprintf(IGDBGameQueryTemplate, query, limit)
	a.logger.Debug("IGDB Query", map[string]any{
		"query": igdbQuery,
	})

	// Call the IGDB API
	gameDetails, err := a.client.SearchGames(igdbQuery)
	if err != nil {
		return nil, err
	}

	// Convert the IGDB Backend Response to what is expected by the Frontend
	var games []*models.Game
	for _, details := range gameDetails {
		games = append(games, &models.Game{
			ID:                    details.ID,
			Name:                  details.Name,
			Summary:               details.Summary,
			CoverURL:              details.CoverURL,
			FirstReleaseDate:      details.FirstReleaseDate,
			Rating:                details.Rating,
			PlatformNames:         details.PlatformNames,
			GenreNames:            details.GenreNames,
			ThemeNames:            details.ThemeNames,
		})
	}

	return games, nil
}

// UpdateToken updates the token in the underlying IGDB client
func (a *IGDBAdapter) UpdateToken(token string) error {
	if a.client == nil {
		return fmt.Errorf("IGDB client is nil")
	}
	a.client.UpdateToken(token)
	return nil
}
