package digital

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
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
	MaxBulkSize = 12 // Maximum number of location IDs allowed in bulk operations
)

// Valid payment methods
var ValidPaymentMethods = map[string]bool{
	"alipay": true,
	"amex": true,
	"diners": true,
	"discover": true,
	"elo": true,
	"generic": true,
	"hiper": true,
	"hipercard": true,
	"jcb": true,
	"maestro": true,
	"mastercard": true,
	"mir": true,
	"paypal": true,
	"unionpay": true,
	"visa": true,
}

// Valid billing cycles
var ValidBillingCycles = map[string]bool{
	"1 month": true,
	"3 month": true,
	"6 month": true,
	"12 month": true,
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
		"is_subscription": location.IsSubscription,
		"has_subscription": location.Subscription != nil,
	})

	var validatedLocation models.DigitalLocation
	var violations []string

	// Copy ID and user ID first - these are required and don't need validation
	validatedLocation.ID = location.ID
	validatedLocation.UserID = location.UserID

	// IMPORTANT: Preserve is_active flag
	validatedLocation.IsActive = location.IsActive

	// IMPORTANT: Preserve is_subscription flag
	validatedLocation.IsSubscription = location.IsSubscription

	// IMPORTANT: Preserve original subscription data


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

	// First, validate payment method for ALL digital locations
	if validatedPaymentMethod, err := v.ValidatePaymentMethod(location.PaymentMethod); err != nil {
		violations = append(violations, err.Error())
	} else {
		validatedLocation.PaymentMethod = validatedPaymentMethod
	}

	if location.IsSubscription {
		// Subscription services must have subscription data
		if location.Subscription == nil {
			violations = append(violations, "Subscription service must have subscription data")
		} else {
			// Validate subscription data IN PLACE (don't create new struct)
			if err := v.validateSubscriptionFields(location.Subscription); err != nil {
				violations = append(violations, err.Error())
			} else {
				// Keep the original subscription data (preserves AnchorDate)
				validatedLocation.Subscription = location.Subscription
			}
		}
	} else {
		// Non-subscription service must NOT have subscription data
		if location.Subscription != nil {
			violations = append(violations, "Non-subscription service cannot have subscription data")
		}
	}

	// Copy other fields that don't need validation
	validatedLocation.CreatedAt = location.CreatedAt
	validatedLocation.UpdatedAt = location.UpdatedAt

	if len(violations) > 0 {
		v.logger.Debug("Digital location validation failed", map[string]any{
			"violations": violations,
		})
		return models.DigitalLocation{}, &validationErrors.ValidationError{
			Field:   "digital_location",
			Message: fmt.Sprintf("Digital location validation failed: %v", violations),
		}
	}

	v.logger.Debug("Digital location validation successful", map[string]any{
		"location": validatedLocation,
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

func (v *DigitalValidator) ValidatePaymentMethod(paymentMethod string) (string, error) {
	if paymentMethod == "" {
		return "", &validationErrors.ValidationError{
			Field:   "payment_method",
			Message: "Payment method is required",
		}
	}

	// Convert to lowercase for validation
	paymentMethod = strings.ToLower(paymentMethod)

	if !ValidPaymentMethods[paymentMethod] {
		return "", &validationErrors.ValidationError{
			Field:   "payment_method",
			Message: fmt.Sprintf("Invalid payment method: %s", paymentMethod),
		}
	}

	return paymentMethod, nil
}

func (v *DigitalValidator) ValidateDigitalLocationsBulk(
	locations []models.DigitalLocation,
) ([]models.DigitalLocation, error) {
	v.logger.Debug("Validating digital locations in bulk", map[string]any{
		"count": len(locations),
	})

	if len(locations) == 0 {
		return nil, &validationErrors.ValidationError{
			Field:   "locations",
			Message: "no locations provided for validation",
		}
	}

	var validatedLocations []models.DigitalLocation
	var violations []string

	for i := 0; i < len(locations); i++ {
		location := locations[i]

		v.logger.Debug("Validating location in bulk", map[string]any{
			"index": i,
			"location_id": location.ID,
		})

		validatedLocation, err := v.ValidateDigitalLocation(location)
		if err != nil {
			// Add the error to violations with index information
			if validationErr, ok := err.(*validationErrors.ValidationError); ok {
				violations = append(violations, fmt.Sprintf("location[%d]: %s", i, validationErr.Message))
			} else {
				violations = append(violations, fmt.Sprintf("location[%d]: %v", i, err))
			}
			continue
		}

		validatedLocations = append(validatedLocations, validatedLocation)
	}

	if len(violations) > 0 {
		v.logger.Debug("Bulk validation failed", map[string]any{
			"violations": violations,
			"validated_count": len(validatedLocations),
			"total_count": len(locations),
		})

		// If no locations were validated successfully, return an error
		if len(validatedLocations) == 0 {
			return nil, &validationErrors.ValidationError{
				Field:   "locations",
				Message: fmt.Sprintf("all locations failed validation: %v", violations),
			}
		}

		// If some locations were validated successfully, return them with a warning
		v.logger.Warn("Partial validation success", map[string]any{
			"validated_count": len(validatedLocations),
			"total_count": len(locations),
			"violations": violations,
		})
	}

	v.logger.Debug("Bulk validation completed", map[string]any{
		"validated_count": len(validatedLocations),
		"total_count": len(locations),
	})

	return validatedLocations, nil
}

// ValidateRemoveDigitalLocation validates location IDs for deletion operations.
// It handles both single and bulk deletion scenarios.
func (v *DigitalValidator) ValidateRemoveDigitalLocation(
	userID string,
	locationIDs []string,
) ([]string, error) {
	v.logger.Debug("Validating location IDs for deletion", map[string]any{
		"userID": userID,
		"locationIDsCount": len(locationIDs),
		"isBulkOperation": len(locationIDs) > 1,
	})

	// Validate userID
	if userID == "" {
		return nil, &validationErrors.ValidationError{
			Field:   "userID",
			Message: "user ID cannot be empty",
		}
	}

	// Validate locationIDs array
	if len(locationIDs) == 0 {
		return nil, &validationErrors.ValidationError{
			Field:   "locationIDs",
			Message: "no location IDs provided for deletion",
		}
	}

	// Enforce maximum number of IDs in bulk operations
	if len(locationIDs) > MaxBulkSize {
		return nil, &validationErrors.ValidationError{
			Field:   "locationIDs",
			Message: fmt.Sprintf("too many location IDs, maximum allowed is %d", MaxBulkSize),
		}
	}

	// Sanitize and deduplicate location IDs
	locationSeen := make(map[string]bool)
	var validatedIDs []string
	for i := 0; i < len(locationIDs); i++ {
		// Skip empty IDs
		if locationIDs[i] == "" {
			continue
		}

		// Sanitize ID
		sanitizedID, err := v.sanitizer.SanitizeString(locationIDs[i])
		if err != nil {
			return nil, &validationErrors.ValidationError{
				Field:   "locationIDs",
				Message: fmt.Sprintf("invalid location ID: %s", locationIDs[i]),
			}
		}

		// Deduplicate
		if !locationSeen[sanitizedID] {
			locationSeen[sanitizedID] = true
			validatedIDs = append(validatedIDs, sanitizedID)
		}
	}

	if len(validatedIDs) == 0 {
    return nil, &validationErrors.ValidationError{
        Field:   "locationIDs",
        Message: "no valid location IDs after sanitization and deduplication",
    }
}

	// Validate that each sanitized ID is a valid UUID
	for i := 0; i < len(validatedIDs); i++ {
		if _, err := uuid.Parse(validatedIDs[i]); err != nil {
			return nil, &validationErrors.ValidationError{
				Field:   "locationIDs",
				Message: fmt.Sprintf("failed to sanitize location ID: %s", validatedIDs[i]),
			}
		}
	}

	// Log if dupes were removed
	if len(validatedIDs) < len(locationIDs) {
		v.logger.Debug("Duplicate location IDs removed", map[string]any{
			"original_count": len(locationIDs),
			"deduplicated_count": len(validatedIDs),
		})
	}

	v.logger.Debug("Location IDs validation successful", map[string]any{
		"userID": userID,
		"count": len(validatedIDs),
	})

	return validatedIDs,nil
}


// ---- validator refactor /
func (v *DigitalValidator) validateSubscriptionFields(subscription *models.Subscription) error {
	var violations []string

	// Validate billing cycle
	if !ValidBillingCycles[subscription.BillingCycle] {
			violations = append(violations, fmt.Sprintf("invalid billing cycle: %s", subscription.BillingCycle))
	}

	// Validate cost per cycle
	if subscription.CostPerCycle <= 0 {
			violations = append(violations, "cost per cycle must be greater than 0")
	} else if subscription.CostPerCycle > MaxCostPerCycle {
			violations = append(violations, fmt.Sprintf("cost per cycle must be less than %.2f", MaxCostPerCycle))
	}

	// Validate payment method
	if subscription.PaymentMethod == "" {
			violations = append(violations, "payment method is required for subscription services")
	} else if !ValidPaymentMethods[subscription.PaymentMethod] {
			violations = append(violations, fmt.Sprintf("Invalid payment method: %s", subscription.PaymentMethod))
	}

	if len(violations) > 0 {
			return &validationErrors.ValidationError{
					Field:   "subscription",
					Message: fmt.Sprintf("Subscription validation failed: %v", violations),
			}
	}

	return nil
}