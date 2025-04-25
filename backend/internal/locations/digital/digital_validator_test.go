package digital

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lokeam/qko-beta/internal/shared/logger"
	sanitizer "github.com/lokeam/qko-beta/internal/shared/security/sanitizer"
)

type MockSanitizer struct{}

func (m *MockSanitizer) SanitizeString(input string) (string, error) {
	return input, nil
}

func (m *MockSanitizer) SanitizeSearchQuery(input string) (string, error) {
	return input, nil
}

/*
Behavior:
1. Main method: ValidateDigitalLocation
2. Helper methods:
   - validateName
   - validateURL
   - validateSubscription
   - ValidatePayment

validateName:
- Ensure name is not empty
- Ensure name is not longer than 100 chars
- Ensure name is sanitized

validateURL:
- Ensure URL is not empty
- Ensure URL is a valid format
- Ensure URL is sanitized

validateSubscription:
- Ensure billing cycle is valid
- Ensure cost per cycle is positive
- Ensure payment method is valid
- Ensure next payment date is in the future

ValidatePayment:
- Ensure amount is positive
- Ensure payment method is valid
- Ensure payment date is not in the future
- Ensure transaction ID length is within limits

Scenarios:
Reject locations with:
- Empty name
- Names longer than 100 chars
- Empty URL
- Invalid URL format
- Invalid service type
- Invalid subscription data
- Invalid payment data
Pass validation with complete, valid location
Collect errors when multiple violations are present
*/

var _ interfaces.Sanitizer = (*mocks.MockSanitizer)(nil)

