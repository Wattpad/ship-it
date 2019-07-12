package service

import (
	"context"

	"ship-it/internal/api/models"

	"github.com/go-kit/kit/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/helm/pkg/proto/hapi/release"
)

func New(l ReleaseLister, travis TravisChecker, logger log.Logger) *Service {
	return &Service{
		lister: l,
		travis: travis,
		logger: logger,
	}
}

type ReleaseLister interface {
	ListAll(namespace string) ([]models.Release, error)
}

type TravisChecker interface {
	GetTravisCIBuildURLForRef(ctx context.Context, repo string, ref string) (string, error)
}

type Service struct {
	lister ReleaseLister
	travis TravisChecker
	logger log.Logger
}

func (s *Service) ListReleases(ctx context.Context) ([]models.Release, error) {
	releases, err := s.lister.ListAll(v1.NamespaceAll)
	if err != nil {
		return nil, err
	}

	for i := range releases {
		r := &releases[i]

		r.Status = s.getReleaseStatus(r).String()
		r.Build.Travis = s.getTravisURL(ctx, *r)

	}

	return releases, err
}

func (s *Service) getReleaseStatus(_ *models.Release) release.Status_Code {
	// The default state of a registered service is 'PENDING_INSTALL'.
	// There's only one possible release status for now.
	return release.Status_PENDING_INSTALL
}

func (s *Service) getTravisURL(ctx context.Context, r models.Release) string {
	url, err := s.travis.GetTravisCIBuildURLForRef(ctx, r.Code.Github, r.Code.Ref)
	if err != nil {
		s.logger.Log("Failed to fetch build URL")
		return ""
	}
	return url
}
