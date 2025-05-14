package igdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lokeam/qko-beta/internal/types"
)

// Helper fn - makeRequest handles the actual HTTP communication with the IGDB API.
// It:
//   1. Creates and sends an HTTP POST request to IGDB
//   2. Adds required headers (Client-ID and Authorization)
//   3. Handles the response and any errors
//   4. Unmarshals the JSON response into the provided result
//
// The function is used by all IGDB API calls to ensure consistent:
//   - Error handling
//   - Logging
//   - Authentication
//   - Response processing
func (c *IGDBClient) makeRequest(endpoint string, query string, result interface{}) error {
    if c == nil {
        return fmt.Errorf("IGDBClient is nil")
    }

    c.logger.Info("igdb client - makeRequest called with endpoint: %s and query: %s",
        map[string]any{"endpoint": endpoint, "query": query})

    // Use the IGDB API URL in igdb_constants file to make requests
    url := fmt.Sprintf("%s/%s", BASE_IGDB_API_URL, endpoint)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(query)))
    if err != nil {
        c.logger.Error("igdb client - makeRequest - failed to create request: %w", map[string]any{"error": err})
        return fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Add("Client-ID", c.clientID)
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
    req.Header.Add("Accept", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        c.logger.Error("igdb client - makeRequest - failed to send request: %w", map[string]any{"error": err})
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    // Context cancellation check
    select {
    case <-req.Context().Done():
        return req.Context().Err()
    default:
    }

    c.logger.Info("igdb client - makeRequest - response: ", map[string]any{"resp": resp})

    // Read the response body ONCE then store it
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        c.logger.Error("igdb client - makeRequest - failed to read response body: %w", map[string]any{"error": err})
        return fmt.Errorf("failed to read response body: %w", err)
    }

    c.logger.Info("igdb client - makeRequest - response received", map[string]any{
        "status": resp.StatusCode,
        "body":   string(body),
    })

    // Use the stored body for status code check
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("IGDB API error (status %d): %s", resp.StatusCode, string(body))
    }

    // Use the same stored body for decoding
    if err := json.Unmarshal(body, result); err != nil {
        c.logger.Error("igdb client - makeRequest - failed to decode response: %w", map[string]any{"error": err})
        return fmt.Errorf("failed to decode response: %w", err)
    }

    return nil
}

// ExecuteQuery executes a query against the IGDB API and returns the results.
// It:
//   1. Takes a QueryBuilder to construct the query
//   2. Builds the query string
//   3. Makes the request to IGDB
//   4. Returns the results as a slice of pointers to IGDBResponse
//
// The function uses makeRequest to handle the actual HTTP communication
// and ensures proper error handling and logging.
func (c *IGDBClient) ExecuteQuery(ctx context.Context, queryBuilder *QueryBuilder) ([]*types.IGDBResponse, error) {
    if c == nil {
        return nil, fmt.Errorf("IGDBClient is nil")
    }

    if queryBuilder == nil {
        return nil, fmt.Errorf("query is nil")
    }

    // Build the query string
    queryStr, err := queryBuilder.Build()
    if err != nil {
        return nil, fmt.Errorf("failed to build query: %w", err)
    }

   // Make the request with a slice of pointers, allows unmarshaler to create pointers directly
    var responses []*types.IGDBResponse
    if err := c.makeRequest("games", queryStr, &responses); err != nil {
        return nil, fmt.Errorf("failed to execute query: %w", err)
    }

    return responses, nil
}


// UpdateToken updates the client's authentication token.
// This is needed because IGDB tokens expire and need to be refreshed
// to maintain API access.
func (c *IGDBClient) UpdateToken(token string) {
    c.token = token
}
