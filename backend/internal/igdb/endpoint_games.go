package igdb

import (
	"fmt"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

func (c *IGDBClient) SearchGames(query string) ([]*models.Game, error) {
    c.logger.Info("igdb client - SearchGames - query: ", map[string]any{"query": query})

    var games []*models.Game
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

    //c.logger.Debug("igdb client - SearchGames - results with details: ", map[string]any{"results": results})

    // Convert GameDetails back to Game and filter out games without cover URLs
    var convertedGames []*models.Game
    for i := 0; i < len(results); i++ {
        detail := results[i]
        // Skip games without cover URLs
        if detail.CoverURL == "" {
            c.logger.Debug("igdb client - SearchGames - skipping game without cover URL", map[string]any{
                "gameID": detail.ID,
                "name": detail.Name,
            })
            continue
        }

        // Create a new GameType with the data from IGDB
        // gameType := types.GameType{
        //     DisplayText:   detail.GameType.Type,
        //     NormalizedText: strings.ToLower(strings.ReplaceAll(detail.GameType.Type, " ", "_")),
        // }
        gameType := types.GameType{
            ID: detail.GameType.ID,
            Type: detail.GameType.Type,
            DisplayText: types.GameTypes[detail.GameType.ID].DisplayText,
            NormalizedText: types.GameTypes[detail.GameType.ID].NormalizedText,
        }

        // Log the game type data for debugging
        c.logger.Debug("Creating new GameType UPDATED", map[string]any{
            "gameID": detail.ID,
            "gameType": gameType,
        })

        convertedGames = append(convertedGames, &models.Game{
            ID:                detail.ID,
            Name:              detail.Name,
            Summary:           detail.Summary,
            CoverURL:          detail.CoverURL,
            FirstReleaseDate:  detail.FirstReleaseDate,
            Rating:            detail.Rating,
            PlatformNames:     detail.PlatformNames,
            GenreNames:        detail.GenreNames,
            ThemeNames:        detail.ThemeNames,
            GameType:          gameType,
        })
    }

    return convertedGames, nil
}

func (c *IGDBClient) GetGameDetailsBySearch(games []*models.Game) ([]*types.GameDetails, error) {
    //c.logger.Info("igdb client - GetGameDetailsBySearch called with games", map[string]any{"games": games})

    // Collect unique IDs for covers, genres, and platforms.
    coverIDsSet := make(map[int64]struct{})
    platformIDsSet := make(map[int64]struct{})
    genreIDsSet := make(map[int64]struct{})
    themeIDsSet := make(map[int64]struct{})

    for _, game := range games {
        if game.CoverID != 0 {
            coverIDsSet[game.CoverID] = struct{}{}
        }
        for _, platformID := range game.Platforms {
            platformIDsSet[platformID] = struct{}{}
        }
        for _, genreID := range game.Genres {
            genreIDsSet[genreID] = struct{}{}
        }
        for _, themeID := range game.Themes {
            themeIDsSet[themeID] = struct{}{}
        }
    }

    // Convert sets to slices.
    coverIDs := make([]int64, 0, len(coverIDsSet))
    for id := range coverIDsSet {
        coverIDs = append(coverIDs, id)
    }

    platformIDs := make([]int64, 0, len(platformIDsSet))
    for id := range platformIDsSet {
        platformIDs = append(platformIDs, id)
    }

    genreIDs := make([]int64, 0, len(genreIDsSet))
    for id := range genreIDsSet {
        genreIDs = append(genreIDs, id)
    }

    themeIDs := make([]int64, 0, len(themeIDsSet))
    for id := range themeIDsSet {
        themeIDs = append(themeIDs, id)
    }

    // Get cover details
    covers, err := c.GetCovers(coverIDs)
    if err != nil {
        c.logger.Warn("igdb client - GetGameDetailsBySearch - failed to get covers", map[string]any{"error": err})
    }

    // Get platform details
    platforms, err := c.GetPlatforms(platformIDs)
    if err != nil {
        c.logger.Warn("igdb client - GetGameDetailsBySearch - failed to get platforms", map[string]any{"error": err})
    }

    // Get genre details
    genres, err := c.GetGenres(genreIDs)
    if err != nil {
        c.logger.Warn("igdb client - GetGameDetailsBySearch - failed to get genres", map[string]any{"error": err})
    } else {
        c.logger.Debug("igdb client - GetGameDetailsBySearch - genres fetched", map[string]any{"genres": genres})
    }

    themes, err := c.GetThemes(themeIDs)
    if err != nil {
        c.logger.Warn("igdb client - GetGameDetailsBySearch - failed to get themes", map[string]any{"error": err})
    } else {
        c.logger.Debug("igdb client - GetGameDetailsBySearch - themes fetched", map[string]any{"themes": themes})
    }

    // Create cover map for quick lookup
    coverMap := make(map[int64]types.Cover)
    for _, cover := range covers {
        coverMap[cover.ID] = cover
    }

    platformMap := make(map[int64]types.Platform)
    for _, platform := range platforms {
        platformMap[platform.ID] = platform
    }

    genreMap := make(map[int64]types.Genre)
    for _, genre := range genres {
        genreMap[genre.ID] = genre
    }

    themeMap := make(map[int64]types.Theme)
    for _, theme := range themes {
        themeMap[theme.ID] = theme
    }

    // Initialize results slice
    var results []*types.GameDetails
    for _, game := range games {
        details := &types.GameDetails{
            ID:                 game.ID,
            Name:               game.Name,
            Summary:            game.Summary,
            FirstReleaseDate:   game.FirstReleaseDate,
            Rating:             game.Rating,
            CoverURL:           game.CoverURL,
            Platforms:          []types.Platform{},
            PlatformNames:      []string{},
            Genres:             []types.Genre{},
            GenreNames:         []string{},
            Themes:             []types.Theme{},
            ThemeNames:         []string{},
            GameType: types.IGDBGameType{
                ID:   game.GameType.ID,
                Type: game.GameType.Type,
            },
        }

        // Log the raw game type data from IGDB response
        c.logger.Debug("Creating new GameType in GetGameDetailsBySearch", map[string]any{
            "id": game.GameType.ID,
            "type": game.GameType.Type,
            "rawGameType": game.GameType,
        })

        // Handle cover mapping
        if game.CoverID != 0 {
            if cover, exists := coverMap[game.CoverID]; exists {
                game.CoverURL = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", cover.ImageID)
                details.CoverURL = game.CoverURL
                c.logger.Debug("Cover URL set: ", map[string]any{
                    "gameID": game.ID,
                    "coverURL": game.CoverURL,
                })
            } else {
                c.logger.Warn("igdb client - GetGameDetailsBySearch - cover not found for game",
                    map[string]any{"gameID": game.ID, "coverID": game.CoverID})
            }
        }

        // Handle platform mapping
        for _, platformID := range game.Platforms {
            if platform, exists := platformMap[platformID]; exists {
                details.Platforms = append(details.Platforms, platform)
                details.PlatformNames = append(details.PlatformNames, platform.Name)
            } else {
                c.logger.Warn("igdb client = GetGameDetailsBySearch - platform not found for game", map[string]any{
                    "gameID": game.ID,
                    "platformID": platform.ID,
                })
            }
        }

        // Handle genre mapping
        for _, genreID := range game.Genres {
            if genre, exists := genreMap[genreID]; exists {
                details.Genres = append(details.Genres, genre)
                details.GenreNames = append(details.GenreNames, genre.Name)
            } else {
                c.logger.Warn("igdb client - GetGameDetailsBySearch - genre not found for game", map[string]any{
                    "gameID":  game.ID,
                    "genreID": genreID,
                })
            }
        }

        // Handle theme mapping
        for _, themeID := range game.Themes {
            if theme, exists := themeMap[themeID]; exists {
                details.Themes = append(details.Themes, theme)
                details.ThemeNames = append(details.ThemeNames, theme.Name)
            } else {
                c.logger.Warn("igdb client - GetGameDetailsBySearch - theme not found for game", map[string]any{
                    "gameID": game.ID,
                    "themeID": themeID,
                })
            }
        }

        results = append(results, details)
    }

    return results, nil
}