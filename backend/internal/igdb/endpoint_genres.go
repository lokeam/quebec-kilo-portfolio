package igdb

import (
	"github.com/lokeam/qko-beta/internal/types"
)

func (c *IGDBClient) GetGenres(genreIDs []int64) (map[int64]types.Genre, error) {
    c.logger.Info("igdb client - GetGenres called with genreIDs", map[string]any{
        "genreIDs": genreIDs,
    })

    if len(genreIDs) == 0 {
        return map[int64]types.Genre{}, nil
    }

    query := buildIDQuery(genreIDs, "id,name,slug")
    var genres []types.Genre
    if err := c.makeRequest("genres", query, &genres); err != nil {
        return nil, err
    }

    genreMap := make(map[int64]types.Genre)
    for _, genre := range genres {
        genreMap[genre.ID] = genre
    }

    return genreMap, nil
}