package middleware

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/metrics"
)

func Timer(h metrics.Histogram) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func(t0 time.Time) {
					h.Observe(millisecondsSince(t0))
				}(time.Now())
				next.ServeHTTP(w, r)
			},
		)
	}
}

func millisecondsSince(t time.Time) float64 {
	return float64(time.Since(t).Seconds() * 1000)
}
