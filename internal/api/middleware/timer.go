package middleware

import (
	"fmt"
	"net/http"
	"time"

	"regexp"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/dogstatsd"
)

func getHist(r *http.Request, t *dogstatsd.Timing) metrics.Histogram {
	str := strings.ReplaceAll(chi.RouteContext(r.Context()).RoutePattern(), "/", ".")
	fmt.Println(str)
	regx := regexp.MustCompile(`\{[^}]*\}`)
	if regx.MatchString(str) {
		str = regx.ReplaceAllString(str, "")
	}
	fmt.Println(str)
	return t.With(str, "method", r.Method)
}

func Timer(t *dogstatsd.Timing) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
				h := getHist(r, t)
				defer func(t0 time.Time) {
					h.Observe(millisecondsSince(t0))
				}(time.Now())
			},
		)
	}
}

func millisecondsSince(t time.Time) float64 {
	return float64(time.Since(t).Seconds() * 1000)
}
