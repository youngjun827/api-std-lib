package rate

import (
	"net/http"
	"sync"
	"time"
)

var (
	rateLimitingWindow = 1 * time.Minute
	requestLimit       = 10
	requests           = make(map[string]int)
	mu                 sync.Mutex
)

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		clientIP := r.RemoteAddr
		if count, exists := requests[clientIP]; exists && count >= requestLimit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Increment the request count for this IP
		requests[clientIP]++

		// Schedule to decrement the count after the rate limiting window
		go func(clientIP string) {
			time.Sleep(rateLimitingWindow)
			mu.Lock()
			requests[clientIP]--
			mu.Unlock()
		}(clientIP)

		next.ServeHTTP(w, r)
	})
}
