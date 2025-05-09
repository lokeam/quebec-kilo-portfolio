package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *IGDBClient) makeRequest(endpoint string, query string, result interface{}) error {
    if c == nil {
        return fmt.Errorf("IGDBClient is nil")
    }

    c.logger.Info("igdb client - makeRequest called with endpoint: %s and query: %s",
        map[string]any{"endpoint": endpoint, "query": query})

    url := fmt.Sprintf("%s/%s", IGDB_API_URL, endpoint)

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

func buildIDQuery(ids []int64, fields string) string {
    var strIDs []string
    for _, id := range ids {
        strIDs = append(strIDs, fmt.Sprintf("%d", id))
    }
    return fmt.Sprintf("fields %s; where id = (%s);", fields, strings.Join(strIDs, ","))
}

// UpdateToken updates the client's authentication token
func (c *IGDBClient) UpdateToken(token string) {
    c.token = token
}
