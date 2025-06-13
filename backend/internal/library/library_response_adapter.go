package library

import (
	"html"

	"github.com/lokeam/qko-beta/internal/types"
)

type LibraryResponseAdapter struct{}

func NewLibraryResponseAdapter() *LibraryResponseAdapter {
	return &LibraryResponseAdapter{}
}

// AdaptToLibraryResponse transforms database results into a library response
func (a *LibraryResponseAdapter) AdaptToLibraryResponse(
	games []types.LibraryGameDBResult,
	physicalLocations []types.LibraryGamePhysicalLocationDBResponse,
	digitalLocations []types.LibraryGameDigitalLocationDBResponse,
) types.LibraryResponse {
	// Create map for quick lookups
	locationsByGame := make(map[int64][]types.GameLocationDBResult)

	// Group all locations by game ID
	for i := 0; i < len(physicalLocations); i++ {
		loc := physicalLocations[i]
		locationsByGame[loc.ID] = append(locationsByGame[loc.ID], types.GameLocationDBResult{
			ID: loc.ID,
			PlatformID: loc.PlatformID,
			Type: "physical",
			LocationID: loc.LocationID,
			LocationName: html.UnescapeString(loc.LocationName),
			LocationType: loc.LocationType,
			SublocationID: loc.SublocationID,
			SublocationName: loc.SublocationName,
			SublocationType: loc.SublocationType,
			SublocationBgColor: loc.SublocationBgColor,
			PlatformName: loc.PlatformName,
		})
	}
	for i := 0; i < len(digitalLocations); i++ {
		loc := digitalLocations[i]
		locationsByGame[loc.ID] = append(locationsByGame[loc.ID], types.GameLocationDBResult{
			ID: loc.ID,
			PlatformID: loc.PlatformID,
			Type: "digital",
			LocationID: loc.LocationID,
			LocationName: html.UnescapeString(loc.LocationName),
			PlatformName: loc.PlatformName,
			IsActive: loc.IsActive,
		})
	}

	// Transform games
	responses := make([]types.LibraryGameResponse, len(games))
	for i := 0; i < len(games); i++ {
		game := games[i]
		responses[i] = types.LibraryGameResponse{
			ID:              game.ID,
			Name:            game.Name,
			CoverURL:        game.CoverURL,
			FirstReleaseDate: game.FirstReleaseDate,
			Rating:          game.Rating,
			ThemeNames:      game.ThemeNames,
			IsInLibrary:     true,
			IsInWishlist:    game.IsInWishlist,
			Favorite:        game.Favorite,
			GameType: struct {
				DisplayText    string `json:"displayText"`
				NormalizedText string `json:"normalizedText"`
			}{
				DisplayText:    game.GameTypeDisplay,
				NormalizedText: game.GameTypeNormalized,
			},
			GamesByPlatformAndLocation: locationsByGame[game.ID],
		}
	}

	return types.LibraryResponse{
		Success: true,
		Games:   responses,
	}
}

// AdaptToSingleGameResponse transforms a single game database result into a response
func (a *LibraryResponseAdapter) AdaptToSingleGameResponse(
	game types.LibraryGameDBResult,
	physicalLocations []types.LibraryGamePhysicalLocationDBResponse,
	digitalLocations []types.LibraryGameDigitalLocationDBResponse,
) types.SingleGameResponse {
	// Create map for quick lookups
	locationsByGame := make(map[int64][]types.GameLocationDBResult)

	// Group all locations by game ID
	for i := 0; i < len(physicalLocations); i++ {
		loc := physicalLocations[i]
		locationsByGame[loc.ID] = append(locationsByGame[loc.ID], types.GameLocationDBResult{
			ID: loc.ID,
			PlatformID: loc.PlatformID,
			Type: "physical",
			LocationID: loc.LocationID,
			LocationName: html.UnescapeString(loc.LocationName),
			LocationType: loc.LocationType,
			SublocationID: loc.SublocationID,
			SublocationName: loc.SublocationName,
			SublocationType: loc.SublocationType,
			SublocationBgColor: loc.SublocationBgColor,
			PlatformName: loc.PlatformName,
		})
	}
	for i := 0; i < len(digitalLocations); i++ {
		loc := digitalLocations[i]
		locationsByGame[loc.ID] = append(locationsByGame[loc.ID], types.GameLocationDBResult{
			ID: loc.ID,
			PlatformID: loc.PlatformID,
			Type: "digital",
			LocationID: loc.LocationID,
			LocationName: html.UnescapeString(loc.LocationName),
			PlatformName: loc.PlatformName,
			IsActive: loc.IsActive,
		})
	}

	return types.SingleGameResponse{
		Success: true,
		Game: types.LibraryGameResponse{
			ID:              game.ID,
			Name:            game.Name,
			CoverURL:        game.CoverURL,
			FirstReleaseDate: game.FirstReleaseDate,
			Rating:          game.Rating,
			ThemeNames:      game.ThemeNames,
			IsInLibrary:     true,
			IsInWishlist:    game.IsInWishlist,
			Favorite:        game.Favorite,
			GameType: struct {
				DisplayText    string `json:"displayText"`
				NormalizedText string `json:"normalizedText"`
			}{
				DisplayText:    game.GameTypeDisplay,
				NormalizedText: game.GameTypeNormalized,
			},
			GamesByPlatformAndLocation: locationsByGame[game.ID],
		},
	}
}
