package connectionutil

import (
	"net"
	"testing"
	"time"
)

const (
	localhost = "127.0.0.1:0"
)

/*
	Behaviors:
		- IsOnline returns true if the network connection is successful
		- IsOnline returns false if the network connection is not successful
		- IsOnline returns false if the network connection is not successful

	Scenarios:
		- If the the local server is listening on a port, then IsOnline should return true
		- If no local server is listening on a port, then IsOnline should return false
*/

func TestIsOnline(t *testing.T) {
	testCases := []struct {
		name             string
		description      string
		host             string
		timeout          time.Duration
		isNetConnected   bool
		setup            func(t *testing.T) (port int, cleanup func())
	}{
		{
			name: "Successful connection",
			description: `
				GIVEN a local server is listening,
				WHEN checking the connection,
				THEN IsOnline returns true.
			`,
			host:       "127.0.0.1",
			timeout:    1 * time.Second,
			isNetConnected: true,
			setup: func(t *testing.T) (int, func()) {
				listener, err := net.Listen("tcp", "127.0.0.1:0")
				if err != nil {
					t.Fatalf("Failed to start listener: %v", err)
				}
				port := listener.Addr().(*net.TCPAddr).Port

				// Create channels to ensure accept goroutine is ready
				acceptReady := make(chan struct{})
				accepted := make(chan struct{})

				go func() {
					close(acceptReady) // Signal that we're about to block on Accept().
					conn, err := listener.Accept()
					if err == nil && conn != nil {
						conn.Close()
					}
					close(accepted) // Explicitly signal that we've accepted a connection.
				}()

				// Wait until the goroutine signals it's ready to accept.
				<-acceptReady

				cleanup := func() {
					listener.Close()
					<-accepted // Wait for accept goroutine to finish
				}
				return port, cleanup
			},
		},
		{
			name: "Failed connection",
			description: `
				GIVEN no local server is listening on a chosen port,
				WHEN checking the internet connection
				THEN IsOnline should return false
			`,
			host:               localhost,
			timeout:            100 * time.Millisecond,
			isNetConnected:     false,
			setup: func(t *testing.T) (int, func()) {
				// Use a random unused port
				return 55555, func() {}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Create port + cleanup from setup function
			port, cleanup := testCase.setup(t)
			defer cleanup()

			// Define actual network connection status
			actualNetworkConnection := IsOnline(testCase.host, port, testCase.timeout)

			// If the actual network connection status it not true, then test should fail
			if actualNetworkConnection != testCase.isNetConnected {
				t.Errorf(
					"Expected network connection to be %v, but got %v",
					testCase.isNetConnected,
					actualNetworkConnection,
				)
			}
		})
	}
}
