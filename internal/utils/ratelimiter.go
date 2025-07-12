package utils

import (
	"net/http"

	"golang.org/x/time/rate"
)

// NewRateLimiter creates a new rate limiter with the given rate and burst.
// rate: how many events are allowed per second.
// burst: how many events are allowed to accumulate in the bucket.
func NewRateLimiter(rateLimit rate.Limit, burst int) *rate.Limiter {
	return rate.NewLimiter(rateLimit, burst)
}

// RateLimitMiddleware is a middleware that applies rate limiting to requests.
func RateLimitMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
