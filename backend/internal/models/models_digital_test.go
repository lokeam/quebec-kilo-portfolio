package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDigitalLocation_IsSubscriptionService(t *testing.T) {
	tests := []struct {
		name         string
		isSubscription bool
		subscription *Subscription
		expected     bool
	}{
		{
			name:         "Basic service without subscription",
			isSubscription: false,
			subscription: nil,
			expected:     false,
		},
		{
			name:         "Subscription service without subscription",
			isSubscription: true,
			subscription: nil,
			expected:     true,
		},
		{
			name:        "Basic service with subscription",
			isSubscription: false,
			subscription: &Subscription{
				BillingCycle: "1 month",
			},
			expected: true,
		},
		{
			name:        "Subscription service with subscription",
			isSubscription: true,
			subscription: &Subscription{
				BillingCycle: "1 month",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := DigitalLocation{
				IsSubscription: tt.isSubscription,
				Subscription: tt.subscription,
			}

			if got := dl.IsSubscriptionService(); got != tt.expected {
				t.Errorf("DigitalLocation.IsSubscriptionService() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDigitalLocation_MarshalJSON(t *testing.T) {
	now := time.Now()
	dl := DigitalLocation{
		ID:          "test-id",
		UserID:      "user-id",
		Name:        "Test Service",
		IsSubscription: true,
		IsActive:    true,
		URL:         "https://example.com",
		CreatedAt:   now,
		UpdatedAt:   now,
		Items:       []Game{},
		Subscription: &Subscription{
			ID:              1,
			LocationID:      "test-id",
			BillingCycle:    "1 month",
			CostPerCycle:    9.99,
			NextPaymentDate: now,
			PaymentMethod:   "visa",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}

	data, err := json.Marshal(dl)
	if err != nil {
		t.Fatalf("Failed to marshal DigitalLocation: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify the isSubscriptionService field was added
	if isSubscription, ok := result["isSubscriptionService"].(bool); !ok || !isSubscription {
		t.Errorf("MarshalJSON did not correctly add isSubscriptionService field, got %v", result["isSubscriptionService"])
	}

	// Verify other fields were preserved
	if id, ok := result["id"].(string); !ok || id != "test-id" {
		t.Errorf("Expected id to be test-id, got %v", result["id"])
	}

	if name, ok := result["name"].(string); !ok || name != "Test Service" {
		t.Errorf("Expected name to be Test Service, got %v", result["name"])
	}

	if isSubscription, ok := result["is_subscription"].(bool); !ok || !isSubscription {
		t.Errorf("Expected is_subscription to be true, got %v", result["is_subscription"])
	}
}
