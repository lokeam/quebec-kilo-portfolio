package formatters

import (
	"fmt"
	"strings"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

// FormatDigitalLocationToFrontend converts a DigitalLocation model to frontend-compatible format
func FormatDigitalLocationToFrontend(dl *models.DigitalLocation) map[string]interface{} {
	// Get logo name from service name
	logoName := getLogoNameFromService(dl.Name)

	result := map[string]interface{}{
		"id":           dl.ID,
		"name":         dl.Name,
		"service_type": dl.ServiceType,
		"is_active":    dl.IsActive,
		"url":          dl.URL,
		"logo":         logoName,
		"label":        getDisplayName(dl.Name),
		"created_at":   dl.CreatedAt.Format(time.RFC3339),
		"updated_at":   dl.UpdatedAt.Format(time.RFC3339),
		"isSubscriptionService": dl.ServiceType == models.ServiceTypeSubscription || dl.Subscription != nil,
	}

	// Create billing object instead of including subscription
	if dl.Subscription != nil {
		// For subscription services with subscription data
		monthlyCost := formatCurrency(dl.Subscription.CostPerCycle)
		quarterlyCost := formatCurrency(dl.Subscription.CostPerCycle * 3)
		annualCost := formatCurrency(dl.Subscription.CostPerCycle * 12)

		// Map backend billing cycle to frontend format
		cycle := FormatBillingCycleToFrontend(dl.Subscription.BillingCycle)

		billingInfo := map[string]interface{}{
			"cycle": cycle,
			"fees": map[string]interface{}{
				"monthly":   monthlyCost,
				"quarterly": quarterlyCost,
				"annual":    annualCost,
			},
			"paymentMethod": dl.Subscription.PaymentMethod,
		}

		// Add renewal date if available
		if !dl.Subscription.NextPaymentDate.IsZero() {
			month := dl.Subscription.NextPaymentDate.Format("January")
			day := dl.Subscription.NextPaymentDate.Day()

			billingInfo["renewalDate"] = map[string]interface{}{
				"month": month,
				"day":   day,
			}
		} else {
			// Default renewal date
			billingInfo["renewalDate"] = map[string]interface{}{
				"month": "January",
				"day":   1,
			}
		}

		result["billing"] = billingInfo
	} else if dl.ServiceType == models.ServiceTypeSubscription {
		// Default billing info for subscription services without subscription data
		result["billing"] = map[string]interface{}{
			"cycle": "1 month",
			"fees": map[string]interface{}{
				"monthly":   "$0.00",
				"quarterly": "$0.00",
				"annual":    "$0.00",
			},
			"paymentMethod": "None",
			"renewalDate": map[string]interface{}{
				"month": "January",
				"day":   1,
			},
		}
	} else {
		// Non-subscription services
		result["billing"] = map[string]interface{}{
			"cycle": "NA",
			"fees": map[string]interface{}{
				"monthly":   "FREE",
				"quarterly": "FREE",
				"annual":    "FREE",
			},
			"paymentMethod": "None",
			"renewalDate": map[string]interface{}{
				"day":   "NA",
				"month": "NA",
			},
		}
	}

	if len(dl.Items) > 0 {
		result["items"] = dl.Items
	} else {
		result["items"] = []models.Game{}
	}

	return result
}

// Helper function to format currency
func formatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// Helper function to get logo name from service name
func getLogoNameFromService(name string) string {
	// Convert to lowercase and trim whitespace
	normalizedName := strings.ToLower(strings.TrimSpace(name))

	// Handle special cases with explicit mappings
	logoMappings := map[string]string{
		"playstation":            "playstation",
		"playstation network":    "playstation",
		"psn":                    "playstation",
		"xbox":                   "xbox",
		"xbox network":           "xbox",
		"xbox game pass":         "xbox",
		"steam":                  "steam",
		"epic games":             "epic",
		"epic games store":       "epic",
		"nintendo":               "nintendo",
		"nintendo switch online": "nintendo",
		"ea play":                "ea",
		"electronic arts":        "ea",
		"ubisoft":                "ubisoft",
		"ubisoft+":               "ubisoft",
		"gog":                    "gog",
		"gog.com":                "gog",
		"humble bundle":          "humble",
		"humble":                 "humble",
		"green man gaming":       "greenman",
		"fanatical":              "fanatical",
		"apple arcade":           "apple",
		"netflix games":          "netflix",
		"geforce now":            "nvidia",
		"nvidia":                 "nvidia",
		"prime gaming":           "prime",
		"amazon luna":            "luna",
		"luna":                   "luna",
		"meta quest":             "meta",
		"meta":                   "meta",
		"google play pass":       "playpass",
		"play pass":              "playpass",
	}

	if logoName, exists := logoMappings[normalizedName]; exists {
		return logoName
	}

	// For other services, remove spaces and special characters
	simplified := strings.ReplaceAll(normalizedName, " ", "")
	simplified = strings.ReplaceAll(simplified, ".", "")
	simplified = strings.ReplaceAll(simplified, "+", "")
	simplified = strings.ReplaceAll(simplified, "-", "")

	return simplified
}

// Helper function to get display name from service name
func getDisplayName(serviceName string) string {
	// Special case mappings for specific services
	displayNameMappings := map[string]string{
		"steam":       "Steam",
		"psn":         "PlayStation Network",
		"playstation": "PlayStation Network",
		"xbox":        "Xbox Network",
		"nintendo":    "Nintendo Switch Online",
		"epic":        "Epic Games Store",
		"epicgames":   "Epic Games Store",
		"ea":          "EA Play",
		"eaplay":      "EA Play",
		"gog":         "GOG.com",
		"ubisoft":     "Ubisoft+",
		"applearcade": "Apple Arcade",
		"netflix":     "Netflix Games",
		"netflixgames": "Netflix Games",
		"nvidia":      "GeForce Now",
		"geforce":     "GeForce Now",
		"prime":       "Prime Gaming",
		"primegaming": "Prime Gaming",
		"playpass":    "Google Play Pass",
		"meta":        "Meta Quest+",
		"quest+":      "Meta Quest+",
	}

	lowercaseName := strings.ToLower(serviceName)
	if displayName, exists := displayNameMappings[lowercaseName]; exists {
		return displayName
	}

	// Convert first letter of each word to uppercase for other services
	words := strings.Fields(strings.ToLower(serviceName))
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[0:1]) + word[1:]
		}
	}

	return strings.Join(words, " ")
}
