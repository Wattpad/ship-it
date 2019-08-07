package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/metrics/discard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called()
}

type mockController struct {
	mock.Mock
}

func (m *mockController) ListReleases(http.ResponseWriter, *http.Request) {
	m.Called()
}

func (m *mockController) Health(http.ResponseWriter, *http.Request) {
	m.Called()
}

func TestRootRoute(t *testing.T) {
	type testCase struct {
		expect     func(*mock.Mock)
		statusCode int
		request    *http.Request
	}

	testCases := []testCase{
		{
			expect:     func(m *mock.Mock) { m.On("ServeHTTP") },
			statusCode: http.StatusOK,
			request:    httptest.NewRequest(http.MethodGet, "/does/not/exist", nil),
		},
		{
			expect:     func(m *mock.Mock) { m.On("ServeHTTP") },
			statusCode: http.StatusNotFound,
			request:    httptest.NewRequest(http.MethodGet, "/does/not/exist", nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.request.URL.Path, func(t *testing.T) {
			var m mockHandler
			tc.expect(&m.Mock)

			rec := httptest.NewRecorder()

			handler := newRouter(&m, new(mockController), discard.NewHistogram())
			handler.ServeHTTP(rec, tc.request)

			assert.Equal(t, rec.Code, http.StatusOK)

			m.AssertExpectations(t)
		})

	}
}

func TestControllerRoutes(t *testing.T) {
	type testCase struct {
		expect     func(*mock.Mock)
		statusCode int
		request    *http.Request
	}

	testCases := []testCase{
		{
			expect:     func(m *mock.Mock) { m.On("Health") },
			statusCode: http.StatusOK,
			request:    httptest.NewRequest(http.MethodGet, "/health", nil),
		},
		{
			expect:     func(m *mock.Mock) { m.On("ListReleases") },
			statusCode: http.StatusOK,
			request:    httptest.NewRequest(http.MethodGet, "/api/releases", nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.request.URL.Path, func(t *testing.T) {
			var m mockController
			tc.expect(&m.Mock)

			rec := httptest.NewRecorder()

			handler := newRouter(new(mockHandler), &m, discard.NewHistogram())
			handler.ServeHTTP(rec, tc.request)

			m.AssertExpectations(t)
			assert.Equal(t, rec.Code, tc.statusCode)
		})
	}
}
