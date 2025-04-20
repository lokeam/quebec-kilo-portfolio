package search

import (
	"fmt"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/services"
)

type searchServiceFactory struct {
	appCtx *appcontext.AppContext
}

// NewSearchServiceFactory creates a new search service factory
func NewSearchServiceFactory(appCtx *appcontext.AppContext) services.SearchServiceFactory {
	return &searchServiceFactory{appCtx: appCtx}
}

// GetService returns a search service for the given domain
func (ssf *searchServiceFactory) GetService(domain string) (services.SearchService, error) {
	switch domain {
	case "games":
			searchService, err := NewGameSearchService(ssf.appCtx)
			if err != nil {
				return nil, err
			}
			return searchService, nil
	default:
			return nil, fmt.Errorf("unsupported domain: %s", domain)
	}
}
