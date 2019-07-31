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
	v1 "k8s.io/api/core/v1"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) ListReleases(ctx context.Context) ([]models.Release, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Release), args.Error(1)
}

func (m *mockService) GetRelease(ctx context.Context, name string) (*models.Release, error) {
	args := m.Called(ctx)

	var ret0 *models.Release
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.(*models.Release)
	}

	return ret0, args.Error(1)
}

func (m *mockService) GetReleaseResources(ctx context.Context, name string) (*models.ReleaseResources, error) {
	args := m.Called(ctx)

	var ret0 *models.ReleaseResources
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.(*models.ReleaseResources)
	}

	return ret0, args.Error(1)
}

func TestListReleasesReturns500ForInternalError(t *testing.T) {
	var m mockService
	m.On("ListReleases", context.Background()).Return([]models.Release{}, errors.New("internal error"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/releases", nil)

	c := newController(&m)
	c.ListReleases(rec, req)

	mock.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusInternalServerError)
}

func TestListReleasesReturns200OnSuccess(t *testing.T) {
	var m mockService
	m.On("ListReleases", context.Background()).Return([]models.Release{}, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/releases", nil)

	c := newController(&m)
	c.ListReleases(rec, req)

	m.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusOK)
}

func withChiRouteContext(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestGetReleaseReturns200OnSuccess(t *testing.T) {
	testRelease := "test-release"

	var m mockService
	m.On("GetRelease", mock.Anything, v1.NamespaceAll, testRelease).Return(&models.Release{}, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/releases/%s", testRelease), nil)

	req = withChiRouteContext(req, "name", testRelease)

	c := newController(&m)
	c.GetRelease(rec, req)

	m.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusOK)
}

func TestGetReleaseReturns400ForInvalidReleaseName(t *testing.T) {
	testRelease := "bad_release_name"

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/releases/%s", testRelease), nil)

	req = withChiRouteContext(req, "name", testRelease)

	m := new(mockService)
	c := newController(m)

	c.GetRelease(rec, req)

	m.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestGetReleaseReturns400ForInvalidReleaseName(t *testing.T) {
	testRelease := "test-release"

	var m mockService
	m.On("GetRelease", mock.Anything, v1.NamespaceAll, testRelease).Return(nil, errors.New("internal error"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/releases/%s", testRelease), nil)

	req = withChiRouteContext(req, "name", testRelease)

	c := newController(&m)
	c.GetRelease(rec, req)

	m.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusInternalServerError)
}

func TestGetReleaseResourcesReturns200OnSuccess(t *testing.T) {
	testRelease := "test-release"

	var m mockService
	m.On("GetReleaseResources", mock.Anything, testRelease).Return(&models.ReleaseResources{}, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/releases/%s/resources", testRelease), nil)

	req = withChiRouteContext(req, "name", testRelease)

	c := newController(&m)
	c.GetReleaseResources(rec, req)

	m.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusOK)
}

func TestGetReleaseResourcesReturns400ForInvalidReleaseName(t *testing.T) {
	invalidRelease := "bad_release_name"

	var m mockService

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/releases/%s/resources", invalidRelease), nil)

	req = withChiRouteContext(req, "name", invalidRelease)

	c := newController(&m)
	c.GetReleaseResources(rec, req)

	m.AssertNotCalled(t, "GetReleaseResources")
	assert.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestGetReleaseResourcesReturns500ForInternalError(t *testing.T) {
	testRelease := "test-release"

	var m mockService
	m.On("GetReleaseResources", mock.Anything, testRelease).Return(nil, errors.New("internal error"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/releases/%s/resources", testRelease), nil)

	req = withChiRouteContext(req, "name", testRelease)

	c := newController(&m)
	c.GetReleaseResources(rec, req)

	m.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusInternalServerError)
}
