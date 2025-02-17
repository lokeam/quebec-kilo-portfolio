package interfaces

import (
	"context"

	"github.com/lokeam/qko-beta/internal/search/searchdef"
)

type IGDBCacheWrapper interface {
	GetCachedResults(ctx context.Context, query searchdef.SearchQuery) (*searchdef.SearchResult, error)
	SetCachedResults(ctx context.Context, query searchdef.SearchQuery, results *searchdef.SearchResult) error
}
