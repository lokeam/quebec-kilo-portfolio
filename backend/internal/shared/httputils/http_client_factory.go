package httputils

import "net/http"

type DefaultHTTPClientFactory struct{}

func (hcf *DefaultHTTPClientFactory) Create() *http.Client {
	return &http.Client{}
}