package middleware

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var tests = []struct {
	input    string
	expected string
}{
	{"/dashboard/", "dashboard."},
	{"/releases/{name}/resources/", "releases.name.resources."},
	{"/releases/{name}/resources/{pod}/", "releases.name.resources.pod."},
}

func TestGetIdentifier(t *testing.T) {
	assert := assert.New(t)
	for _, test := range tests {
		assert.Equal(test.expected, getIdentifier(test.input))
	}
}
