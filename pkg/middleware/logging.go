package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseRecorder is a wrapper around http.ResponseWriter to capture the status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Create a response writer to capture the status code
		rr := &responseRecorder{w, http.StatusOK}
		next.ServeHTTP(rr, r)

		log.Printf("Completed %s %s in %v with status %d", r.Method, r.URL.Path, time.Since(start), rr.statusCode)
	})
}
