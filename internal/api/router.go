package api

import (
	"context"
	"net/http"

	"ship-it/internal/api/middleware"
	"ship-it/internal/models"
	"ship-it/internal/service"

	"github.com/go-chi/chi"
)

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type Service interface {
	ListReleases(context.Context) ([]models.Release, error)
}

// New returns an 'http.Handler' that serves the ship-it API.
func New(s *service.Service) http.Handler {
	c := newController(s)

	r := chi.NewRouter()

	r.Get("/health", health)

	r.Get("/releases", c.ListReleases)
	hist := s.DDTimer.With("method", "get", "releases")
	r.Use(middleware.Timer(hist))

	r.Mount("/dashboard", http.FileServer(http.Dir("")))
	r.Mount("/static", http.FileServer(http.Dir("dashboard")))

	return r
}
