package analytics

import (
	"fmt"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/shared/utils"
)

// Helper function to format currency
func formatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// FormatStorageStats formats storage statistics for the frontend
func FormatStorageStats(stats *StorageStats) {
	// Log data before formatting
	fmt.Printf("\n[Formatter] Data before formatting:\n")

	// DEBUG - Log digital locations
	// for _, loc := range stats.DigitalLocations {
	// 	fmt.Printf("Location: %s\n", loc.Name)
	// 	fmt.Printf("  Billing Cycle: %v\n", loc.BillingCycle)
	// 	fmt.Printf("  Cost Per Cycle: %v\n", loc.CostPerCycle)
	// 	fmt.Printf("  Monthly Cost: %v\n", loc.MonthlyCost)
	// }

	// DEBUG - Log physical locations
	for _, loc := range stats.PhysicalLocations {
		fmt.Printf("Location: %s\n", loc.Name)
		fmt.Printf("  Location Type: %v\n", loc.LocationType)
		fmt.Printf("  Sublocations Count: %v\n", len(loc.Sublocations))
	}

	// Format digital locations
	for i := range stats.DigitalLocations {
		loc := &stats.DigitalLocations[i]
		loc.LocationType = "digital"

		// Digital Location Debug logs before formatting
		// fmt.Printf("\n[Formatter] Before formatting %s:\n", loc.Name)
		// fmt.Printf("  Billing Cycle: %v\n", loc.BillingCycle)
		// fmt.Printf("  Cost Per Cycle: %v\n", loc.CostPerCycle)
		// fmt.Printf("  Monthly Cost (from DB): %v\n", loc.MonthlyCost)

		// Keep the monthly cost as is since it's already calculated in the database
		// No need to recalculate since we're using the standardized format

		// Digital Location Debug logs after formatting
		// fmt.Printf("\n[Formatter] After formatting %s:\n", loc.Name)
		// fmt.Printf("  Billing Cycle: %v\n", loc.BillingCycle)
		// fmt.Printf("  Cost Per Cycle: %v\n", loc.CostPerCycle)
		// fmt.Printf("  Monthly Cost (final): %v\n", loc.MonthlyCost)
	}

	// Format physical locations
	for i := range stats.PhysicalLocations {
		loc := &stats.PhysicalLocations[i]
		loc.LocationType = "physical"

		fmt.Printf("\n[Formatter] Before formatting %s:\n", loc.Name)
		fmt.Printf("  Location Type: %v\n", loc.LocationType)
		fmt.Printf("  Map Coordinates: %+v\n", loc.MapCoordinates)
		fmt.Printf("  Sublocations Count: %v\n", len(loc.Sublocations))

		// Generate Google Maps link if coordinates exist
		if loc.MapCoordinates.Coords != "" {
			fmt.Printf("  Attempting to parse coordinates: %q\n", loc.MapCoordinates.Coords)
			lat, lng, err := utils.ParseCoordinates(loc.MapCoordinates.Coords)
			if err != nil {
				fmt.Printf("  Failed to parse coordinates: %v\n", err)
			} else {
				fmt.Printf("  Successfully parsed coordinates: lat=%f, lng=%f\n", lat, lng)
				// Create a new PhysicalMapCoordinates struct with both fields
				loc.MapCoordinates = models.PhysicalMapCoordinates{
					Coords:         loc.MapCoordinates.Coords,
					GoogleMapsLink: utils.BuildGoogleMapsURL(lat, lng),
				}
				fmt.Printf("  Generated Google Maps link: %s\n", loc.MapCoordinates.GoogleMapsLink)
			}
		} else {
			fmt.Printf("  No coordinates to parse\n")
		}

		// Format sublocations if any exist
		for j := range loc.Sublocations {
			subloc := &loc.Sublocations[j]
			subloc.LocationType = "sublocation" // Ensure sublocation type is set
		}

		fmt.Printf("\n[Formatter] After formatting %s:\n", loc.Name)
		fmt.Printf("  Location Type: %v\n", loc.LocationType)
		fmt.Printf("  Map Coordinates: %+v\n", loc.MapCoordinates)
		fmt.Printf("  Sublocations Count: %v\n", len(loc.Sublocations))
	}
}