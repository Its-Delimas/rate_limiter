package main

import (
	"fmt"
	"net"
	"net/http"

	limiter "github.com/Its-Delimas/rate_limiter"
	"github.com/Its-Delimas/rate_limiter/middleware"
)

func main() {
	l := limiter.NewTokenBucketLimiter(5, 1)

	rl := middleware.RateLimit(l, func(r *http.Request) string {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return r.RemoteAddr
		}
		return host
	})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "request allowed")
	})

	http.Handle("/", rl(handler))

	fmt.Println("listenin on :8080")
	http.ListenAndServe(":8080", nil)
}
