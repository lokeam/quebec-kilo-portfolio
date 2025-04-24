package formatters

import (
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

// FormatDigitalLocationToFrontend converts a DigitalLocation model to frontend-compatible format
func FormatDigitalLocationToFrontend(dl *models.DigitalLocation) map[string]interface{} {
	result := map[string]interface{}{
		"id":           dl.ID,
		"name":         dl.Name,
		"service_type": dl.ServiceType,
		"is_active":    dl.IsActive,
		"url":          dl.URL,
		"created_at":   dl.CreatedAt.Format(time.RFC3339),
		"updated_at":   dl.UpdatedAt.Format(time.RFC3339),
	}

	if dl.Subscription != nil {
		result["subscription"] = map[string]interface{}{
			"id":            dl.Subscription.ID,
			"location_id":   dl.Subscription.LocationID,
			"billing_cycle": FormatBillingCycleToFrontend(dl.Subscription.BillingCycle),
			"cost":          dl.Subscription.CostPerCycle,
			"created_at":    dl.Subscription.CreatedAt.Format(time.RFC3339),
			"updated_at":    dl.Subscription.UpdatedAt.Format(time.RFC3339),
		}
	}

	return result
}