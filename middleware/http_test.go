package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	limiter "github.com/Its-Delimas/rate_limiter"
)

func TestRateLimit_BlocksOverLimiter(t *testing.T) {
	l := limiter.NewTokenBucketLimiter(1, 1)
	handlerCalled := 0

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled++
		w.WriteHeader(http.StatusOK)
	})

	rl := RateLimit(l, func(r *http.Request) string {
		return "test-client"
	})
	wrapped := rl(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// first request - should pass
	w1 := httptest.NewRecorder()
	wrapped.ServeHTTP(w1, req)
	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200 on first request, got %d", w1.Code)
	}

	// second request - bucket exhausted - should be blocked
	w2 := httptest.NewRecorder()
	wrapped.ServeHTTP(w2, req)
	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 on second request, got %d", w2.Code)
	}

	if handlerCalled != 1 {
		t.Fatalf("expected next handler called exactly once, got %d", handlerCalled)
	}

}
