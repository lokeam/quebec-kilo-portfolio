package connectionutil

import (
	"context"
	"fmt"
	"net"
	"time"
)

// CheckInternetConnection attempts a TCP connection to the specified host and port
// within the given timeout. It returns true if the connection is successful, false otherwise.
//
// Parameters:
//   - host: The host to connect to.
//   - port: The port to connect to.
//   - timeout: The duration to wait for connection to be established.
//
// Returns:
//   - bool: true if the connection is successful, false otherwise.
func IsOnline(host string, port int, timeout time.Duration) bool {
	// Create a context with the specified timeout.
	// Allows for more flexible control than using net.DialTimeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	d := net.Dialer{}
	// Allow for future cancellations or deadline extensions
	conn, err := d.DialContext(
		ctx,
		"tcp",
		// Combine host + port into a proper address string, avoiding manual fmt.Sprintf("%s:%d", host, port) every time
		net.JoinHostPort(host, fmt.Sprint(port)),
	)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
