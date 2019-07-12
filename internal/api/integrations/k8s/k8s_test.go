package k8s

import (
	"os"
	"testing"

	"ship-it/internal"

	"github.com/go-kit/kit/log"

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
	tests := []struct {
		inputMap map[string]interface{}
		repo     string
		expected internal.Image
	}{
		{
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
				},
				"image": map[string]interface{}{
					"repository": "foo/bar",
					"tag":        "baz",
				},
			},
			"bar",
			internal.Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "baz",
			},
		}, {
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
				},
				"image": map[string]interface{}{
					"repository": "foo/oof",
					"tag":        "baz",
				},
			},
			"bar",
			internal.Image{},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, GetImageForRepo(test.repo, test.inputMap, log.NewJSONLogger(log.NewSyncWriter(os.Stdout))))
	}
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
