package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"ship-it/internal/api/models"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) ListReleases(ctx context.Context) ([]models.Release, error) {
	args := m.Called(ctx)

	var ret0 []models.Release
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.([]models.Release)
	}

	return ret0, args.Error(1)
}

func (m *mockService) GetRelease(ctx context.Context, name string) (*models.Release, error) {
	args := m.Called(ctx, name)

	var ret0 *models.Release
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.(*models.Release)
	}

	return ret0, args.Error(1)
}

func (m *mockService) GetReleaseResources(ctx context.Context, name string) (*models.ReleaseResources, error) {
	args := m.Called(ctx, name)

	var ret0 *models.ReleaseResources
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.(*models.ReleaseResources)
	}

	return ret0, args.Error(1)
}

func withRouteContext(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestListReleases(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/releases", nil)

	t.Run("endpoint returns 200 on success", func(t *testing.T) {
		var m mockService
		m.On("ListReleases", context.Background()).Return(nil, errors.New("internal error"))

		rec := httptest.NewRecorder()

		c := newController(&m)
		c.ListReleases(rec, req)

		m.AssertExpectations(t)
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("endpoint returns 500 for internal error", func(t *testing.T) {
		var m mockService
		m.On("ListReleases", context.Background()).Return([]models.Release{}, nil)

		rec := httptest.NewRecorder()

		c := newController(&m)
		c.ListReleases(rec, req)

		m.AssertExpectations(t)
		assert.Equal(t, rec.Code, http.StatusOK)
	})
}

func TestGetReleaseResources(t *testing.T) {
	testRelease := "test-release"
	invalidRelease := "bad$release#name"

	t.Run("endpoint returns 200 on success", func(t *testing.T) {
		var m mockService
		m.On("GetReleaseResources", mock.Anything, testRelease).Return(&models.ReleaseResources{}, nil)

		rec := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/releases/%s/resources", testRelease), nil)
		req = withRouteContext(req, "name", testRelease)

		c := newController(&m)
		c.GetReleaseResources(rec, req)

		m.AssertExpectations(t)
		assert.Equal(t, rec.Code, http.StatusOK)
	})

	t.Run("endpoint returns 400 for invalid request", func(t *testing.T) {
		var m mockService

		rec := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/releases/%s/resources", invalidRelease), nil)
		req = withRouteContext(req, "name", invalidRelease)

		c := newController(&m)
		c.GetReleaseResources(rec, req)

		m.AssertNotCalled(t, "GetReleaseResources")
		assert.Equal(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("endpoint returns 500 for internal error", func(t *testing.T) {
		var m mockService
		m.On("GetReleaseResources", mock.Anything, testRelease).Return(nil, errors.New("internal error"))

		rec := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/releases/%s/resources", testRelease), nil)
		req = withRouteContext(req, "name", testRelease)

		c := newController(&m)
		c.GetReleaseResources(rec, req)

		m.AssertExpectations(t)
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
	})
}
