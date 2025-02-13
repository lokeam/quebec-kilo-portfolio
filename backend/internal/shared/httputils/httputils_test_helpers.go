//lint:file-ignore U1000 These helper functions + types are used by the other httputilstest files.
package httputils

import (
	"bytes"
	"errors"
	"net/http"
)

type testResponseWriter struct {
	header        http.Header
	buf           *bytes.Buffer
	body          *bytes.Buffer
	status        int
	writtenFlag   bool
}

func newTestResponseWriter() *testResponseWriter {
	return &testResponseWriter{
		header: make(http.Header),
		buf:    &bytes.Buffer{},
	}
}

// Header returns the HTTP header map.
func (trw *testResponseWriter) Header() http.Header {
	return trw.header
}

// Write writes the provided bytes into the internal buffer.
// It also sets a flag indicating that something was written.
func (trw *testResponseWriter) Write(b []byte) (int, error) {
	trw.writtenFlag = true
	return trw.buf.Write(b)
}

// WriteHeader sets the response status code.
func (trw *testResponseWriter) WriteHeader(statusCode int) {
	trw.status = statusCode
}

// Written returns whether the Write method was called.
func (trw *testResponseWriter) Written() bool {
	return trw.writtenFlag
}

// errorResponseWriter simulates a ResponseWriter that always fails on Write.
type errorResponseWriter struct{}

// newErrorResponseWriter creates a new errorResponseWriter.
func newErrorResponseWriter() *errorResponseWriter {
	return &errorResponseWriter{}
}

// Header returns an empty header.
func (e *errorResponseWriter) Header() http.Header {
	return make(http.Header)
}

// Write always returns an error.
func (e *errorResponseWriter) Write(b []byte) (int, error) {
	return 0, errors.New("simulated write error")
}

// WriteHeader is a no-op.
func (e *errorResponseWriter) WriteHeader(statusCode int) {}
