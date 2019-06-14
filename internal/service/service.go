package service

import (
	"context"
	"errors"

	"ship-it/internal/models"
)

var ErrNotImplemented = errors.New("not implemented")

func New() *Service {
	return &Service{}
}

type Service struct{}

func (s *Service) ListDeployments(ctx context.Context) ([]models.DeploymentDetail, error) {
	return nil, ErrNotImplemented
}
