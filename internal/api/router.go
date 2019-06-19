package api

import (
	"context"
	"net/http"

	"ship-it/internal/models"

	"github.com/go-chi/chi"
)

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type Service interface {
	ListReleases(context.Context) ([]models.Release, error)
}

// New returns an 'http.Handler' that serves the ship-it API.
func New(s Service) http.Handler {
	c := newController(s)

	r := chi.NewRouter()

	r.Get("/health", health)

	r.Route("/api", func(r chi.Router) {
		r.Get("/releases", c.ListReleases)
	})

	r.Mount("/", http.FileServer(http.Dir("dashboard")))

	return r
}
