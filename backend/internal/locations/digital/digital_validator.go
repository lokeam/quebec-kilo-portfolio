package digital

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/shared/logger"
	validationErrors "github.com/lokeam/qko-beta/internal/shared/validation"
)

const (
	MaxNameLength = 100
	MaxURLLength = 2048
	MaxTransactionIDLength = 100
	MaxCostPerCycle = 1000000.0 // Maximum allowed cost per cycle
)

// Valid service types from the catalog
var ValidServiceTypes = map[string]bool{
	"basic": true,
	"subscription": true,
}

// Valid payment methods
var ValidPaymentMethods = map[string]bool{
	"Alipay": true,
	"Amex": true,
	"Diners": true,
	"Discover": true,
	"Elo": true,
	"Generic": true,
	"Hiper": true,
	"Hipercard": true,
	"Jcb": true,
	"Maestro": true,
	"Mastercard": true,
	"Mir": true,
	"Paypal": true,
	"Unionpay": true,
	"Visa": true,
}

// Valid billing cycles
var ValidBillingCycles = map[string]bool{
	"monthly": true,
	"quarterly": true,
	"bi-annually": true,
	"annually": true,
}

type DigitalValidator struct {
	sanitizer interfaces.Sanitizer
	timeSource func() time.Time
	logger    logger.Logger
}

func NewDigitalValidator(sanitizer interfaces.Sanitizer) (*DigitalValidator, error) {
	log, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &DigitalValidator{
		sanitizer:  sanitizer,
		timeSource: time.Now,
		logger:     *log,
	}, nil
}

func (v *DigitalValidator) ValidateDigitalLocation(
	location models.DigitalLocation,
) (models.DigitalLocation, error) {
	v.logger.Debug("Validating digital location", map[string]any{
		"location": location,
		"incoming_is_active": location.IsActive,
		"service_type": location.ServiceType,
		"has_subscription": location.Subscription != nil,
	})

	var validatedLocation models.DigitalLocation
	var violations []string

	// Copy ID and user ID first - these are required and don't need validation
	validatedLocation.ID = location.ID
	validatedLocation.UserID = location.UserID

	// IMPORTANT: Preserve is_active flag
	validatedLocation.IsActive = location.IsActive

	// Validate name
	if sanitizedName, err := v.validateName(location.Name); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedLocation.Name = sanitizedName
	}

	// Validate URL
	if sanitizedURL, err := v.validateURL(location.URL); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedLocation.URL = sanitizedURL
	}

	// Validate service type - simple check, no auto-correction
	if location.ServiceType != "basic" && location.ServiceType != "subscription" {
			violations = append(violations, "Invalid service type. Must be 'basic' or 'subscription'")
	} else {
			validatedLocation.ServiceType = location.ServiceType

			// SIMPLE validation: If it's subscription type, it MUST have subscription data
			if location.ServiceType == "subscription" && location.Subscription == nil {
					violations = append(violations, "Subscription service must have subscription data")
			}

			// SIMPLE validation: If it's not subscription type, it must NOT have subscription data
			if location.ServiceType != "subscription" && location.Subscription != nil {
					violations = append(violations, "Non-subscription service cannot have subscription data")
			}
	}

	// Validate subscription if present (only for subscription type)
	if location.Subscription != nil && location.ServiceType == "subscription" {
			if validatedSubscription, err := v.validateSubscription(*location.Subscription); err != nil {
					violations = append(violations, err.Error())
			} else {
					validatedLocation.Subscription = &validatedSubscription
			}
	}

	// Copy other fields that don't need validation
	validatedLocation.CreatedAt = location.CreatedAt
	validatedLocation.UpdatedAt = location.UpdatedAt

	if len(violations) > 0 {
		v.logger.Debug("Validation failed", map[string]any{
			"violations": violations,
		})
		return models.DigitalLocation{}, &validationErrors.ValidationError{
			Field:   "location",
			Message: fmt.Sprintf("Digital location validation failed: %v", violations),
		}
	}

	v.logger.Debug("Validation successful", map[string]any{
		"location": validatedLocation,
		"validated_is_active": validatedLocation.IsActive,
	})
	return validatedLocation, nil
}

func (v *DigitalValidator) validateName(name string) (string, error) {
	// Check if name is empty
	if name == "" {
		return "", &validationErrors.ValidationError{
			Field:   "name",
			Message: "name cannot be empty",
		}
	}

	// Check name length
	if len(name) > MaxNameLength {
		return "", &validationErrors.ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("name must be less than %d characters", MaxNameLength),
		}
	}

	// Sanitize name
	sanitized, err := v.sanitizer.SanitizeString(name)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("invalid name content: %v", err),
		}
	}

	return sanitized, nil
}

