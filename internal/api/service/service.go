package service

import (
	"context"

	"ship-it/internal/api/models"

	v1 "k8s.io/api/core/v1"
	"k8s.io/helm/pkg/proto/hapi/release"
)

func New(l ReleaseLister) *Service {
	return &Service{
		lister: l,
	}
}

type ReleaseLister interface {
	ListAll(namespace string) ([]models.Release, error)
}

type Service struct {
	lister ReleaseLister
}

func (s *Service) ListReleases(ctx context.Context) ([]models.Release, error) {
	releases, err := s.lister.ListAll(v1.NamespaceAll)
	if err != nil {
		return nil, err
	}

	for i := range releases {
		r := &releases[i]

		r.Status = s.getReleaseStatus(r).String()
	}

	return releases, err
}

func (s *Service) getReleaseStatus(_ *models.Release) release.Status_Code {
	// The default state of a registered service is 'PENDING_INSTALL'.
	// There's only one possible release status for now.
	return release.Status_PENDING_INSTALL
}
