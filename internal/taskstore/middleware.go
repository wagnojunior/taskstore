package taskstore

import (
	"log"
	"net/http"
	"time"
)

// `LoggingMiddleware` logs the request method, URL, and elapsed time
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}
