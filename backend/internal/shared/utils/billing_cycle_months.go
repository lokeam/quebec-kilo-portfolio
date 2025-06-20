package utils

import "fmt"

// Converts a billing cycle string to the number of months
func GetBillingCycleMonths(billingCycle string) (int, error) {
	switch billingCycle {
	case "1 month":
		return 1, nil
	case "3 month":
		return 3, nil
	case "6 month":
		return 6, nil
	case "12 month":
		return 12, nil
	default:
		return 0, fmt.Errorf("invalid billing cycle: %s. Must be one of: 1 month, 3 month, 6 month, 12 month", billingCycle)
	}
}