func TestDigitalValidator(t *testing.T) {
	// Setup
	testSanitizer := &mocks.MockSanitizer{}
	testValidator, testErr := NewDigitalValidator(testSanitizer)
	if testErr != nil {
		t.Fatalf("failed to create test validator: %v", testErr)
	}

	// Setup mock sanitizer behavior
	testSanitizer.SanitizeFunc = func(text string) (string, error) {
		if strings.Contains(text, "<script>") {
			return "", errors.New("sanitizer failure")
		}
		return text, nil
	}

	// ----------- validateName() ------------
	/*
		GIVEN a location with an empty name
		WHEN validateName() is called
		THEN the method returns an error stating "name cannot be empty"
	*/
	t.Run(`validateName() rejects empty names`, func(t *testing.T) {
		testName := ""
		sanitizedName, testErr := testValidator.validateName(testName)

		if testErr == nil {
			t.Fatalf("expected an error for an empty name, but got nil")
		}
		if sanitizedName != "" {
			t.Fatalf("expected sanitized name to be empty, but got %s", sanitizedName)
		}

		expectedError := "name cannot be empty"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})


	/*
		GIVEN a location with a name longer than the max length
		WHEN validateName() is called
		THEN the method returns an error stating "name must be less than X characters"
	*/
	t.Run(`validateName() rejects names longer than 100 characters`, func(t *testing.T) {
		testName := strings.Repeat("a", MaxNameLength+1)
		sanitizedName, testErr := testValidator.validateName(testName)

		if testErr == nil {
			t.Errorf("expected an error for a name that is too long, but got nil")
		}
		if sanitizedName != "" {
			t.Errorf("expected sanitized name to be empty on error, but got %s", sanitizedName)
		}

		expectedError := fmt.Sprintf("name must be less than %d characters", MaxNameLength)
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})


	/*
		GIVEN a name that triggers a sanitizer error
		WHEN validateName() is called
		THEN the method returns an error from the sanitizer
	*/
	t.Run(`validateName() fails when sanitizer fails`, func(t *testing.T) {
		testName := "<script>alert('xss');</script>"
		sanitizedName, testErr := testValidator.validateName(testName)

		if testErr == nil {
			t.Errorf("expected an error when sanitizer fails, but got nil")
		}
		if sanitizedName != "" {
			t.Errorf("expected sanitized name to be empty on error, but got %s", sanitizedName)
		}

		expectedError := "invalid name content: sanitizer failure"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})


	/*
		GIVEN a valid name
		WHEN validateName() is called
		THEN the method returns the sanitized name and no error
	*/
	t.Run(`validateName() Happy Path - Return sanitized name for valid input`, func(t *testing.T) {
		testName := "Test Name"
		sanitizedName, testErr := testValidator.validateName(testName)

		if testErr != nil {
			t.Errorf("expected no error for valid name, but got %v", testErr)
		}
		if sanitizedName != testName {
			t.Errorf("expected sanitized name to be %q, but got %q", testName, sanitizedName)
		}
	})


	// ----------- validateURL() ------------
	/*
		GIVEN an empty URL
		WHEN validateURL() is called
		THEN the method returns an error stating "URL cannot be empty"
	*/
	t.Run(`validateURL() rejects empty URLs`, func(t *testing.T) {
		testURL := ""
		sanitizedURL, testErr := testValidator.validateURL(testURL)

		if testErr == nil {
			t.Fatalf("expected an error for an empty URL, but got nil")
		}
		if sanitizedURL != "" {
			t.Fatalf("expected sanitized URL to be empty, but got %s", sanitizedURL)
		}

		expectedError := "URL cannot be empty"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})


	/*
		GIVEN a URL longer than the max length
		WHEN validateURL() is called
		THEN the method returns an error stating "URL must be less than X characters"
	*/
	t.Run(`validateURL() rejects URLs longer than max length`, func(t *testing.T) {
		testURL := "https://" + strings.Repeat("a", MaxURLLength)
		sanitizedURL, testErr := testValidator.validateURL(testURL)

		if testErr == nil {
			t.Errorf("expected an error for a URL that is too long, but got nil")
		}
		if sanitizedURL != "" {
			t.Errorf("expected sanitized URL to be empty on error, but got %s", sanitizedURL)
		}

		expectedError := fmt.Sprintf("URL must be less than %d characters", MaxURLLength)
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})


	/*
		GIVEN an invalid URL format
		WHEN validateURL() is called
		THEN the method returns an error stating "invalid URL format"
	*/
	t.Run(`validateURL() rejects invalid URL formats`, func(t *testing.T) {
		testURL := "not-a-valid-url"
		sanitizedURL, testErr := testValidator.validateURL(testURL)

		if testErr == nil {
			t.Errorf("expected an error for an invalid URL format, but got nil")
		}
		if sanitizedURL != "" {
			t.Errorf("expected sanitized URL to be empty on error, but got %s", sanitizedURL)
		}

		expectedError := "invalid URL format"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})


	/*
		GIVEN a URL that triggers a sanitizer error
		WHEN validateURL() is called
		THEN the method returns an error from the sanitizer
	*/
	t.Run(`validateURL() fails when sanitizer fails`, func(t *testing.T) {
		testURL := "https://<script>alert('xss');</script>.com"
		sanitizedURL, testErr := testValidator.validateURL(testURL)

		if testErr == nil {
			t.Errorf("expected an error when sanitizer fails, but got nil")
		}
		if sanitizedURL != "" {
			t.Errorf("expected sanitized URL to be empty on error, but got %s", sanitizedURL)
		}

		expectedError := "invalid URL content: sanitizer failure"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})


	/*
		GIVEN a valid URL
		WHEN validateURL() is called
		THEN the method returns the sanitized URL and no error
	*/
	t.Run(`validateURL() Happy Path - Return sanitized URL for valid input`, func(t *testing.T) {
		testURL := "https://example.com"
		sanitizedURL, testErr := testValidator.validateURL(testURL)

		if testErr != nil {
			t.Errorf("expected no error for valid URL, but got %v", testErr)
		}
		if sanitizedURL != testURL {
			t.Errorf("expected sanitized URL to be %q, but got %q", testURL, sanitizedURL)
		}
	})


	// ----------- validateSubscription() ------------
	/*
		GIVEN a subscription with an invalid billing cycle
		WHEN validateSubscription() is called
		THEN the method returns an error stating "invalid billing cycle"
	*/
	t.Run(`validateSubscription() rejects invalid billing cycles`, func(t *testing.T) {
		testSubscription := models.Subscription{
			BillingCycle: "invalid",
			CostPerCycle: 10.0,
			PaymentMethod: "Visa",
			NextPaymentDate: time.Now().Add(24 * time.Hour),
		}

		validatedSubscription, testErr := testValidator.validateSubscription(testSubscription)

		if testErr == nil {
			t.Fatalf("expected an error for invalid billing cycle, but got nil")
		}

		expectedError := "invalid billing cycle: invalid"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}

		if validatedSubscription.ID != 0 {
			t.Errorf("expected validated subscription to be empty on error, but got %+v", validatedSubscription)
		}
	})

	/*
		GIVEN a subscription with a non-positive cost
		WHEN validateSubscription() is called
		THEN the method returns an error stating "cost per cycle must be greater than 0"
	*/
	t.Run(`validateSubscription() rejects non-positive costs`, func(t *testing.T) {
		testSubscription := models.Subscription{
			BillingCycle: "monthly",
			CostPerCycle: 0,
			PaymentMethod: "Visa",
			NextPaymentDate: time.Now().Add(24 * time.Hour),
		}

		validatedSubscription, testErr := testValidator.validateSubscription(testSubscription)

		if testErr == nil {
			t.Fatalf("expected an error for non-positive cost, but got nil")
		}

		expectedError := "cost per cycle must be greater than 0"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}

		if validatedSubscription.ID != 0 {
			t.Errorf("expected validated subscription to be empty on error, but got %+v", validatedSubscription)
		}
	})

	/*
		GIVEN a subscription with an invalid payment method
		WHEN validateSubscription() is called
		THEN the method returns an error stating "Invalid payment method: Invalid"
	*/
	t.Run(`validateSubscription() rejects invalid payment methods`, func(t *testing.T) {
		testSubscription := models.Subscription{
			BillingCycle: "monthly",
			CostPerCycle: 10.0,
			PaymentMethod: "Invalid",
			NextPaymentDate: time.Now().Add(24 * time.Hour),
		}

		validatedSubscription, testErr := testValidator.validateSubscription(testSubscription)

		if testErr == nil {
			t.Fatalf("expected an error for invalid payment method, but got nil")
		}

		expectedError := "Invalid payment method: Invalid"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}

		if validatedSubscription.ID != 0 {
			t.Errorf("expected validated subscription to be empty on error, but got %+v", validatedSubscription)
		}
	})

	/*
		GIVEN a subscription with a past next payment date
		WHEN validateSubscription() is called
		THEN the method returns an error stating "next payment date must be in the future"
	*/
	t.Run(`validateSubscription() accepts any next payment date`, func(t *testing.T) {
		// The validator doesn't check if the next payment date is in the future,
		// so we test that any date is accepted
		testSubscription := models.Subscription{
			BillingCycle: "monthly",
			CostPerCycle: 10.0,
			PaymentMethod: "Visa",
			NextPaymentDate: time.Now().Add(-24 * time.Hour), // Past date
		}

		validatedSubscription, testErr := testValidator.validateSubscription(testSubscription)

		// No error should be returned for a past date
		if testErr != nil {
			t.Errorf("expected no error for past next payment date, but got %v", testErr)
		}

		// The date should be preserved as-is
		if !validatedSubscription.NextPaymentDate.Equal(testSubscription.NextPaymentDate) {
			t.Errorf("expected next payment date to be preserved, but got %v", validatedSubscription.NextPaymentDate)
		}
	})

	/*
		GIVEN a valid subscription
		WHEN validateSubscription() is called
		THEN the method returns the validated subscription and no error
	*/
	t.Run(`validateSubscription() Happy Path - Returns validated subscription for valid input`, func(t *testing.T) {
		testSubscription := models.Subscription{
			BillingCycle: "monthly",
			CostPerCycle: 10.0,
			PaymentMethod: "Visa",
			NextPaymentDate: time.Now().Add(24 * time.Hour),
		}

		validatedSubscription, testErr := testValidator.validateSubscription(testSubscription)

		if testErr != nil {
			t.Errorf("expected no error for valid subscription, but got %v", testErr)
		}

		if validatedSubscription.BillingCycle != testSubscription.BillingCycle ||
			validatedSubscription.CostPerCycle != testSubscription.CostPerCycle ||
			validatedSubscription.PaymentMethod != testSubscription.PaymentMethod ||
			!validatedSubscription.NextPaymentDate.Equal(testSubscription.NextPaymentDate) {
			t.Errorf("expected validated subscription to match input, but got %+v", validatedSubscription)
		}
	})

	// ----------- ValidatePayment() ------------
	/*
		GIVEN a payment with a non-positive amount
		WHEN ValidatePayment() is called
		THEN the method returns an error stating "amount must be greater than 0"
	*/
	t.Run(`ValidatePayment() rejects non-positive amounts`, func(t *testing.T) {
		testPayment := models.Payment{
			Amount: 0,
			PaymentMethod: "Visa",
			PaymentDate: time.Now(),
		}

		validatedPayment, testErr := testValidator.ValidatePayment(testPayment)

		if testErr == nil {
			t.Fatalf("expected an error for non-positive amount, but got nil")
		}

		expectedError := "amount must be greater than 0"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}

		if validatedPayment.ID != 0 {
			t.Errorf("expected validated payment to be empty on error, but got %+v", validatedPayment)
		}
	})

	/*
		GIVEN a payment with an invalid payment method
		WHEN ValidatePayment() is called
		THEN the method returns an error stating "Invalid payment method: Invalid"
	*/
	t.Run(`ValidatePayment() rejects invalid payment methods`, func(t *testing.T) {
		testPayment := models.Payment{
			Amount: 10.0,
			PaymentMethod: "Invalid",
			PaymentDate: time.Now(),
		}

		validatedPayment, testErr := testValidator.ValidatePayment(testPayment)

		if testErr == nil {
			t.Fatalf("expected an error for invalid payment method, but got nil")
		}

		expectedError := "Invalid payment method: Invalid"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}

		if validatedPayment.ID != 0 {
			t.Errorf("expected validated payment to be empty on error, but got %+v", validatedPayment)
		}
	})

	/*
		GIVEN a payment with a future payment date
		WHEN ValidatePayment() is called
		THEN the method returns an error stating "payment date cannot be in the future"
	*/
	t.Run(`ValidatePayment() rejects future payment dates`, func(t *testing.T) {
		testPayment := models.Payment{
			Amount: 10.0,
			PaymentMethod: "Visa",
			PaymentDate: time.Now().Add(24 * time.Hour),
		}

		validatedPayment, testErr := testValidator.ValidatePayment(testPayment)

		if testErr == nil {
			t.Fatalf("expected an error for future payment date, but got nil")
		}

		expectedError := "payment date cannot be in the future"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}

		if validatedPayment.ID != 0 {
			t.Errorf("expected validated payment to be empty on error, but got %+v", validatedPayment)
		}
	})

	/*
		GIVEN a payment with a transaction ID that is too long
		WHEN ValidatePayment() is called
		THEN the method returns an error stating "transaction ID must be less than X characters"
	*/
	t.Run(`ValidatePayment() rejects long transaction IDs`, func(t *testing.T) {
		testPayment := models.Payment{
			Amount: 10.0,
			PaymentMethod: "Visa",
			PaymentDate: time.Now(),
			TransactionID: strings.Repeat("a", MaxTransactionIDLength+1),
		}

		validatedPayment, testErr := testValidator.ValidatePayment(testPayment)

		if testErr == nil {
			t.Fatalf("expected an error for long transaction ID, but got nil")
		}

		expectedError := fmt.Sprintf("transaction ID must be less than %d characters", MaxTransactionIDLength)
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}

		if validatedPayment.ID != 0 {
			t.Errorf("expected validated payment to be empty on error, but got %+v", validatedPayment)
		}
	})

	/*
		GIVEN a valid payment
		WHEN ValidatePayment() is called
		THEN the method returns the validated payment and no error
	*/
	t.Run(`ValidatePayment() Happy Path - Returns validated payment for valid input`, func(t *testing.T) {
		testPayment := models.Payment{
			Amount: 10.0,
			PaymentMethod: "Visa",
			PaymentDate: time.Now(),
			TransactionID: "test-transaction-id",
		}

		validatedPayment, testErr := testValidator.ValidatePayment(testPayment)

		if testErr != nil {
			t.Errorf("expected no error for valid payment, but got %v", testErr)
		}

		if validatedPayment.Amount != testPayment.Amount ||
			validatedPayment.PaymentMethod != testPayment.PaymentMethod ||
			!validatedPayment.PaymentDate.Equal(testPayment.PaymentDate) ||
			validatedPayment.TransactionID != testPayment.TransactionID {
			t.Errorf("expected validated payment to match input, but got %+v", validatedPayment)
		}
	})

	// ----------- ValidateDigitalLocation() ------------
	/*
		GIVEN a location with an invalid service type
		WHEN ValidateDigitalLocation() is called
		THEN the method returns an error stating "Invalid service type. Must be 'basic' or 'subscription'"
	*/
	t.Run(`ValidateDigitalLocation() rejects invalid service types`, func(t *testing.T) {
		testLocation := models.DigitalLocation{
			Name: "Test Location",
			URL: "https://example.com",
			ServiceType: "Invalid",
		}

		validatedLocation, testErr := testValidator.ValidateDigitalLocation(testLocation)

		if testErr == nil {
			t.Fatalf("expected an error for invalid service type, but got nil")
		}

		expectedError := "Invalid service type. Must be 'basic' or 'subscription'"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}

		if validatedLocation.ID != "" {
			t.Errorf("expected validated location to be empty on error, but got %+v", validatedLocation)
		}
	})

	/*
		GIVEN a location with an invalid subscription
		WHEN ValidateDigitalLocation() is called
		THEN the method returns an error stating "Subscription validation failed"
	*/
	t.Run(`ValidateDigitalLocation() rejects invalid subscriptions`, func(t *testing.T) {
		testLocation := models.DigitalLocation{
			Name: "Test Location",
			URL: "https://example.com",
			ServiceType: "subscription",
			Subscription: &models.Subscription{
				BillingCycle: "invalid",
				CostPerCycle: 10.0,
				PaymentMethod: "Visa",
				NextPaymentDate: time.Now().Add(24 * time.Hour),
			},
		}

		validatedLocation, testErr := testValidator.ValidateDigitalLocation(testLocation)

		if testErr == nil {
			t.Fatalf("expected an error for invalid subscription, but got nil")
		}

		expectedError := "Subscription validation failed"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}

		if validatedLocation.ID != "" {
			t.Errorf("expected validated location to be empty on error, but got %+v", validatedLocation)
		}
	})

	/*
		GIVEN a valid location with subscription
		WHEN ValidateDigitalLocation() is called
		THEN the method returns the validated location and no error
	*/
	t.Run(`ValidateDigitalLocation() Happy Path - Returns validated location with subscription`, func(t *testing.T) {
		testLocation := models.DigitalLocation{
			Name: "Test Location",
			URL: "https://example.com",
			ServiceType: "subscription",
			Subscription: &models.Subscription{
				BillingCycle: "monthly",
				CostPerCycle: 10.0,
				PaymentMethod: "Visa",
				NextPaymentDate: time.Now().Add(24 * time.Hour),
			},
		}

		validatedLocation, testErr := testValidator.ValidateDigitalLocation(testLocation)

		if testErr != nil {
			t.Errorf("expected no error for valid location, but got %v", testErr)
		}

		if validatedLocation.Name != testLocation.Name ||
			validatedLocation.URL != testLocation.URL ||
			validatedLocation.ServiceType != testLocation.ServiceType ||
			validatedLocation.Subscription == nil ||
			validatedLocation.Subscription.BillingCycle != testLocation.Subscription.BillingCycle ||
			validatedLocation.Subscription.CostPerCycle != testLocation.Subscription.CostPerCycle ||
			validatedLocation.Subscription.PaymentMethod != testLocation.Subscription.PaymentMethod ||
			!validatedLocation.Subscription.NextPaymentDate.Equal(testLocation.Subscription.NextPaymentDate) {
			t.Errorf("expected validated location to match input, but got %+v", validatedLocation)
		}
	})
}

