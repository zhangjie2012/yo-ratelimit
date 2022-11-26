package ratelimit

import (
	"sync"
	"time"
)

// TokenBucket token-bucket algorithm implemention
type TokenBucket struct {
	mu sync.Mutex

	// span unit is second, burst/second, per second token increase
	Span  float64
	Burst float64

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

func (l *TokenBucket) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 1. update bucket tokens by passed time
	current := time.Now()
	timePassed := current.Sub(l.last)
	l.last = current

	l.tokens += (float64(timePassed) / float64(time.Second)) * (l.Burst / l.Span)
	if l.tokens > l.Burst {
		l.tokens = l.Burst
	}

	// 2. logic
	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}
