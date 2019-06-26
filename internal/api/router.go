package api

import (
	"context"
	"net/http"

	"ship-it/internal/api/middleware"
	"ship-it/internal/api/models"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/metrics"
)

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type Service interface {
	ListReleases(context.Context) ([]models.Release, error)
}

// New returns an 'http.Handler' that serves the ship-it API.
func New(s Service, t metrics.Histogram) http.Handler {
	c := newController(s)

	r := chi.NewRouter()
	r.Use(middleware.Timer(t))

	r.Get("/health", health)

	r.Route("/api", func(r chi.Router) {
		r.Get("/releases", c.ListReleases)
	})

	r.Mount("/", http.FileServer(http.Dir("dashboard")))

	return r
}
