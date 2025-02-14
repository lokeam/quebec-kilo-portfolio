package cache

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/lokeam/qko-beta/internal/shared/logger"
)

/*
	Behaviors:
	1. Rueidis client is able to initialize with a valid configuration
	2. Rueidis client is able to make requests to Redis
	3. Rueidis client is able to handle errors gracefully
	4. Rueidis client is able handle multiple requests concurrently

	Scenarios:
		- Valid config + ready client
		- Invalid config or Redis isn't reachable
		- Client able to make successful SET and GET requests
		- Client should be able to handle errors gracefully such as trying to GET a key that doesn't exist
		- Client should be able to convert a redis error to a regular one used by the codebase
		- Client should be able to handle timeouts by returning the appropriate error when timeouts occur
		- Client should be able to notify that it isn't ready
		- Client should either succeed or fail in an expected manner when handling multiple concurrent requests
		- Client should close gracefully without leftover resources or errors
*/

// Helper to create and return a new RueidisClient for integration testing
func getTestClient(t *testing.T) *RueidisClient {
	t.Helper()

	// NOTE: Ensure that a Redis server is running on localhost:6379
	testConfig := &RueidisConfig{
		Host: "localhost",
		Port: 6379,
		Password: os.Getenv("TEST_REDIS_PASSWORD"),
	}
	testLogger, _ := logger.NewLogger()
	testRueidisClient, err := NewRueidisClient(testConfig, testLogger)
	if err != nil {
		t.Fatalf("GIVEN a valid config, WHEN we call NewRueidisClient(), THEN we expect no error but we got this: %v", err)
	}
	if !testRueidisClient.IsReady() {
		t.Fatalf("GIVEN a valid config, WHEN we call NewRueidisClient(), THEN we expect the client should be ready")
	}

	return testRueidisClient
}

func TestRueidisClientIntegration(t *testing.T) {
	// --------- Testing Initialization + Client Readiness ---------
	t.Run(
		`Initialization and Client Readiness`,
		func(t *testing.T) {
			/*
				GIVEN a valid Redis config
				WHEN we call RueidisClient()
				THEN the client should connect to Redis and report that its ready
			*/
			testRueidisClient := getTestClient(t)
			rClientIsReady := testRueidisClient.IsReady()
			if !rClientIsReady {
				t.Fatalf("Expected client to be ready, but it isn't")
			}

			// Cleanup
			_ = testRueidisClient.Close()
		},
	)

	// --------- Testing Successful SET + GET ---------
	t.Run(
		`Successful SET + GET`,
		func(t *testing.T) {
			/*
				GIVEN an initialized client
				WHEN the client executes a SET and GET command
				THEN the GET command returns the expected value
			*/
			testRueidisClient := getTestClient(t)
			defer testRueidisClient.Close()

			testIntTestKey := "integration-test-key"
			expectedValue := "integration-test-value"

			testError := testRueidisClient.Set(context.Background(), testIntTestKey, expectedValue, 2 * time.Minute)
			if testError != nil {
				t.Fatalf("GIVEN a valid key-value, WHEN we call Set(), THEN we expect no error but instead got this: %v", testError)
			}

			actualValue, testError := testRueidisClient.Get(context.Background(), testIntTestKey)
			if testError != nil {
				t.Fatalf("GIVEN a valid Redis key, WHEN we call Get(), THEN we expect to receive no error but instead got this: %v", testError)
			}

			if actualValue != expectedValue {
				t.Fatalf("GIVEN a valid key-value, WHEN we call Set(), THEN we expect Get() to return the expected value: %q but instead got this: %q", expectedValue, actualValue)
			}
		},
	)

	// --------- GET for a key that doesn't exist ---------
	t.Run(
		`GET a key that doesn't exist`,
		func(t *testing.T) {
			/*
				GIVEN an initialized client
				WHEN the client executes a GET command for a key that doesn't exist
				THEN the client should return the error constant ErrorKeyNotFound
			*/
			testRueidisClient := getTestClient(t)
			defer testRueidisClient.Close()

			nonExistentKey := "non-existent-key"
			_, testError := testRueidisClient.Get(context.Background(), nonExistentKey)
			if testError == nil {
				t.Fatalf("GIVEN a missing Redis key, WHEN we call Get(), THEN we expect to receive ErrorKeyNotFound but got nil")
			}
			if !errors.Is(testError, ErrorKeyNotFound) {
				t.Fatalf("Expected to receive ErrorKeyNotFound: %q but instead got this: %v", ErrorKeyNotFound.Error(), testError)
			}
		},
	)

	// --------- Handle Context Deadline Exceeded Error ---------
	t.Run(
		`Context Deadline Exceeded Error`,
		func(t *testing.T) {
			/*
				GIVEN an initialized client
				WHEN the client executes a SET command with an expired context
				THEN the client should return the error constant ErrorContextDeadlineExceeded
			*/
			testRueidisClient := getTestClient(t)
			defer testRueidisClient.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Nanosecond)
			defer cancel()

			// Allow the context to expire
			time.Sleep(10 * time.Millisecond)

			testError := testRueidisClient.Set(ctx, "timeout-key", "timeout-value", 0)
			if testError == nil {
				t.Fatalf("GIVEN an expired context, WHEN we call Set(), THEN we expect to receive an ErrorTimeout error but instead got nil")
			}
			if !errors.Is(testError, ErrorTimeout) {
				t.Fatalf("Expected error to be ErrorTimeout: %q but instead got this: %v", ErrorTimeout.Error(), testError)
			}
		},
	)

	// --------- Handle Concurrent Requests ---------
	t.Run(
		`Concurrent Command Execution`,
		func(t *testing.T) {
			/*
				GIVEN an initialized client
				WHEN the client executes many commands concurrently
				THEN the client should handle the commands without errors
			*/
			testRueidisClient := getTestClient(t)
			defer testRueidisClient.Close()

			concurrentRequests := 10
			var testWaitGroup sync.WaitGroup

			for i := 0; i < concurrentRequests; i++ {
				testWaitGroup.Add(1)
				go func(i int) {
					defer testWaitGroup.Done()
					testKey := fmt.Sprintf("concurrent-key-%d", i)
					expectedValue := fmt.Sprintf("concurrent-value-%d", i)

					if testError := testRueidisClient.Set(context.Background(), testKey, expectedValue, 2 * time.Minute); testError != nil {
						t.Errorf("We called SET concurrently for key %q and expected no error, but instead got this: %v", testKey, testError)
					}
					actualValue, testError := testRueidisClient.Get(context.Background(), testKey)
					if testError != nil {
						t.Errorf("We called GET concurrently for key %q and expected no error, but instead got this: %v", testKey, testError)
					}
					if actualValue != expectedValue {
						t.Errorf("We called GET concurrently and expected to get this value: %q, but instead got this: %q", expectedValue, actualValue)
					}
				}(i)
			}
			testWaitGroup.Wait()
		},
	)

	// --------- Graceful Shutdown + Resource Cleanup ---------
	t.Run(
		`Graceful shutdown and resource cleanup`,
		func(t *testing.T) {
			/*
				GIVEN an initialized client
				WHEN executing the Close() shutdown procedure
				THEN the client should report that its no longer read``
			*/
			testRueidisClient := getTestClient(t)

			testError := testRueidisClient.Close()
			if testError != nil {
				t.Fatalf("We called Close() and expected no error, but instead got this: %v", testError)
			}
			if testRueidisClient.IsReady() {
				t.Fatalf("We called Close() and expected the client to no longer be ready, but its still reporting that its ready")
			}
		},
	)
}