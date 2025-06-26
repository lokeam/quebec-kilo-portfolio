package digital

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/types"
)

// Takes the initial POST request and transforms it to a digital location model for db adapter
func TransformCreateRequestToModel(
	request types.DigitalLocationRequest,
	userID string,
) (models.DigitalLocation, error) {
	now := time.Now()
	locationUUID := uuid.New().String()

	digitalLocationModel := models.DigitalLocation{
		ID:             locationUUID,
		UserID:         userID,
		Name:           request.Name,
		IsSubscription: request.IsSubscription,
		IsActive:       request.IsActive,
		URL:            request.URL,
		PaymentMethod:  request.PaymentMethod,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Transform subscription if present
	if err := transformSubscriptionFieldToModel(
		request,
		&digitalLocationModel,
		now,
	); err != nil {
		return models.DigitalLocation{}, err
	}

	return digitalLocationModel, nil
}


// Takes the initial PUT request and transforms it to a digital location model for db adapter
func TransformUpdateRequestToModel(
	request types.DigitalLocationRequest,
	locationID string,
	existingLocation models.DigitalLocation,
) (models.DigitalLocation, error) {
	now := time.Now()

	digitalLocationModel := models.DigitalLocation{
		ID:             locationID,
		UserID:         existingLocation.UserID,
		Name:           request.Name,
		IsSubscription: request.IsSubscription,
		IsActive:       request.IsActive,
		URL:            request.URL,
		PaymentMethod:  request.PaymentMethod,
		CreatedAt:      existingLocation.CreatedAt,
		UpdatedAt:      time.Now(),
	}

	// Transform subscription if present
	if err := transformSubscriptionFieldToModel(
		request,
		&digitalLocationModel,
		now,
	); err != nil {
		return models.DigitalLocation{}, err
	}

	return digitalLocationModel, nil
}


// Converts the optional subscription field to a subscription model for db adapter consumption
// Pointer to types.DigitalLocationRequestSubscription is used to handle this optional field
func TransformSubscriptionRequestToModel(
	requestSubscription *types.DigitalLocationRequestSubscription,
	locationID string,
	now time.Time,
) (*models.Subscription, error) {

	anchorDate, err := time.Parse("2006-01-02T15:04:05Z", requestSubscription.AnchorDate)
	if err != nil {
		return nil, fmt.Errorf("invalid anchor_date format: %w", err)
	}

	return &models.Subscription{
		LocationID:     locationID,
		BillingCycle:   requestSubscription.BillingCycle,
		CostPerCycle:   requestSubscription.CostPerCycle,
		AnchorDate:     anchorDate,
		PaymentMethod:  requestSubscription.PaymentMethod,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func transformSubscriptionFieldToModel(
	request types.DigitalLocationRequest,
	digitalLocationModel *models.DigitalLocation,
	now time.Time,
) error {
	if request.IsSubscription && request.Subscription != nil {
		transformedSubscription, err := TransformSubscriptionRequestToModel(
			request.Subscription,
			digitalLocationModel.ID,
			now,
		)
		if err != nil {
			return err
		}
		digitalLocationModel.Subscription = transformedSubscription
	}

	return nil
}