package search

import (
	"fmt"

	"github.com/lokeam/qko-beta/internal/appcontext"
)

type searchServiceFactory struct {
	appCtx *appcontext.AppContext
}

type SearchServiceFactory interface {
	GetService(domain string) (SearchService, error)
}

func NewSearchServiceFactory(appCtx *appcontext.AppContext) SearchServiceFactory {
	return &searchServiceFactory{appCtx: appCtx}
}

func (ssf *searchServiceFactory) GetService(domain string) (SearchService, error) {
	switch domain {
	case "games":
			return NewGameSearchService(ssf.appCtx)
	default:
			return nil, fmt.Errorf("unsupported domain: %s", domain)
	}
}
