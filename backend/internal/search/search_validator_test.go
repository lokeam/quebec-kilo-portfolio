package search

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/search/searchdef"
	"github.com/lokeam/qko-beta/internal/testutils/mocks"
)

/*
	Behavior:
	1 Main method: ValidateQuery
	2 Helper methods: validateQueryString, validateResultLimit

	validateQueryString:
		- Ensure the query string is at least 2 characters long
		- Ensure the query string is no longer than 100 characters
		- Ensure the query string is sanitized (we already have tests for sanitizer.go)

	validateResultLimit:
		- Ensure that the query at least returns 1 result
		- Ensure that the query does not return more than 50 results
		- Ensure that the query does not return more than 500 pages of results (max offset)

	Scenarios:
		- Reject queries that are too short
		- Reject queries that are too long
		- Return a sanitized query for a valid input
		- Fail when sanitizer fails
		- Reject a limit value that is less than the min result limit (example: 1)
		- Reject a limit value that is greater than the max result limit (example: > 50)
		- Reject an offset value that is negative
		- Reject an offset value that is greater than the max result offset (example: > 500)
		- Accept valid limits and offset values
		- Pass validation for a complete, valid query
		- Collect errors when multiple violations are present
*/

// Make sure MockSanitizer implements interfaces.Sanitizer
var _ interfaces.Sanitizer = (*mocks.MockSanitizer)(nil)

