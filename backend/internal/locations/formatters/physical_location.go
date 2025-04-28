package formatters

import (
	"html"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

// FormatPhysicalLocationToFrontend converts a physical location model to a frontend-compatible format
func FormatPhysicalLocationToFrontend(pl *models.PhysicalLocation) map[string]interface{} {
	result := map[string]interface{}{
		"id":             pl.ID,
		"name":           html.UnescapeString(pl.Name),
		"label":          html.UnescapeString(pl.Label),
		"location_type":  pl.LocationType,
		"map_coordinates": pl.MapCoordinates,
		"created_at":     pl.CreatedAt.Format(time.RFC3339),
		"updated_at":     pl.UpdatedAt.Format(time.RFC3339),
	}

	// Format sublocations if they exist
	if pl.SubLocations != nil {
		var formattedSublocations []map[string]interface{}
		for _, subloc := range *pl.SubLocations {
			formattedSubloc := map[string]interface{}{
				"id":                  subloc.ID,
				"name":                html.UnescapeString(subloc.Name),
				"location_type":       subloc.LocationType,
				"bg_color":            subloc.BgColor,
				"stored_items":        subloc.StoredItems,
				"created_at":          subloc.CreatedAt.Format(time.RFC3339),
				"updated_at":          subloc.UpdatedAt.Format(time.RFC3339),
			}

			// Format items if they exist
			if subloc.Items != nil {
				formattedSubloc["items"] = subloc.Items
			} else {
				formattedSubloc["items"] = []models.Game{}
			}

			formattedSublocations = append(formattedSublocations, formattedSubloc)
		}
		result["sublocations"] = formattedSublocations
	} else {
		result["sublocations"] = []map[string]interface{}{}
	}

	return result
}

// FormatPhysicalLocationsToFrontend converts a slice of physical locations to frontend-compatible format
func FormatPhysicalLocationsToFrontend(locations []models.PhysicalLocation) []map[string]interface{} {
	var result []map[string]interface{}
	for _, location := range locations {
		formattedLocation := FormatPhysicalLocationToFrontend(&location)
		result = append(result, formattedLocation)
	}
	return result
}