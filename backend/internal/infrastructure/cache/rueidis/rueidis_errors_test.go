package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/lokeam/qko-beta/internal/testutils"
)

/*
	Behaviors:
	NewRueidisError()
		- When constructing a custom error using this method, the error string and its unwrapped error match what is expected.

	IsRedisError()
		- Returns the expected boolean flag given varioius import errors such as nil and non-Redis errors

	ConvertRueidisError()
		- NOTE: need to simulate different branches of code that would cause this error
		  Maybe use a mock Rueidis client to test conversion logic?

		- Converts an error into one of our custom error constants based on conditions:
			* If error is "redisnil"
			* If error is a timeout
			* If client is not ready
			* If error is a generic Redis error

	Scenarios:
		NewRueidisError()
			- When creating a new error with key + error, output of Error() method must include all supplied context
			* Need to compare returned string against expected format + confirm that errors.Unwrap() finds orignal error

		IsRedisError()
			- If no error or non-Redis error is passed, return false
			- When passing an error that represents a Redis error, return true
			  NOTE: Inject or simulate behavior of rueidis.IsRedisErr -> fake error type? dep injection?

		ConvertRueidisError()
			- If error is nil, passes through as nil
			- If error is recongized as "Redis nil", returns ErrorKeyNotFound
			- If error is context.DeadlineExceeded, returns ErrorTimeout
			- If client is "not ready", returns ErrorClientNotReady
			- If error is a generic Redis error, return ErrorConnectionFaile
			- If none of these conditions above are met, return the original error

*/

