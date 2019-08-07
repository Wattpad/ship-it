package api

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
)

// importing "k8s.io/helm/pkg/tiller" (specifically its transitive dependency
// on 'k8s.io/kubernetes' pkgs) breaks the build horribly.
// https://github.com/helm/helm/blob/master/pkg/tiller/release_server.go#L82
var tillerValidName = regexp.MustCompile("^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])+$")

type controller struct {
	svc Service
}

func newController(s Service) *controller {
	return &controller{
		svc: s,
	}
}

func (c *controller) ListReleases(w http.ResponseWriter, r *http.Request) {
	releases, err := c.svc.ListReleases(r.Context())
	if err != nil {
		Error500(w, err)
		return
	}

	Success200(w, releases)
}

func (c *controller) GetRelease(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := validateReleaseName(name); err != nil {
		Error400(w, err)
		return
	}

	release, err := c.svc.GetRelease(r.Context(), name)
	if err != nil {
		Error500(w, err)
		return
	}

	Success200(w, release)
}

func (c *controller) GetReleaseResources(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := validateReleaseName(name); err != nil {
		Error400(w, err)
		return
	}

	status, err := c.svc.GetReleaseResources(r.Context(), name)
	if err != nil {
		Error500(w, err)
		return
	}

	Success200(w, status)
}

func validateReleaseName(name string) error {
	if name == "" {
		return errors.New("missing release name")
	}

	// https: //github.com/helm/helm/blob/master/pkg/tiller/release_server.go#L50
	releaseNameMaxLen := 53

	if !tillerValidName.MatchString(name) || len(name) > releaseNameMaxLen {
		return errors.New("invalid release name")
	}

	return nil
}
