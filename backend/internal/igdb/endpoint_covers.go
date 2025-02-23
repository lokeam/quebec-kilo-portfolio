package igdb

import (
	"fmt"

	"github.com/lokeam/qko-beta/internal/types"
)

// GetCovers fetches cover details in bulk by cover IDs.
func (c *IGDBClient) GetCovers(ids []int64) (map[int64]types.Cover, error) {
	c.logger.Info("igdb client - GetCovers called with ids: ", map[string]any{"ids": ids})

	if len(ids) == 0 {
		return map[int64]types.Cover{}, nil
}

	query := buildIDQuery(ids, "id,image_id,url")
	c.logger.Debug("igdb client - GetCovers - query: ", map[string]any{"query": query})

	var covers []types.Cover
	if err := c.makeRequest("covers", query, &covers); err != nil {
		c.logger.Error("igdb client - GetCovers - failed to make request: %w", map[string]any{"error": err})
		return nil, err
	}

	c.logger.Debug("igdb client - GetCovers - raw covers fetched", map[string]any{"covers": covers})

	coverMap := make(map[int64]types.Cover)
	for _, cover := range covers {
		coverMap[cover.ID] = cover
}
	c.logger.Debug("igdb client - GetCovers - coverMap constructed: ", map[string]any{"coverMap": coverMap})
	return coverMap, nil
}

func (c *IGDBClient) GetCoverURL(coverID int) (string, error) {
    if coverID == 0 {
        return "", nil // No cover ID, return empty URL
    }

    // Fetch the cover details
    var covers []types.Cover
    query := fmt.Sprintf("fields image_id; where id = %d;", coverID)
    if err := c.makeRequest("covers", query, &covers); err != nil {
        return "", err
    }

    if len(covers) == 0 {
        return "", fmt.Errorf("cover not found for ID: %d", coverID)
    }

    // Construct the URL
    imageID := covers[0].ImageID
    return fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", imageID), nil
}