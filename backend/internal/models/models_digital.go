package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	ServiceTypeBasic        = "basic"
	ServiceTypeSubscription = "subscription"
)

type DigitalLocation struct {
	ID          string            `json:"id" db:"id"`
	UserID      string            `json:"user_id" db:"user_id"`
	Name        string            `json:"name" db:"name"`
	ServiceType string            `json:"service_type" db:"service_type"`
	IsActive    bool              `json:"is_active" db:"is_active"`
	URL         string            `json:"url" db:"url"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
	Items       []Game            `json:"items" db:"items"`
	Subscription *Subscription     `json:"subscription,omitempty"`
}

// Payment model
type Payment struct {
	ID            int64     `json:"id" db:"id"`
	LocationID    string    `json:"location_id" db:"digital_location_id"`
	Amount        float64     `json:"amount" db:"amount"`
	PaymentDate   time.Time `json:"payment_date" db:"payment_date"`
	PaymentMethod string    `json:"payment_method" db:"payment_method"`
	TransactionID string    `json:"transaction_id,omitempty" db:"transaction_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// IsSubscriptionService returns true if the location's service type is 'subscription'
func (dl *DigitalLocation) IsSubscriptionService() bool {
	return dl.ServiceType == "subscription" || dl.Subscription != nil
}

// MarshalJSON implements json.Marshaler for DigitalLocation
// to include computed fields like isSubscriptionService
func (dl DigitalLocation) MarshalJSON() ([]byte, error) {
	type Alias DigitalLocation
	return json.Marshal(&struct {
		Alias
		IsSubscriptionService bool `json:"isSubscriptionService"`
	}{
		Alias:                Alias(dl),
		IsSubscriptionService: dl.IsSubscriptionService(),
	})
}

// ToFrontendDigitalLocation converts a backend DigitalLocation to a frontend-compatible format
// This is used for API responses that need to match frontend expectations
func (dl *DigitalLocation) ToFrontendDigitalLocation() map[string]interface{} {

	// Get logo name from service name
	logoName := getLogoNameFromService(dl.Name)

	result := map[string]interface{}{
		"id":                    dl.ID,
		"user_id":               dl.UserID,
		"name":                  dl.Name,
		"service_type":          dl.ServiceType,
		"is_active":             dl.IsActive,
		"url":                   dl.URL,
		"logo":                  logoName,
		"label":                 getDisplayName(dl.Name),
		"created_at":            dl.CreatedAt.Format(time.RFC3339),
		"updated_at":            dl.UpdatedAt.Format(time.RFC3339),
		"isSubscriptionService": dl.IsSubscriptionService(),
	}

	// Remove subscription object inclusion
	// if dl.Subscription != nil {
	// 	result["subscription"] = dl.Subscription
	// }

	// Add billing information
  // For non-subscription services, provide defaults
  if !dl.IsSubscriptionService() {
			result["billing"] = map[string]interface{}{
					"cycle": "NA",
					"fees": map[string]interface{}{
							"monthly": "FREE",
							"quarterly": "FREE",
							"annual": "FREE",
					},
					"renewalDate": map[string]interface{}{
							"day": "NA",
							"month": "NA",
					},
					"paymentMethod": "None",
			}
	} else if dl.Subscription != nil {
			// For subscription services with subscription data
			monthlyCost := formatCurrency(dl.Subscription.CostPerCycle)

			// Calculate quarterly and annual costs based on monthly
			quarterlyCost := formatCurrency(dl.Subscription.CostPerCycle * 3)
			annualCost := formatCurrency(dl.Subscription.CostPerCycle * 12)

			// Map backend billing cycle to frontend format
			cycle := mapBillingCycleToFrontend(dl.Subscription.BillingCycle)

			billingInfo := map[string]interface{}{
					"cycle": cycle,
					"fees": map[string]interface{}{
							"monthly": monthlyCost,
							"quarterly": quarterlyCost,
							"annual": annualCost,
					},
					"paymentMethod": dl.Subscription.PaymentMethod,
			}

			// Add renewal date if available
			if !dl.Subscription.NextPaymentDate.IsZero() {
					month := dl.Subscription.NextPaymentDate.Format("January")
					day := dl.Subscription.NextPaymentDate.Day()

					billingInfo["renewalDate"] = map[string]interface{}{
							"month": month,
							"day": day,
					}
			} else {
					// Default renewal date
					billingInfo["renewalDate"] = map[string]interface{}{
							"month": "January",
							"day": 1,
					}
			}

			result["billing"] = billingInfo
	}

	if len(dl.Items) > 0 {
		result["items"] = dl.Items
	}

	return result
}

// Subscription model
type Subscription struct {
	ID              int64     `json:"id" db:"id"`
	LocationID      string    `json:"location_id" db:"digital_location_id"`
	BillingCycle    string    `json:"billing_cycle" db:"billing_cycle"`
	CostPerCycle    float64     `json:"cost_per_cycle" db:"cost_per_cycle"`
	NextPaymentDate time.Time `json:"next_payment_date" db:"next_payment_date"`
	PaymentMethod   string    `json:"payment_method" db:"payment_method"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// UnmarshalJSON implements json.Unmarshaler for Subscription
func (s *Subscription) UnmarshalJSON(data []byte) error {
	type Alias Subscription
	aux := &struct {
		NextPaymentDate string `json:"next_payment_date"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.NextPaymentDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05Z", aux.NextPaymentDate)
		if err != nil {
			return err
		}
		s.NextPaymentDate = t
	}
	return nil
}

