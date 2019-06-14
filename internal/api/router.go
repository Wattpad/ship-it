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
	ListDeployments(context.Context) ([]models.DeploymentDetail, error)
}

// New returns an 'http.Handler' that serves the ship-it API.
func New(s Service) http.Handler {
	c := newController(s)

	r := chi.NewRouter()

	r.Get("/health", health)

	r.Get("/deployments", c.ListDeployments)

	r.Mount("/dashboard", http.FileServer(http.Dir("")))
	r.Mount("/static", http.FileServer(http.Dir("dashboard")))

	return r
}
