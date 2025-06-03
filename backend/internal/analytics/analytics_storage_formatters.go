package analytics

import (
	"fmt"
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

		// Log before formatting
		fmt.Printf("\n[Formatter] Before formatting %s:\n", loc.Name)
		fmt.Printf("  Location Type: %v\n", loc.LocationType)
		fmt.Printf("  Sublocations Count: %v\n", len(loc.Sublocations))

		// Format sublocations if any exist
		for j := range loc.Sublocations {
			subloc := &loc.Sublocations[j]
			subloc.LocationType = "sublocation" // Ensure sublocation type is set
		}

		// Log after formatting
		fmt.Printf("\n[Formatter] After formatting %s:\n", loc.Name)
		fmt.Printf("  Location Type: %v\n", loc.LocationType)
		fmt.Printf("  Sublocations Count: %v\n", len(loc.Sublocations))
	}
}