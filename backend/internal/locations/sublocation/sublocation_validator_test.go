package sublocation

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
)

/*
	Behavior:
	1. Main method: ValidateSublocation()
	2. Helpers:
		- validateName()
		- validateLocationType()
		- validateStoredItems()
		- validatePhysicalLocationID()

	validateName()
		- Ensure name is:
			- not empty
			- not longer than 100 chars
			- sanitized

	validateLocationType()
		- Ensure location type is:
			- one of allowed values
			- sanitized

	validateStoredItems()
		- Ensure stored items is:
			- not negative
			- within reasonable bounds (0-1000)

	validatePhysicalLocationID()
		- Ensure physical location ID is:
			- not empty
			- valid UUID
			- sanitized

	Scenarios:
		Reject sublocations with:
			- empty name
			- name longer than 100 chars
			- invalid location types
			- negative stored items
			- stored items greater than max allowed
			- invalid physical location ID
		Pass validation with complete, valid sublocation
		Collect errors when multiple violations are met
*/

// Ensure MockSanitizer implements interfaces.Sanitizer
var _ interfaces.Sanitizer = (*mocks.MockSanitizer)(nil)

func TestSublocationValidator(t *testing.T) {
	testSanitizer := &mocks.MockSanitizer{}
	testCacheWrapper := &mocks.MockSublocationCacheWrapper{}
	testLogger := testutils.NewTestLogger()
	testValidator, testErr := NewSublocationValidator(testSanitizer, testCacheWrapper, testLogger)
	if testErr != nil {
		t.Fatalf("failed to create test validator: %v", testErr)
	}

	// Setup mock sanitizer behavior
	testSanitizer.SanitizeFunc = func(text string) (string, error) {
		if strings.Contains(text, "<script>") {
			return "", fmt.Errorf("sanitizer failure")
		}
		return text, nil
	}

	// Setup mock cache wrapper behavior
	testCacheWrapper.GetCachedSublocationsFunc = func(ctx context.Context, userID string) ([]models.Sublocation, error) {
		return []models.Sublocation{}, nil
	}

	// ----------- validateName() ------------
	t.Run(`validateName() rejects empty names`, func(t *testing.T) {
		testName := ""
		testErr := testValidator.validateName(testName)

		if testErr == nil {
			t.Fatalf("expected an error for an empty name, but got nil")
		}

		expectedError := "name must be at least 1 characters long"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	t.Run(`validateName() rejects names longer than 100 characters`, func(t *testing.T) {
		testName := strings.Repeat("a", MaxNameLength+1)
		testErr := testValidator.validateName(testName)

		if testErr == nil {
			t.Errorf("expected an error for a name that is too long, but got nil")
		}

		expectedError := fmt.Sprintf("name must not exceed %d characters", MaxNameLength)
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	t.Run(`validateName() fails when sanitizer fails`, func(t *testing.T) {
		testName := "<script>alert('xss');</script>"
		testErr := testValidator.validateName(testName)

		if testErr == nil {
			t.Errorf("expected an error when sanitizer fails, but got nil")
		}

		expectedError := "invalid name content: sanitizer failure"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	t.Run(`validateName() Happy Path - Accepts valid input`, func(t *testing.T) {
		testName := "Test Shelf"
		testErr := testValidator.validateName(testName)

		if testErr != nil {
			t.Errorf("expected no error for valid name, but got %v", testErr)
		}
	})

	// ----------- validateLocationType() ------------
	t.Run(`validateLocationType() rejects invalid location types`, func(t *testing.T) {
		testType := "InvalidType"
		testErr := testValidator.validateLocationType(testType)

		if testErr == nil {
			t.Errorf("expected an error for invalid location type, but got nil")
		}

		expectedError := "location type must be one of"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	t.Run(`validateLocationType() Happy Path - Accepts valid input`, func(t *testing.T) {
		testType := "shelf"
		testErr := testValidator.validateLocationType(testType)

		if testErr != nil {
			t.Errorf("expected no error for valid location type, but got %v", testErr)
		}
	})

	// ----------- validateStoredItems() ------------
	t.Run(`validateStoredItems() rejects negative values`, func(t *testing.T) {
		testItems := -1
		testErr := testValidator.validateStoredItems(testItems)

		if testErr == nil {
			t.Errorf("expected an error for negative stored items, but got nil")
		}

		expectedError := "stored_items cannot be negative"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	t.Run(`validateStoredItems() rejects values above max`, func(t *testing.T) {
		testItems := MaxStoredItems + 1
		testErr := testValidator.validateStoredItems(testItems)

		if testErr == nil {
			t.Errorf("expected an error for stored items above max, but got nil")
		}

		expectedError := fmt.Sprintf("stored_items must not exceed %d", MaxStoredItems)
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	t.Run(`validateStoredItems() Happy Path - Accepts valid input`, func(t *testing.T) {
		testItems := 10
		testErr := testValidator.validateStoredItems(testItems)

		if testErr != nil {
			t.Errorf("expected no error for valid stored items, but got %v", testErr)
		}
	})

	// ----------- validatePhysicalLocationID() ------------
	t.Run(`validatePhysicalLocationID() rejects empty IDs`, func(t *testing.T) {
		testID := ""
		testErr := testValidator.validatePhysicalLocationID(testID)

		if testErr == nil {
			t.Errorf("expected an error for empty physical_location_id, but got nil")
		}

		expectedError := "physical_location_id cannot be empty"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	t.Run(`validatePhysicalLocationID() rejects invalid UUIDs`, func(t *testing.T) {
		testID := "not-a-uuid"
		testErr := testValidator.validatePhysicalLocationID(testID)

		if testErr == nil {
			t.Errorf("expected an error for invalid UUID, but got nil")
		}

		expectedError := "physical_location_id must be a valid UUID"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	t.Run(`validatePhysicalLocationID() fails when sanitizer fails`, func(t *testing.T) {
		testID := "<script>alert('xss');</script>"
		testErr := testValidator.validatePhysicalLocationID(testID)

		if testErr == nil {
			t.Errorf("expected an error when sanitizer fails, but got nil")
		}

		expectedError := "invalid physical_location_id content: sanitizer failure"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	t.Run(`validatePhysicalLocationID() Happy Path - Accepts valid UUID`, func(t *testing.T) {
		testID := "123e4567-e89b-12d3-a456-426614174000" // Valid UUID format
		testErr := testValidator.validatePhysicalLocationID(testID)

		if testErr != nil {
			t.Errorf("expected no error for valid UUID, but got %v", testErr)
		}
	})

	// ----------- ValidateSublocation() ------------
	t.Run(`ValidateSublocation() - Collects errors when multiple violations are present`, func(t *testing.T) {
		testSublocation := models.Sublocation{
			UserID:             "test-user-id",
			Name:               "", // Empty name
			LocationType:       "Invalid", // Invalid type
			StoredItems:        -10, // Negative stored items
			PhysicalLocationID: "not-a-uuid", // Invalid UUID
		}

		validatedSublocation, testErr := testValidator.ValidateSublocation(testSublocation)

		if testErr == nil {
			t.Errorf("expected an error for sublocation with multiple issues, but got nil")
		}

		// Check that the error contains all expected validation messages
		errorStr := testErr.Error()
		expectedErrors := []string{
			"name must be at least 1 characters long",
			"location type must be one of",
			"stored_items cannot be negative",
			"physical_location_id must be a valid UUID",
		}

		for _, expectedError := range expectedErrors {
			if !strings.Contains(errorStr, expectedError) {
				t.Errorf("expected error to contain %q, but got %q", expectedError, errorStr)
			}
		}

		// Check that the returned sublocation is empty
		if validatedSublocation.Name != "" || validatedSublocation.LocationType != "" ||
			validatedSublocation.StoredItems != 0 {
			t.Errorf("expected validated sublocation to be empty on error, but got %+v", validatedSublocation)
		}
	})

	t.Run(`ValidateSublocation() - Happy Path - Passes validation for valid sublocation`, func(t *testing.T) {
		testSublocation := models.Sublocation{
			UserID:             "test-user-id",
			Name:               "Game Shelf",
			LocationType:       "shelf",
			StoredItems:        50,
			PhysicalLocationID: "123e4567-e89b-12d3-a456-426614174000", // Valid UUID
		}

		validatedSublocation, testErr := testValidator.ValidateSublocation(testSublocation)

		if testErr != nil {
			t.Errorf("expected no error for valid sublocation, but got %v", testErr)
		}

		// Check that the validated sublocation matches the input
		if validatedSublocation.Name != testSublocation.Name ||
			validatedSublocation.LocationType != testSublocation.LocationType ||
			validatedSublocation.StoredItems != testSublocation.StoredItems ||
			validatedSublocation.PhysicalLocationID != testSublocation.PhysicalLocationID ||
			validatedSublocation.UserID != testSublocation.UserID {
			t.Errorf("expected validated sublocation to match input, but got %+v", validatedSublocation)
		}
	})
}