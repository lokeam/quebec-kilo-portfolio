package search

import (
	"context"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
)

type IGDBCacheAdapter struct {
	cacheWrapper interfaces.CacheWrapper
}

func NewIGDBCacheAdapter(
	cacheWrapper interfaces.CacheWrapper,
) (interfaces.IGDBCacheWrapper, error) {
	return &IGDBCacheAdapter{
		cacheWrapper: cacheWrapper,
	}, nil
}

func (ica *IGDBCacheAdapter) GetCachedResults(
	ctx context.Context,
	query searchdef.SearchQuery,
	) (*searchdef.SearchResult, error) {
		cacheKey := query.ToCacheKey()

		var result searchdef.SearchResult
		cacheHit, err := ica.cacheWrapper.GetCachedResults(ctx, cacheKey, &result)
		if err != nil {
			return nil, err
		}

		if cacheHit {
			result.Meta.CacheHit = true
			return &result, nil
		}

		return nil, nil
	}


func (ica *IGDBCacheAdapter) SetCachedResults(
	ctx context.Context,
	query searchdef.SearchQuery,
	result *searchdef.SearchResult,
) error {
	cacheKey := query.ToCacheKey()

	return ica.cacheWrapper.SetCachedResults(ctx, cacheKey, result)
}
