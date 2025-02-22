package interfaces

import "net/http"

type HTTPClientFactory interface {
	NewHTTPClient() *http.Client
}