func (v *DigitalValidator) validateURL(urlStr string) (string, error) {
	// Check URL is empty
	if urlStr == "" {
		return "", &validationErrors.ValidationError{
			Field:     "url",
			Message:   "URL cannot be empty",
		}
	}

	// Check URL length
	if len(urlStr) > MaxURLLength {
		return "", &validationErrors.ValidationError{
			Field:     "url",
			Message:   fmt.Sprintf("URL must be less than %d characters", MaxURLLength),
		}
	}

	// Validate URL format
	parsedURL, err := url.Parse(urlStr)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", &validationErrors.ValidationError{
			Field:      "url",
			Message:    "invalid URL format",
		}
	}

	// Ensure URL has http or https scheme
	if !strings.HasPrefix(parsedURL.Scheme, "http") {
		return "", &validationErrors.ValidationError{
			Field:     "url",
			Message:   "URL must use http or https protocol",
		}
	}

	// Sanitize URL
	sanitized, err := v.sanitizer.SanitizeString(urlStr)
	if err != nil {
		return "", &validationErrors.ValidationError{
			Field:     "url",
			Message:   fmt.Sprintf("invalid URL content: %v", err),
		}
	}

	return sanitized, nil
}

func (v *DigitalValidator) validateSubscription(subscription models.Subscription) (models.Subscription, error) {
	v.logger.Debug("Validating subscription details", map[string]any{
		"subscription": subscription,
	})

	var violations []string
	var validatedSubscription models.Subscription

	// Validate billing cycle
	if !ValidBillingCycles[subscription.BillingCycle] {
		violations = append(violations, fmt.Sprintf("invalid billing cycle: %s", subscription.BillingCycle))
	} else {
		validatedSubscription.BillingCycle = subscription.BillingCycle
	}

	// Validate cost per cycle
	if subscription.CostPerCycle <= 0 {
		violations = append(violations, "cost per cycle must be greater than 0")
	} else if subscription.CostPerCycle > MaxCostPerCycle {
		violations = append(violations, fmt.Sprintf("cost per cycle must be less than %.2f", MaxCostPerCycle))
	} else {
		validatedSubscription.CostPerCycle = subscription.CostPerCycle
	}

	// Validate payment method
	if !ValidPaymentMethods[subscription.PaymentMethod] {
		violations = append(violations, fmt.Sprintf("Invalid payment method: %s", subscription.PaymentMethod))
	} else {
		validatedSubscription.PaymentMethod = subscription.PaymentMethod
	}

	// Copy other fields that don't need validation
	validatedSubscription.ID = subscription.ID
	validatedSubscription.LocationID = subscription.LocationID
	validatedSubscription.CreatedAt = subscription.CreatedAt
	validatedSubscription.UpdatedAt = subscription.UpdatedAt
	validatedSubscription.NextPaymentDate = subscription.NextPaymentDate

	if len(violations) > 0 {
		v.logger.Debug("Subscription validation failed", map[string]any{
			"violations": violations,
		})
		return models.Subscription{}, &validationErrors.ValidationError{
			Field:   "subscription",
			Message: fmt.Sprintf("Subscription validation failed: %v", violations),
		}
	}

	v.logger.Debug("Subscription validation successful", map[string]any{
		"subscription": validatedSubscription,
	})
	return validatedSubscription, nil
}

func (v *DigitalValidator) ValidatePayment(payment models.Payment) (models.Payment, error) {
	var validatedPayment models.Payment
	var violations []string

	// Validate amount
	if payment.Amount <= 0 {
		violations = append(violations, "amount must be greater than 0")
	} else {
		validatedPayment.Amount = payment.Amount
	}

	// Validate payment method
	if !ValidPaymentMethods[payment.PaymentMethod] {
		violations = append(violations, fmt.Sprintf("Invalid payment method: %s", payment.PaymentMethod))
	} else {
		validatedPayment.PaymentMethod = payment.PaymentMethod
	}

	// Validate payment date
	if payment.PaymentDate.After(time.Now()) {
		violations = append(violations, "payment date cannot be in the future")
	} else {
		validatedPayment.PaymentDate = payment.PaymentDate
	}

	// Validate transaction ID if present
	if payment.TransactionID != "" {
		if len(payment.TransactionID) > MaxTransactionIDLength {
			violations = append(violations, fmt.Sprintf("transaction ID must be less than %d characters", MaxTransactionIDLength))
		} else {
			validatedPayment.TransactionID = payment.TransactionID
		}
	}

	// Copy other fields that don't need validation
	validatedPayment.ID = payment.ID
	validatedPayment.CreatedAt = payment.CreatedAt

	if len(violations) > 0 {
		return models.Payment{}, &validationErrors.ValidationError{
			Field:   "payment",
			Message: fmt.Sprintf("Payment validation failed: %v", violations),
		}
	}

	return validatedPayment, nil
}
