package go_utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// BlowUpOnClose provides a closer that always returns an error, for testing purposes
type BlowUpOnClose struct{}

// Close on this mock always returns an error
func (rc BlowUpOnClose) Close() error { return fmt.Errorf("ka-boom") }

// Read on this mock always reads 0 bytes and returns no error
func (rc BlowUpOnClose) Read(_ []byte) (n int, err error) {
	return 0, nil
}

// BlowUpOnRead provides a reader that always returns an error, for testing purposes
type BlowUpOnRead struct{}

// Close on this mock always succeeds with no error
func (rc BlowUpOnRead) Close() error { return nil }

// Read on this mock always returns an error
func (rc BlowUpOnRead) Read(_ []byte) (n int, err error) { return 0, fmt.Errorf("boom") }

// MockHTTPTransportFunc defines the signature of a function that can be assigned to a HTTP client Transport value
type MockHTTPTransportFunc func(req *http.Request) (*http.Response, error)

// RoundTrip always applies the func on which it is a receiver
func (f MockHTTPTransportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// MockHTTPClient returns *http.Client with Transport replaced to avoid making real calls
func MockHTTPClient(fn MockHTTPTransportFunc) *http.Client {
	return &http.Client{
		Transport: MockHTTPTransportFunc(fn),
	}
}

// NewErrorResponseWriter returns an initialized ErrorResponseWriter
func NewErrorResponseWriter(err error) *ErrorResponseWriter {
	return &ErrorResponseWriter{
		err: err,
		rec: httptest.NewRecorder(),
	}
}

// ErrorResponseWriter is a http.ResponseWriter that always errors on attempted writes.
//
// It is necessary for tests.
type ErrorResponseWriter struct {
	err error
	rec *httptest.ResponseRecorder
}

// Header delegates reading of headers to the underlying response writer
func (w *ErrorResponseWriter) Header() http.Header {
	return w.rec.Header()
}

// Write always returns the supplied error on any attempt to write.
func (w *ErrorResponseWriter) Write([]byte) (int, error) {
	return 0, w.err
}

// WriteHeader delegates writing of headers to the underlying response writer
func (w *ErrorResponseWriter) WriteHeader(statusCode int) {
	w.rec.WriteHeader(statusCode)
}
