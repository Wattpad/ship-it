package service

import (
	"context"

	"ship-it/internal/api/models"

	v1 "k8s.io/api/core/v1"
	"k8s.io/helm/pkg/proto/hapi/release"
)

func New(l ReleaseLister, git GitCommands, repo string, ref string) *Service {
	return &Service{
		lister:     l,
		client:     git,
		repository: repo,
		ref:        ref,
	}
}

type ReleaseLister interface {
	ListAll(namespace string) ([]models.Release, error)
}

type GitCommands interface {
	GetTravisCIBuildURLForRef(ctx context.Context, repo string, ref string) (string, error)
}

type Service struct {
	lister     ReleaseLister
	client     GitCommands
	repository string
	ref        string
}

func (s *Service) ListReleases(ctx context.Context) ([]models.Release, error) {
	releases, err := s.lister.ListAll(v1.NamespaceAll)
	if err != nil {
		return nil, err
	}

	for i := range releases {
		r := &releases[i]

		r.Status = s.getReleaseStatus(r).String()
		r.Build.Travis = s.getTravisURL(ctx)

	}

	return releases, err
}

func (s *Service) getReleaseStatus(_ *models.Release) release.Status_Code {
	// The default state of a registered service is 'PENDING_INSTALL'.
	// There's only one possible release status for now.
	return release.Status_PENDING_INSTALL
}

func (s *Service) getTravisURL(ctx context.Context) string {
	url, err := s.client.GetTravisCIBuildURLForRef(ctx, s.repository, s.ref)
	if err != nil {
		return ""
	}
	return url
}
