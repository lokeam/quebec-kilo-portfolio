package library

import (
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
	// Create maps for quick lookups
	physicalLocationsByGame := make(map[int64][]types.LibraryGamePhysicalLocationDBResponse)
	digitalLocationsByGame := make(map[int64][]types.LibraryGameDigitalLocationDBResponse)

	// Group locations by game ID
	for _, loc := range physicalLocations {
		physicalLocationsByGame[loc.GameID] = append(physicalLocationsByGame[loc.GameID], loc)
	}
	for _, loc := range digitalLocations {
		digitalLocationsByGame[loc.GameID] = append(digitalLocationsByGame[loc.GameID], loc)
	}

	// Transform games
	responses := make([]types.LibraryGameResponse, len(games))
	for i, game := range games {
		responses[i] = types.LibraryGameResponse{
			ID:              game.ID,
			Name:            game.Name,
			CoverURL:        game.CoverURL,
			FirstReleaseDate: game.FirstReleaseDate,
			Rating:          game.Rating,
			ThemeNames:      game.ThemeNames,
			IsInLibrary:     true, // We know this is true since it's from the library
			IsInWishlist:    game.IsInWishlist,
			Favorite:        game.Favorite,
			GameType: struct {
				DisplayText    string `json:"displayText"`
				NormalizedText string `json:"normalizedText"`
			}{
				DisplayText:    game.GameTypeDisplay,
				NormalizedText: game.GameTypeNormalized,
			},
			PhysicalLocations: a.adaptPhysicalLocations(physicalLocationsByGame[game.ID]),
			DigitalLocations:  a.adaptDigitalLocations(digitalLocationsByGame[game.ID]),
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
			PhysicalLocations: a.adaptPhysicalLocations(physicalLocations),
			DigitalLocations:  a.adaptDigitalLocations(digitalLocations),
		},
	}
}

// Helper function to adapt physical locations
func (a *LibraryResponseAdapter) adaptPhysicalLocations(
	locations []types.LibraryGamePhysicalLocationDBResponse,
) []types.PhysicalLocationResponse {
	responses := make([]types.PhysicalLocationResponse, len(locations))
	for i, loc := range locations {
		responses[i] = types.PhysicalLocationResponse{
			Name:    loc.LocationName,
			Type:    loc.LocationType,
			Sublocation: types.SublocationResponse{
				Name:     loc.SublocationName,
				Type:     loc.SublocationType,
				BgColor:  loc.SublocationBgColor,
			},
			Platform: loc.PlatformName,
		}
	}
	return responses
}

// Helper function to adapt digital locations
func (a *LibraryResponseAdapter) adaptDigitalLocations(
	locations []types.LibraryGameDigitalLocationDBResponse,
) []types.DigitalLocationResponse {
	responses := make([]types.DigitalLocationResponse, len(locations))
	for i, loc := range locations {
		responses[i] = types.DigitalLocationResponse{
			Name:          loc.LocationName,
			NormalizedName: loc.NormalizedName,
			Platform:      loc.PlatformName,
		}
	}
	return responses
}