package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Henry-Sarabia/igdb"
	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/sony/gobreaker"
)

// IGDBAdapter retrieves data directly from the IGDB endpoint.
// It uses a circuit breaker to prevent cascading failures when the IGDB API is slow or failing.
type IGDBAdapter struct {
	apiURL     string                     // Base URL for IGDB API.
	apiKey     string                     // API key for authentication.
	httpClient *http.Client               // HTTP client configured with timeouts.
	breaker    *gobreaker.CircuitBreaker  // Circuit breaker wrapping IGDB calls.
}

// NewIGDBAdapter creates a new IGDBAdapter instance with circuit breaker protection.
func NewIGDBAdapter(config *config.Config, logger interfaces.Logger) (*IGDBAdapter, error) {

	apiURL := config.IGDB.ClientID
	apiKey := config.IGDB.ClientSecret

	logger.Info("Creating IGDBAdapter", map[string]any{
		"apiURL": apiURL,
	})

	// Configure the circuit breaker settings.
	cbSettings := gobreaker.Settings{
		Name: "IGDBAPI",
		Timeout: 5 * time.Second,  // Duration after which the circuit breaker resets.
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Open the circuit after 3 consecutive failures.
			return counts.ConsecutiveFailures >= 3
		},

		// Interval: time window to reset counts if the failure rate improves.
		Interval: 1 * time.Minute,
	}
	return &IGDBAdapter{
		apiURL: apiURL,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Set an appropriate timeout for IGDB calls.
		},
		breaker: gobreaker.NewCircuitBreaker(cbSettings),
	}, nil
}

// SearchGames retrieves games from the IGDB API using the circuit breaker.
// It returns a slice of pointers to igdb.Game or an error.
func (a *IGDBAdapter) SearchGames(ctx context.Context, query string, limit int) ([]*igdb.Game, error) {
	// Wrap the API call in a circuit breaker.
	result, err := a.breaker.Execute(func() (interface{}, error) {
		// NOTE: Build the request URL. In production, use url.Values for query parameters.
		urlStr := a.apiURL + "/games?search=" + url.QueryEscape(query) + "&limit=" + fmt.Sprintf("%d", limit)
		req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
		if err != nil {
			return nil, err
		}
		// Include the API key in the request (header or query parameter, as required).
		req.Header.Set("Authorization", "Bearer "+a.apiKey)

		resp, err := a.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("IGDB API returned non-200 status: " + resp.Status)
		}

		var games []*igdb.Game
		// Decode the JSON response into games.
		if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
			return nil, err
		}
		return games, nil
	})
	if err != nil {
		return nil, err
	}
	games, ok := result.([]*igdb.Game)
	if !ok {
		return nil, errors.New("unexpected result type from IGDB API")
	}
	return games, nil
}
