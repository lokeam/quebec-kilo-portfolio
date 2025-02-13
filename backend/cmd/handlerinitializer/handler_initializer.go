package handlerinitializer

import (
	"net/http"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/health"
	"github.com/lokeam/qko-beta/internal/search"
)

type HandlerInitializer struct {
	Search  http.Handler
	Health  http.Handler
}

func NewHandlerInitializer(appCtx *appcontext.AppContext) *HandlerInitializer {
	return &HandlerInitializer{
		Search: search.NewSearchHandler(appCtx),
		Health: health.NewHealthHandler(appCtx.Config, appCtx.Logger),
	}
}
