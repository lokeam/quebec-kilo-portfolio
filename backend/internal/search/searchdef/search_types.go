package searchdef

import (
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

const (
	DefaultPageSize = 30
	MaxPageSize     = 50
)

// PROCESSED DATA FOR FRONTEND
type SearchRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

// SearchQuery represents the search parameters. This type should include any fields
// needed to generate a unique cache key.
// For this example, we assume a simple query string.
type SearchQuery struct {
	Query string
	Limit int
}

func (sq SearchQuery) ToCacheKey() string {
	return fmt.Sprintf("search:%s:%d", sq.Query, sq.Limit)
}

// Game represents a title from IGDB
// type Game struct {
// 	ID               int64   `json:"id"`
// 	Name             string  `json:"name"`
// 	Summary          string  `json:"summary,omitempty"`
// 	Cover            string  `json:"cover,omitempty"`
// 	FirstReleaseDate int64   `json:"first_release_date,omitempty"`
// 	Rating           float64 `json:"rating,omitempty"`
// }
// SearchResponse represents the overall search response from IGDB.
type SearchResponse struct {
	Games []models.Game `json:"games"`
	Total int    `json:"total"`
}

// Search Meta contains info about the search request
type SearchMeta struct {
	// Pagination info
	Total          int `json:"total"`            // Total number of search results available
	CurrentPage    int `json:"current_page"`     // Current page number
	ResultsPerPage int `json:"results_per_page"` // Number of results per page

	// Request info
	Query     string `json:"query"`       // Original search query
	TimeTaken int    `json:"time_taken"`  // How long the search took

	// Pagination (required for frontend)
	HasNextPage     bool `json:"has_next_page"`
	HasPreviousPage bool `json:"has_previous_page"`
	TotalPages      int  `json:"total_pages"`

	// Request Tracking (Required for monitoring + debugging)
	RequestID    string `json:"request_id"`    // Trace ID for request
	TimestampUTC string `json:"timestamp_utc"` // When we performed the search

	// Cache info (required for monitoring + debugging)
	CacheHit bool          `json:"cache_hit"` // Was result from L2 cache?
	CacheTTL time.Duration `json:"cache_ttl"` // How long result stays cached
}

// Search Result wraps the games data + metadata
type SearchResult struct {
	Games []models.Game     `json:"games"` // The actual game results
	Meta  SearchMeta `json:"meta"`  // Metadata about the search
	Error error      `json:"-"`     // Any error that occurred (not serialized to JSON)
}

// Constructor for NewSearchResult, creates empty result, uses builder pattern to create results step by step
func NewSearchResult() *SearchResult {
	return &SearchResult{
		Games: make([]models.Game, 0),
		Meta:  SearchMeta{
			ResultsPerPage: 30,
		},
	}
}

// WithError adds an error the result if needed
func (r *SearchResult) WithError(err error) *SearchResult {
	r.Error = err
	return r
}

// WithMeta adds metadata to the result if needed
func (r *SearchResult) WithMeta(meta SearchMeta) *SearchResult {
	r.Meta = meta
	return r
}

// WithGames add games to the result if needed
func (r * SearchResult) WithGames(games []models.Game, currentPage int, limit int) *SearchResult {
	r.Games = games
	r.Meta.Total = len(games)

	// Calculate pagination
	r.calculatePagination(currentPage, limit)
	return r
}


// Calculation helper to centralize pagination logic
func (r *SearchResult) calculatePagination(currentPage int, limit int) {
	r.Meta.CurrentPage = currentPage
	r.Meta.ResultsPerPage = limit

	// Calculate total pages
	if (r.Meta.Total > 0) {
		r.Meta.TotalPages = (r.Meta.Total + limit - 1) / limit
	}

	// Calculate page indicators
	r.Meta.HasPreviousPage = currentPage > 1
	r.Meta.HasNextPage = currentPage < r.Meta.TotalPages
}
