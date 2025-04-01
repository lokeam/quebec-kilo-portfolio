package search

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/infrastructure/cache"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
	"github.com/lokeam/qko-beta/internal/shared/worker"
)

// GameSearchService processes search requests by validating and sanitizing the query,
// then delegating the retrieval to the IGDBCacheWrapper.
type GameSearchService struct {
	adapter         interfaces.IGDBAdapter       // Retrieves data directly from IGDB.
	config          *config.Config
	logger          interfaces.Logger
	validator       interfaces.SearchValidator
	sanitizer       interfaces.Sanitizer
	cacheWrapper    interfaces.IGDBCacheWrapper  // Only handles caching.
	appContext      *appcontext.AppContext
}

type SearchService interface {
	Search(ctx context.Context, req searchdef.SearchRequest) (*searchdef.SearchResult, error)
}

// NewGameSearchService wires up the GameSearchService with its dependencies.
func NewGameSearchService(appContext *appcontext.AppContext) (*GameSearchService, error) {

	// Create an adapter to search IGDB
	appContext.Logger.Info("Game Search Service - Creating IGDBAdapter", map[string]any{
		"appContext": appContext,
	})
	adapter, err := NewIGDBAdapter(appContext)
	if err != nil {
		return nil, err
	}

	// Create sanitizer to feed into validator
	sanitizer, err := security.NewSanitizer()
	if err != nil {
		return nil, err
	}

	// Create validator to validate search queries
	validator, err := NewSearchValidator(sanitizer)
	if err != nil {
		return nil, err
	}

	// Create cache wrapper to handle Redis caching
	cacheWrapper, err := cache.NewCacheWrapper(
		appContext.RedisClient,
		appContext.Config.Redis.RedisTTL,
		appContext.Config.Redis.RedisTimeout,
		appContext.Logger,
	)
	if err != nil {
		return nil, err
	}

	igdbCacheAdapter, err := NewIGDBCacheAdapter(cacheWrapper)
	if err != nil {
		return nil, err
	}

	return &GameSearchService{
		adapter:      adapter,
		config:       appContext.Config,
		logger:       appContext.Logger,
		validator:    validator,
		sanitizer:    sanitizer,
		cacheWrapper: igdbCacheAdapter,
		appContext:   appContext,
	}, nil
}

// Search processes the search request: it sanitizes and validates the query,
// then first checks the cache. If the cache misses, it fetches fresh data,
// caches it, and returns the result.
func (s *GameSearchService) Search(ctx context.Context, req searchdef.SearchRequest) (*searchdef.SearchResult, error) {
	// 1. Sanitize the query.
	sanitized, err := s.sanitizer.SanitizeSearchQuery(req.Query)
	if err != nil {
		s.logger.Error("Search query sanitization failed", map[string]any{"error": err})
		return nil, err
	}
	req.Query = sanitized

	// 2. Validate the query.
	if err := s.validator.ValidateQuery(searchdef.SearchQuery(req)); err != nil {
		s.logger.Error("Search validation failed", map[string]any{"error": err})
		return nil, err
	}

	// 3. Build a SearchQuery from the request.
	// Note: SearchQuery must provide the ToCacheKey() method.
	sq := searchdef.SearchQuery(req)

	// 3. Attempt to retrieve cached results.
	cachedResult, err := s.cacheWrapper.GetCachedResults(ctx, sq)
	if err == nil && cachedResult != nil {
		s.logger.Debug("Cache hit", map[string]any{"query": req.Query})
		return cachedResult, nil
	}

	s.logger.Debug("Cache miss; performing IGDB search", map[string]any{"query": req.Query, "limit": req.Limit})

	// 4. If cache miss, fetch data using adapter with retry for auth errors
	games, err := s.searchWithTokenRefresh(
		ctx,
		req.Query,
		req.Limit,
	)
	if err != nil {
		return nil, err
	}

	// Convert adapter results.
	convertedGames := make([]models.Game, len(games))
	for i, gamePtr := range games {
		convertedGames[i] = convertIGDBGame(*gamePtr)
	}

	result := &searchdef.SearchResult{Games: convertedGames}
	result.Meta.CacheHit = false
	result.Meta.CacheTTL = s.config.Redis.RedisTTL
	result.Meta.TimestampUTC = time.Now().UTC().Format(time.RFC3339)

	// 5. Cache the fresh data.
	if err := s.cacheWrapper.SetCachedResults(ctx, sq, result); err != nil {
		s.logger.Error("Failed to cache fresh search results", map[string]any{"error": err})
	}

	return result, nil
}

func convertIGDBGame(g models.Game) models.Game {
	return models.Game{
		ID:                  int64(g.ID),
		Name:                g.Name,
		Summary:             g.Summary,
		CoverURL:            g.CoverURL,
		FirstReleaseDate:    int64(g.FirstReleaseDate),
		Rating:              g.Rating,
		PlatformNames:       g.PlatformNames,
		GenreNames:          g.GenreNames,
		ThemeNames:          g.ThemeNames,
	}
}

// Attempts to search IGDB and automatically refreshes the token if a 401 error occurs
func (s *GameSearchService) searchWithTokenRefresh(
	ctx context.Context,
	query string,
	limit int,
) ([]*models.Game, error) {
	// Attempt to search IGDB
	games, err := s.adapter.SearchGames(ctx, query, limit)

	// If we get an error, check if it is related to auth (401)
	if err != nil && IAuthError(err) {
		s.logger.Warn("Received authentication error from IGDB, attempting token refresh", map[string]any{
			"error": err,
		})

		// Attempt to refresh token
		if refreshErr := s.refreshToken(ctx); refreshErr != nil {
			s.logger.Error("Failed to refresh token", map[string]any{
				"error": refreshErr,
			})
			return nil, fmt.Errorf("failed to refresh token: %w", refreshErr)
		}

		s.logger.Info("Token refreshed successfully, retrying search request", nil)

		// Retry the search with the new token
		return s.adapter.SearchGames(ctx, query, limit)
	}

	return games, err
}

// TODO: Move this to error handling package
func IAuthError(err error) bool {
	// Check for specific error types or messages that indicate auth failure
	if err == nil {
			return false
	}

	// Look for 401 in the error message
	return strings.Contains(err.Error(), "401") ||
				 strings.Contains(strings.ToLower(err.Error()), "unauthorized") ||
				 strings.Contains(strings.ToLower(err.Error()), "authentication")
}

// Trigger a token refresh using existing worker jobs
func (s *GameSearchService) refreshToken(ctx context.Context) error {
	// Get config values
	clientID := s.config.IGDB.ClientID
	clientSecret := s.config.IGDB.ClientSecret
	authURL := s.config.IGDB.AuthURL
	redisKey := s.config.IGDB.AccessTokenKey

	// Get the token via retry logic
	token, err := worker.GetTwitchTokenWithRetry(
		ctx,
		clientID,
		clientSecret,
		authURL,
		s.logger,
	)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	// Save the token with refresh logic
	if err := worker.UpdateTwitchTokenJob(
		ctx,
		redisKey,
		s.appContext.RedisClient,
		s.appContext.MemCache,
		*token,
		s.logger,
	); err != nil {
		return fmt.Errorf("failed to save refreshed token: %w", err)
	}

	s.logger.Info("Successfully refreshed and saved Twitch token", nil)

	return nil
}