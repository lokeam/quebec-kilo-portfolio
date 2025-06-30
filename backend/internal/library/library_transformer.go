package library

import (
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

// transformLocationDBToResponse converts a LibraryLocationDB to a LibraryGamesByPlatformAndLocationItemFINAL
func TransformLocationDBToResponse(db models.GameLocationDatabaseEntry) types.LibraryGamesByPlatformAndLocationItemFINAL {
	return types.LibraryGamesByPlatformAndLocationItemFINAL{
		ID: db.GameID,
		PlatformID: db.PlatformID,
		PlatformName: db.PlatformName,
		IsPC: db.Category == "pc",
		IsMobile: db.Category == "mobile",
		DateAdded: db.CreatedAt.Unix(),
		ParentLocationID: db.ParentLocationID.String,
		ParentLocationName: db.ParentLocationName.String,
		ParentLocationType: db.ParentLocationType.String,
		ParentLocationBgColor: db.ParentLocationBgColor.String,
		SublocationID: db.SublocationID.String,
		SublocationName: db.SublocationName.String,
		SublocationType: db.SublocationType.String,
	}
}

// transformGameDBToResponse converts a LibraryGameDB to a LibraryGameItemBFFResponseFINAL
func TransformGameDBToResponse(db models.LibraryGameDB, locations []types.LibraryGamesByPlatformAndLocationItemFINAL) types.LibraryGameItemBFFResponseFINAL {
	return types.LibraryGameItemBFFResponseFINAL{
		ID: db.ID,
		Name: db.Name,
		CoverURL: db.CoverURL,
		GameTypeDisplayText: db.GameTypeDisplayText,
		GameTypeNormalizedText: db.GameTypeNormalizedText,
		IsFavorite: db.IsFavorite,
		GamesByPlatformAndLocation: locations,
	}
}

// -- REFACTORED LIBRARY RESPONSE TRANSFORMERS --
// transformToRefactoredResponse transforms database results to the refactored BFF response
func (la *LibraryDbAdapter) TransformToRefactoredResponse(
	games []models.LibraryGameRefactoredDB,
	physicalLocations []models.PhysicalLocationDB,
	digitalLocations []models.DigitalLocationDB,
) types.LibraryBFFRefactoredResponse {
	// Group locations by game ID
	physicalByGame := make(map[int64][]models.PhysicalLocationDB)
	digitalByGame := make(map[int64][]models.DigitalLocationDB)

	for _, loc := range physicalLocations {
			physicalByGame[loc.GameID] = append(physicalByGame[loc.GameID], loc)
	}

	for _, loc := range digitalLocations {
			digitalByGame[loc.GameID] = append(digitalByGame[loc.GameID], loc)
	}

	// Transform games to response format
	libraryItems := make([]types.SingleLibraryGameBFFResponse, 0, len(games))
	recentlyAdded := make([]types.SingleLibraryGameBFFResponse, 0)

	// Calculate 6-month cutoff for recently added
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)

	for _, game := range games {
			// Transform physical locations
			physicalLocations := la.TransformPhysicalLocations(physicalByGame[game.ID])

			// Transform digital locations
			digitalLocations := la.TransformDigitalLocations(digitalByGame[game.ID])

			// Calculate version counts
			totalPhysicalVersions := len(physicalByGame[game.ID])
			totalDigitalVersions := len(digitalByGame[game.ID])

			gameResponse := types.SingleLibraryGameBFFResponse{
					ID:                    game.ID,
					Name:                  game.Name,
					CoverURL:              game.CoverURL,
					IsInWishlist:          game.IsInWishlist,
					FirstReleaseDate:      game.FirstReleaseDate,
					GenreNames:            game.GenreNames,
					GameType: types.GameTypeResponse{
							DisplayText:    game.GameTypeDisplayText,
							NormalizedText: game.GameTypeNormalizedText,
					},
					Favorite:              game.Favorite,
					TotalPhysicalVersions: totalPhysicalVersions,
					TotalDigitalVersions:  totalDigitalVersions,
					PhysicalLocations:     physicalLocations,
					DigitalLocations:      digitalLocations,
			}

			libraryItems = append(libraryItems, gameResponse)

			// Check if game is recently added
			if game.CreatedAt.After(sixMonthsAgo) {
					recentlyAdded = append(recentlyAdded, gameResponse)
			}
	}

	return types.LibraryBFFRefactoredResponse{
			LibraryItems:  libraryItems,
			RecentlyAdded: recentlyAdded,
	}
}

// transformPhysicalLocations transforms physical location database results to response format
func (la *LibraryDbAdapter) TransformPhysicalLocations(
	locations []models.PhysicalLocationDB,
) []types.LibraryBFFSinglePhysicalLocationResponse {
	if len(locations) == 0 {
			return []types.LibraryBFFSinglePhysicalLocationResponse{}
	}

	// Group by location (parent + sublocation combination)
	locationGroups := make(map[string][]models.PhysicalLocationDB)
	for _, loc := range locations {
			key := fmt.Sprintf("%s_%s", loc.ParentLocationID, loc.SublocationID)
			locationGroups[key] = append(locationGroups[key], loc)
	}

	result := make([]types.LibraryBFFSinglePhysicalLocationResponse, 0, len(locationGroups))
	for _, group := range locationGroups {
			if len(group) == 0 {
					continue
			}

			loc := group[0] // Use first item for location details
			platformVersions := make([]types.PlatformVersionResponse, len(group))
			for i, platform := range group {
					platformVersions[i] = types.PlatformVersionResponse{
							PlatformName: platform.PlatformName,
							PlatformId:   platform.PlatformID,
					}
			}

			result = append(result, types.LibraryBFFSinglePhysicalLocationResponse{
					ParentLocationName:    loc.ParentLocationName,
					ParentLocationId:      loc.ParentLocationID,
					ParentLocationType:    loc.ParentLocationType,
					ParentLocationBgColor: loc.ParentLocationBgColor,
					SublocationName:       loc.SublocationName,
					SublocationId:         loc.SublocationID,
					SublocationType:       loc.SublocationType,
					GamePlatformVersions:  platformVersions,
			})
	}

	return result
}

// transformDigitalLocations transforms digital location database results to response format
func (la *LibraryDbAdapter) TransformDigitalLocations(
	locations []models.DigitalLocationDB,
) []types.LibraryBFFSingleDigitalLocationResponse {
	if len(locations) == 0 {
			return []types.LibraryBFFSingleDigitalLocationResponse{}
	}

	// Group by digital location
	locationGroups := make(map[string][]models.DigitalLocationDB)
	for _, loc := range locations {
			locationGroups[loc.DigitalLocationID] = append(locationGroups[loc.DigitalLocationID], loc)
	}

	result := make([]types.LibraryBFFSingleDigitalLocationResponse, 0, len(locationGroups))
	for _, group := range locationGroups {
			if len(group) == 0 {
					continue
			}

			loc := group[0] // Use first item for location details
			platformVersions := make([]types.PlatformVersionResponse, len(group))
			for i, platform := range group {
					platformVersions[i] = types.PlatformVersionResponse{
							PlatformName: platform.PlatformName,
							PlatformId:   platform.PlatformID,
					}
			}

			result = append(result, types.LibraryBFFSingleDigitalLocationResponse{
					DigitalLocationName:  loc.DigitalLocationName,
					DigitalLocationId:    loc.DigitalLocationID,
					GamePlatformVersions: platformVersions,
			})
	}

	return result
}