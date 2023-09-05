package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

var (
	rateLimitingWindow = 1 * time.Minute
	requestLimit       = 10
	requests           = sync.Map{}
	cleanupInterval    = 5 * time.Minute
)

func init() {
	go periodicCleanup()
}

func periodicCleanup() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		requests.Range(func(key, value interface{}) bool {
			clientIP := key.(string)
			count := value.(int)

			if count == 0 {
				requests.Delete(clientIP)
			}

			return true
		})
	}
}

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		count, _ := requests.LoadOrStore(clientIP, 0)
		requestCount := count.(int)

		// Log request
		slog.Info("Rate limiting request: clientIP=%s, requestCount=%d", clientIP, requestCount)

		if requestCount >= requestLimit {
			// Return JSON error
			jsonErr := map[string]string{
				"error": "Rate limit exceeded",
			}
			jsonBytes, err := json.Marshal(jsonErr)
			if err != nil {
				slog.Error("Failed to marshal JSON error: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write(jsonBytes)
			return
		}

		// Store request
		requests.Store(clientIP, requestCount+1)

		// Serve next handler
		next.ServeHTTP(w, r)
	})
}
