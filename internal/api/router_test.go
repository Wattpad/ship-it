package api

import (
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	expectedRoutes := map[string]interface{}{
		"/*":      nil,
		"/health": nil,
		"/api/*": map[string]interface{}{
			"/releases": nil,
		},
	}

	handler := New(nil, nil)

	if routes, ok := handler.(chi.Routes); ok {
		assertRoutes(t, expectedRoutes, routes)
	}
}

func assertRoutes(t *testing.T, expected map[string]interface{}, r chi.Routes) {
	for _, route := range r.Routes() {
		assert.Contains(t, expected, route.Pattern)
		if nested, ok := expected[route.Pattern].(map[string]interface{}); ok {
			assertRoutes(t, nested, route.SubRoutes)
		}
	}
}
