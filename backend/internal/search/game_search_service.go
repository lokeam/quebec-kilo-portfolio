package search

import (
	"context"
	"time"

	"github.com/lokeam/qko-beta/config"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	security "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
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
}

type SearchService interface {
	Search(ctx context.Context, req searchdef.SearchRequest) (*searchdef.SearchResult, error)
}

// NewGameSearchService wires up the GameSearchService with its dependencies.
func NewGameSearchService(appContext *appcontext.AppContext) (*GameSearchService, error) {

	// Create an adapter to search IGDB
	adapter, err := NewIGDBAdapter(appContext.Config, appContext.Logger)
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
	cacheWrapper, err := NewIGDBCacheWrapper(
		appContext.RedisClient,
		appContext.Config.Redis.RedisTTL,
		appContext.Config.Redis.RedisTimeout,
		appContext.Logger,
	)
	if err != nil {
		return nil, err
	}

	return &GameSearchService{
		adapter:      adapter,
		config:       appContext.Config,
		logger:       appContext.Logger,
		validator:    validator,
		sanitizer:    sanitizer,
		cacheWrapper: cacheWrapper,
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
	// 4. If cache miss, fetch fresh data using the adapter.
	games, err := s.adapter.SearchGames(ctx, req.Query, req.Limit)
	if err != nil {
		s.logger.Error("IGDB adapter error", map[string]any{"error": err})
		return nil, err
	}

	// Convert adapter results.
	convertedGames := make([]searchdef.Game, len(games))
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
