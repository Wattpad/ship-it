package service

import (
	"context"

	"ship-it/internal/api/models"

	v1 "k8s.io/api/core/v1"
)

type ReleaseLister interface {
	List(ctx context.Context, namespace string) ([]models.Release, error)
}

type ResourcesGetter interface {
	Get(namespace, release string) (string, error)
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

func (s *Service) GetReleaseResources(ctx context.Context, name string) (*models.ReleaseResources, error) {
	resources, err := s.resources.Get(s.Namespace, name)
	if err != nil {
		return nil, err
	}

	return &models.ReleaseResources{
		Name:      name,
		Resources: resources,
	}, nil
}
