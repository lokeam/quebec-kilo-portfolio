package handlerinitializer

import (
	"net/http"
)

type HandlerInitializer struct {
	Search  http.Handler
	Health  http.Handler
}

// func NewHandlerInitializer(appCtx *appcontext.AppContext) *HandlerInitializer {
// 	searchServiceFactory := search.NewSearchServiceFactory(appCtx)

// 	return &HandlerInitializer{
// 		Search: search.NewSearchHandler(appCtx, searchServiceFactory),
// 		Health: health.NewHealthHandler(appCtx.Config, appCtx.Logger),
// 	}
// }
