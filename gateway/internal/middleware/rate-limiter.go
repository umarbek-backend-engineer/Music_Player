package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// each visitor info
type visitor struct {
	count    int
	lastSeen time.Time
}

// rate limiter struct
type rateLimiter struct {
	mu        sync.Mutex
	visitor   map[string]*visitor
	limit     int
	resetTime time.Duration
}

// initializer
func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitor:   make(map[string]*visitor),
		limit:     limit,
		resetTime: resetTime,
	}
	// start the reset routine
	go rl.Cleaup()
	return rl
}

// cleanup old visitors
func (r *rateLimiter) Cleaup() {
	for {
		// sleep the duration and reset the again
		time.Sleep(time.Minute)
		r.mu.Lock()
		for ip, v := range r.visitor {
			if time.Since(v.lastSeen) > r.resetTime*2 {
				delete(r.visitor, ip)
			}
		}
		r.mu.Unlock()
	}
}

func (rl *rateLimiter) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get the user ip  address
		ip := ctx.ClientIP()

		rl.mu.Lock()

		v, exists := rl.visitor[ip]
		if !exists {
			v = &visitor{
				count:    1,
				lastSeen: time.Now(),
			}
		} else {
			// reset if the reset time is expired
			if time.Since(v.lastSeen) > rl.resetTime {
				v.count = 1
				v.lastSeen = time.Now()
			} else {
				v.count++
			}
		}

		// if the v.count is has reached the limit the Middleware will give an errro too many requests
		if v.count > rl.limit {
			rl.mu.Unlock()
			ctx.AbortWithStatusJSON(429, gin.H{
				"error": "too many requests",
			})
			return
		}

		rl.mu.Unlock()
		ctx.Next()
	}
}
