package ratelimit

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTokenBucketLimiter(t *testing.T) {
	// single thread
	{
		b := NewTokenBucket(5, 25) // 5 qps

		for i := 0; i < 100; i++ {
			time.Sleep(100 * time.Millisecond) // 10 qps
			t.Log(i, b.Allow())
		}
	}

	// multiple thread
	{
		b := NewTokenBucket(5, 25) // 5 qps

		wg := sync.WaitGroup{}
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				time.Sleep(100 * time.Millisecond)
				for j := 0; j < 100; j++ {
					time.Sleep(200 * time.Millisecond) // 5 qps
					t.Log(idx, j, b.Allow())
				}
			}(i)
		}
		wg.Wait()
	}
}

func TestSlidingWindow(t *testing.T) {
	// single thread
	{
		l := NewSlidingWindow(1, 5) // 5 qps
		time.Sleep(100 * time.Millisecond)
		t.Log(l.Allow())

		time.Sleep(200 * time.Millisecond)
		t.Log(l.Allow())

		time.Sleep(300 * time.Millisecond)
		t.Log(l.Allow())

		// move one window (100+200+300+500 - 1000)/1000=10%
		time.Sleep(500 * time.Millisecond)
		t.Log(l.Allow())

		// 100+100
		time.Sleep(100 * time.Millisecond)
		t.Log(l.Allow())
		time.Sleep(10 * time.Millisecond)
		t.Log(l.Allow())
		time.Sleep(10 * time.Millisecond)
		t.Log(l.Allow())
		time.Sleep(10 * time.Millisecond)
		t.Log(l.Allow())
		time.Sleep(200 * time.Millisecond)
		t.Log(l.Allow())
		time.Sleep(10 * time.Millisecond)
		t.Log(l.Allow())
		time.Sleep(10 * time.Millisecond)
		t.Log(l.Allow())
		time.Sleep(10 * time.Millisecond)
		t.Log(l.Allow())
		time.Sleep(2000 * time.Millisecond)
		t.Log(l.Allow())
	}

	// multiple thread
	{
		b := NewSlidingWindow(5, 25) // 5 qps

		wg := sync.WaitGroup{}
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				time.Sleep(100 * time.Millisecond)
				for j := 0; j < 100; j++ {
					time.Sleep(200 * time.Millisecond) // 5 qps
					t.Log(idx, j, b.Allow())
				}
			}(i)
		}
		wg.Wait()
	}
}

func TestRateLimiterPool(t *testing.T) {
	pool := NewRateLimiterPool(
		WithNewLimiterFunc(NewTokenBucket),
		WithGCPeriod(2*time.Second),
		WithExpiredPeriod(3*time.Second),
	)

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			fingerprint := fmt.Sprintf("IP:127.0.0.%d", idx)
			time.Sleep(1 * time.Second)
			for j := 0; j < 100; j++ {
				time.Sleep(150 * time.Millisecond)
				// t.Log(fingerprint, j, pool.Allow(fingerprint, 1, 5))
				pool.Allow(fingerprint, 1, 5)
			}
		}(i)
	}
	wg.Wait()
	time.Sleep(10 * time.Second)
}
