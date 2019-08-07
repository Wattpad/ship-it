package api

import (
	"context"
	"net/http"

	"ship-it/internal/api/models"
)

type Service interface {
	ListReleases(context.Context) ([]models.Release, error)
}

type controller struct {
	svc Service
}

func NewController(s Service) Controller {
	return &controller{
		svc: s,
	}
}

func (c *controller) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c *controller) ListReleases(w http.ResponseWriter, r *http.Request) {
	deps, err := c.svc.ListReleases(r.Context())
	if err != nil {
		Error500(w, err)
		return
	}

	Success200(w, deps)
}
