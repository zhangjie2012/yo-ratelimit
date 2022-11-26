package ratelimit

import (
	"sync"
	"time"
)

// RateLimiterPool a limiter pool
// every limiter has unique key call "fingerprint", the fingerprint may be ip, userid, api, etc.
// when limiter expired(long time unused), it will auto recycle by "gc" function
//
// "gcPeriod", "expiredPeriod" must far bigger than "span", suggest at least 10x, 20x ...
// "gc" will lock pool, don't set too short value
type RateLimiterPool struct {
	mu             sync.Mutex
	limiters       map[string]RateLimiter
	newLimiterFunc func(span, burst float64) RateLimiter
	gcPeriod       time.Duration
	expiredPeriod  time.Duration
	logger         Logger
}

type PoolOptionFunc func(*RateLimiterPool)

func WithNewLimiterFunc(f func(span, burst float64) RateLimiter) PoolOptionFunc {
	return func(p *RateLimiterPool) {
		p.newLimiterFunc = f
	}
}

func WithGCPeriod(gcPeriod time.Duration) PoolOptionFunc {
	return func(p *RateLimiterPool) {
		p.gcPeriod = gcPeriod
	}
}

func WithExpiredPeriod(expiredPeriod time.Duration) PoolOptionFunc {
	return func(p *RateLimiterPool) {
		p.expiredPeriod = expiredPeriod
	}
}

func WithLogger(logger Logger) PoolOptionFunc {
	return func(p *RateLimiterPool) {
		p.logger = logger
	}
}

func NewRateLimiterPool(options ...PoolOptionFunc) *RateLimiterPool {
	pool := &RateLimiterPool{
		limiters:       make(map[string]RateLimiter),
		newLimiterFunc: NewTokenBucket,
		gcPeriod:       1 * time.Hour,
		expiredPeriod:  6 * time.Hour,
		logger:         NewXLog(),
	}

	for _, opt := range options {
		opt(pool)
	}

	go pool.gcLoop()

	return pool
}

// Allow "fingerprint" not exist, will auto create a limiter, you don't need care about it.
//   "span", "burst" only used in auto create.
func (pool *RateLimiterPool) Allow(fingerprint string, span, burst float64) bool {
	limiter, ok := pool.limiters[fingerprint]
	if !ok {
		limiter = pool.add(fingerprint, span, burst)
	}
	return limiter.Allow()
}

func (pool *RateLimiterPool) add(fingerprint string, span, burst float64) RateLimiter {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// double check
	if limiter, ok := pool.limiters[fingerprint]; ok {
		return limiter
	}

	limiter := pool.newLimiterFunc(span, burst)
	pool.limiters[fingerprint] = limiter

	pool.logger.Infof("add limiter to pool, fingerprint=\"%s\"\n", fingerprint)

	return limiter
}

func (pool *RateLimiterPool) gcLoop() {
	ticker := time.NewTicker(pool.gcPeriod)
	for {
		t := <-ticker.C
		pool.gc(t)
	}
}

func (pool *RateLimiterPool) gc(t time.Time) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// pool.logger.Debugf("[%s] start gc total=%d\n", t, len(pool.limiters))
	for fingerprint, limiter := range pool.limiters {
		if t.Sub(limiter.Last()) > pool.expiredPeriod {
			delete(pool.limiters, fingerprint)
			pool.logger.Infof("gc(%s): delete fingerprint, \"%s\"\n", t, fingerprint)
		}
	}
}