func TestSearchValidator(t *testing.T) {

	testSanitizer := &mocks.MockSanitizer{}
	testValidator, testErr := NewSearchValidator(testSanitizer)
	if testErr != nil {
		t.Fatalf("failed to create test validator: %v", testErr)
	}

	// validateQueryString() helper method
	// --------- Reject queries that are too short ---------
	t.Run(
		`validateQueryString() rejects queries that are too short`,
		func(t *testing.T) {
			/*
				GIVEN a query string that is shorter than the min query length (example: 1 char)
				WHEN validateQueryString() is called
				THEN the method returns an error stating "query must be at least 2 characters"
			*/
			testQuery := "a"
			testSanitizedQuery, testErr := testValidator.validateQueryString(testQuery)

			if testErr == nil {
				t.Errorf("expected an error for a query that is too short, but instead got nil")
			}
			if testSanitizedQuery != "" {
				t.Errorf("expected the sanitized query to be empty, but instead got %s", testSanitizedQuery)
			}

			expectedError := fmt.Sprintf("query must be at least %d characters", MinQueryLength)
			if !strings.Contains(testErr.Error(), expectedError) {
				t.Errorf("expected error to contain %q, but instead got %q", expectedError, testErr.Error())
			}
		},
	)

	// --------- Reject queries that are too long ---------
	t.Run(
		`validateQueryString() rejects queries that are too long`,
		func(t *testing.T) {
			/*
				GIVEN a query string that is longer than the max query length (example: 100 chars)
				WHEN validateQueryString(query) is called
				THEN the method returns an error stating "search query must be less than 100 characters"
			*/
			testQuery := strings.Repeat("a", MaxQueryLength + 1)
			testSanitizedQuery, testErr := testValidator.validateQueryString(testQuery)

			if testErr == nil {
				t.Errorf("expected an error for a query that is too long, but instead got nil")
			}
			if testSanitizedQuery != "" {
				t.Errorf("expected sanitized result to be empty on error, but instead got %s", testSanitizedQuery)
			}

			expectedError := fmt.Sprintf("search query must be less than %d characters", MaxQueryLength)
			if !strings.Contains(testErr.Error(), expectedError) {
				t.Errorf("expected error to contain this query: %q, but instead got %q", expectedError, testErr.Error())
			}
		},
	)

	// --------- Fail when sanitizer fails ---------
	t.Run(
		`validateQueryString() fails when sanitizer fails`,
		func(t *testing.T) {
			/*
				GIVEN a query string that is valid in length but which, passed to the sanitizer, results in an error
				WHEN validateQueryString(query) is called
				THEN the method returns an error stating "invalid query content: <error message from sanitizer>"
			*/
			testQuery := "<script>alert('trigger xss sanitizer error');</script>"
			sanitizedQuery, testErr := testValidator.validateQueryString(testQuery)

			if testErr == nil {
				t.Errorf("expected an error for a query triggers a sanitizer error, but instead got nil")
			}
			if sanitizedQuery != "" {
				t.Errorf("expected sanitized query to be empty on error, but instead got: %s", sanitizedQuery)
			}

			expectedError := "invalid query content: sanitizer failure"
			if !strings.Contains(testErr.Error(), expectedError) {
				t.Errorf("expected error to contain this query: %q, but instead got %q", expectedError, testErr.Error())
			}
		},
	)

	// --------- Happy Path:Return a sanitized query for a valid input ---------
	t.Run(
		`validateQueryString() Happy Path: Return a sanitized query for a valid input `,
		func(t *testing.T) {
			/*
				GIVEN a valid query string that meets the length criteria
				AND a sanitizer that returns a transformed (or identical) query string without error
				WHEN validateQueryString(query) is called
				THEN the method returns the sanitized query and no error
			*/
			testQuery := "a valid query"
			testSanitizedQuery, testErr := testValidator.validateQueryString(testQuery)

			if testErr != nil {
				t.Errorf("expected no error for a valid query, but instead got %v", testErr)
			}
			if testSanitizedQuery != testQuery {
				t.Errorf("expected sanitized query to be idential to the input query: %q, but instead got %q", testQuery, testSanitizedQuery)
			}
		},
	)

	// validateResultLimit() helper method
	// --------- Reject a limit value that is less than the min result limit (example: 1) ---------
	t.Run(
		`validateResultLimit() rejects limit that is below min result limit`,
		func(t *testing.T) {
			/*
				GIVEN a limit value less than the min result limit (example: 0)
				WHEN validateResultLimit(limit, offset) is called with a valid offset (example: 0)
				THEN the method returns an error stating "result limit must be between 1 and 50"
			*/
			testLimit := MinResultLimit - 1
			testOffset := 0
			testErr := testValidator.validateResultLimit(testLimit, testOffset)

			if testErr == nil {
				t.Errorf("expected an error for a limit that is below min result limit, but instead got nil")
			}
			expectedError := fmt.Sprintf("result limit must be between %d and %d", MinResultLimit, MaxResultLimit)
			if !strings.Contains(testErr.Error(), expectedError) {
				t.Errorf("expected error to contain this query: %q, but instead got %q", expectedError, testErr.Error())
			}
		},
	)

	// --------- Reject a limit value that is greater than the max result limit (example: > 50) ---------
	t.Run(
		`validateResultLimit() rejects limit that is above max result limit`,
		func (t *testing.T) {
			/*
				GIVEN a limit value greater than the max result limit (example: 51)
				WHEN validateResultLimit(limit, offset) is called
				THEN the method returns an error stating "result limit must be between 1 and 50"
			*/
			testLimit := MaxResultLimit + 1
			testOffset := 0
			testErr := testValidator.validateResultLimit(testLimit, testOffset)

			if testErr == nil {
				t.Errorf("expected an error for a limit that is above max result limit, but instead got nil")
			}
			expectedError := fmt.Sprintf("result limit must be between %d and %d", MinResultLimit, MaxResultLimit)
			if !strings.Contains(testErr.Error(), expectedError) {
				t.Errorf("expected error to contain this query: %q, but got: %q", expectedError, testErr.Error())
			}
		},
	)

	// --------- Reject an offset value that is negative ---------
	t.Run(
		`validateResultLimit() rejects offset that is negative`,
		func(t *testing.T) {
			/*
				GIVEN an offset value that is negative (example: -1)
				WHEN validateResultLimit(limit, offset) is called
				THEN the method returns an error stating "result offset must be less than 500"
			*/
			testLimit := 10
			testOffset := -1
			testErr := testValidator.validateResultLimit(testLimit, testOffset)

			if testErr == nil {
				t.Errorf("expected an error for an offset that is negative, but instead got nil")
			}
			expectedError := fmt.Sprintf("result offset must be less than %d", MaxResultOffset)
			if !strings.Contains(testErr.Error(), expectedError) {
				t.Errorf("expected error to contain this query: %q, but instead got %q", expectedError, testErr.Error())
			}
		},
	)

	// --------- Reject an offset value that is above the max result offset (example: > 500) ---------
	t.Run(
		`validateResultLimit() rejects offset that is above max result offset`,
		func(t *testing.T) {
			/*
				GIVEN an offset value greater than the max result offset (example: 501)
				WHEN validateResultLimit(limit, offset) is called
				THEN the method returns an error stating "result offset must be less than 500"
			*/
			testLimit := 10
			testOffset := MaxResultOffset + 1
			testErr := testValidator.validateResultLimit(testLimit, testOffset)

			if testErr == nil {
				t.Errorf("expected an error for an offset that is above the max result offset, but instead got nil")
			}
			expectedError := fmt.Sprintf("result offset must be less than %d", MaxResultOffset)
			if !strings.Contains(testErr.Error(), expectedError) {
				t.Errorf("expected error to contain this query: %q, but instead got %q", expectedError, testErr.Error())
			}
		},
	)

	// --------- Happy Path: Accept valid limit and offset values ---------
	t.Run(
		`validateResultLimit() Happy Path: Accept valid limit and return nil`,
		func(t *testing.T) {
			/*
				GIVEN a valid limit value between 1 and 50 and an offset between 0 and 500
				WHEN validateResultLimit(limit, offset) is called
				THEN the method returns nil (no error)
			*/
			testLimit := 10
			testOffset := 100
			testErr := testValidator.validateResultLimit(testLimit, testOffset)

			if testErr != nil {
				t.Errorf("expected no error for a valid limit and offset, but instead got %v", testErr)
			}
		},
	)

	// ValidateQuery() main method
	// --------- Collect errors when multiple violations are present ---------
	t.Run(
		`ValidateQuery() collects errors when multiple violations are present`,
		func(t *testing.T) {
			/*
				GIVEN a SearchQuery with multiple violations (example: query too short/long AND result limit too low/high)
				WHEN ValidateQuery(query) is called
				THEN the method returns an error with all the violations
			*/
			testSearchQuery := searchdef.SearchQuery{
				Query: "a", // too short (min 2 chars)
				Limit: 0,   // too low (min result limit is 1)
			}
			testErr := testValidator.ValidateQuery(testSearchQuery)

			if testErr == nil {
				t.Errorf("expected an error for a query with multiple validation errors, but instead got nil")
			} else {
				errorStr := testErr.Error()
				// Validate that the error contains BOTH expected messages
				if !strings.Contains(errorStr, fmt.Sprintf("query must be at least %d characters", MinQueryLength)) {
					t.Errorf("expected validation error to mention min query length, but instead got %s", errorStr)
				}
				if !strings.Contains(errorStr, fmt.Sprintf("result limit must be between %d and %d", MinResultLimit, MaxResultLimit)) {
					t.Errorf("expected validation error to mention max quer length, but instead got %s", errorStr)
				}
			}
		},
	)

	// --------- Happy Path: Accept valid query, limit, and offset values ---------
	t.Run(
		`ValidateQuery() Happy Path: Pass validation for a complete, valid query`,
		func(t *testing.T) {
			/*
				GIVEN a valid SearchQuery (string meets length criteria, sanitizer returns valid output, limit within range)
				WHEN ValidateQuery(query) is called
				THEN the method returns nil (no error)
			*/
			testSearchQuery := searchdef.SearchQuery{
				Query: "this is a valid query",
				Limit: 50,
			}

			testErr := testValidator.ValidateQuery(testSearchQuery)

			if testErr != nil {
				t.Errorf("expected no validation errors for a valid query, but instead got %v", testErr)
			}
		},
	)
}
