package server

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		originalURL := r.URL.String()

		h.ServeHTTP(w, r)
		log.Printf("logger: new request %s -> %s took %s", originalURL, r.URL.String(), time.Since(start).String())
	})
}
