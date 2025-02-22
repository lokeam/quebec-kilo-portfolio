package types

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	authenticationFailedMessage = "authentication failed"
	invalidRequestMessage       = "invalid request"
	clientErrorMessage          = "client error with status code: %d"
	serverErrorMessage          = "server error with status code: %d"

	ScenarioDefault        ErrorScenario = "default"
	ScenarioIrregularJSON  ErrorScenario = "irregular_json"
	ScenarioCircuitBreaker ErrorScenario = "circuit_breaker"
	ScenarioContextCancel  ErrorScenario = "context_cancel"
	ScenarioNon200Status   ErrorScenario = "non_200_status"
)

type ErrorScenario string

type HTTPError struct {
	StatusCode    int
	Message       string
	Err           error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %d - %s", e.StatusCode, e.Message)
}

func NewHTTPError(statusCode int, message string, err error) *HTTPError {
	return &HTTPError{
			StatusCode: statusCode,
			Message:    message,
			Err:        err,
	}
}

type DomainError struct {
	Domain string
	Err    error
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("unsupported search domain: %s", e.Domain)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

type HTTPErrorInterpreter struct {
	DefaultClientErrorMessage string
	DefaultServerErrorMessage string
	CustomHandlers            map[int]func(statusCode int) error
}

func (het *HTTPErrorInterpreter) WithCustomHandler(
	statusCode int,
	handler func(statusCode int) error,
) *HTTPErrorInterpreter {
	het.CustomHandlers[statusCode] = handler
	return het
}

func (het *HTTPErrorInterpreter) InterpretStatusCode(response *http.Response) error {
	// Check for custom handler
	if handler, exists := het.CustomHandlers[response.StatusCode]; exists {
		return handler(response.StatusCode)
	}

	switch {
	case response.StatusCode == http.StatusUnauthorized:
		return NewHTTPError(
			http.StatusUnauthorized,
			authenticationFailedMessage,
			errors.New("invalid credentials"),
		)
	case response.StatusCode == http.StatusBadRequest:
		return NewHTTPError(
			http.StatusBadRequest,
			invalidRequestMessage,
			errors.New("malformed request parameters"),
		)

	case response.StatusCode >= 400 && response.StatusCode < 500:
		return NewHTTPError(
			response.StatusCode,
			het.DefaultClientErrorMessage,
			fmt.Errorf(clientErrorMessage, response.StatusCode),
		)

	case response.StatusCode >= 500:
		return NewHTTPError(
			response.StatusCode,
			het.DefaultServerErrorMessage,
			fmt.Errorf(serverErrorMessage, response.StatusCode),
		)
	}

	return nil
}

const (
	DefaultClientErrorMessage = "An error occurred while processing your request."
	DefaultServerErrorMessage = "An internal server error occurred."
)

func NewHTTPErrorInterpreter() *HTTPErrorInterpreter {
	return &HTTPErrorInterpreter{
		DefaultClientErrorMessage: DefaultClientErrorMessage,
		DefaultServerErrorMessage: DefaultServerErrorMessage,
		CustomHandlers:            map[int]func(statusCode int) error{},
	}
}