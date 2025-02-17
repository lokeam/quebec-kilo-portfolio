package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Henry-Sarabia/igdb"
	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/sony/gobreaker"
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
	testLogger := testutils.NewTestLogger()

	t.Run(
		`SearchGames() returns a non 200 HTTP status code and error`,
		func(t *testing.T) {
			/*
				GIVEN a properly configured IGDB adapter
				AND a simulated API endpoint that returns a non 200 HTTP status code
				WHEN SearchGames() is called with a valid context + query
				THEN SearchGames() returns and error indicating the IGDB API returned a non 200 status code
			*/
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			})
			testServer := httptest.NewServer(testHandler)
			defer testServer.Close()

			testConfig := &config.Config{
				IGDB: &config.IGDBConfig{
						ClientID:       testServer.URL,
						ClientSecret:   "dummySecret",
						AuthURL:        "dummyAuthURL",
						BaseURL:        "dummyBaseURL",
						TokenTTL:       24 * time.Hour,
						// Simulate error by leaving out AccessTokenKey:
						AccessTokenKey: "",
				},
			}

			testIGDBAdapter, testErr := NewIGDBAdapter(testConfig, testLogger)
			if testErr != nil {
				t.Fatalf("Failed to create IGDB Adapter: %v", testErr)
			}
			testIGDBAdapter.httpClient = testServer.Client()

			testContext := context.Background()
			_, testErr = testIGDBAdapter.SearchGames(testContext, "Dark Souls", 1)
			if testErr == nil {
				t.Fatal("expected error for non 200 response but got nil")
			}
			if !strings.Contains(testErr.Error(), "non-200") {
				t.Errorf("expected error to mention non-200 status, but got: %v", testErr)
			}
		},
	)

	t.Run(
		`SearchGames() returns an error immediately after HTTP request error`,
		func(t *testing.T) {
			/*
				GIVEN the IGDB Adapter + HTTP client is configured to simulate a network error
				WHEN SearchGames() is called with a valid context + query
				THEN SearchGames() immediately returns an error and doesn't proceed to JSON encoding
			*/
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			testServer := httptest.NewServer(handler)
			defer testServer.Close()

			testConfig := &config.Config{
				IGDB: &config.IGDBConfig{
					ClientID:       testServer.URL,
					ClientSecret:   "dummySecret",
					AuthURL:        "dummyAuthURL",
					BaseURL:        "dummyBaseURL",
					TokenTTL:       24 * time.Hour,
					AccessTokenKey: "abc123",
				},
			}

			testIGDBAdapter, testErr := NewIGDBAdapter(testConfig, testLogger)
			if testErr != nil {
				t.Fatalf("failed to create IGDB adapter: %v", testErr)
			}
			testIGDBAdapter.httpClient = &http.Client{
				Transport: &testutils.ErrorRoundTripper{},
			}

			testContext := context.Background()
			_, testErr = testIGDBAdapter.SearchGames(testContext, "DarkSouls", 1)
			if testErr == nil {
				t.Fatal("expected error for network error but got nil")
			}
			if !strings.Contains(testErr.Error(), "network error") {
				t.Errorf("expected error to mention simulated network error, but got: %v", testErr)
			}
		},
	)

	t.Run(
		`SearchGames() returns an irregular JSON response `,
		func(t *testing.T) {
			/*
				GIVEN a properly configured IGDB Adapter
				AND a simulated IGDB endpoint that returns HTTP 200 but with malformed/invalid JSON
				WHEN SearchGames() is called with valid context + query
				THEN SearchGames() returns a JSON decoding error
			*/
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "invalid JSON")
		})
			testServer := httptest.NewServer(testHandler)
			defer testServer.Close()

			testConfig := &config.Config{
				IGDB: &config.IGDBConfig{
					ClientID:       testServer.URL,
					ClientSecret:   "dummySecret",
					AuthURL:        "dummyAuthURL",
					BaseURL:        "dummyBaseURL",
					TokenTTL:       24 * time.Hour,
					AccessTokenKey: "abc123",
				},
			}
			testIGDBAdapter, testErr := NewIGDBAdapter(testConfig, testLogger)
			if testErr != nil {
				t.Fatalf("failed to create IGDB adapter: %v", testErr)
			}
			testIGDBAdapter.httpClient = testServer.Client()

			testContext := context.Background()
			_, testErr = testIGDBAdapter.SearchGames(testContext, "Dark Souls", 1)
			if testErr == nil {
				t.Fatal("expected JSON decoding error but got nil")
			}
			if !strings.Contains(testErr.Error(), "invalid character") {
				t.Errorf("expected a JSON decoding error, but instead got: %v", testErr)
			}
		},
	)

	t.Run(
		`Circuit breaker trips after N-configured number of consecutive failures`,
		func(t *testing.T) {
			/*
				GIVEN an IGDB Adapter configured with a circuit breaker (set to trip after 3 consective failures)
				AND a simulated IGDB endpoint that consistentyly fails (or returns errors/timeouts)
				WHEN SearchGames() is called at least the number of times required to trip the breaker
				THEN After the configured number failures, SearchGames() should trip
				AND subsequent calls should quickly return an error w/o waiting for HTTP client to timeout
				(thus indicating that the circuit breaker is still open)
			*/
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			})
			testServer := httptest.NewServer(testHandler)
			defer testServer.Close()

			testConfig := &config.Config{
				IGDB: &config.IGDBConfig{
					ClientID:       testServer.URL,
					ClientSecret:   "dummySecret",
					AuthURL:        "dummyAuthURL",
					BaseURL:        "dummyBaseURL",
					TokenTTL:       24 * time.Hour,
					AccessTokenKey: "abc123",
				},
			}
			testIGDBAdapter, testErr := NewIGDBAdapter(testConfig, testLogger)
			if testErr != nil {
				t.Fatalf("failed to create IGDB adapter: %v", testErr)
			}
			testIGDBAdapter.httpClient = testServer.Client()

			testContext := context.Background()
			var finalErr error
			for i := 0; i < 5; i++ {
				_, finalErr = testIGDBAdapter.SearchGames(testContext, "Dark Souls", 1)
			}

			if finalErr == nil {
				t.Fatalf("expected a circuit breaker error after multiple failures")
			}
			// Check if circuit breaker is still open
			if !errors.Is(finalErr, gobreaker.ErrOpenState) {
				t.Errorf("expected circuit breaker error, got: %v", finalErr)
		}
		},
	)

	t.Run(
		`Context is cancelled before the HTTP request is complete`,
		func(t *testing.T) {
			/*
				GIVEN a properly configured IGDB Adapter
				AND a context that is cancelled/has a deadline in the past
				WHEN SearchGames() is called with this cancelled context
				THEN SearchGames() should return an error related to the context being cancelled/expired
			*/
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(50 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode([]*igdb.Game{{ID: 1, Name: "Dark Souls"}})
			})
			testServer := httptest.NewServer(testHandler)
			defer testServer.Close()

			testConfig := &config.Config{
				IGDB: &config.IGDBConfig{
					ClientID:       testServer.URL,
					ClientSecret:   "dummySecret",
					AuthURL:        "dummyAuthURL",
					BaseURL:        "dummyBaseURL",
					TokenTTL:       24 * time.Hour,
					AccessTokenKey: "abc123",
				},
			}
			testIGDBAdapter, testErr := NewIGDBAdapter(testConfig, testLogger)
			if testErr != nil {
				t.Fatalf("failed to create IGDB adapter: %v", testErr)
			}
			testIGDBAdapter.httpClient = testServer.Client()

		testContext, cancel := context.WithCancel(context.Background())
			cancel()

			_, testErr = testIGDBAdapter.SearchGames(testContext, "Dark Souls", 1)
			if testErr == nil {
				t.Fatalf("expected an error due to context cancellation but instead got: %v", testErr)
			}

			// Check for a context cancellation error
			if !errors.Is(testErr, context.Canceled) && !strings.Contains(testErr.Error(), "context canceled") {
				t.Errorf("expected context cancellation error, but instead got: %v", testErr)
			}
		},
	)

	t.Run(
		`Happy Path - Successful game search`,
		func(t *testing.T) {
			/*
				GIVEN a properly configured IGDB Adapter with a a valid API URL+key and HTTP client
				AND a running (or simulated) API endpoint that returns HTTP 200 w/ a well formed JSON response for a list of games
				WHEN SearchGames() is called with a valid context, query and limit
				THEN the method returns a slice of *igdb.Game objects containing the expected data
			*/
			expectedGameSearchResult := []*igdb.Game{
				{ ID: 1, Name: "Dark Souls 1" },
				{ ID: 2, Name: "Dark Souls 2" },
				{ ID: 3, Name: "Dark Souls 3" },
			}

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(expectedGameSearchResult)
			})
			testServer := httptest.NewServer(testHandler)
			defer testServer.Close()

			testConfig := &config.Config{
				IGDB: &config.IGDBConfig{
					ClientID:       testServer.URL,
					ClientSecret:   "dummySecret",
					AuthURL:        "dummyAuthURL",
					BaseURL:        "dummyBaseURL",
					TokenTTL:       24 * time.Hour,
					AccessTokenKey: "abc123",
				},
			}
			testIGDBAdapter, testErr := NewIGDBAdapter(testConfig, testLogger)
			if testErr != nil {
				t.Fatalf("failed to create IGDB adapter: %v", testErr)
			}
			testIGDBAdapter.httpClient = testServer.Client()

			testContext := context.Background()
			actualGames, testErr := testIGDBAdapter.SearchGames(testContext, "Dark Souls", 3)
			if testErr != nil {
				t.Fatalf("expected no error on successful search, but got: %v", testErr)
			}
			if len(actualGames) != len(expectedGameSearchResult) {
				t.Errorf("expected %d games, but got %d", len(expectedGameSearchResult), len(actualGames))
			}

			// Check each field
			for i, game := range actualGames {
				if game.ID != expectedGameSearchResult[i].ID || game.Name != expectedGameSearchResult[i].Name {
						t.Errorf("there was a mismatch in game details; we expected %+v, but got %+v", expectedGameSearchResult[i], game)
				}
			}
		},
	)
}
