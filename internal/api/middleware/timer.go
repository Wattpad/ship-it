package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/metrics"
)

func getIdentifier(ctx context.Context) string {
	str := strings.ReplaceAll(chi.RouteContext(ctx).RoutePattern(), "/", ".")
	regx := regexp.MustCompile(`[{}]`)

	str = regx.ReplaceAllString(str, "")
	str = strings.Replace(str, ".", "", 1)
	return str
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
