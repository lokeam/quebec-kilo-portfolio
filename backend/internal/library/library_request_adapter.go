package library

import (
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

type LibraryRequestAdapter struct{}

func NewLibraryRequestAdapter() *LibraryRequestAdapter {
	return &LibraryRequestAdapter{}
}

func (a *LibraryRequestAdapter) AdaptCreateRequestToLibraryGameModel(
	req types.CreateLibraryGameRequest,
	) models.GameToSave {
	return models.GameToSave{
		GameID:               req.GameID,
		GameName:             req.GameName,
		GameCoverURL:         req.GameCoverURL,
		GameFirstReleaseDate: req.GameFirstReleaseDate,
		GameType:             models.GameToSaveIGDBType{
			DisplayText:        req.GameType.DisplayText,
			NormalizedText:     req.GameType.NormalizedText,
		},
		GameThemeNames:       req.GameThemeNames,
		PlatformLocations:    a.transformPlatformLocations(req.GamesByPlatformAndLocation),
	}
}

func (a *LibraryRequestAdapter) AdaptUpdateRequestToLibraryGameModel(
	req types.UpdateLibraryGameRequest,
) models.GameToSave {
	return models.GameToSave{
		PlatformLocations: a.transformPlatformLocations(req.GamesByPlatformAndLocation),
	}
}

func (a *LibraryRequestAdapter) transformPlatformLocations(
	locations []types.LibraryRequestGameLocation,
) []models.GameToSaveLocation {
	platformLocations := make([]models.GameToSaveLocation, len(locations))

	for i := 0; i < len(locations); i++ {
		platformLocations[i] = models.GameToSaveLocation{
			PlatformID:   locations[i].PlatformID,
			PlatformName: locations[i].PlatformName,
			Type:         locations[i].Type,
			Location: models.GameToSaveLocationDetails{
				SublocationID:     locations[i].Location.SublocationID,
				DigitalLocationID: locations[i].Location.DigitalLocationID,
			},
		}
	}
	return platformLocations
}

func (a *LibraryRequestAdapter) getLocationID(
	location types.GameLocation,
) string {
	if location.DigitalLocationID != "" {
		return location.DigitalLocationID
	}
	return location.SublocationID
}