package service

import (
	"context"

	"ship-it/internal/models"

	v1 "k8s.io/api/core/v1"
)

func New(k K8sInterface) (*Service, error) {
	return &Service{k}, nil
}

type K8sInterface interface {
	ListAll(namespace string) ([]models.Release, error)
}

type Service struct {
	kube K8sInterface
}

func (s *Service) ListReleases(ctx context.Context) ([]models.Release, error) {
	return s.kube.ListAll(v1.NamespaceAll)
}
