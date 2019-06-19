package middleware

import (
	"fmt"
	"net/http"
	"time"

	"strings"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/dogstatsd"
)

func getHist(route string, t *dogstatsd.Timing) metrics.Histogram {
	str := strings.ReplaceAll(route, "/", ".")
	fmt.Println(str)
	return t.With(str)
}

func Timer(t *dogstatsd.Timing) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				h := getHist(chi.RouteContext(r.Context()).RoutePattern(), t)
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
