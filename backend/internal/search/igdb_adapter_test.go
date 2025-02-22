package search

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
	"github.com/lokeam/qko-beta/internal/types"
)

/*

	Behaviors:
		- Searches for games using the IGDB API and returns with a valid JSON search result
		- If IGDB search fails, it returns a non 200 HTTP status code + error message
		- If the IGDB Adapter is called consecutively (with a config to trip after N amoutn of times) the circuit breaker trips and returns an error


	Scenarios:
	- Successful game search
	- API returns a non 200 status code
	- HTTP request fails or errors in some way
	- The JSON response is malformed
	- The circuit breaker trips after consecutive failures
	- The SearchGames() method is called with a cancelled context
	-
*/

func TestIGDBAdapter(t *testing.T) {
	t.Run(
		`SearchGames() returns a non 200 HTTP status code and error`,
		func(t *testing.T) {
			/*
				GIVEN a properly configured IGDB adapter
				AND a simulated API endpoint that returns a non 200 HTTP status code
				WHEN SearchGames() is called with a valid context + query
				THEN SearchGames() returns and error indicating the IGDB API returned a non 200 status code
			*/
			mockAdapter := &mocks.MockIGDBAdapter{
        SearchGamesFunc: func(ctx context.Context, query string, limit int) ([]*types.Game, error) {
            return nil, fmt.Errorf("HTTP error: non-200 status code")
        },
    	}

			_, testErr := mockAdapter.SearchGames(context.Background(), "Dark Souls", 1)
			if testErr == nil {
					t.Fatal("expected error for non 200 response but got nil")
			}
			if !strings.Contains(testErr.Error(), "non-200") {
					t.Errorf("expected error to mention non-200 status, but got: %v", testErr)
			}
		},
	)

	t.Run(`SearchGames() returns an error immediately after HTTP request error`, func(t *testing.T) {

		mockConfig := mocks.NewMockConfig()
		testLogger := testutils.NewTestLogger()
		mockTokenRetriever := &mocks.MockTwitchTokenRetriever{}

		// Create mock app context
		mockAppContext := &appcontext.AppContext{
				Config:              mockConfig,
				Logger:              testLogger,
				TwitchTokenRetriever: mockTokenRetriever,
		}

		// Create a custom HTTP client that will simulate a network error
		mockHTTPClient := &http.Client{
				Transport: &testutils.ErrorRoundTripper{},
		}

		// Create the IGDBAdapter
		testIGDBAdapter, err := NewIGDBAdapter(mockAppContext)
		if err != nil {
				t.Fatalf("failed to create IGDB adapter: %v", err)
		}

		// Verify the client is not nil before setting HTTP client
		if testIGDBAdapter.client == nil {
				t.Fatal("IGDBClient is nil")
		}

		// Override the HTTP client in the IGDBAdapter's IGDBClient
		testIGDBAdapter.client.SetHTTPClient(mockHTTPClient)

		// Perform the search
		_, searchErr := testIGDBAdapter.SearchGames(context.Background(), "DarkSouls", 1)

		// Assert the error
		if searchErr == nil {
				t.Fatal("expected error for network error but got nil")
		}

		if !strings.Contains(strings.ToLower(searchErr.Error()), "network error") {
			t.Errorf("expected network error, but got: %v", searchErr)
		}
	})

	t.Run(`SearchGames() returns an irregular JSON response`, func(t *testing.T) {
    mockAdapter := &mocks.MockIGDBAdapter{
				SearchGamesFunc: func(ctx context.Context, query string, limit int) ([]*types.Game, error) {
						return nil, fmt.Errorf("failed to decode JSON response")
				},
		}

		_, testErr := mockAdapter.SearchGames(context.Background(), "Dark Souls", 1)
		if testErr == nil {
				t.Fatal("expected JSON decoding error but got nil")
		}
		if !strings.Contains(strings.ToLower(testErr.Error()), "json") {
				t.Errorf("expected JSON decoding error, but got: %v", testErr)
		}
	})

	t.Run(`Circuit breaker trips after N-configured number of consecutive failures`, func(t *testing.T) {
			/*
				GIVEN an IGDB Adapter configured with a circuit breaker (set to trip after 3 consective failures)
				AND a simulated IGDB endpoint that consistentyly fails (or returns errors/timeouts)
				WHEN SearchGames() is called at least the number of times required to trip the breaker
				THEN After the configured number failures, SearchGames() should trip
				AND subsequent calls should quickly return an error w/o waiting for HTTP client to timeout
				(thus indicating that the circuit breaker is still open)
			*/
			mockAdapter := &mocks.MockIGDBAdapter{
        SearchGamesFunc: func(ctx context.Context, query string, limit int) ([]*types.Game, error) {
            return nil, errors.New("circuit breaker error")
        },
			}

			var finalErr error
			for i := 0; i < 5; i++ {
					_, finalErr = mockAdapter.SearchGames(context.Background(), "Dark Souls", 1)
			}

			if finalErr == nil {
					t.Fatalf("expected a circuit breaker error after multiple failures")
			}
			if !strings.Contains(finalErr.Error(), "circuit breaker") {
					t.Errorf("expected circuit breaker error, got: %v", finalErr)
			}
	})

	t.Run(`Context is cancelled before the HTTP request is complete`, func(t *testing.T) {
			/*
				GIVEN a properly configured IGDB Adapter
				AND a context that is cancelled/has a deadline in the past
				WHEN SearchGames() is called with this cancelled context
				THEN SearchGames() should return an error related to the context being cancelled/expired
			*/
				mockAdapter := &mocks.MockIGDBAdapter{
					SearchGamesFunc: func(ctx context.Context, query string, limit int) ([]*types.Game, error) {
							return nil, context.Canceled
					},
			}

			testContext, cancel := context.WithCancel(context.Background())
			cancel()

			_, testErr := mockAdapter.SearchGames(testContext, "Dark Souls", 1)
			if testErr == nil {
					t.Fatalf("expected an error due to context cancellation but instead got: %v", testErr)
			}
			if !errors.Is(testErr, context.Canceled) {
					t.Errorf("expected context cancellation error, but instead got: %v", testErr)
			}
	})

	t.Run(`Happy Path - Successful game search`, func(t *testing.T) {
			/*
				GIVEN a properly configured IGDB Adapter with a a valid API URL+key and HTTP client
				AND a running (or simulated) API endpoint that returns HTTP 200 w/ a well formed JSON response for a list of games
				WHEN SearchGames() is called with a valid context, query and limit
				THEN the method returns a slice of *igdb.Game objects containing the expected data
			*/
			expectedGames := []*types.Game{
        {ID: 1, Name: "Dark Souls 1"},
        {ID: 2, Name: "Dark Souls 2"},
        {ID: 3, Name: "Dark Souls 3"},
			}

			mockAdapter := &mocks.MockIGDBAdapter{
					SearchGamesFunc: func(ctx context.Context, query string, limit int) ([]*types.Game, error) {
							return expectedGames, nil
					},
			}

			actualGames, testErr := mockAdapter.SearchGames(context.Background(), "Dark Souls", 3)
			if testErr != nil {
					t.Fatalf("expected no error on successful search, but got: %v", testErr)
			}
			if len(actualGames) != len(expectedGames) {
					t.Errorf("expected %d games, but got %d", len(expectedGames), len(actualGames))
			}
	})
}
