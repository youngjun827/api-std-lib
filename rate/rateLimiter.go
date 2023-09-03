package rate

import (
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

		if requestCount >= requestLimit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		requests.Store(clientIP, requestCount+1)

		next.ServeHTTP(w, r)
	})
}
