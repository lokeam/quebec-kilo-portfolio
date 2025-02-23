package igdb

import (
	"github.com/lokeam/qko-beta/internal/types"
)

func (c *IGDBClient) GetPlatforms(platformIDs []int64) (map[int64]types.Platform, error) {
    if len(platformIDs) == 0 {
        return map[int64]types.Platform{}, nil
    }

    query := buildIDQuery(platformIDs, "id,name,platform_logo")
    var platforms []types.Platform
    if err := c.makeRequest("platforms", query, &platforms); err != nil {
        return nil, err
    }

    platformMap := make(map[int64]types.Platform)
    for _, platform := range platforms {
        platformMap[platform.ID] = platform
    }

    return platformMap, nil
}

func (c *IGDBClient) GetPlatformNames(platformIDs []int64) ([]string, error) {
	platformMap, err := c.GetPlatforms(platformIDs)
	if err != nil {
			return nil, err
	}

	var platformNames []string
	for _, platformID := range platformIDs {
			if platform, exists := platformMap[platformID]; exists {
					platformNames = append(platformNames, platform.Name)
			}
	}

	return platformNames, nil
}