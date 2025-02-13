package httputils

import "net/http"

// This is complete bullshit that I have to write this
type ResponseWriter interface {
	http.ResponseWriter
	Written() bool
}

type responseWriterAdapter struct {
	http.ResponseWriter
	written bool
}

// Write sets the written flag and delegates to the underlying ResponseWriter.
func (r *responseWriterAdapter) Write(b []byte) (int, error) {
	r.written = true
	return r.ResponseWriter.Write(b)
}

// Written returns whether the response has been written.
func (r *responseWriterAdapter) Written() bool {
	return r.written
}

// NewResponseWriterAdapter wraps an http.ResponseWriter.
func NewResponseWriterAdapter(w http.ResponseWriter) ResponseWriter {
	return &responseWriterAdapter{ResponseWriter: w}
}