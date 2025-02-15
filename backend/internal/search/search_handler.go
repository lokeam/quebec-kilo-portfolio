package search

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/lokeam/qko-beta/internal/appcontext" // Concrete implementation of SearchService.
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

// SearchHandler is wired with explicit domain-specific services.
type SearchHandler struct {
	appContext        *appcontext.AppContext
	gameSearchService SearchService
}

// NewSearchHandler creates a new SearchHandler using ONLY the provided AppContext.
// It resolves the search service dependency internally to avoid circular dependencies.
func NewSearchHandler(appCtx *appcontext.AppContext) *SearchHandler {
	// Instantiate the concrete search service.
	gameService, err := NewGameSearchService(appCtx)
	if err != nil {
		panic(err)
	}
	return &SearchHandler{
		appContext:        appCtx,
		gameSearchService: gameService,
	}
}

// ServeHTTP implements the http.Handler interface by sending the request to HandleSearch.
func (s *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := s.HandleSearch(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleSearch delegates the search based on the "domain" query parameter.
func (s *SearchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) error {
	// 1. Get the search query parameter.
	requestID := r.Header.Get(httputils.XRequestIDHeader)
	query := r.URL.Query().Get("query")
	if query == "" {
		err := errors.New("search query is required")
		httputils.RespondWithError(
			httputils.NewResponseWriterAdapter(w),
			s.appContext.Logger,
			requestID,
			err,
		)
		return err
	}

	// 2. Retrieve the IGDB access token key.
	twitchAccessTokenKey, err := s.appContext.Config.IGDB.GetAccessTokenKey()
	if err != nil || twitchAccessTokenKey == ""{
		err := errors.New("failed to retrieve token")
		 httputils.RespondWithError(
			httputils.NewResponseWriterAdapter(w),
			s.appContext.Logger,
			requestID,
			err,
		)
		return err
	}

	s.appContext.Logger.Info("Successfully retrieved Twitch token", map[string]any{
		"request_id": requestID,
		"token_key":  twitchAccessTokenKey,
	})

	// 3. Get the domain parameter. Default to "games" if not provided.
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		domain = "games"
	}

	// 4. Optional limit parameter. Default to 50.
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	s.appContext.Logger.Info("Handling search", map[string]any{
		"query":      query,
		"limit":      limit,
		"domain":     domain,
		"request_id": requestID,
	})

	// 5. Build the search request.
	req := searchdef.SearchRequest{Query: query, Limit: limit}
	var result *searchdef.SearchResult

	// 6. Dispatch to the appropriate service.
	switch domain {
	case "games":
		result, err = s.gameSearchService.Search(r.Context(), req)
		// TODO: Add more domains here.
	default:
		err := errors.New("unsupported search domain")
		httputils.RespondWithError(
			httputils.NewResponseWriterAdapter(w),
			s.appContext.Logger,
			requestID,
			err,
		)
		return err
	}
	if err != nil {
		httputils.RespondWithError(
			httputils.NewResponseWriterAdapter(w),
			s.appContext.Logger,
			requestID,
			err,
		)
		return err
	}

	// 7. Construct a unified response.
	response := searchdef.SearchResponse{
		Games: result.Games,
		Total: len(result.Games),
	}

	// 8. Return the search response as JSON.
	return httputils.RespondWithJSON(w, s.appContext.Logger, http.StatusOK, response)
}