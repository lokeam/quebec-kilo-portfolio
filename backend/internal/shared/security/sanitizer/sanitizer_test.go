package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	inputContainsInvalidChars = "sanitization error: input contains invalid characters"
)

/*
	Behaviors:
	- Sanitizer returns an unchanged string if no unsafe content is found
	- Sanitizer returns an error if if unsafe HTML special characters are found
	- Sanitizer returns an error if a query with a potential XSS attach is found

	Scenarios:
		- Sanitizer returns an unchanged string if no unsafe content is found
		- Sanitizer returns an error if if unsafe HTML special characters are found
		- Sanitizer returns an error if a query with a potential XSS attach is found
*/

func TestSanitizeSearchQuerySearchQuery(t *testing.T) {
	testCases := []struct{
		name                  string
		description           string
		searchQuery           string
		expectedOutput        string
		sanitizerShouldError  bool
		expectedErrorMsg      string
	}{
		{
			name: "Valid search query",
			description: `
				GIVEN a safe search query,
				WHEN the query is sanitized,
				THEN it returns the unchanged query
			`,
			searchQuery: "After fighting, everything in else in your life got the volume turned down",
			expectedOutput: "After fighting, everything in else in your life got the volume turned down",
			sanitizerShouldError: false,
		},
		{
			name: "Query with HTML Special Characters",
			description: `
				GIVEN a query containing HTML special characters,
				WHEN the query is sanitized,
				THEN it returns an error indicating invalid characters
			`,
			searchQuery: "The things you <bold>own</bold> end up owning <bold>you</bold>",
			expectedOutput: "",
			sanitizerShouldError: true,
			expectedErrorMsg: inputContainsInvalidChars,
		},
		{
			name: "Query with potential XSS attack",
			description: `
				GIVEN a query with a potential XSS attack,
				WHEN the query is sanitized,
				THEN it returns an error indicating invalid characters
			`,
			searchQuery: "<script>iAmJacksSmirking('revenge')</script>",
			expectedOutput: "",
			sanitizerShouldError: true,
			expectedErrorMsg: inputContainsInvalidChars,
		},
	}

	// Test runner loop
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Log test case description
			t.Log(testCase.description)

			// Given a new sanitizer instance
			sanitizer, err := NewSanitizer()
			require.NoError(t, err, "Failed to create sanitizer")

			// When we sanitize the provided search query
			santizerOutput, err := sanitizer.SanitizeSearchQuery(testCase.searchQuery)

			// Then check for error or sanitized output
			if testCase.sanitizerShouldError {
				require.Error(t, err, "Expected an error for input: %s", testCase.searchQuery)
				assert.Equal(t, err.Error(), testCase.expectedErrorMsg, "Error message mismatch")

			} else {
				require.NoError(t, err, "Didn't expect an error for input: %s", testCase.searchQuery)
				assert.Equal(t, testCase.expectedOutput, santizerOutput, "Sanitized output mismatch")
			}
		})
	}
}

func BenchmarkSanitizeSearchQuery(b *testing.B) {
	sanitizer, err := NewSanitizer()
	require.NoError(b, err, "Failed to create sanitizer")
	input := "Do you know what a duvee is?"

	for i := 0; i < b.N; i++ {
		sanitizer.SanitizeSearchQuery(input)
	}
}