func TestRueidisErrorBehaviors(t *testing.T) {
	// --------- Test NewRueidisError and Unwrap ---------
	t.Run(
		`NewRueidisError()`,
		func(t *testing.T) {
			/*
			  GIVEN an original non-Redis error AND a Redis operation (GET) AND a KEY
			  WHEN we call NewRueidisError() with test values "GET", "user:123" and the original error
				THEN the error msg should include the operation, key and original error AND errors.Unwrap(err) should return the original error
			*/
			t.Run(
				`WITH provided key`,
				func(t *testing.T) {
					originalError := errors.New("original error - GET")
					testError := NewRueidisError("GET", "user:123", originalError)

					expectedErrorMsg := fmt.Sprintf("redis %s operation failed for key '%s': '%v'", "GET", "user:123", originalError)
					if actualErrorMsg := strings.TrimSpace(testError.Error()); actualErrorMsg != expectedErrorMsg {
						t.Fatalf("expected error message to be: %q, but instead got: %q", expectedErrorMsg, actualErrorMsg)
					}
					if !errors.Is(testError, originalError) {
						t.Fatalf("expected error message to be: %v, but instead it wasn't thrown at all", originalError)
					}
				},
			)

			/*
				GIVEN an original non-Redis error AND a Redis operation (SET) with an empty key
				WHEN we call NewRueidisError() with test values "SET", "", and the original error
				THEN the error msg should the SET operation and original error (without the key), matching the format expected
			*/
			t.Run(
				`WITHOUT provided key`,
				func(t *testing.T) {
					originalError := errors.New("original error - SET")
					testError := NewRueidisError("SET", "", originalError)
					expectedErrorMsg := fmt.Sprintf("redis %s operation failed: '%v'", "SET", originalError)
					if actualErrorMsg := strings.TrimSpace(testError.Error()); actualErrorMsg != expectedErrorMsg {
						t.Fatalf("expected error message to be: %q, but instead got: %q", expectedErrorMsg, actualErrorMsg)
					}
				},
			)
		},
	)

	// --------- Test IsRedisError ---------
	t.Run(
		`IsRedisError()`,
		func(t *testing.T) {
			/*
				GIVEN a nil error
				WHEN we call IsRedisError()
				THEN IsRedisError() should return false
			*/
			t.Run(
				`Nil error returns false`,
				func(t *testing.T) {
					if IsRedisError(nil) {
						t.Fatalf("expected IsRedisError() to return false, but instead it returned true")
					}
				},
			)

			/*
			  GIVEN a generic non-Redis error
			  WHEN we pass this error into IsRedisError()
			  THEN the function should return false since it's not recognixed as a Redis error
		  */
			t.Run(
				`Generic non-Redis error returns false`,
				func(t *testing.T) {
					testError := errors.New("i am not a Redis error")
					if IsRedisError(testError) {
						t.Fatalf("expected IsRedisError(generic error) to return false, but instead it returned true")
					}
				},
			)

			/*
				GIVEN a simulated Redis error such as "ruedis: some test error" that rueidis.IsRedisErr() would interpret as a Redis error
				WHEN we pass this error into IsRedisError()
				THEN the function should return true
			*/
			t.Run(
				`Simulated Redis error returns true`,
				func(t *testing.T) {
					simulatedRedisError := &simulatedRedisError{}
					if !IsRedisError(simulatedRedisError) {
						t.Fatalf("expected IsRedisError(simulated redis error) to return true, but instead it returned false")
					}
				},
			)
		},
	)

	// --------- Test ConvertRueidisError ---------
	t.Run(
		`ConvertRueidisError()`,
		func(t *testing.T) {
			// NOTE: we need to prep a mock rueidis client for error converstion testing
			testLogger := testutils.NewTestLogger()
			mockRueidisClient := &testMockRueidisClient{ready: true, logger: testLogger}

			/*
				GIVEN a nil error
				WHEN we call ConvertRueidisError() with nil, "GET" on a mock rueidis client
				THEN the function return nil
			*/
			t.Run(
				`Nil rueidis error returns nil`,
				func(t *testing.T) {
					if mockRueidisClient.ConvertRueidisError(nil, "GET") != nil {
						t.Fatalf("expected ConvertRueidisError to return nil for a nil Redis error")
					}
				},
			)

			/*
				GIVEN an error that simulates a Redis nil scenario (example: error w/ msg: "rueidis: nil")
				WHEN we call ConvertRueidisError(nilErr, "GET")
				THEN the function should return the custom error constant ErrorKeyNotFound
			*/
			t.Run(
				`Redis nil error returns ErrorKeyNotFound`,
				func(t *testing.T) {
					testNilRedisError := &simulatedRedisNilError{}
					actualTestError := mockRueidisClient.ConvertRueidisError(testNilRedisError, "GET")
					if actualTestError != ErrorKeyNotFound {
					t.Fatalf("expected ConvertRueidisError to return ErrorKeyNotFound: %q, but instead it returned: %q",
							ErrorKeyNotFound.Error(), actualTestError.Error())
					}
				},
			)

			/*
				GIVEN a deadlined exceeded error (example: context.DeadlineExceeded)
				WHEN we call ConvertRueidisError(context.DeadlineExceeded, "SET")
				THEN the function should return the custom error constant ErrorTimeout
			*/
			t.Run(
				`DeadlineExceeded error returns ErrorTimeout`,
				func(t *testing.T) {
					actualTestError := mockRueidisClient.ConvertRueidisError(context.DeadlineExceeded, "SET")
					if !errors.Is(actualTestError, ErrorTimeout) {
						t.Fatalf("expected ConvertRueidisError to return ErrorTimeout: %q, but instead it returned: %q", ErrorTimeout.Error(), actualTestError.Error())
					}
				},
			)

			/*
				GIVEN a mock rueidis client that indicates it is not ready (example: IsReady() returns false)
				AND a generic Redis error that doesnt' match conversion rules
				WHEN calling ConvertRueidisError(genericRedisError, "DEL") on the mock rueidis client that isn't ready
				THEN the function should return the custom error constant ErrorClientNotReady
			*/
			t.Run(
				`Client not ready returns ErrorClientNotReady`,
				func(t *testing.T) {
					notReadyMockRueidisClient := &testMockRueidisClient{ready: false, logger: testLogger}
					genericRedisError := errors.New("some generic Redis error")
					actualTestError := notReadyMockRueidisClient.ConvertRueidisError(genericRedisError, "DEL")
					if actualTestError != ErrorClientNotReady {
						t.Fatalf("expected ConvertRueidisError to return ErrorClientNotReady: %q, but instead it returned: %q", ErrorClientNotReady.Error(), actualTestError.Error())
					}
				},
			)

			/*
				GIVEN a simulated Redis error (example: rueidis: some test error) that is recognized as a Redis error by IsRedisError()
				WHEN calling ConvertRueidisError(simulatedRedisError, "HSET") on a mock rueidis client that is ready
				THEN the function should return the custom error constant ErrorConnectionFailed
			*/
			t.Run(
				`Generic Redis error returns ErrorConnectionFailed`,
				func(t *testing.T) {
					simulatedRedisError := &simulatedRedisError{}
					actualTestError := mockRueidisClient.ConvertRueidisError(simulatedRedisError, "HSET")
					if actualTestError != ErrorConnectionFailed {
						t.Fatalf("expected ConvertRueidisError to return ErrorConnectionFailed: %q, but instead it returned: %q", ErrorConnectionFailed.Error(), actualTestError.Error())
					}
				},
			)

			/*
				GIVEN a generic error that isn't a Redis error
				WHEN calling ConvertRueidisError(genericError, "INCR")
				THEN the function should return the original, unchanged error
			*/
			t.Run(
				`Error that isn't a Redis error returns the original error`,
				func(t *testing.T) {
					nonRedisError := errors.New("some non-Redis error")
					actualTestError := mockRueidisClient.ConvertRueidisError(nonRedisError, "INCR")
					if actualTestError != nonRedisError {
						t.Fatalf("expected ConvertRueidisError to return the original error: %q, but instead it returned: %q", nonRedisError.Error(), actualTestError.Error())
					}
				},
			)
		},
	)
}

// Simulated Redis errors
type simulatedRedisError struct{}

func (mre *simulatedRedisError) Error() string {
	return "rueidis: simulated redis error"
}

// NOTE: need a marker method so that IsRedisError() can recognize this error as a Redis error
func (mrc *simulatedRedisError) RedisErrorMarker() {}

var errRedisNil = errors.New("redis: nil")

type simulatedRedisNilError struct{}

func (e *simulatedRedisNilError) Error() string {
	return errRedisNil.Error()
}

func (e *simulatedRedisNilError) Unwrap() error {
	return errRedisNil
}

// Test implementations for testing ConvertRueidisError
type testMockRueidisClient struct {
	ready     bool
	logger    *testutils.TestLogger
}

func (mrc *testMockRueidisClient) IsReady() bool {
	return mrc.ready
}

func (mrc *testMockRueidisClient) ConvertRueidisError(err error, rueidisOperation string) error {
	if err == nil {
			return nil
	}

	mrc.logger.Debug("redis error", map[string]any{
			"error":     err,
			"operation": rueidisOperation,
	})

	switch {
	case isRedisNil(err):  // NOTE: Must use helper to detect simulated nil errors
			return ErrorKeyNotFound
	case errors.Is(err, context.DeadlineExceeded):
			return ErrorTimeout
	case !mrc.IsReady():
			return ErrorClientNotReady
	case IsRedisError(err):
			return ErrorConnectionFailed
	default:
			return err
	}
}
