package search

import "fmt"

// SearchQuery represents the parameters for a search request.
type SearchQuery struct {
	Query string // The search term.
	Limit int    // Maximum number of results.
}

// ToCacheKey returns a unique key for caching based on the query and limit.
func (sq SearchQuery) ToCacheKey() string {
	return fmt.Sprintf("search:%s:%d", sq.Query, sq.Limit)
}