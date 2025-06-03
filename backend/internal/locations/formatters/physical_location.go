package formatters

import (
	"html"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

// FormatPhysicalLocationToFrontend converts a physical location model to a frontend-compatible format
func FormatPhysicalLocationToFrontend(pl *models.PhysicalLocation) map[string]any {
	result := map[string]any{
		"id":                pl.ID,
		"name":              html.UnescapeString(pl.Name),
		"label":             html.UnescapeString(pl.Label),
		"location_type":     pl.LocationType,
		"map_coordinates":   pl.MapCoordinates,
		"bg_color":          pl.BgColor,
		"created_at":        pl.CreatedAt.Format(time.RFC3339),
		"updated_at":        pl.UpdatedAt.Format(time.RFC3339),
	}

	// Format sublocations if they exist
	if pl.SubLocations != nil {
		var formattedSublocations []map[string]any
		for i := 0; i < len(*pl.SubLocations); i++ {
			formattedSubloc := map[string]any{
				"id":                  (*pl.SubLocations)[i].ID,
				"name":                html.UnescapeString((*pl.SubLocations)[i].Name),
				"location_type":       (*pl.SubLocations)[i].LocationType,
				"stored_items":        (*pl.SubLocations)[i].StoredItems,
				"created_at":          (*pl.SubLocations)[i].CreatedAt.Format(time.RFC3339),
				"updated_at":          (*pl.SubLocations)[i].UpdatedAt.Format(time.RFC3339),
			}

			// Format items if they exist
			if (*pl.SubLocations)[i].Items != nil {
				formattedSubloc["items"] = (*pl.SubLocations)[i].Items
			} else {
				formattedSubloc["items"] = []models.Game{}
			}

			formattedSublocations = append(formattedSublocations, formattedSubloc)
		}
		// for _, subloc := range *pl.SubLocations {
		// 	formattedSubloc := map[string]any{
		// 		"id":                  subloc.ID,
		// 		"name":                html.UnescapeString(subloc.Name),
		// 		"location_type":       subloc.LocationType,
		// 		"stored_items":        subloc.StoredItems,
		// 		"created_at":          subloc.CreatedAt.Format(time.RFC3339),
		// 		"updated_at":          subloc.UpdatedAt.Format(time.RFC3339),
		// 	}

		// 	// Format items if they exist
		// 	if subloc.Items != nil {
		// 		formattedSubloc["items"] = subloc.Items
		// 	} else {
		// 		formattedSubloc["items"] = []models.Game{}
		// 	}

		// 	formattedSublocations = append(formattedSublocations, formattedSubloc)
		// }
		result["sublocations"] = formattedSublocations
	} else {
		result["sublocations"] = []map[string]any{}
	}

	return result
}

// FormatPhysicalLocationsToFrontend converts a slice of physical locations to frontend-compatible format
func FormatPhysicalLocationsToFrontend(locations []models.PhysicalLocation) []map[string]any{
	var result []map[string]any
	for i := 0; i < len(locations); i++ {
		formattedLocation := FormatPhysicalLocationToFrontend(&locations[i])
		result = append(result, formattedLocation)
	}
	return result
}