package digital

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
)

/*
Behavior:
1. Main method: ValidateDigitalLocation
2. Helper methods:
   - validateName
   - validateURL

validateName:
- Ensure name is not empty
- Ensure name is not longer than 100 chars
- Ensure name is sanitized

validateURL:
- Ensure URL is not empty
- Ensure URL is a valid format
- Ensure URL is sanitized

Scenarios:
Reject locations with:
- Empty name
- Names longer than 100 chars
- Empty URL
- Invalid URL format
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


	// ----------- ValidateDigitalLocation() ------------
	/*
		GIVEN a location with multiple validation issues
		WHEN ValidateDigitalLocation() is called
		THEN the method returns an error with all violations
	*/
	t.Run(`ValidateDigitalLocation() - Collects errors when multiple violations are present`, func(t *testing.T) {
		testLocation := models.DigitalLocation{
			Name:    "",
			URL:    "not-a-valid-url",
		}

		validatedLocation, testErr := testValidator.ValidateDigitalLocation(testLocation)

		if testErr == nil {
			t.Errorf("expected an error for location with multiple issues, but got nil")
		}

		// Check that the error contains all expected validation messages
		errorStr := testErr.Error()
		expectedErrors := []string{
			"name cannot be empty",
			"invalid URL format",
		}

		for _, expectedError := range expectedErrors {
			if !strings.Contains(errorStr, expectedError) {
				t.Errorf("expected error to contain %q, but got %q", expectedError, errorStr)
			}
		}

		// Check that the returned location is empty
		if validatedLocation.Name != "" || validatedLocation.URL != "" {
			t.Errorf("expected validated location to be empty on error, but got %+v", validatedLocation)
		}
	})


	/*
		GIVEN a valid location
		WHEN ValidateDigitalLocation() is called
		THEN the method returns the validated location and no error
	*/
	t.Run(`ValidateDigitalLocation() - Happy Path - Passes validation for valid location`, func(t *testing.T) {
		testLocation := models.DigitalLocation{
			Name:     "My Website",
			URL:      "https://example.com",
			IsActive: true,
		}

		validatedLocation, testErr := testValidator.ValidateDigitalLocation(testLocation)

		if testErr != nil {
			t.Errorf("expected no error for valid location, but got %v", testErr)
		}

		// Check that the validated location matches the input
		if validatedLocation.Name != testLocation.Name ||
			validatedLocation.URL != testLocation.URL ||
			validatedLocation.IsActive != testLocation.IsActive {
			t.Errorf("expected validated location to match input, but got %+v", validatedLocation)
		}
	})
}
