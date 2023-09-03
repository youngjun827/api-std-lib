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
)

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

		go func(clientIP string) {
			time.Sleep(rateLimitingWindow)

			requests.LoadAndDelete(clientIP)
		}(clientIP)

		next.ServeHTTP(w, r)
	})
}
