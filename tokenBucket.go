package limiter

import (
	"sync"
	"time"
)

// structs
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

// constructor
func NewTokenBucketLimiter(capacity, refillRate float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		buckets:    make(map[string]*TokenBucket),
		capacity:   capacity,
		refillRate: refillRate,
	}
}

func (l *TokenBucketLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	bucket, exists := l.buckets[key]
	if !exists {
		bucket = &TokenBucket{
			tokens:         l.capacity,
			capacity:       l.capacity,
			refillRate:     l.refillRate,
			lastRefillTime: time.Now(),
		}
		l.buckets[key] = bucket
	}

	now := time.Now()
	elapsed := now.Sub(bucket.lastRefillTime).Seconds()
	bucket.tokens += elapsed * bucket.refillRate
	if bucket.tokens > bucket.capacity {
		bucket.tokens = bucket.capacity
	}
	bucket.lastRefillTime = now

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}
	return false
}
