package httputils

import (
	"testing"
)

/*
	Behaviors:
		1. Making Write Operations
			- When we call the Write() method, it should make a call to http.ResponseWrite

		2. Tracking Write Status
			- The response writer adapter should be able to know THAT any data HAS BEEN written
			- The response writer adapter should be able to know WHEN any data is written

	Scenarios:
		- Before any write ops, adapter should report that no data has been written
		- When data is written using the Write() method, adapter should:
			* Call the Write() method on the http.ResponseWriter
			* Set an internal boolean written flag to true
		- When some known data is written, the adapter should known that http.ResponseWriter RECEIVED the data
		- When some known data is written, the adapter should know THAT the Written() method returns true
		- If the http.ResponseWriter fails, the adapter should:
			* Mark itself as written (to keep track that we made a write attempt)
			* Send the error to the caller of the Write() method
*/

func TestNewResponseWriterAdapter(t *testing.T) {

			testResponseWriter := newTestResponseWriter()
			testRWAdapter := NewResponseWriterAdapter(testResponseWriter)

	// Scenario 1: No written data
	t.Run(
		`GIVEN a valid ResponseWriter adapter`,
		func(t *testing.T) {
			t.Run(
				`THEN before any write operations, Written() should return false`,
				func(t *testing.T) {
					if testRWAdapter.Written() {
						t.Errorf("We expected that Written() would return false, but it returned true")
					}
				},
			)

			t.Run(
				`WHEN the Write() method is called with some known data`,
				func(t *testing.T) {
					dataToWrite := []byte("You are not special. You are not some unique snowflake. You're the same decaying organic matter as everyone else. We're all part of the same compost heap")
					bytesWritten, err := testRWAdapter.Write(dataToWrite)

					t.Run(
						`THEN the Write() method should succeed`,
						func(t *testing.T) {
							if err != nil {
								t.Fatalf("We expected the Write() method to succeed, but it failed with an error: %v", err)
							}
							if bytesWritten != len(dataToWrite) {
								t.Errorf("We expected %d bytes to be written, but instead %d bytes were written", len(dataToWrite), bytesWritten)
							}
						},
					)

					t.Run(
						`THEN the Written() method should return true`,
						func(t *testing.T) {
							if !testRWAdapter.Written() {
								t.Errorf("We expected that Written() would return true after a write operation")
							}
						},
					)

					t.Run(
						`THEN the http.ResponseWriter should contain the written data`,
						func(t *testing.T) {
							if actualWrittenData := testResponseWriter.buf.String(); actualWrittenData != string(dataToWrite) {
								t.Errorf("We expected the http.ResponseWriter to contain this written data: %q, but instead it contained: %q", string(dataToWrite), actualWrittenData)
							}
						},
					)
				},
			)
		},
	)

	// Scenario 2: Write() error
	t.Run(
		`GIVEN a ResponseWriter adapter wrapping a http.ResponseWriter that returns an error when writing`,
		func(t *testing.T) {
			testError := &errorResponseWriter{}
			testRWAdapter := NewResponseWriterAdapter(testError)

			t.Run(
				`WHEN the we call the Write() method`,
				func(t *testing.T) {
					dataToWrite := []byte("I am Jack's inflamed sense of rejection")
					bytesWritten, err := testRWAdapter.Write(dataToWrite)

					t.Run(
						`THEN the Write() method should return an error`,
						func(t *testing.T) {
							if err == nil {
								t.Errorf("We expected the Write() method to return an error, instead it returned nil")
							}
							if bytesWritten != 0 {
								t.Errorf("We expected 0 bytes to be written on error, instead %d bytes were written", bytesWritten)
							}
						},
					)

					t.Run(
						`THEN the Written() method should still return true`,
						func(t *testing.T) {
							// Note: Even on error, the adapter sets the "written" flag to be true
							if !testRWAdapter.Written() {
								t.Errorf("We expected the Written() method to return true, even if Write() returned an error")
							}
						},
					)
				},
			)
		},
	)
}
