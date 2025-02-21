package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	memorycache "github.com/lokeam/qko-beta/internal/infrastructure/cache/memorycache"
	cache "github.com/lokeam/qko-beta/internal/infrastructure/cache/rueidis"
	"github.com/lokeam/qko-beta/internal/interfaces"
)

// TokenResponse represents the structure of a Twitch token response.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// ErrInvalidResponse indicates a non-200 response from the Twitch API
type ErrInvalidResponse struct {
	StatusCode int
	Body       string
}

func (e *ErrInvalidResponse) Error() string {
	return fmt.Sprintf("received non-200 status code %d: %s", e.StatusCode, e.Body)
}

// RefreshToken sends a POST request to Twitch's token endpoint.
// It returns a TokenResponse on success or an error otherwise.
var RefreshToken = func(
	ctx context.Context,
	clientID string,
	clientSecret string,
	authURL string,
	logger interfaces.Logger,
) (*TokenResponse, error) {

	// Populate correct form data for request
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request with context
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		authURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		logger.Error("Failed to create request", map[string]any{"error": err})
		return nil, err
	}

	// Set headers (break this up into http utils)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	logger.Info("Sending request to Twitch token endpoint", map[string]any{
		"url": authURL,
		"client_id": clientID,
	})

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("HTTP request error", map[string]any{"error": err})
		return nil, err
	}
	defer resp.Body.Close()



	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		logger.Error("Non-200 response from Twitch", map[string]any{
			"status": resp.StatusCode,
			"body":   string(body),
		})
		return nil, &ErrInvalidResponse{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}

	// Parse response
	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		logger.Error("Failed to parse response", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Validate response fields
	if tokenResponse.AccessToken == "" {
		return nil, fmt.Errorf("received empty access token")
	}

	logger.Info("Successfully retrieved Twitch token", map[string]any{
		"access_token": tokenResponse.AccessToken,
		"expires_in": tokenResponse.ExpiresIn,
		"token_type": tokenResponse.TokenType,
	})

	return &tokenResponse, nil
}

// TwitchTokenRetriever handles retrieving and caching a Twitch token.
type TwitchTokenRetriever struct {
	memCache    *memorycache.MemoryCache
	redisClient *cache.RueidisClient
	twitchClientID  string
	twitchSecret    string
	twitchAuthURL   string
}



// NewTwitchTokenRetriever creates a new TwitchTokenRetriever.
func NewTwitchTokenRetriever(
	memCache *memorycache.MemoryCache,
	redisClient *cache.RueidisClient,
	twitchClientID,
	twitchSecret,
	twitchAuthURL string,
) *TwitchTokenRetriever {
	return &TwitchTokenRetriever{
		memCache:    memCache,
		redisClient: redisClient,
		twitchClientID:  twitchClientID,
		twitchSecret:    twitchSecret,
		twitchAuthURL:   twitchAuthURL,
	}
}

// GetToken checks cache layers (memory first, then Redis) for an existing token.
// If not found, it refreshes the token using RefreshToken and caches the result.
func (r *TwitchTokenRetriever) GetToken(
	ctx context.Context,
	clientID,
	clientSecret,
	authURL string,
	logger interfaces.Logger,
) (string, error) {
	const cacheKey = "twitch:access_token"

	// Attempt to get token from the memory cache.
	token, err := r.memCache.Get(ctx, cacheKey)
	if err != nil {
		logger.Error("Failed to get token from memory cache", map[string]any{"error": err})
		return "", err
	}
	if token != "" {
		logger.Debug("Token found in memory cache", map[string]any{
			"token": token,
		})
		return token, nil
	}

	// Attempt to get token from Redis.
	token, err = r.redisClient.Get(ctx, cacheKey)
	if err == nil && token != "" {
		logger.Debug("Token found in Redis; caching in memcache", map[string]any{
			"token": token,
		})

		// Cache it in memory for a faster lookup next time.
		r.memCache.Set(ctx, cacheKey, token, 10*time.Minute)
		return token, nil
	}



	// Otherwise, refresh a new token.
	// tokenResponse, err := RefreshToken(ctx, clientID, clientSecret, authURL, logger)
	// if err != nil {
	// 	return "", err
	// }
	logger.Info("Token not found in cache; refreshing Twitch token", nil)
	tokenResponse, err := RefreshToken(
		ctx,
		r.twitchClientID,
		r.twitchSecret,
		r.twitchAuthURL,
		logger,
	)
	if err != nil {
		logger.Error("Failed to refresh token", map[string]any{"error": err})
		return "", err
	}

	expiration := time.Duration(tokenResponse.ExpiresIn) * time.Second
	logger.Debug("Twitch token refreshed", map[string]any{
		"access_token": tokenResponse.AccessToken,
		"expires_in": tokenResponse.ExpiresIn,
		"token_type": tokenResponse.TokenType,
	})


	// Cache new token in both memory and redis.
	r.memCache.Set(ctx, cacheKey, tokenResponse.AccessToken, expiration)
	_ = r.redisClient.Set(ctx, cacheKey, tokenResponse.AccessToken, expiration)
	logger.Info("Token cached in memory and Redis", nil)

	return tokenResponse.AccessToken, nil
}
