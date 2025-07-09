package httputils

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/lokeam/qko-beta/internal/shared/core"
)

/*
	Behaviors:
	- RespondWithJSON should log the response start time and the state of the response
	- RespondWithJSON should ensure that Content-Type header is set to "application/json" if not already set
	- RespondWithJSON should write the HTTP status code if not already set
	- RespondWithJSON should encode the provided data as JSON and write it to the response
	- RespondWithJSON should log errors that occur during encoding and return them

	- RespondWithError should map error types to HTTP status codes
		Example:
		* If error is core.ErrValidation -> HTTP status code 400 (bad request)
		* If error is core.ErrAuthentication -> HTTP status code 401 (unauthorized)
	- RespondWithError should create a templated error response object with an error message and request ID
	- RespondWithError should let RespondWithJSON handle actually writing the response

	Scenarios:
*/

type errorResponse struct {
	Error      string  `json:"error"`
	RequestID  string  `json:"requestId"`
}

// Mindlessly satisfies the logger.LoggerInterface
type dummyLogger struct{}

func (d dummyLogger) Debug(msg string, fields map[string]any) {}
func (d dummyLogger) Error(msg string, fields map[string]any) {}
func (d dummyLogger) Info(msg string, fields map[string]any) {}
func (d dummyLogger) Warn(msg string, fields map[string]any) {}

func TestRespondWithJSON(t *testing.T) {
	testLogger := dummyLogger{}
	requestID := "test-abc-123"

	// Generic error
	t.Run(
		`GIVEN a generic error`,
		func(t *testing.T) {
			testResponseWriter := newTestResponseWriter()
			genericError := errors.New("generic error")

			// Sanity check: make sure the buffer is not nil
			if testResponseWriter.buf == nil {
				t.Fatal("warning - writer.buf is nil -- plz check the newTestResponseWriter constructor")
			}

			t.Run(
				`WHEN RespondWithError is called, THEN status code should be 500 and response properly encoded`,
				func(t *testing.T) {
					err := RespondWithError(testResponseWriter, testLogger, requestID, genericError)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					// Response writer hasn't been used, WriteHeader should be called
					if testResponseWriter.status != http.StatusInternalServerError {
						t.Errorf("We expected the status code to be %d, but instead it is %d", http.StatusInternalServerError, testResponseWriter.status)
					}

					// Decode JSON error response
					var errResponse errorResponse
					if err := json.Unmarshal(testResponseWriter.buf.Bytes(), &errResponse); err != nil {
						t.Fatalf("failed to decode json error response: %v", err)
					}

					if errResponse.Error != genericError.Error() {
						t.Errorf("We expected error msg to be '%s', but instead its '%s'", genericError.Error(), errResponse.Error)
					}

					if errResponse.RequestID != requestID {
						t.Errorf("We expected requestID '%s', but instead its '%s'", requestID, errResponse.RequestID)
					}
				},
 			)
		},
	)

	// Validation error
	t.Run(
		`GIVEN a validation error`,
		func (t *testing.T) {
			testResponseWriter := newTestResponseWriter()
			validationError := core.ErrValidation

			t.Run(
				"WHEN RespondWithError is called THEN status should be 400",
				func(t *testing.T) {
					err := RespondWithError(testResponseWriter, testLogger, requestID, validationError)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					if testResponseWriter.status != http.StatusBadRequest {
						t.Errorf("We expected the status code to be %d, but instead it is %d", http.StatusBadRequest, testResponseWriter.status)
					}

					// Decode JSON error response
					var errResponse errorResponse
					if err := json.Unmarshal(testResponseWriter.buf.Bytes(), &errResponse); err != nil {
						t.Fatalf("failed to decode error response: %v", err)
					}

					if errResponse.Error != validationError.Error() {
						t.Errorf("We expected the error msg to be '%s', but instead its '%s'", validationError.Error(), errResponse.Error)
					}

					if errResponse.RequestID != requestID {
						t.Errorf("We expected requestID '%s', but instead its '%s'", requestID, errResponse.RequestID)
					}
				},
			)
		},
	)

	// Authentication error
	t.Run(
		`GIVEN an authentication error`,
		func(t *testing.T) {
			testResponseWriter := newTestResponseWriter()
			authError := core.ErrAuthentication

			t.Run(
				`WHEN RespondWithError is called THEN status should be 401`,
				func(t *testing.T) {
					err := RespondWithError(testResponseWriter, testLogger, requestID, authError)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					if testResponseWriter.status != http.StatusUnauthorized {
						t.Errorf("We expected the status code to be %d, but instead its %d", http.StatusUnauthorized, testResponseWriter.status)
					}

					var errResponse errorResponse
					if err := json.Unmarshal(testResponseWriter.buf.Bytes(), &errResponse); err != nil {
						t.Fatalf("failed to decode error response: %v", err)
					}

					if errResponse.Error != authError.Error() {
						t.Errorf("We expected the error msg to be '%s', but instead its '%s'", authError.Error(), errResponse.Error)
					}

					if errResponse.RequestID != requestID {
						t.Errorf("We expected requestID to be '%s', but instead its '%s'", requestID, errResponse.RequestID)
					}
				},
			)
		},
	)

}

