package ratelimit

import (
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
			t.Log(i, b.Allow(), b.tokens)
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
					t.Log(idx, j, b.Allow(), b.tokens)
				}
			}(i)
		}
		wg.Wait()
	}
}
