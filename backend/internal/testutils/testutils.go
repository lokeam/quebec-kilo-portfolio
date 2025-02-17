package testutils

import (
	"errors"
	"net/http"
	"sync"
)

// --------- Test Logger ---------
type TestLogger struct {
	Mu         sync.Mutex
	InfoCalls  []string
	DebugCalls []string
	ErrorCalls []string
	WarnCalls  []string
}

func NewTestLogger() *TestLogger {
	return &TestLogger{}
}

func (tl *TestLogger) Info(msg string, fields map[string]any) {
	tl.Mu.Lock()
	defer tl.Mu.Unlock()
	tl.InfoCalls = append(tl.InfoCalls, msg)
}

func (tl *TestLogger) Debug(msg string, fields map[string]any) {
	tl.Mu.Lock()
	defer tl.Mu.Unlock()
	tl.DebugCalls = append(tl.DebugCalls, msg)
}

func (tl *TestLogger) Error(msg string, fields map[string]any) {
	tl.Mu.Lock()
	defer tl.Mu.Unlock()
	errStr := msg
	if fields != nil {
		if err, ok := fields["error"].(error); ok {
			errStr += ": " + err.Error()
		}
	}
	tl.ErrorCalls = append(tl.ErrorCalls, errStr)
}

func (tl *TestLogger) Warn(msg string, fields map[string]any) {
	tl.Mu.Lock()
	defer tl.Mu.Unlock()
	tl.WarnCalls = append(tl.WarnCalls, msg)
}


// --------- IGDB Config ---------
type MockIGDBConfig struct {
	tokenKey   string
	err        error
}

func (mic *MockIGDBConfig) GetTokenKey() (string, error) {
	if mic.err != nil {
		return "", mic.err
	}

	return mic.tokenKey, nil
}


// --------- Round Tripper (simulate network error) ---------
type ErrorRoundTripper struct {}

func (ert *ErrorRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("network error")
}