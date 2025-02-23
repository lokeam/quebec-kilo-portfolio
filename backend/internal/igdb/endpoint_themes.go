package igdb

import (
	"github.com/lokeam/qko-beta/internal/types"
)

func (c *IGDBClient) GetThemes(themeIDs []int64) (map[int64]types.Theme, error) {
	c.logger.Info("igdb client - Get Themes called with themeIDs", map[string]any{
		"themeIDs": themeIDs,
	})

	if len(themeIDs) == 0 {
		return map[int64]types.Theme{}, nil
	}

	query := buildIDQuery(themeIDs, "id,name,slug")
	var themes []types.Theme
	if err := c.makeRequest("themes", query, &themes); err != nil {
		return nil, err
	}

	themeMap := make(map[int64]types.Theme)
	for _, theme := range themes {
		themeMap[theme.ID] = theme
	}

	return themeMap, nil
}
