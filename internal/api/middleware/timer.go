package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/metrics"
)

func getIdentifier(ctx context.Context) string {
	r := strings.NewReplacer("{", "", "}", "", "/", ".")
	str := r.Replace(chi.RouteContext(ctx).RoutePattern())

	return str[1:]
}

func Timer(h metrics.Histogram) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
				id := getIdentifier(r.Context())
				defer func(t0 time.Time) {
					h.With("endpoint", id, "method", r.Method).Observe(millisecondsSince(t0))
				}(time.Now())
			},
		)
	}
}

func millisecondsSince(t time.Time) float64 {
	return float64(time.Since(t).Seconds() * 1000)
}
