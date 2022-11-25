package ratelimit

import (
	"sync"
	"time"
)

// TokenBucket
type TokenBucket struct {
	// span unit is second, burst/second, per second token increase
	Span  float64
	Burst float64

	mu sync.Mutex
	// tokens is current unused token count
	tokens float64
	last   time.Time
}

func NewTokenBucket(span, burst float64) *TokenBucket {
	return &TokenBucket{
		Span:   span,
		Burst:  burst,
		tokens: burst,
		last:   time.Now(),
	}
}

func (limiter *TokenBucket) Allow() bool {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	// 1. update bucket tokens by passed time
	current := time.Now()
	timePassed := current.Sub(limiter.last)
	limiter.last = current

	limiter.tokens += (float64(timePassed) / float64(time.Second)) * (limiter.Burst / limiter.Span)
	if limiter.tokens > limiter.Burst {
		limiter.tokens = limiter.Burst
	}

	// 2. logic
	if limiter.tokens >= 1 {
		limiter.tokens--
		return true
	}
	return false
}
