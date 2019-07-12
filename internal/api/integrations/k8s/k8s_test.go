package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnnotationFor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"word-counts", "helmreleases.k8s.wattpad.com/word-counts"},
		{"", "helmreleases.k8s.wattpad.com/"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, annotationFor(test.input))
	}
}

func TestGetImageForRepo(t *testing.T) {

}

func TestFindRepo(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://github.com/Wattpad/highlander/tree/master/wattpad/src/services/word-counts", "highlander"},
		{"https://github.com/Wattpad/miranda/", "miranda"},
		{"https://github.com/Wattpad/", ""},
		{"https://github.com/Wattpad", ""},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, findVCSRepo(test.input))
	}
}
