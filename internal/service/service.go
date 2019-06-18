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

func (s *Service) ListReleases(ctx context.Context) ([]models.Release, error) {
	return nil, ErrNotImplemented
}
