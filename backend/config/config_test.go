package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Suite(t *testing.T) {
	// Set test cases w/ BDD style naming
	testCases := []struct {
		name         string
		envVars      map[string]string
		wantConfig   *Config
		wantErr      error
		description  string
	}{
		{
			name: "Given no environment variables, When loading config, Then return default config",
			description: `
			  This test ensures that the app loads with default values when no environment variables are set.
				Crucial for dev environments and first-time setup.
			`,
			envVars: map[string]string{},
			wantConfig: &Config{
				Server: ServerConfig{
					Port: 8080,
					Host: "localhost",
				},
				Env: "dev",
			},
		},
		{
			name: "Given production environment, When loading config, Then return prod config",
			description: `
				Validates that prod environment settings are recognized and applied.
				Critical for ensuring proper production deployments.
			`,
			envVars: map[string]string{
				"ENV": "prod",
			},
			wantConfig: &Config{
				Server: ServerConfig{
					Port: 8080,
					Host: "localhost",
				},
				Env: "prod",
			},
		},
		{
			name: "Given invalid port number, When loading config, Then return error",
			description: `
			  Validates custom port and host configuration.
			`,
			envVars: map[string]string{
				"PORT": "9000",
				"HOST": "custom.host",
			},
			wantConfig: &Config{
				Server: ServerConfig{
					Port: 9000,
					Host: "custom.host",
				},
				Env: "dev",
			},
		},
		{
			name: "Given all environment variables set, When loading config, Then complete custom config",
			description: `
			  Validates that all environment variables are properly processed together.
				Tests the complete configuration parsing capability.
			`,
			envVars: map[string]string{
				"ENV": "prod",
				"PORT": "9000",
				"HOST": "custom.host",
			},
			wantConfig: &Config{
				Server: ServerConfig{
					Port: 9000,
					Host: "custom.host",
				},
				Env: "prod",
			},
		},
	}

	// Test runner loop
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// GIVEN
			// Generate clean env for each test
			t.Cleanup(func() {
				cleanupEnv(t, testCase.envVars)
			})

			// Setup test env
			if err := setupTest(t, testCase.envVars); err != nil {
				t.Fatalf("failed to setup test env: %v", err)
			}

			// WHEN
			got, err := Load()

			// THEN
			if testCase.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.wantErr.Error())
				return
			}

			// Assert success case
			require.NoError(t, err)
			assertConfigEquals(t, testCase.wantConfig, got)
		})
	}
}


/* Helper Functions */
// Load Edge Cases
func TestLoad_EdgeCases(t *testing.T) {
	edgeCases := []struct {
		name    string
		envVars map[string]string
		wantErr error
		description string
	}{
		{
			name: "Given port number below valid range (0), When loading config, Then return error",
			envVars: map[string]string{
				"PORT": "0",
			},
			wantErr: fmt.Errorf("port number must be between 1 and 65535"),
			description: "Validates minimum port boundary",
		},
		{
			name: "Given port number above valid range (65535), When loading config, Then return error",
			envVars: map[string]string{
				"PORT": "65536",
			},
			wantErr: fmt.Errorf("port number must be between 1 and 65535"),
			description: "Validates maximum port boundary",
		},
		{
			name: "Given invalid port number format, When loading config, Then return error",
			envVars: map[string]string{
				"PORT": "not-a-number",
			},
			wantErr: fmt.Errorf("invalid port number format"),
			description: "Validates port number parsing",
		},
		{
			name: "Given empty host value, When loading config, Then return error",
			envVars: map[string]string{
				"HOST": "",
			},
			wantErr: fmt.Errorf("host cannot be empty"),
			description: "Validates host value is not empty",
		},
		{
			name: "Given invalid env value, When loading config, Then return error",
			envVars: map[string]string{
				"ENV": "invalid-env",
			},
			wantErr: fmt.Errorf("invalid environment: must be one of dev, test or prod"),
			description: "Validates env value constraints",
		},
	}

	// Test runner loop
	for _, testCase := range edgeCases {
		t.Run(testCase.name, func(t *testing.T) {
			// GIVEN
			// Generate clean env for each test
			t.Cleanup(func() {
				cleanupEnv(t, testCase.envVars)
			})

			// Setup test env
			if err := setupTest(t, testCase.envVars); err != nil {
				t.Fatalf("failed to setup test env: %v", err)
			}

			// WHEN
			_, err := Load()

			// THEN
			if err == nil {
				t.Error("Expected an error but none")
				return
			}

			if testCase.wantErr != nil {
				assert.Contains(
					t,
					err.Error(),
					testCase.wantErr.Error(),
					"Error message doesn't match expected error",
				)
			}
		})
	}
}

// Setup tests, cleanup env, custom assertions
func setupTest(test *testing.T, envVars map[string]string) error {
	test.Helper()

	for envVar, value := range envVars {
		if err := os.Setenv(envVar, value); err != nil {
			test.Errorf("failed to set env var %s: %v", envVar, err)
			return err
		}
	}
	return nil
}

func cleanupEnv(test *testing.T, envVars map[string] string) {
	test.Helper()

	for envVar := range envVars {
		if err := os.Unsetenv(envVar); err != nil {
			test.Errorf("failed to unset env var %s: %v", envVar, err)
		}
	}
}

func assertConfigEquals(test *testing.T, want, got *Config) {
	test.Helper()

	assert.Equal(test, want.Env, got.Env, "Environment mismatch")
	assert.Equal(test, want.Server.Port, got.Server.Port, "Port mismatch")
	assert.Equal(test, want.Server.Host, got.Server.Host, "Host mismatch")
}

// Benchmarks
func BenchmarkLoad(benchmark *testing.B) {
	// Reset env before each benchmark
	os.Clearenv()

	benchmark.ResetTimer()
	for index := 0; index < benchmark.N; index++ {
		_, err := Load()
		if err != nil {
			benchmark.Fatalf("benchmark failed: %v", err)
		}
	}
}