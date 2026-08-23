package limiter

import (
	"sync"
	"testing"
	"time"
)

func TestSlidingWindow_ExhaustsLimit(t *testing.T) {
	l := NewSlidingWindowLimiter(3, time.Second)

	for i := 0; i < 3; i++ {
		if !l.Allow("client1") {
			t.Fatalf("request %d should have been allowed", i+1)
		}
	}

	if l.Allow("client1") {
		t.Fatal("4th request should have been denied, limit reached")
	}
}

func TestSlidingWindow_AllowAfterWindowPasses(t *testing.T) {
	l := NewSlidingWindowLimiter(1, 500*time.Millisecond)

	if !l.Allow("client1") {
		t.Fatal("first request should be allowed")
	}
	if l.Allow("client1") {
		t.Fatal("Second immediate request should be denied")
	}
	time.Sleep(600 * time.Millisecond) //let the window pass

	if !l.Allow("client1") {
		t.Fatal("request after window passed should be allowed")
	}
}

func TestSlidingWindow_SeparateKeysIndependent(t *testing.T) {
	l := NewSlidingWindowLimiter(1, time.Second)

	if !l.Allow("client1") {
		t.Fatal("client1 first request should be allowed")
	}
	if !l.Allow("client2") {
		t.Fatal("client2 should have its own log, unaffected by client1")
	}
}

func TestSlidinfWindow_ConcurrentAccess(t *testing.T) {
	l := NewSlidingWindowLimiter(100, time.Second)
	var wg sync.WaitGroup
	var mu sync.Mutex
	allowedCount := 0

	for i := 0; i < 150; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.Allow("client1") {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if allowedCount != 100 {
		t.Fatalf("expected exactly 100 requests, got %d", allowedCount)
	}
}
