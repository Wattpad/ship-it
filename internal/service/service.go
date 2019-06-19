package service

import (
	"context"
	"errors"

	"ship-it/internal/models"

	"github.com/go-kit/kit/metrics/dogstatsd"
)

var ErrNotImplemented = errors.New("not implemented")

func New(t *dogstatsd.Timing) *Service {
	return &Service{
		DDTimer: t,
	}
}

type Service struct {
	DDTimer *dogstatsd.Timing
}

func (s *Service) ListReleases(ctx context.Context) ([]models.Release, error) {
	return nil, ErrNotImplemented
}
