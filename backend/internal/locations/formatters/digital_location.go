package formatters

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

func FormatDigitalLocationToFrontend(dl *models.DigitalLocation) map[string]any {
	// Debug logging for payment method
	fmt.Printf("\n[DEBUG] Formatter input for %s:\n", dl.Name)
	fmt.Printf("  Is Subscription: %v\n", dl.IsSubscription)
	fmt.Printf("  Location Payment Method: %v\n", dl.PaymentMethod)
	if dl.Subscription != nil {
		fmt.Printf("  Subscription Payment Method: %v\n", dl.Subscription.PaymentMethod)
	}

	// Get logo name from service name
	logoName := getLogoNameFromService(dl.Name)

	result := map[string]any{
			"id":           dl.ID,
			"name":         dl.Name,
			"is_active":    dl.IsActive,
			"url":          dl.URL,
			"logo":         logoName,
			"label":        getDisplayName(dl.Name),
			"created_at":   dl.CreatedAt.Format(time.RFC3339),
			"updated_at":   dl.UpdatedAt.Format(time.RFC3339),
			"isSubscriptionService": dl.IsSubscription,
			"paymentMethod": dl.PaymentMethod,
	}

	// Create billing object
	if dl.Subscription != nil {
			var monthlyCost, quarterlyCost, annualCost string

			switch dl.Subscription.BillingCycle {
			case "1 month":
				monthlyCost = formatCurrency(dl.Subscription.CostPerCycle)
				quarterlyCost = formatCurrency(dl.Subscription.CostPerCycle * 3)
				annualCost = formatCurrency(dl.Subscription.CostPerCycle * 12)
			case "3 month":
				monthlyCost = formatCurrency(dl.Subscription.CostPerCycle / 3)
				quarterlyCost = formatCurrency(dl.Subscription.CostPerCycle)
				annualCost = formatCurrency(dl.Subscription.CostPerCycle * 4)
			case "6 month":
				monthlyCost = formatCurrency(dl.Subscription.CostPerCycle / 6)
				quarterlyCost = formatCurrency(dl.Subscription.CostPerCycle / 2)
				annualCost = formatCurrency(dl.Subscription.CostPerCycle * 2)
			case "12 month":
				monthlyCost = formatCurrency(dl.Subscription.CostPerCycle / 12)
				quarterlyCost = formatCurrency(dl.Subscription.CostPerCycle / 4)
				annualCost = formatCurrency(dl.Subscription.CostPerCycle)
			default:
				// Log unknown billing cycle and default to monthly
				fmt.Printf("Unknown billing cycle: %s, defaulting to monthly calculations\n", dl.Subscription.BillingCycle)
				monthlyCost = formatCurrency(dl.Subscription.CostPerCycle)
				quarterlyCost = formatCurrency(dl.Subscription.CostPerCycle * 3)
				annualCost = formatCurrency(dl.Subscription.CostPerCycle * 12)
			}

			billingInfo := map[string]any{
					"cycle": dl.Subscription.BillingCycle,
					"fees": map[string]any{
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

					billingInfo["renewalDate"] = map[string]any{
							"month": month,
							"day":   day,
					}
			} else {
					// Default renewal date
					billingInfo["renewalDate"] = map[string]any{
							"month": "January",
							"day":   1,
					}
			}

			result["billing"] = billingInfo
	} else {
			// Non-subscription services or subscription services without data (should never happen after validation)
			result["billing"] = map[string]any{
					"cycle": "NA",
					"fees": map[string]any{
							"monthly":   "FREE",
							"quarterly": "FREE",
							"annual":    "FREE",
					},
					"paymentMethod": dl.PaymentMethod,
					"renewalDate": map[string]any{
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
	// Round to 2 decimal places
	rounded := math.Round(amount*100) / 100
	return fmt.Sprintf("$%.2f", rounded)
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
		"psn":         "PlayStation Plus",
		"playstation": "PlayStation Plus",
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
