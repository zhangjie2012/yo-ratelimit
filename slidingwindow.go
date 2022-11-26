package ratelimit

import (
	"sync"
	"time"
)

// SlidingWindow slide-windows algorithm implemention
type SlidingWindow struct {
	mu   sync.Mutex
	last time.Time // just store last allow time

	Span  float64 // second
	Burst float64

	// |  previous span   |   current span    |
	// |------------------|---------------....|
	// |    prevCount     |     currCount     |
	// |                  |<--currStart      |
	prevCount float64
	currCount float64
	currStart time.Time
}

func NewSlidingWindow(span, burst float64) RateLimiter {
	return &SlidingWindow{
		Span:      span,
		Burst:     burst,
		prevCount: 0,
		currCount: 0,
		currStart: time.Now(),
		last:      time.Now(),
	}
}

func (l *SlidingWindow) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	current := time.Now()
	l.last = current

	// reach next(or multiple next) window, update current window
	if float64(current.Sub(l.currStart))/float64(time.Second) > l.Span {
		// find current window
		windowStart := l.currStart
		for windowStart.Add(time.Duration(l.Span) * time.Second).Before(current) {
			windowStart = windowStart.Add(time.Duration(l.Span) * time.Second) // move next
		}
		l.currStart = windowStart
		// Always treat the previous window amount as the last one,
		//   no matter how many time windows the period spans.
		// If you want to change this behavior, add a counter in for loop, then then make a judgment.
		l.prevCount = l.currCount
		l.currCount = 0
	}

	preWindowPercent := (l.Span - float64(current.Sub(l.currStart))/float64(time.Second)) / l.Span
	counter := l.prevCount*preWindowPercent + l.currCount

	if counter > l.Burst {
		return false
	}
	l.currCount++
	return true
}

func (l *SlidingWindow) Last() time.Time {
	return l.last
}
