package server

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

type WriterWithStatus struct {
	http.ResponseWriter
	status int
}

func (w *WriterWithStatus) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.status = statusCode
}

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		originalURL := r.URL.String()

		writer := &WriterWithStatus{w, http.StatusOK}
		h.ServeHTTP(writer, r)

		log.Printf("logger: [%d] %s -> %s %s", writer.status, originalURL, r.URL.String(), time.Since(start).String())
	})
}
