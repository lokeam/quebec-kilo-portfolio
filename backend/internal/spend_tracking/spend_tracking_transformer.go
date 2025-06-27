package spend_tracking

import (
	"fmt"
	"time"

	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

// TransformCreateRequestToModel converts a SpendTrackingRequest to a SpendTrackingOneTimePurchaseDB model
func TransformCreateRequestToModel(
	request types.SpendTrackingRequest,
	userID string,
) (models.SpendTrackingOneTimePurchaseDB, error) {
	now := time.Now()

	// Parse purchase date
	purchaseDate, err := time.Parse("2006-01-02T15:04:05Z", request.PurchaseDate)
	if err != nil {
		return models.SpendTrackingOneTimePurchaseDB{}, fmt.Errorf("invalid purchase_date format: %w", err)
	}

	// Set default values for optional fields
	isDigital := false
	if request.IsDigital != nil {
		isDigital = *request.IsDigital
	}

	isWishlisted := false
	if request.IsWishlisted != nil {
		isWishlisted = *request.IsWishlisted
	}

	// Determine media type based on spending category
	mediaType := "misc"
	if request.SpendingCategoryID > 0 {
		// This could be enhanced to map category IDs to media types
		// For now, using a simple mapping
		switch request.SpendingCategoryID {
		case 1:
			mediaType = "hardware"
		case 2:
			mediaType = "dlc"
		case 3:
			mediaType = "in_game_purchase"
		case 4:
			mediaType = "physical_game"
		case 5:
			mediaType = "digital_game"
		default:
			mediaType = "misc"
		}
	}

	return models.SpendTrackingOneTimePurchaseDB{
		UserID:            userID,
		Title:             request.Title,
		Amount:            request.Amount,
		PurchaseDate:      purchaseDate,
		PaymentMethod:     request.PaymentMethod,
		CategoryID:        request.SpendingCategoryID,
		DigitalLocationID: request.DigitalLocationID,
		IsDigital:         isDigital,
		IsWishlisted:      isWishlisted,
		MediaType:         mediaType,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}