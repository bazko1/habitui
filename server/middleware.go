package server

import (
	"log"
	"net/http"
)

func logRequestMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ssrw := NewStatusSaveResponseWriter(w)
		next.ServeHTTP(ssrw, r)
		log.Printf("method: %s, path: %s | %d", r.Method, r.URL.Path, ssrw.StatusCode)
	}
}

type statusSaveResponseWriter struct {
	responseWriter http.ResponseWriter
	StatusCode     int
}

func NewStatusSaveResponseWriter(w http.ResponseWriter) *statusSaveResponseWriter {
	return &statusSaveResponseWriter{w, http.StatusOK}
}

func (w *statusSaveResponseWriter) Write(b []byte) (int, error) {
	return w.responseWriter.Write(b)
}

func (w *statusSaveResponseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

func (w *statusSaveResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.responseWriter.WriteHeader(statusCode)
}
