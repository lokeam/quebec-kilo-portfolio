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
var IsOnline = func (host string, port int, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	d := net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", net.JoinHostPort(host, fmt.Sprint(port)))
	if err != nil {
			// Optionally log err here if a logger is available.
			return false
	}
	_ = conn.Close()
	return true
}
