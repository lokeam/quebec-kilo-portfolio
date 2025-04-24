package formatters

import (
	"testing"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
)

func TestFormatDigitalLocationToFrontend(t *testing.T) {
    now := time.Now()

    /*
        GIVEN a location with subscription data
        WHEN FormatDigitalLocationToFrontend() is called
        THEN the method returns a location with billing object instead of subscription
    */
    t.Run(`FormatDigitalLocationToFrontend() converts subscription to billing object`, func(t *testing.T) {
        input := &models.DigitalLocation{
            ID:          "test-id",
            Name:        "Steam",
            ServiceType: "subscription",
            IsActive:    true,
            URL:         "https://test.com",
            CreatedAt:   now,
            UpdatedAt:   now,
            Subscription: &models.Subscription{
                ID:            1,
                LocationID:    "test-id",
                BillingCycle:  BackendMonthly,
                CostPerCycle:  9.99,
                PaymentMethod: "credit_card",
                NextPaymentDate: now.AddDate(0, 1, 0), // 1 month in the future
                CreatedAt:     now,
                UpdatedAt:     now,
            },
        }

        got := FormatDigitalLocationToFrontend(input)

        // Check basic fields
        if got["id"] != input.ID {
            t.Errorf("FormatDigitalLocationToFrontend() id = %v, want %v", got["id"], input.ID)
        }
        if got["name"] != input.Name {
            t.Errorf("FormatDigitalLocationToFrontend() name = %v, want %v", got["name"], input.Name)
        }
        if got["service_type"] != input.ServiceType {
            t.Errorf("FormatDigitalLocationToFrontend() service_type = %v, want %v", got["service_type"], input.ServiceType)
        }
        if got["is_active"] != input.IsActive {
            t.Errorf("FormatDigitalLocationToFrontend() is_active = %v, want %v", got["is_active"], input.IsActive)
        }
        if got["url"] != input.URL {
            t.Errorf("FormatDigitalLocationToFrontend() url = %v, want %v", got["url"], input.URL)
        }
        if got["created_at"] != input.CreatedAt.Format(time.RFC3339) {
            t.Errorf("FormatDigitalLocationToFrontend() created_at = %v, want %v", got["created_at"], input.CreatedAt.Format(time.RFC3339))
        }
        if got["updated_at"] != input.UpdatedAt.Format(time.RFC3339) {
            t.Errorf("FormatDigitalLocationToFrontend() updated_at = %v, want %v", got["updated_at"], input.UpdatedAt.Format(time.RFC3339))
        }

        // Check new fields
        if got["logo"] != "steam" {
            t.Errorf("FormatDigitalLocationToFrontend() logo = %v, want %v", got["logo"], "steam")
        }
        if got["label"] != "Steam" {
            t.Errorf("FormatDigitalLocationToFrontend() label = %v, want %v", got["label"], "Steam")
        }
        if got["isSubscriptionService"] != true {
            t.Errorf("FormatDigitalLocationToFrontend() isSubscriptionService = %v, want %v", got["isSubscriptionService"], true)
        }

        // Check billing object exists
        billing, ok := got["billing"].(map[string]interface{})
        if !ok {
            t.Fatal("FormatDigitalLocationToFrontend() did not create billing object")
        }

        // Check billing fields
        if billing["cycle"] != FrontendMonthly {
            t.Errorf("FormatDigitalLocationToFrontend() billing.cycle = %v, want %v", billing["cycle"], FrontendMonthly)
        }

        fees, ok := billing["fees"].(map[string]interface{})
        if !ok {
            t.Fatal("FormatDigitalLocationToFrontend() did not create fees object")
        }

        if fees["monthly"] != "$9.99" {
            t.Errorf("FormatDigitalLocationToFrontend() fees.monthly = %v, want %v", fees["monthly"], "$9.99")
        }
        if fees["quarterly"] != "$29.97" {
            t.Errorf("FormatDigitalLocationToFrontend() fees.quarterly = %v, want %v", fees["quarterly"], "$29.97")
        }
        if fees["annual"] != "$119.88" {
            t.Errorf("FormatDigitalLocationToFrontend() fees.annual = %v, want %v", fees["annual"], "$119.88")
        }

        if billing["paymentMethod"] != "credit_card" {
            t.Errorf("FormatDigitalLocationToFrontend() paymentMethod = %v, want %v", billing["paymentMethod"], "credit_card")
        }

        // Check renewal date
        renewalDate, ok := billing["renewalDate"].(map[string]interface{})
        if !ok {
            t.Fatal("FormatDigitalLocationToFrontend() did not create renewalDate object")
        }

        expectedMonth := input.Subscription.NextPaymentDate.Format("January")
        expectedDay := input.Subscription.NextPaymentDate.Day()

        if renewalDate["month"] != expectedMonth {
            t.Errorf("FormatDigitalLocationToFrontend() renewalDate.month = %v, want %v", renewalDate["month"], expectedMonth)
        }
        if renewalDate["day"] != expectedDay {
            t.Errorf("FormatDigitalLocationToFrontend() renewalDate.day = %v, want %v", renewalDate["day"], expectedDay)
        }

        // Verify subscription object is not present
        if _, ok := got["subscription"]; ok {
            t.Error("FormatDigitalLocationToFrontend() should not include subscription object")
        }
    })

    /*
        GIVEN a location with subscription type but no subscription data
        WHEN FormatDigitalLocationToFrontend() is called
        THEN the method returns a location with default billing object
    */
    t.Run(`FormatDigitalLocationToFrontend() creates default billing for subscription services without data`, func(t *testing.T) {
        input := &models.DigitalLocation{
            ID:          "test-id",
            Name:        "Steam",
            ServiceType: "subscription",
            IsActive:    true,
            URL:         "https://test.com",
            CreatedAt:   now,
            UpdatedAt:   now,
        }

        got := FormatDigitalLocationToFrontend(input)

        // Check isSubscriptionService flag
        if got["isSubscriptionService"] != true {
            t.Errorf("FormatDigitalLocationToFrontend() isSubscriptionService = %v, want %v", got["isSubscriptionService"], true)
        }

        // Check billing object exists
        billing, ok := got["billing"].(map[string]interface{})
        if !ok {
            t.Fatal("FormatDigitalLocationToFrontend() did not create billing object")
        }

        // Check default subscription billing values
        if billing["cycle"] != "1 month" {
            t.Errorf("FormatDigitalLocationToFrontend() billing.cycle = %v, want %v", billing["cycle"], "1 month")
        }

        fees, ok := billing["fees"].(map[string]interface{})
        if !ok {
            t.Fatal("FormatDigitalLocationToFrontend() did not create fees object")
        }

        if fees["monthly"] != "$0.00" {
            t.Errorf("FormatDigitalLocationToFrontend() fees.monthly = %v, want %v", fees["monthly"], "$0.00")
        }
        if fees["quarterly"] != "$0.00" {
            t.Errorf("FormatDigitalLocationToFrontend() fees.quarterly = %v, want %v", fees["quarterly"], "$0.00")
        }
        if fees["annual"] != "$0.00" {
            t.Errorf("FormatDigitalLocationToFrontend() fees.annual = %v, want %v", fees["annual"], "$0.00")
        }

        // Check renewal date
        renewalDate, ok := billing["renewalDate"].(map[string]interface{})
        if !ok {
            t.Fatal("FormatDigitalLocationToFrontend() did not create renewalDate object")
        }

        if renewalDate["month"] != "January" {
            t.Errorf("FormatDigitalLocationToFrontend() renewalDate.month = %v, want %v", renewalDate["month"], "January")
        }
        if renewalDate["day"] != 1 {
            t.Errorf("FormatDigitalLocationToFrontend() renewalDate.day = %v, want %v", renewalDate["day"], 1)
        }
    })

    /*
        GIVEN a location without subscription data
        WHEN FormatDigitalLocationToFrontend() is called
        THEN the method returns a location with default billing object
    */
    t.Run(`FormatDigitalLocationToFrontend() creates default billing for non-subscription services`, func(t *testing.T) {
        input := &models.DigitalLocation{
            ID:          "test-id",
            Name:        "Epic Games Store",
            ServiceType: "basic",
            IsActive:    true,
            URL:         "https://test.com",
            CreatedAt:   now,
            UpdatedAt:   now,
            Items: []models.Game{{Name: "Test Game", ID: 1}},
        }

        got := FormatDigitalLocationToFrontend(input)

        // Check basic fields
        if got["id"] != input.ID {
            t.Errorf("FormatDigitalLocationToFrontend() id = %v, want %v", got["id"], input.ID)
        }
        if got["name"] != input.Name {
            t.Errorf("FormatDigitalLocationToFrontend() name = %v, want %v", got["name"], input.Name)
        }
        if got["service_type"] != input.ServiceType {
            t.Errorf("FormatDigitalLocationToFrontend() service_type = %v, want %v", got["service_type"], input.ServiceType)
        }
        if got["is_active"] != input.IsActive {
            t.Errorf("FormatDigitalLocationToFrontend() is_active = %v, want %v", got["is_active"], input.IsActive)
        }
        if got["url"] != input.URL {
            t.Errorf("FormatDigitalLocationToFrontend() url = %v, want %v", got["url"], input.URL)
        }

        // Check new fields
        if got["logo"] != "epic" {
            t.Errorf("FormatDigitalLocationToFrontend() logo = %v, want %v", got["logo"], "epic")
        }
        if got["label"] != "Epic Games Store" {
            t.Errorf("FormatDigitalLocationToFrontend() label = %v, want %v", got["label"], "Epic Games Store")
        }
        if got["isSubscriptionService"] != false {
            t.Errorf("FormatDigitalLocationToFrontend() isSubscriptionService = %v, want %v", got["isSubscriptionService"], false)
        }

        // Check billing object exists
        billing, ok := got["billing"].(map[string]interface{})
        if !ok {
            t.Fatal("FormatDigitalLocationToFrontend() did not create billing object")
        }

        // Check default billing values
        if billing["cycle"] != "NA" {
            t.Errorf("FormatDigitalLocationToFrontend() billing.cycle = %v, want %v", billing["cycle"], "NA")
        }

        fees, ok := billing["fees"].(map[string]interface{})
        if !ok {
            t.Fatal("FormatDigitalLocationToFrontend() did not create fees object")
        }

        if fees["monthly"] != "FREE" {
            t.Errorf("FormatDigitalLocationToFrontend() fees.monthly = %v, want %v", fees["monthly"], "FREE")
        }
        if fees["quarterly"] != "FREE" {
            t.Errorf("FormatDigitalLocationToFrontend() fees.quarterly = %v, want %v", fees["quarterly"], "FREE")
        }
        if fees["annual"] != "FREE" {
            t.Errorf("FormatDigitalLocationToFrontend() fees.annual = %v, want %v", fees["annual"], "FREE")
        }

        if billing["paymentMethod"] != "None" {
            t.Errorf("FormatDigitalLocationToFrontend() paymentMethod = %v, want %v", billing["paymentMethod"], "None")
        }

        // Check renewal date
        renewalDate, ok := billing["renewalDate"].(map[string]interface{})
        if !ok {
            t.Fatal("FormatDigitalLocationToFrontend() did not create renewalDate object")
        }

        if renewalDate["month"] != "NA" {
            t.Errorf("FormatDigitalLocationToFrontend() renewalDate.month = %v, want %v", renewalDate["month"], "NA")
        }
        if renewalDate["day"] != "NA" {
            t.Errorf("FormatDigitalLocationToFrontend() renewalDate.day = %v, want %v", renewalDate["day"], "NA")
        }

        // Check items are included
        items, ok := got["items"].([]models.Game)
        if !ok {
            t.Fatal("FormatDigitalLocationToFrontend() did not include items")
        }
        if len(items) != 1 || items[0].ID != 1 {
            t.Errorf("FormatDigitalLocationToFrontend() items = %v, want %v", items, input.Items)
        }
    })
}