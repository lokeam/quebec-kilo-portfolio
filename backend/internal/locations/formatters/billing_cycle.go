package formatters

// Billing cycle formatter constants
const (
	// Frontend to backend mapping
	FrontendMonthly     = "1 month"
	FrontendQuarterly   = "3 months"
	FrontendBiAnnually  = "6 months"
	FrontendAnnually    = "1 year"

	// Backend to frontend mapping
	BackendMonthly    = "monthly"
	BackendQuarterly  = "quarterly"
	BackendAnnually   = "annually"
)

// Billing cycle formatter maps
var (
	frontendToBackendBillingCycle = map[string]string{
		FrontendMonthly:    BackendMonthly,
		FrontendQuarterly:  BackendQuarterly,
		FrontendBiAnnually: BackendQuarterly, // Map bi-annually to quarterly for backend
		FrontendAnnually:   BackendAnnually,
	}

	backendToFrontendBillingCycle = map[string]string{
		BackendMonthly:   FrontendMonthly,
		BackendQuarterly: FrontendQuarterly,
		BackendAnnually:  FrontendAnnually,
		// Direct mappings for when frontend formats are used in backend
		FrontendMonthly:   FrontendMonthly,
		FrontendQuarterly: FrontendQuarterly,
		FrontendBiAnnually: FrontendBiAnnually,
		FrontendAnnually:  FrontendAnnually,
	}
)

// FormatBillingCycleToBackend converts frontend billing cycle format to backend format
func FormatBillingCycleToBackend(frontendCycle string) string {
	if backendCycle, exists := frontendToBackendBillingCycle[frontendCycle]; exists {
		return backendCycle
	}
	// Default to monthly if unknown
	return BackendMonthly
}

// FormatBillingCycleToFrontend converts backend billing cycle format to frontend format
func FormatBillingCycleToFrontend(backendCycle string) string {
	if frontendCycle, exists := backendToFrontendBillingCycle[backendCycle]; exists {
		return frontendCycle
	}

	// Check if the backend cycle already matches a frontend format (like "1 year")
	// This is for cases where the frontend format was preserved in the backend
	switch backendCycle {
	case "1 month", "1 year", "3 months", "6 months":
		return backendCycle
	}

	// Default to monthly if unknown
	return FrontendMonthly
}