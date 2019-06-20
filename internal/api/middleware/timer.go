package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/metrics"
)

func getIdentifier(route string) string {
	r := strings.NewReplacer("{", "", "}", "", "/", ".")
	str := r.Replace(route)
	fmt.Println(str[1:])
	return str[1:]
}

func Timer(h metrics.Histogram) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				t0 := time.Now()
				next.ServeHTTP(w, r)
				id := getIdentifier(chi.RouteContext(r.Context()).RoutePattern())
				defer func(t0 time.Time) {
					h.With("endpoint", id, "method", r.Method).Observe(millisecondsSince(t0))
				}(t0)
			},
		)
	}
}

func millisecondsSince(t time.Time) float64 {
	return float64(time.Since(t).Seconds() * 1000)
}
