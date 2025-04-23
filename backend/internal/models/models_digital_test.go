package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDigitalLocation_IsSubscriptionService(t *testing.T) {
	tests := []struct {
		name         string
		serviceType  string
		subscription *Subscription
		expected     bool
	}{
		{
			name:         "Basic service type without subscription",
			serviceType:  "basic",
			subscription: nil,
			expected:     false,
		},
		{
			name:         "Subscription service type without subscription",
			serviceType:  "subscription",
			subscription: nil,
			expected:     true,
		},
		{
			name:        "Basic service type with subscription",
			serviceType: "basic",
			subscription: &Subscription{
				BillingCycle: "monthly",
			},
			expected: true,
		},
		{
			name:        "Subscription service type with subscription",
			serviceType: "subscription",
			subscription: &Subscription{
				BillingCycle: "monthly",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := DigitalLocation{
				ServiceType:  tt.serviceType,
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
		ServiceType: "subscription",
		IsActive:    true,
		URL:         "https://example.com",
		CreatedAt:   now,
		UpdatedAt:   now,
		Items:       []Game{},
		Subscription: &Subscription{
			ID:              1,
			LocationID:      "test-id",
			BillingCycle:    "monthly",
			CostPerCycle:    9.99,
			NextPaymentDate: now,
			PaymentMethod:   "Visa",
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

	if serviceType, ok := result["service_type"].(string); !ok || serviceType != "subscription" {
		t.Errorf("Expected service_type to be subscription, got %v", result["service_type"])
	}
}

func TestDigitalLocation_ToFrontendDigitalLocation(t *testing.T) {
	now := time.Now()
	dl := DigitalLocation{
		ID:          "test-id",
		UserID:      "user-id",
		Name:        "Test Service",
		ServiceType: "subscription",
		IsActive:    true,
		URL:         "https://example.com",
		CreatedAt:   now,
		UpdatedAt:   now,
		Items:       []Game{},
		Subscription: &Subscription{
			ID:              1,
			LocationID:      "test-id",
			BillingCycle:    "monthly",
			CostPerCycle:    9.99,
			NextPaymentDate: now,
			PaymentMethod:   "Visa",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}

	frontendDL := dl.ToFrontendDigitalLocation()

	// Verify key fields
	if id, ok := frontendDL["id"].(string); !ok || id != "test-id" {
		t.Errorf("Expected id to be test-id, got %v", frontendDL["id"])
	}

	if name, ok := frontendDL["name"].(string); !ok || name != "Test Service" {
		t.Errorf("Expected name to be Test Service, got %v", frontendDL["name"])
	}

	if serviceType, ok := frontendDL["service_type"].(string); !ok || serviceType != "subscription" {
		t.Errorf("Expected service_type to be subscription, got %v", frontendDL["service_type"])
	}

	if isSubscription, ok := frontendDL["isSubscriptionService"].(bool); !ok || !isSubscription {
		t.Errorf("Expected isSubscriptionService to be true, got %v", frontendDL["isSubscriptionService"])
	}

	// Verify billing field is present instead of subscription
	if billing, ok := frontendDL["billing"].(map[string]interface{}); !ok || billing == nil {
		t.Errorf("Expected billing to be present, got %v", frontendDL["billing"])
	}

	// Verify subscription field is NOT present
	if _, ok := frontendDL["subscription"]; ok {
		t.Errorf("Expected subscription field to be absent, but it was found")
	}
}