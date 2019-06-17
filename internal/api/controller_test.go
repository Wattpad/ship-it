package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ship-it/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) ListReleases(ctx context.Context) ([]models.Release, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Release), args.Error(1)
}

func TestListReleasesReturns500ForInternalError(t *testing.T) {
	mock := new(mockService)
	mock.On("ListReleases", context.Background()).Return([]models.Release{}, errors.New("internal error"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/releases", nil)

	c := newController(mock)
	c.ListReleases(rec, req)

	mock.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusInternalServerError)
}

func TestListReleasesReturns200OnSuccess(t *testing.T) {
	mock := new(mockService)
	mock.On("ListReleases", context.Background()).Return([]models.Release{}, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/releases", nil)

	c := newController(mock)
	c.ListReleases(rec, req)

	mock.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusOK)
}
