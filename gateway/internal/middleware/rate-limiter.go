package auth

import (
	"net/http"
	"sync"
	"time"
)

// visitor struct
// it will save information of every client in this structure
type rateLimiter struct {
	mu        sync.Mutex
	visitor   map[string]int
	limit     int
	resetTime time.Duration
}

// initializer
func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitor:   make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
	// start the reset routine
	go rl.ResetVisitorCount()
	return rl
}

func (r *rateLimiter) ResetVisitorCount() {
	for {
		// sleep the duration and reset the again
		time.Sleep(r.resetTime)
		r.mu.Lock()
		r.visitor = make(map[string]int)
		r.mu.Unlock()
	}
}

func (rl *rateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		// get the user ip  address
		ip := r.RemoteAddr
		// set the ip address and increment the number of that object of the map
		rl.visitor[ip]++
		// check if the user with the ip has not riched the limit of requests
		if rl.visitor[ip] > rl.limit {
			http.Error(w, "Too many request", http.StatusTooManyRequests)
			return
		}
		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
