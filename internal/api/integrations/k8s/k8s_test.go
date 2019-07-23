package k8s

import (
	"testing"

	"ship-it/internal"

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
		expected *internal.Image
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
			&internal.Image{
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
			nil,
		},
	}

	for _, test := range tests {
		actualImg, _ := GetImageForRepo(test.repo, test.inputMap)
		assert.Equal(t, test.expected, actualImg)
	}
}

func TestFindRepo(t *testing.T) {

	type testCase struct {
		input    string
		expected string
	}

	tests := map[string]testCase{
		"valid repo with extra content appended": {
			"https://github.com/Wattpad/highlander/tree/master/wattpad/src/services/word-counts",
			"highlander",
		},
		"valid with repo with no extra content": {
			"https://github.com/Wattpad/miranda/",
			"miranda",
		},
		"no repo with trailing /": {
			"https://github.com/Wattpad/",
			"",
		},
		"no repo without trailing /": {
			"https://github.com/Wattpad",
			"",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repo, _ := findVCSRepo(test.input)
			assert.Equal(t, test.expected, repo)
		})
	}
}
