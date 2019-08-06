package service

import (
	"context"

	"ship-it/internal/api/models"

	v1 "k8s.io/api/core/v1"
)

type ReleaseLister interface {
	Get(ctx context.Context, namespace, release string) (*models.Release, error)
	List(ctx context.Context, namespace string) ([]models.Release, error)
}

type ResourcesGetter interface {
	Get(release string) (string, error)
}

func New(rl ReleaseLister, rg ResourcesGetter) *Service {
	return &Service{
		Namespace: v1.NamespaceAll,
		releases:  rl,
		resources: rg,
	}
}

type Service struct {
	Namespace string
	releases  ReleaseLister
	resources ResourcesGetter
}

func (s *Service) ListReleases(ctx context.Context) ([]models.Release, error) {
	return s.releases.List(ctx, s.Namespace)
}

func (s *Service) GetRelease(ctx context.Context, name string) (*models.Release, error) {
	return s.releases.Get(ctx, s.Namespace, name)
}

func (s *Service) GetReleaseResources(ctx context.Context, name string) (*models.ReleaseResources, error) {
	resources, err := s.resources.Get(name)
	if err != nil {
		return nil, err
	}

	return &models.ReleaseResources{
		Name:      name,
		Resources: resources,
	}, nil
}
