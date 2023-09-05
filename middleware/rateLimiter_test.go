package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRateLimiter(t *testing.T) {
	requestLimit = 2 

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	limiter := RateLimiter(handler)

	t.Run("under rate limit", func(t *testing.T) {
		requests = sync.Map{}
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		req.RemoteAddr = "192.168.0.1"

		limiter.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("unexpected status: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("over rate limit", func(t *testing.T) {
		requests = sync.Map{}
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.0.1"
	
		limiter.ServeHTTP(rr, req)
		rr = httptest.NewRecorder()  
		limiter.ServeHTTP(rr, req)
		rr = httptest.NewRecorder()  
		limiter.ServeHTTP(rr, req) 
	
		if status := rr.Code; status != http.StatusTooManyRequests {
			t.Errorf("unexpected status: got %v want %v", status, http.StatusTooManyRequests)
		}
	})

	t.Run("different clients", func(t *testing.T) {
		requests = sync.Map{}
		rr1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/", nil)
		req1.RemoteAddr = "192.168.0.2"

		rr2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/", nil)
		req2.RemoteAddr = "192.168.0.3"

		limiter.ServeHTTP(rr1, req1)
		limiter.ServeHTTP(rr2, req2)

		if status := rr1.Code; status != http.StatusOK {
			t.Errorf("unexpected status for client 1: got %v want %v", status, http.StatusOK)
		}

		if status := rr2.Code; status != http.StatusOK {
			t.Errorf("unexpected status for client 2: got %v want %v", status, http.StatusOK)
		}
	})
}
