package limiter

import (
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	mu          sync.Mutex
	logs        map[string][]time.Time
	maxRequests int
	windowSize  time.Duration
}

func NewSlidingWindowLimiter(maxRequests int, windowSize time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		logs:        make(map[string][]time.Time),
		maxRequests: maxRequests,
		windowSize:  windowSize,
	}
}

func (l *SlidingWindowLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-l.windowSize)

	timestamps := l.logs[key]

	// pruning timestamps outside the window
	valid := timestamps[:0]
	for _, t := range timestamps {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= l.maxRequests {
		l.logs[key] = valid
		return false
	}

	valid = append(valid, now)
	l.logs[key] = valid
	return true
}