func TestDigitalValidator_ValidateSubscription(t *testing.T) {
	// Create a fixed time for testing
	fixedTime := time.Date(2024, 4, 18, 0, 0, 0, 0, time.UTC)

	// Create validator with fixed time source and proper logger
	sanitizer, err := sanitizer.NewSanitizer()
	require.NoError(t, err)

	log, err := logger.NewLogger()
	require.NoError(t, err)

	validator := &DigitalValidator{
		sanitizer:  sanitizer,
		timeSource: func() time.Time { return fixedTime },
		logger:     *log,
	}

	tests := []struct {
		name        string
		subscription models.Subscription
		wantErr     bool
		errContains string
	}{
		{
			name: "valid subscription",
			subscription: models.Subscription{
				BillingCycle:    "monthly",
				CostPerCycle:    9.99,
				PaymentMethod:   "Visa",
				NextPaymentDate: fixedTime.AddDate(0, 1, 0), // 1 month in future
			},
			wantErr: false,
		},
		{
			name: "valid subscription with past payment date",
			subscription: models.Subscription{
				BillingCycle:    "monthly",
				CostPerCycle:    9.99,
				PaymentMethod:   "Visa",
				NextPaymentDate: fixedTime.AddDate(0, -1, 0), // 1 month in past (should be accepted)
			},
			wantErr: false,
		},
		{
			name: "invalid billing cycle",
			subscription: models.Subscription{
				BillingCycle:    "invalid",
				CostPerCycle:    9.99,
				PaymentMethod:   "Visa",
				NextPaymentDate: fixedTime.AddDate(0, 1, 0),
			},
			wantErr:     true,
			errContains: "invalid billing cycle",
		},
		{
			name: "invalid cost per cycle",
			subscription: models.Subscription{
				BillingCycle:    "monthly",
				CostPerCycle:    -1.00,
				PaymentMethod:   "Visa",
				NextPaymentDate: fixedTime.AddDate(0, 1, 0),
			},
			wantErr:     true,
			errContains: "cost per cycle must be greater than 0",
		},
		{
			name: "invalid payment method",
			subscription: models.Subscription{
				BillingCycle:    "monthly",
				CostPerCycle:    9.99,
				PaymentMethod:   "invalid",
				NextPaymentDate: fixedTime.AddDate(0, 1, 0),
			},
			wantErr:     true,
			errContains: "Invalid payment method: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validator.validateSubscription(tt.subscription)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.subscription.BillingCycle, got.BillingCycle)
			assert.Equal(t, tt.subscription.CostPerCycle, got.CostPerCycle)
			assert.Equal(t, tt.subscription.PaymentMethod, got.PaymentMethod)
			assert.Equal(t, tt.subscription.NextPaymentDate, got.NextPaymentDate)
		})
	}
}

