package middleware

import (
	"net/http"
	"time"

	"regexp"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/metrics"
)

func getHist(r *http.Request, t metrics.Histogram) metrics.Histogram {
	str := strings.ReplaceAll(chi.RouteContext(r.Context()).RoutePattern(), "/", ".")
	regx := regexp.MustCompile(`[{}]`)
	if regx.MatchString(str) {
		str = regx.ReplaceAllString(str, "")
	}

	str = strings.Replace(str, ".", "", 1)

	return t.With(str, "method", r.Method)
}

func Timer(t metrics.Histogram) func(http.Handler) http.Handler {
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
