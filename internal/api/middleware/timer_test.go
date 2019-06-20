package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
