package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/types"
)

type IGDBClient struct {
	baseURL        string
	clientID       string
	token          string
	httpClient     *http.Client
	logger         interfaces.Logger
	appContext     *appcontext.AppContext
}

// NewIGDBClient creates a new IGDB client.
//func NewIGDBClient(clientID, token string) *Client {
func NewIGDBClient(appContext *appcontext.AppContext, token string) *IGDBClient {
	if appContext == nil {
		panic("appContext is nil")
	}
	if token == "" {
			panic("token is empty")
	}
	if appContext.Config.IGDB.ClientID == "" {
			panic("ClientID is empty")
	}
	if appContext.Config.IGDB.BaseURL == "" {
			panic("BaseURL is empty")
	}


	return &IGDBClient{
		//clientID: clientID,
		appContext:     appContext,
		clientID:       appContext.Config.IGDB.ClientID,
		token:          token,
		httpClient:     &http.Client{},
		baseURL:        appContext.Config.IGDB.BaseURL,
		logger:         appContext.Logger,
	}
}

// SearchGames performs a search against the games endpoint.
// The query parameter should be a valid IGDB query string.
func (c *IGDBClient) SearchGames(query string) ([]*types.Game, error) {
	//fmt.Println("igdb client - query: ", query)
	//c.logger.Info("igdb client - query: ", map[string]any{"query": query})

	var games []*types.Game
	if err := c.makeRequest("games", query, &games); err != nil {
		return nil, err
	}

	c.logger.Debug("igdb client - SearchGames - raw games fetched", map[string]any{
		"games": games,
	})

	// Get Game details
	results, err := c.GetGameDetailsBySearch(games)
	if err != nil {
		c.logger.Error("igdb client - SearchGames - failed to get game details: %w", map[string]any{"error": err})
		return nil, err
	}

	c.logger.Debug("igdb client - SearchGames - results with details: ", map[string]any{"results": results})

	return games, nil
}


