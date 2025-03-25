package mocks

import (
	"bytes"
	"io"
	"net/http"
)

type MockIGDBRoundTripper struct {
	response *http.Response
	err      error
}

func (m *MockIGDBRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
			return nil, m.err
	}

	if m.response == nil {
			// Create a mock successful response if no specific response is set
			body := io.NopCloser(bytes.NewReader([]byte(`[{"id": 1, "name": "Test Game"}]`)))
			m.response = &http.Response{
					StatusCode: http.StatusOK,
					Body:       body,
			}
	}

	return m.response, nil
}
