# limiter

A pluggable rate limiting library for Go, with net/http middleware.

## Features
- Token bucket algorithm (lazy refill, no background goroutines)
- Interface-based design — swap algorithms without touching middleware
- Thread-safe, verified under `-race`
- net/http middleware with pluggable key extraction

## Requirements
- Go 1.25.0+

## Installation
- go get [github.com](https://github.com/Its-Delimas/rate_limiter)


## Usage

```go
import (
    "net"
    "net/http"

    "github.com/Its-Delimas/rate_limiter"
    "github.com/Its-Delimas/rate_limiter/middleware"
)

l := limiter.NewTokenBucketLimiter(5, 1) // capacity 5, refill 1/sec

rl := middleware.RateLimit(l, func(r *http.Request) string {
    host, _, _ := net.SplitHostPort(r.RemoteAddr)
    return host
})

http.Handle("/", rl(myHandler))
```

## Design notes
- Lazy refill: tokens computed on-demand from elapsed time, not a ticking goroutine per bucket — avoids goroutine leaks at scale.

- Key extraction is caller-supplied, not hardcoded to IP — works behind proxies, auth systems, etc.

- New keys start with a full bucket, not empty, so first-time clients aren't penalized.

## Testing
go test -race ./
-   Output:
```bash
ok      github.com/Its-Delimas/rate_limiter     (cached)
```
## License

Apache 2.0 - see [LICENSE](LICENSE)