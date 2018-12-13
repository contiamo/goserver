package http

import (
	"net/http"
)

// InstrumentedResponseWriter is a http ResponseWriter that supports recording the status code and
// response size.  This is used in logging, metrics, and tracing
type InstrumentedResponseWriter struct {
	http.ResponseWriter
	length     int
	statusCode int
}

func (w *InstrumentedResponseWriter) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)
	w.length += n
	return
}

func (w *InstrumentedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
	return
}

func (w *InstrumentedResponseWriter) Length() int {
	return w.length
}

func (w *InstrumentedResponseWriter) StatusCode() int {
	return w.statusCode
}
