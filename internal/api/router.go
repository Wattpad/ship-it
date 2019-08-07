package api

import (
	"net/http"

	"ship-it/internal/api/middleware"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/metrics"
)

type Controller interface {
	Health(http.ResponseWriter, *http.Request)
	ListReleases(http.ResponseWriter, *http.Request)
}

// New returns an 'http.Handler' that serves the ship-it API.
func NewRouter(root http.Handler, c Controller, t metrics.Histogram) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Timer(t))

	r.Get("/health", c.Health)

	r.Route("/api", func(r chi.Router) {
		r.Get("/releases", c.ListReleases)
	})

	r.Mount("/", root)

	return r
}
