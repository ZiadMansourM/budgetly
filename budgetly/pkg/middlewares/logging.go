package middlewares

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// responseRecorder is a wrapper around http.ResponseWriter to capture the status code and size of the response
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// WriteHeader captures the status code
func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

// Write captures the size of the response body
func (rr *responseRecorder) Write(b []byte) (int, error) {
	size, err := rr.ResponseWriter.Write(b)
	rr.size += size
	return size, err
}

// LoggingMiddleware logs requests with more detailed information
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response recorder to capture the status code and size
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// Get client IP address and port
		clientIP := r.RemoteAddr

		// Strip any IPv6 zone identifiers if present
		if strings.Contains(clientIP, "%") {
			clientIP = strings.Split(clientIP, "%")[0]
		}

		// Serve the request
		next.ServeHTTP(rr, r)

		// Calculate duration
		duration := time.Since(start)

		// Log the request
		slog.Info(
			"HTTP Request",
			"method", r.Method,
			"url", r.URL.String(),
			"proto", r.Proto,
			"client_ip", clientIP,
			"status", rr.statusCode,
			"response_size", rr.size,
			"duration", duration,
		)
	})
}
