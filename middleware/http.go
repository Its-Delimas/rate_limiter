package middleware

import (
	"net/http"

	limiter "github.com/Its-Delimas/rate_limiter"
)

type KeyFunc func(r *http.Request) string

func RateLimit(l limiter.Limiter, KeyFunc KeyFunc) func(http, http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := KeyFunc(r)
			if !l.Allow(key) {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
