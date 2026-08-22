package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens         float64
	capacity       float64
	refillRate     float64 //tokens per second
	lastRefillTime time.Time
}

type TokenBucketLimiter struct {
	mu         sync.Mutex
	buckets    map[string]*TokenBucket
	capacity   float64
	refillRate float64
}

func NewTokenBucketLimiter(capacity, refillRate float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		buckets:    make(map[string]*TokenBucket),
		capacity:   capacity,
		refillRate: refillRate,
	}
}
