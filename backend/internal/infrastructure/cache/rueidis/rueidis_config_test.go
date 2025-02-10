package cache

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

/*
  Behaviors:

	LoadFromEnv() method
		1. Loading stuff from the environment variables
		2. Converting environment variables to the correct types
		3. Storing these values in the config struct

		* GIVEN no REDIS_HOST when LoadFromEnv() is called then error
		* GIVEN no REDIS_PORT when LoadFromEnv() is called then error
		* GIVEN non-numeric string for REDIS_PORT WHEN LoadFromEnv() is called then error
		* GIVEN invalid duration strings for timeout settings WHEN loaded, then error


	Validate() method
		1. Checks that the loaded config is within expected bounds.
			(Throw errors when):
			* GIVEN an empty variable WHEN Validate() is called THEN error
			* GIVEN a port value out of range WHEN Validate() is called THEN error
			* GIVEN a non-positive timeout value or pool size WHEN Validate() is called THEN error
			* GIVEN all correct values WHEN Validate() is called THEN no error, gtg, yay

	GetRueidisOptions() method
		1. Converts config into a rueidis.ClientOption struct used to create the Redis client
			* GIVEN a valid config WHEN GetRueidisOptions() is called THEN the returned structure should have the right values
*/

func TestRueidisConfigBehaviors(t *testing.T) {
	// --------- Test LoadFromEnv() ---------
	t.Run(
		`LoadFromEnv() behaviors:`,
		func(t *testing.T) {
			// GIVEN no REDIS_HOST WHEN LoadFromEnv() is called THEN throw an error
			t.Run(
				`Missing REDIS_HOST`,
				func(t *testing.T) {
					t.Setenv(REDIS_HOST_ENV_VAR, "")

					// We need to set a valid port even if we're testing a missing REDIS_HOST
					t.Setenv(REDIS_PORT_ENV_VAR, "6379")

					testConfig := NewRueidisConfig()
					testError := testConfig.LoadFromEnv()
					expectedErr := "REDIS_HOST environment variable is required"
					if testError == nil {
						t.Fatalf("expected to see an error for missing REDIS_HOST, instead got: %v", testError)
					}
					if strings.TrimSpace(testError.Error()) != expectedErr {
						t.Fatalf("expected error %q for missing REDIS_HOST, got %q (len(expected)=%d, len(actual)=%d)", expectedErr, testError.Error(), len(expectedErr), len(testError.Error()))
					}

				},
			)

			// GIVEN no REDIS_PORT WHEN LoadFromEnv() is called THEN throw an error
			t.Run(
				`Missing REDIS_PORT`,
				func(t *testing.T) {
					t.Setenv(REDIS_HOST_ENV_VAR, "127.0.0.1")
					t.Setenv(REDIS_PORT_ENV_VAR, "")

					testConfig := NewRueidisConfig()
					testError := testConfig.LoadFromEnv()
					if testError == nil || testError.Error() != "REDIS_PORT environment variable is required" {
						t.Fatalf("expected to see an error for missing REDIS_PORT, instead got: %v", testError)
					}
				},
			)

			// GIVEN non-numeric REDIS_PORT WHEN LoadFromEnv() is called THEN throw an error
			t.Run(
				`Non-numeric REDIS_PORT`,
				func(t *testing.T) {
					t.Setenv(REDIS_HOST_ENV_VAR, "127.0.0.1")
					t.Setenv(REDIS_PORT_ENV_VAR, "i-am-not-a-number")

					testConfig := NewRueidisConfig()
					testErr := testConfig.LoadFromEnv()
					if testErr == nil {
						t.Fatalf("expected to see an error for a non-numeric REDIS_PORT, but got nil")
					}
				},
			)

			// GIVEN invalid duration strings for timeout settings WHEN loaded, THEN throw an error
			t.Run(
				`Invalid timeout duration`,
				func(t *testing.T) {
					t.Setenv(REDIS_HOST_ENV_VAR, "127.0.0.1")
					t.Setenv(REDIS_PORT_ENV_VAR, "6379")

					// Set an invalid value for the write timeout
					t.Setenv(REDIS_CONN_WRITE_TIMEOUT_ENV_VAR, "i-am-not-a-time-duration")
					testConfig := NewRueidisConfig()
					testErr := testConfig.LoadFromEnv()
					if testErr == nil {
						t.Fatal("expected to see an error for an invalid timeout duration, but got nil")
					}
				},
			)

			// GIVEN REDIS_DB is set to an invalid value WHEN LoadFromEnv() is called THEN throw an error
			t.Run(
				`REDIS_DB invalid; out of range`,
				func(t *testing.T) {
					t.Setenv(REDIS_HOST_ENV_VAR, "127.0.0.1")
					t.Setenv(REDIS_PORT_ENV_VAR, "6379")

					// Set an invalid DB number such as 9999 which is out of the default range for Redis DBs
					t.Setenv(REDIS_DB_ENV_VAR, "9999")
					testConfig := NewRueidisConfig()
					testErr := testConfig.LoadFromEnv()
					expectedMsg := "REDIS_DB must be between 0 and " + strconv.Itoa(MAX_REDIS_DB)
					if testErr == nil || testErr.Error() != expectedMsg {
						t.Fatalf("expected to see an error for an invalid DB number: %q, but instead got: %v", expectedMsg, testErr)
					}
				},
			)

		},
	)

	// --------- Test Validate() ---------
	t.Run(
		`Validate() behaviors:`,
		func(t *testing.T) {
			// GIVEN an empty host WHEN Validate() method is called, THEN throw an error
			t.Run(
				`Empty host`,
				func(t * testing.T) {
					testConfig := NewRueidisConfig()
					testConfig.Host = ""
					testConfig.Port = 6379
					testErr := testConfig.Validate()
					if testErr == nil || testErr.Error() != "redis host cannot be empty" {
						t.Fatalf("expected to see an error for an empty host, but instead got: %v", testErr)
					}
				},
			)

			// GIVEN a port value out of range WHEN Validate() method is called, THEN throw an error
			t.Run(
				`Port out of range`,
				func(t *testing.T) {
					testConfig := NewRueidisConfig()
					testConfig.Host = "127.0.0.1"
					testConfig.Port = 70000
					testErr := testConfig.Validate()
					expectedErr := fmt.Sprintf("invalid Redis port: %d", testConfig.Port)
					if testErr == nil || testErr.Error() != expectedErr {
						t.Fatalf("expected to see an error for a port out of range, but instead got: %v", testErr)
					}
				},
			)

			t.Run(
				`Non-positive timeout/pool size`,
				func(t *testing.T) {
					testConfig := NewRueidisConfig()
					testConfig.Host = "127.0.0.1"
					testConfig.Port = 6379

					// Set a non-positive value for the read timeout
					testConfig.TimeoutConfig.Read = 0
					testErr := testConfig.Validate()
					if testErr == nil || testErr.Error() != "TimeoutConfig.Read must be positive" {
						t.Fatalf("expected to see an error for a non-positive read timeout, but instead got: %v", testErr)
					}
				},
			)

			// GIVEN all correct values WHEN Validate() method is called, THEN no error is returned
			t.Run(
				`Valid configuration, no problems expected`,
				func(t *testing.T) {
					testConfig := NewRueidisConfig()
					testConfig.Host = "127.0.0.1"
					testConfig.Port = 6379

					testErr := testConfig.Validate()
					if testErr != nil {
						t.Fatalf("expected there to be no errors but instead got: %v", testErr)
					}
				},
			)
		},
	)

	// --------- Test GetRueidisOptions() ---------
	t.Run(
		`GetRueidisOptions() behaviors:`,
		func(t *testing.T) {
			testConfig := NewRueidisConfig()
			testConfig.Host = "192.168.0.1"
			testConfig.Port = 6380
			testConfig.Password = os.Getenv("TEST_REDIS_PASSWORD")
			testConfig.DB = 3
			testConfig.BlockingOpsPoolSize = 15
			testConfig.ConnWriteTimeout = 250 * time.Millisecond

			testOptions := testConfig.GetRueidisOptions()
			expectedAddress := fmt.Sprintf("%s:%d", testConfig.Host, testConfig.Port)

			if len(testOptions.InitAddress) == 0 || testOptions.InitAddress[0] != expectedAddress {
				t.Errorf("expected the InitAddress to be: %q, but instead got: %v", expectedAddress, testOptions.InitAddress)
			}
			if testOptions.Password != testConfig.Password {
				t.Errorf("expected the Password to be: %q, but instead got: %q", testConfig.Password, testOptions.Password)
			}
			if testOptions.SelectDB != testConfig.DB {
				t.Errorf("expected the SelectDB to be: %d, but instead got: %d,", testConfig.DB, testOptions.SelectDB)
			}
			if testOptions.Dialer.Timeout != testConfig.ConnWriteTimeout {
				t.Errorf("expected the Dialer.Timeout to be:%v, but instead got: %v", testConfig.ConnWriteTimeout, testOptions.Dialer.Timeout)
			}
			if testOptions.BlockingPoolSize != testConfig.BlockingOpsPoolSize {
				t.Errorf("expected the BlockingPoolSize to be: %d, but instead got: %d", testConfig.BlockingOpsPoolSize, testOptions.BlockingPoolSize)
			}
		},
	)
}