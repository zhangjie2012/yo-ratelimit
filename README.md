# yo-ratelimit

A go token-limit and sliding window rate limit implementation. additional have a ratelimit pool for complex scene. 
btw, it's thread-safed.

- `span` time span, or window width.
- `burst` max request in span, bigger will not `Allow()`.

## usage

The simple scenario, pure use `NewTokenBucket` or `NewSlidingWindow`.

``` go
b := NewTokenBucket(60, 60000) 

if b.Allow() {
    // do next
} else {
    // block request
}
```

_Above example 1 minute max 60000 request, in fact, not equal `qps = 6000`._

if you want complex sence, for single userid, IP, or api request set a ratelimit. `RateLimiterPool` may be helpful for you.

```
pool := NewRateLimiterPool(
    WithNewLimiterFunc(NewTokenBucket),
    WithGCPeriod(1*time.Hour),
    WithExpiredPeriod(3*time.Hour),
)
if pool.Allow("userid:10001", 60, 1000) {
    // do next
} else {
    // block request
}
```

which means: the `10001` user 1 minute max request 1000 times.

more example checkout: <./ratelimiter_test.go>

