package physical

import (
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
	1. Main method: ValidatePhysicalLocation
	2. Helper methods:
		- validateName
		- validateLabel
		- validateLocationType
		- validateMapCoordinates

	validateName:
		- Ensure name is not emptyu
		- Ensure name is not longer than 100 chars
		- Ensure name is sanitized

	validateLabel:
		- Ensure the label is not longer than 50 chars
		- Ensure label is sanitized (may be empty)

	validatelocationType:
		- Ensure location type is one of allowed values
		- Ensure location type is sanitized

	validateMapCoordinates:
		- Ensure map coordinates are in correct format (lat/long)
		- Ensure coordinates are within valid ranges (lat: -90 to 90, long: -180 to 180)
		- Ensure coordinates are sanitized

	Scenarios:
		Reject locations with:
			- Empty name
			- names longer than 100 chars
			- labels longer than 50 chars
			- invalid location types
			- invalid map coordintes
			- coordinates outside of valid ranges
		Pass validation with complete, valid location
		Collect erorrs when multiple violations are met
*/

// Ensure MockSanitizer implements interfaces.Sanitizer
var _ interfaces.Sanitizer = (*mocks.MockSanitizer)(nil)

// Valid location types
var ValidLocationsTypes = []string{
	"house",
	"apartment",
	"office",
	"warehouse",
}

func TestPhysicalValidator(t *testing.T) {
	// Setup
	testSanitizer := &mocks.MockSanitizer{}
	testCacheWrapper := mocks.DefaultPhysicalCacheWrapper()
	testLogger := testutils.NewTestLogger()
	testValidator, testErr := NewPhysicalValidator(testSanitizer, testCacheWrapper, testLogger)
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

	// ----------- validateName() ------------
	/*
		GIVEN a label longer than the max length
		WHEN validateLabel() is called
		THEN the method returns an error stating "label must be less than X characters"
	*/
	t.Run(`validateLabel() rejects labels that are too long`, func(t *testing.T) {
		testLabel := strings.Repeat("a", MaxLabelLength + 1)
		sanitizedLabel, testErr := testValidator.validateLabel(testLabel)

		if testErr == nil {
			t.Errorf("expected an error for a label that is too long, but got nil")
		}
		if sanitizedLabel != "" {
			t.Errorf("expected sanitized label to be empty on error, but got %s", sanitizedLabel)
		}

		expectedError := fmt.Sprintf("label must be less than %d characters", MaxLabelLength)
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	/*
		GIVEN an empty label
		WHEN validateLabel() is called
		THEN the method returns an empty string and no error
	*/
	t.Run(`validateLabel() Happy Path - Accepts empty label`, func(t *testing.T) {
		testLabel := ""
		sanitizedLabel, testErr := testValidator.validateLabel(testLabel)

		if testErr != nil {
			t.Errorf("expected no error for empty label, but got %v", testErr)
		}
		if sanitizedLabel != "" {
			t.Errorf("expected sanitized label to be empty, but got %s", sanitizedLabel)
		}
	})


	/*
		GIVEN a valid label
		WHEN validateLabel() is called
		THEN the method returns the sanitized label and no error
	*/
	t.Run(`validateLabel() Happy Path - Returns sanitized label for valid input`, func(t *testing.T) {
		testLabel := "Primary"
		sanitizedLabel, testErr := testValidator.validateLabel(testLabel)

		if testErr != nil {
			t.Errorf("expected no error for valid label, but got %v", testErr)
		}
		if sanitizedLabel != testLabel {
			t.Errorf("expected sanitized label to be %q, but got %q", testLabel, sanitizedLabel)
		}
	})


	// ----------- validateLocationType() ------------
	/*
		GIVEN a location with multiple validation issues
		WHEN ValidatePhysicalLocation() is called
		THEN the method returns an error with all violations
	*/
	t.Run(`ValidatePhysicalLocation() - Collects errors when multiple violations are present`, func(t *testing.T) {
		testLocation := models.PhysicalLocation{
			Name:           "", // Empty name
			Label:          strings.Repeat("a", MaxLabelLength+1), // Label too long
			LocationType:   "Invalid", // Invalid type
			MapCoordinates: "invalid-format", // Invalid coordinates
		}

		validatedLocation, testErr := testValidator.ValidatePhysicalLocation(testLocation)

		if testErr == nil {
			t.Errorf("expected an error for location with multiple issues, but got nil")
		}

		// Check that the error contains all expected validation messages
		errorStr := testErr.Error()
		expectedErrors := []string{
			"name cannot be empty",
			fmt.Sprintf("label must be less than %d characters", MaxLabelLength),
			"location type must be one of",
			"map coordinates must be in format latitude,longitude",
		}

		for _, expectedError := range expectedErrors {
			if !strings.Contains(errorStr, expectedError) {
				t.Errorf("expected error to contain %q, but got %q", expectedError, errorStr)
			}
		}

		// Check that the returned location is empty
		if validatedLocation.Name != "" || validatedLocation.Label != "" ||
			 validatedLocation.LocationType != "" || validatedLocation.MapCoordinates != "" {
			t.Errorf("expected validated location to be empty on error, but got %+v", validatedLocation)
		}
	})


	/*
		GIVEN a valid location
		WHEN ValidatePhysicalLocation() is called
		THEN the method returns the validated location and no error
	*/
	t.Run(`ValiatePhysicalLocation() - Happy Path - Passes validation for valid location`, func(t *testing.T) {
		testLocation := models.PhysicalLocation{
			Name:           "My Home",
			Label:          "Primary",
			LocationType:   "Home",
			MapCoordinates: "40.7128,-74.0060",
		}

		validatedLocation, testErr := testValidator.ValidatePhysicalLocation(testLocation)

		if testErr != nil {
			t.Errorf("expected no error for valid location, but got %v", testErr)
		}

		// Check that the validated location matches the input
		if validatedLocation.Name != testLocation.Name ||
			 validatedLocation.Label != testLocation.Label ||
			 validatedLocation.LocationType != testLocation.LocationType ||
			 validatedLocation.MapCoordinates != testLocation.MapCoordinates {
			t.Errorf("expected validated location to match input, but got %+v", validatedLocation)
		}
	})
}