// makeRequest sends an HTTP POST request to the specified endpoint with the given query.
func (c *IGDBClient) makeRequest(endpoint string, query string, result interface{}) error {
	// Early return if client is nil to prevent nil pointer dereference
	if c == nil {
		return fmt.Errorf("IGDBClient is nil")
	}

// Move logging after nil check
	c.logger.Info("igdb client - makeRequest called with endpoint: %s and query: %s",
		map[string]any{"endpoint": endpoint, "query": query})

	url := fmt.Sprintf("%s/%s", IGDB_API_URL, endpoint)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(query)))
	if err != nil {
			c.logger.Error("igdb client - makeRequest - failed to create request: %w", map[string]any{"error": err})
			return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Client-ID", c.clientID)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
			c.logger.Error("igdb client - makeRequest - failed to send request: %w", map[string]any{"error": err})
			return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Context cancellation check (first, before any other processing)
	select {
	case <-req.Context().Done():
			return req.Context().Err()
	default:
	}

	c.logger.Info("igdb client - makeRequest - response: ", map[string]any{"resp": resp})

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			c.logger.Error("igdb client - makeRequest - failed to decode response: %w", map[string]any{"error": err})
			return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// buildIDQuery builds a query string to fetch records by a list of IDs with specific fields.
func buildIDQuery(ids []int, fields string) string {
	// Convert slice of ints to comma-separated string.
	var strIDs []string
	for _, id := range ids {
		strIDs = append(strIDs, fmt.Sprintf("%d", id))
	}
	return fmt.Sprintf("fields %s; where id = (%s);", fields, strings.Join(strIDs, ","))
}

// GetCovers fetches cover details in bulk by cover IDs.
func (c *IGDBClient) GetCovers(ids []int) (map[int]types.Cover, error) {
	c.logger.Info("igdb client - GetCovers called with ids: ", map[string]any{"ids": ids})

	if len(ids) == 0 {
		return map[int]types.Cover{}, nil
}

	query := buildIDQuery(ids, "id,image_id,url")
	c.logger.Debug("igdb client - GetCovers - query: ", map[string]any{"query": query})

	var covers []types.Cover
	if err := c.makeRequest("covers", query, &covers); err != nil {
		c.logger.Error("igdb client - GetCovers - failed to make request: %w", map[string]any{"error": err})
		return nil, err
	}

	c.logger.Debug("igdb client - GetCovers - raw covers fetched", map[string]any{"covers": covers})

	coverMap := make(map[int]types.Cover)
	for _, cover := range covers {
		coverMap[cover.ID] = cover
}
	c.logger.Debug("igdb client - GetCovers - coverMap constructed: ", map[string]any{"coverMap": coverMap})
	return coverMap, nil
}

func (c *IGDBClient) GetCoverURL(imageID string) string {
	if imageID == "" {
		c.logger.Error("igdb client - GetCoverURL - imageID is empty", map[string]any{"imageID": imageID})
			return ""
	}
	return fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", imageID)
}

// GetGenres fetches genre details in bulk by genre IDs.
func (c *IGDBClient) GetGenres(ids []int) (map[int]types.Genre, error) {
	if len(ids) == 0 {
		return map[int]types.Genre{}, nil
	}
	query := buildIDQuery(ids, "id,name")
	var genres []types.Genre
	if err := c.makeRequest("genres", query, &genres); err != nil {
		return nil, err
	}
	genreMap := make(map[int]types.Genre)
	for _, genre := range genres {
		genreMap[genre.ID] = genre
	}
	return genreMap, nil
}

// GetPlatforms fetches platform details in bulk by platform IDs.
func (c *IGDBClient) GetPlatforms(ids []int) (map[int]types.Platform, error) {
	if len(ids) == 0 {
		return map[int]types.Platform{}, nil
	}
	query := buildIDQuery(ids, "id,name")
	var platforms []types.Platform
	if err := c.makeRequest("platforms", query, &platforms); err != nil {
		return nil, err
	}
	platformMap := make(map[int]types.Platform)
	for _, p := range platforms {
		platformMap[p.ID] = p
	}
	return platformMap, nil
}

// GetGameDetailsBySearch combines the calls to search for games and then expand
// the game details by fetching related covers, genres, and platforms.
func (c *IGDBClient) GetGameDetailsBySearch(games []*types.Game) ([]types.GameDetails, error) {
	c.logger.Info("igdb client - GetGameDetailsBySearch called with games", map[string]any{"games": games})

	// Collect unique IDs for covers, genres, and platforms.
	coverIDsSet := make(map[int]struct{})
	genreIDsSet := make(map[int]struct{})
	platformIDsSet := make(map[int]struct{})

	for _, game := range games {
			if game.Cover != 0 {
					coverIDsSet[game.Cover] = struct{}{}
			}
			// for _, id := range game.Genres {
			// 		genreIDsSet[id] = struct{}{}
			// }
			// for _, id := range game.Platforms {
			// 		platformIDsSet[id] = struct{}{}
			// }
	}

	// Convert sets to slices.
	coverIDs := make([]int, 0, len(coverIDsSet))
	for id := range coverIDsSet {
			coverIDs = append(coverIDs, id)
	}

	genreIDs := make([]int, 0, len(genreIDsSet))
	for id := range genreIDsSet {
			genreIDs = append(genreIDs, id)
	}

	platformIDs := make([]int, 0, len(platformIDsSet))
	for id := range platformIDsSet {
			platformIDs = append(platformIDs, id)
	}

	// Get cover details
	covers, err := c.GetCovers(coverIDs)
	if err != nil {
			c.logger.Warn("igdb client - GetGameDetailsBySearch - failed to get covers", map[string]any{"error": err})
	}

	// Create cover map for quick lookup
	coverMap := make(map[int]types.Cover)
	for _, cover := range covers {
			coverMap[cover.ID] = cover
	}

	// Initialize results slice
	var results []types.GameDetails
	for _, game := range games {
			details := types.GameDetails{
					ID:      game.ID,
					Name:    game.Name,
					Summary: game.Summary,
					Rating:  game.Rating,
					CoverURL: game.CoverURL,
			}

			// Handle cover mapping
			if game.Cover != 0 {
					if cover, exists := coverMap[game.Cover]; exists {
							game.CoverURL = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", cover.ImageID)
							c.logger.Debug("Cover URL set: ", map[string]any{
								"gameID": game.ID,
								"coverURL": game.CoverURL,
							})
					} else {
							c.logger.Warn("igdb client - GetGameDetailsBySearch - cover not found for game",
									map[string]any{"gameID": game.ID, "coverID": game.Cover})
					}
			}

			results = append(results, details)
	}

	c.logger.Debug("igdb client - GetGameDetailsBySearch - results constructed",
			map[string]any{"results": results})

	return results, nil
}


// SetHTTPClient allows setting a custom HTTP client for testing purposes.
func (c *IGDBClient) SetHTTPClient(client *http.Client) {
	if client == nil {
			// Fallback to default client if nil is passed
			client = &http.Client{}
	}
	c.httpClient = client
}