// Helper function to get logo name from service name
func getLogoNameFromService(name string) string {
	// Convert to lowercase and trim whitespace
	normalizedName := strings.ToLower(strings.TrimSpace(name))

	// Handle special cases with explicit mappings
	logoMappings := map[string]string{
			"playstation": "playstation",
			"playstation network": "playstation",
			"psn": "playstation",
			"xbox": "xbox",
			"xbox network": "xbox",
			"xbox game pass": "xbox",
			"steam": "steam",
			"epic games": "epic",
			"epic games store": "epic",
			"nintendo": "nintendo",
			"nintendo switch online": "nintendo",
			"ea play": "ea",
			"electronic arts": "ea",
			"ubisoft": "ubisoft",
			"ubisoft+": "ubisoft",
			"gog": "gog",
			"gog.com": "gog",
			"humble bundle": "humble",
			"humble": "humble",
			"green man gaming": "greenman",
			"fanatical": "fanatical",
			"apple arcade": "apple",
			"netflix games": "netflix",
			"geforce now": "nvidia",
			"nvidia": "nvidia",
			"prime gaming": "prime",
			"amazon luna": "luna",
			"luna": "luna",
			"meta quest": "meta",
			"meta": "meta",
			"google play pass": "playpass",
			"play pass": "playpass",
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
			"steam": "Steam",
			"psn": "PlayStation Network",
			"playstation": "PlayStation Network",
			"xbox": "Xbox Network",
			"nintendo": "Nintendo Switch Online",
			"epic": "Epic Games Store",
			"epicgames": "Epic Games Store",
			"ea": "EA Play",
			"eaplay": "EA Play",
			"gog": "GOG.com",
			"ubisoft": "Ubisoft+",
			"applearcade": "Apple Arcade",
			"netflix": "Netflix Games",
			"netflixgames": "Netflix Games",
			"nvidia": "GeForce Now",
			"geforce": "GeForce Now",
			"prime": "Prime Gaming",
			"primegaming": "Prime Gaming",
			"playpass": "Google Play Pass",
			"meta": "Meta Quest+",
			"quest+": "Meta Quest+",
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

// Helper function to map backend billing cycle to frontend format
func mapBillingCycleToFrontend(backendCycle string) string {
	cycleMappings := map[string]string{
			"monthly": "1 month",
			"quarterly": "3 months",
			"semi_annual": "6 months",
			"annual": "1 year",
			"yearly": "1 year",
	}

	if frontendCycle, exists := cycleMappings[backendCycle]; exists {
			return frontendCycle
	}

	// If no mapping exists, return the original
	return backendCycle
}

// Helper function to format currency
func formatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}