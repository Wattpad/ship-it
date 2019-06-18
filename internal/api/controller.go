package api

import (
	"net/http"
)

type controller struct {
	svc Service
}

func newController(s Service) *controller {
	return &controller{
		svc: s,
	}
}

func (c *controller) ListReleases(w http.ResponseWriter, r *http.Request) {
	deps, err := c.svc.ListReleases(r.Context())
	if err != nil {
		Error500(w, err)
		return
	}

	Success200(w, deps)
}
