package taskstore

import (
	"log"
	"net/http"
	"time"
)

// `LoggingMiddleware` logs the request method, URL, and elapsed time
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	}
}
