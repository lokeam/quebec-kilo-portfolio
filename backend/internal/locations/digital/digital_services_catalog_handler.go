package digital

import (
	"net/http"
	"strings"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

// Returns a handler fn that servces the digital catalog with optional filtering by name
func NewDigitalServicesCatalogHandler(appContext *appcontext.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add cache headers
		w.Header().Set("Cache-Control", "public, max-age=3600")

		q := r.URL.Query().Get("q")

		if q == "" {
			httputils.RespondWithJSON(w, appContext.Logger, http.StatusOK, DigitalServicesCatalog)
			return
		}

		// Only allocate memory for results when filtering
		var filtered []DigitalServiceCatalogItem
		lowercaseQuery := strings.ToLower(q)

		for _, service := range DigitalServicesCatalog {
			if strings.Contains(strings.ToLower(service.Name), lowercaseQuery) {
				filtered = append(filtered, service)
			}
		}

		httputils.RespondWithJSON(w, appContext.Logger, http.StatusOK, filtered)
	}
}