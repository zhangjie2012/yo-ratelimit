# yo-ratelimit

A go token-limit and sliding window rate limit implementation

## 术语

流量速度通过 `span` 和 `limit` 来定义。

- `span` 时间跨度，单位为秒
- `burst` 时间跨度的最大流量
- `burst/span` 单位时间内增量
