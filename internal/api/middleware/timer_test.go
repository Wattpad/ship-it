package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/metrics"
	"github.com/stretchr/testify/assert"
)

type mockHistogram struct {
	labels        []string
	observeCalled bool
}

func (m *mockHistogram) Observe(t float64) {
	m.observeCalled = true
}

func (m *mockHistogram) With(kvs ...string) metrics.Histogram {
	m.labels = kvs
	return m
}

func withRouteContext(req *http.Request) *http.Request {
	rctx := chi.NewRouteContext()
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestGetIdentifier(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"/", "root"},
		{"/dashboard/", "dashboard"},
		{"/releases/{name}/resources/", "releases.name.resources"},
		{"/releases/{name}/resources/{pod}/", "releases.name.resources.pod"},
	}
	assert := assert.New(t)
	for _, test := range tests {
		assert.Equal(test.expected, getIdentifier(test.input))
	}
}

func TestTimer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var histo mockHistogram

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	req = withRouteContext(req)

	Timer(&histo)(handler).ServeHTTP(rec, req)

	t.Run("it adds method and endpoint labels", func(t *testing.T) {
		assert.Contains(t, histo.labels, "method")
		assert.Contains(t, histo.labels, "endpoint")
	})

	t.Run("it observes the duration of the HTTP handler", func(t *testing.T) {
		assert.True(t, histo.observeCalled)
	})
}
