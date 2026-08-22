package limiter

import (
	"testing"
	"time"
)

// Exhaustion
// Refill over time
// key isolation

func TestAllow_ExhaustsBucket(t *testing.T) {
	l := NewTokenBucketLimiter(3, 1) //capacity = 3, refill 1 token/sec

	for i := 0; i < 3; i++ {
		if !l.Allow("client1") {
			t.Fatalf("request %d should have been allowed", i+1)
		}
	}

	if l.Allow("client1") {
		t.Fatal("4th request should havebeen denied, bucket exhausted")
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	l := NewTokenBucketLimiter(1, 1)

	if !l.Allow("client1") {
		t.Fatal("first request should be allowed")
	}
	if l.Allow("client1") {
		t.Fatal("Second immediate request should be denied")
	}

	time.Sleep(1100 * time.Millisecond) //wait for refill

	if !l.Allow("client1") {
		t.Fatal("request after refill wait should be allowed")
	}
}

func TestAllow_SeparateKeysIndependent(t *testing.T) {
	l := NewTokenBucketLimiter(1, 1)

	if !l.Allow("client") {
		t.Fatal("client1 first request should be allowed")
	}
	if !l.Allow("client2") {
		t.Fatalf("client2 should have its own bucket, unaffected by client1")
	}
}
