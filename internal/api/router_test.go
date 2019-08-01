package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"ship-it/internal/api/models"
	"testing"

	"github.com/go-kit/kit/metrics/discard"
	"github.com/stretchr/testify/assert"
)

type nopService struct{}

func (m *nopService) ListReleases(ctx context.Context) ([]models.Release, error) {
	return nil, nil
}

func TestRoutes(t *testing.T) {
	requests := []*http.Request{
		httptest.NewRequest(http.MethodGet, "/index.html", nil),
		httptest.NewRequest(http.MethodGet, "/health", nil),
		httptest.NewRequest(http.MethodGet, "/api/releases", nil),
	}

	handler := New(&nopService{}, discard.NewHistogram())

	for _, req := range requests {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.NotEqual(t, rec.Code, http.StatusNotFound)
	}
}
