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
	for _, loc := range stats.DigitalLocations {
		fmt.Printf("Location: %s\n", loc.Name)
		fmt.Printf("  Billing Cycle: %v\n", loc.BillingCycle)
		fmt.Printf("  Cost Per Cycle: %v\n", loc.CostPerCycle)
		fmt.Printf("  Monthly Cost (from DB): %v\n", loc.MonthlyCost)
	}

	// Format digital locations
	for i := range stats.DigitalLocations {
		loc := &stats.DigitalLocations[i]
		loc.LocationType = "digital"

		// Log before formatting
		fmt.Printf("\n[Formatter] Before formatting %s:\n", loc.Name)
		fmt.Printf("  Billing Cycle: %v\n", loc.BillingCycle)
		fmt.Printf("  Cost Per Cycle: %v\n", loc.CostPerCycle)
		fmt.Printf("  Monthly Cost (from DB): %v\n", loc.MonthlyCost)

		// Keep the monthly cost as is since it's already calculated in the database
		// No need to recalculate since we're using the standardized format

		// Log after formatting
		fmt.Printf("\n[Formatter] After formatting %s:\n", loc.Name)
		fmt.Printf("  Billing Cycle: %v\n", loc.BillingCycle)
		fmt.Printf("  Cost Per Cycle: %v\n", loc.CostPerCycle)
		fmt.Printf("  Monthly Cost (final): %v\n", loc.MonthlyCost)
	}
}