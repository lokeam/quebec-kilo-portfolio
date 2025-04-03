package sublocation

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
)

/*
	Behavior:
	1. Main method: ValidateSublocation()
	2. Helpers:
		- validateName()
		- validateLabel()
		- validateLocationType()
		- validateBgColor()
		- validateCapacity()

	validateName()
		- Ensure name is:
			- not empty
			- not longer than 100 chars
			- sanitized

	validateLocationType()
		- Ensure location type is:
			- one of allowed values
			- sanitized

	validateBgColor()
		- Ensure bgColor is:
			- one of allowed string values
			- sanitized

	validateCapacity()
		- Ensure capacity is:
			- a positive number
			- within reasonable bounds (1-1000)


	Scenarios:
		Reject sublocations with:
			- empty name
			- name longer than 100 chars
			- invalid bg colors
			- negative or 0 capacity
			- capacity grreater than max allowed
		Pass validation with complete, valid sublocation
		Collect errors when multiple violations are met
*/

// Ensure MockSanitizer implements interfaces.Sanitizer
var _ interfaces.Sanitizer = (*mocks.MockSanitizer)(nil)

func TestSublocationValidator(t *testing.T) {
	testSanitizer := &mocks.MockSanitizer{}
	testValidator, testErr := NewSublocationValidator(testSanitizer)
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
		GIVEN a sublocation with an empty name
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
		GIVEN a sublocation with a name longer than the max length
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
		testName := "Test Shelf"
		sanitizedName, testErr := testValidator.validateName(testName)

		if testErr != nil {
			t.Errorf("expected no error for valid name, but got %v", testErr)
		}
		if sanitizedName != testName {
			t.Errorf("expected sanitized name to be %q, but got %q", testName, sanitizedName)
		}
	})


	// ----------- validateLocationType() ------------
	/*
		GIVEN an invalid location type
		WHEN validateLocationType() is called
		THEN the method returns an error stating "location type must be one of [valid types]"
	*/
	t.Run(`validateLocationType() rejects invalid location types`, func(t *testing.T) {
		testType := "InvalidType"
		sanitizedType, testErr := testValidator.validateLocationType(testType)

		if testErr == nil {
			t.Errorf("expected an error for invalid location type, but got nil")
		}
		if sanitizedType != "" {
			t.Errorf("expected sanitized type to be empty on error, but got %s", sanitizedType)
		}

		expectedError := "location type must be one of"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})


	/*
		GIVEN a location type that triggers a sanitizer error
		WHEN validateLocationType() is called
		THEN the method returns an error from the sanitizer
	*/
	t.Run(`validateLocationType() fails when sanitizer fails`, func(t *testing.T) {
		testType := "<script>alert('xss');</script>"
		sanitizedType, testErr := testValidator.validateLocationType(testType)

		if testErr == nil {
			t.Errorf("expected an error when sanitizer fails, but got nil")
		}
		if sanitizedType != "" {
			t.Errorf("expected sanitized type to be empty on error, but got %s", sanitizedType)
		}

		expectedError := "invalid location type content: sanitizer failure"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})


	/*
		GIVEN a valid location type
		WHEN validateLocationType() is called
		THEN the method returns the sanitized type and no error
	*/
	t.Run(`validateLocationType() Happy Path - Returns sanitized type for valid input`, func(t *testing.T) {
		testType := "shelf"
		sanitizedType, testErr := testValidator.validateLocationType(testType)

		if testErr != nil {
			t.Errorf("expected no error for valid location type, but got %v", testErr)
		}
		if sanitizedType != testType {
			t.Errorf("expected sanitized type to be %q, but got %q", testType, sanitizedType)
		}
	})


	// ----------- validateBgColor() ------------
	/*
		GIVEN an invalid background color
		WHEN validateBgColor() is called
		THEN the method returns an error stating "background color must be one of [valid colors]"
	*/
	t.Run(`validateBgColor() rejects invalid background colors`, func(t *testing.T) {
		testColor := "not-a-color"
		sanitizedColor, testErr := testValidator.validateBgColor(testColor)

		if testErr == nil {
			t.Errorf("expected an error for invalid background color, but got nil")
		}
		if sanitizedColor != "" {
			t.Errorf("expected sanitized color to be empty on error, but got %s", sanitizedColor)
		}

		expectedError := "background color must be one of"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	/*
		GIVEN a background color that triggers a sanitizer error
		WHEN validateBgColor() is called
		THEN the method returns an error from the sanitizer
	*/
	t.Run(`validateBgColor() fails when sanitizer fails`, func(t *testing.T) {
		testColor := "<script>alert('xss');</script>"
		sanitizedColor, testErr := testValidator.validateBgColor(testColor)

		if testErr == nil {
			t.Errorf("expected an error when sanitizer fails, but got nil")
		}
		if sanitizedColor != "" {
			t.Errorf("expected sanitized color to be empty on error, but got %s", sanitizedColor)
		}

		expectedError := "invalid background color content: sanitizer failure"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	/*
		GIVEN a valid background color
		WHEN validateBgColor() is called
		THEN the method returns the sanitized color and no error
	*/
	t.Run(`validateBgColor() Happy Path - Returns sanitized color for valid input`, func(t *testing.T) {
		testColor := "blue"
		sanitizedColor, testErr := testValidator.validateBgColor(testColor)

		if testErr != nil {
			t.Errorf("expected no error for valid background color, but got %v", testErr)
		}
		if sanitizedColor != testColor {
			t.Errorf("expected sanitized color to be %q, but got %q", testColor, sanitizedColor)
		}
	})


	// ----------- validateCapacity() ------------
	/*
		GIVEN a negative capacity
		WHEN validateCapacity() is called
		THEN the method returns an error stating "capacity must be a positive number"
	*/
	t.Run(`validateCapacity() rejects negative capacity`, func(t *testing.T) {
		testCapacity := -5
		validatedCapacity, testErr := testValidator.validateCapacity(testCapacity)

		if testErr == nil {
			t.Errorf("expected an error for negative capacity, but got nil")
		}
		if validatedCapacity != 0 {
			t.Errorf("expected validated capacity to be 0 on error, but got %d", validatedCapacity)
		}

		expectedError := "capacity must be a positive number"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})

	/*
		GIVEN a zero capacity
		WHEN validateCapacity() is called
		THEN the method returns an error stating "capacity must be a positive number"
	*/
	t.Run(`validateCapacity() rejects zero capacity`, func(t *testing.T) {
		testCapacity := 0
		validatedCapacity, testErr := testValidator.validateCapacity(testCapacity)

		if testErr == nil {
			t.Errorf("expected an error for zero capacity, but got nil")
		}
		if validatedCapacity != 0 {
			t.Errorf("expected validated capacity to be 0 on error, but got %d", validatedCapacity)
		}

		expectedError := "capacity must be a positive number"
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})


	/*
		GIVEN a capacity exceeding the maximum limit
		WHEN validateCapacity() is called
		THEN the method returns an error stating "capacity must not exceed X"
	*/
	t.Run(`validateCapacity() rejects capacity exceeding maximum limit`, func(t *testing.T) {
		testCapacity := MaxCapacity + 1
		validatedCapacity, testErr := testValidator.validateCapacity(testCapacity)

		if testErr == nil {
			t.Errorf("expected an error for capacity exceeding maximum, but got nil")
		}
		if validatedCapacity != 0 {
			t.Errorf("expected validated capacity to be 0 on error, but got %d", validatedCapacity)
		}

		expectedError := fmt.Sprintf("capacity must not exceed %d", MaxCapacity)
		if !strings.Contains(testErr.Error(), expectedError) {
			t.Errorf("expected error to contain %q, but got %q", expectedError, testErr.Error())
		}
	})


	/*
		GIVEN a valid capacity
		WHEN validateCapacity() is called
		THEN the method returns the validated capacity and no error
	*/
	t.Run(`validateCapacity() Happy Path - Returns validated capacity for valid input`, func(t *testing.T) {
		testCapacity := 50
		validatedCapacity, testErr := testValidator.validateCapacity(testCapacity)

		if testErr != nil {
			t.Errorf("expected no error for valid capacity, but got %v", testErr)
		}
		if validatedCapacity != testCapacity {
			t.Errorf("expected validated capacity to be %d, but got %d", testCapacity, validatedCapacity)
		}
	})


	// ----------- ValidateSublocation() ------------
	/*
		GIVEN a sublocation with multiple validation issues
		WHEN ValidateSublocation() is called
		THEN the method returns an error with all violations
	*/
	t.Run(`ValidateSublocation() - Collects errors when multiple violations are present`, func(t *testing.T) {
		testSublocation := models.Sublocation{
			Name:         "", // Empty name
			LocationType: "Invalid", // Invalid type
			BgColor:      "not-a-color", // Invalid color
			Capacity:     -10, // Negative capacity
		}

		validatedSublocation, testErr := testValidator.ValidateSublocation(testSublocation)

		if testErr == nil {
			t.Errorf("expected an error for sublocation with multiple issues, but got nil")
		}

		// Check that the error contains all expected validation messages
		errorStr := testErr.Error()
		expectedErrors := []string{
			"name cannot be empty",
			"location type must be one of",
			"background color must be one of",
			"capacity must be a positive number",
		}

		for _, expectedError := range expectedErrors {
			if !strings.Contains(errorStr, expectedError) {
				t.Errorf("expected error to contain %q, but got %q", expectedError, errorStr)
			}
		}

		// Check that the returned sublocation is empty
		if validatedSublocation.Name != "" || validatedSublocation.LocationType != "" ||
			validatedSublocation.BgColor != "" || validatedSublocation.Capacity != 0 {
			t.Errorf("expected validated sublocation to be empty on error, but got %+v", validatedSublocation)
		}
	})


	/*
		GIVEN a valid sublocation
		WHEN ValidateSublocation() is called
		THEN the method returns the validated sublocation and no error
	*/
	t.Run(`ValidateSublocation() - Happy Path - Passes validation for valid sublocation`, func(t *testing.T) {
		testSublocation := models.Sublocation{
			Name:         "Game Shelf",
			LocationType: "shelf",
			BgColor:      "blue",
			Capacity:     50,
		}

		validatedSublocation, testErr := testValidator.ValidateSublocation(testSublocation)

		if testErr != nil {
			t.Errorf("expected no error for valid sublocation, but got %v", testErr)
		}

		// Check that the validated sublocation matches the input
		if validatedSublocation.Name != testSublocation.Name ||
			validatedSublocation.LocationType != testSublocation.LocationType ||
			validatedSublocation.BgColor != testSublocation.BgColor ||
			validatedSublocation.Capacity != testSublocation.Capacity {
			t.Errorf("expected validated sublocation to match input, but got %+v", validatedSublocation)
		}
	})
}