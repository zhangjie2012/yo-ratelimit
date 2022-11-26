package ratelimit

import "time"

// RateLimiter thread-safed rate limiter
type RateLimiter interface {
	Allow() bool     // allow request pass
	Last() time.Time // last check time
}