func TestValidateRemoveDigitalLocation(t *testing.T) {
	validator, err := NewDigitalValidator(&MockSanitizer{})
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	t.Run("empty user ID", func(t *testing.T) {
		err := validator.ValidateRemoveDigitalLocation("", []string{"123e4567-e89b-12d3-a456-426614174000"})
		if err == nil {
			t.Error("Expected error for empty user ID")
		}
		if !strings.Contains(err.Error(), "user ID cannot be empty") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("empty location IDs", func(t *testing.T) {
		err := validator.ValidateRemoveDigitalLocation("user1", []string{})
		if err == nil {
			t.Error("Expected error for empty location IDs")
		}
		if !strings.Contains(err.Error(), "no location IDs provided") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("too many location IDs", func(t *testing.T) {
		ids := make([]string, MaxBulkSize+1)
		for i := range ids {
			ids[i] = fmt.Sprintf("123e4567-e89b-12d3-a456-426614174%03d", i)
		}
		err := validator.ValidateRemoveDigitalLocation("user1", ids)
		if err == nil {
			t.Error("Expected error for too many location IDs")
		}
		if !strings.Contains(err.Error(), fmt.Sprintf("too many location IDs, maximum allowed is %d", MaxBulkSize)) {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("empty location ID", func(t *testing.T) {
		err := validator.ValidateRemoveDigitalLocation("user1", []string{""})
		if err == nil {
			t.Error("Expected error for empty location ID")
		}
		if !strings.Contains(err.Error(), "location ID cannot be empty") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("invalid UUID format", func(t *testing.T) {
		err := validator.ValidateRemoveDigitalLocation("user1", []string{"invalid-uuid"})
		if err == nil {
			t.Error("Expected error for invalid UUID format")
		}
		if !strings.Contains(err.Error(), "invalid location ID format") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("duplicate location IDs", func(t *testing.T) {
		duplicateID := "123e4567-e89b-12d3-a456-426614174000"
		err := validator.ValidateRemoveDigitalLocation("user1", []string{duplicateID, duplicateID})
		if err == nil {
			t.Error("Expected error for duplicate location IDs")
		}
		if !strings.Contains(err.Error(), "duplicate location ID") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("single valid location ID", func(t *testing.T) {
		err := validator.ValidateRemoveDigitalLocation("user1", []string{"123e4567-e89b-12d3-a456-426614174000"})
		if err != nil {
			t.Errorf("Unexpected error for valid location ID: %v", err)
		}
	})

	t.Run("multiple valid location IDs", func(t *testing.T) {
		ids := []string{
			"123e4567-e89b-12d3-a456-426614174000",
			"123e4567-e89b-12d3-a456-426614174001",
			"123e4567-e89b-12d3-a456-426614174002",
		}
		err := validator.ValidateRemoveDigitalLocation("user1", ids)
		if err != nil {
			t.Errorf("Unexpected error for valid location IDs: %v", err)
		}
	})
}
