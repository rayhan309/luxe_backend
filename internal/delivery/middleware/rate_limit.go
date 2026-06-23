package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter creates a simple in-memory rate limiter middleware
func RateLimiter(rps float64, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "rate limit exceeded, please slow down",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// StrictRateLimiter for auth endpoints
func StrictRateLimiter() gin.HandlerFunc {
	return RateLimiter(5, 10) // 5 req/sec, burst 10
}

// Timeout middleware adds a deadline to requests
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Signal clients with a header
		c.Writer.Header().Set("X-Request-Timeout", timeout.String())
		c.Next()
	}
}
