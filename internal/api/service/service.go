package service

import (
	"context"

	"ship-it/internal/api/models"

	v1 "k8s.io/api/core/v1"
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
	return s.lister.ListAll(v1.NamespaceAll)
}
