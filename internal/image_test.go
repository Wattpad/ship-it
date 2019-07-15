package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseImage(t *testing.T) {
	type testCase struct {
		repo     string
		tag      string
		expected *Image
	}

	testCases := map[string]testCase{
		"valid image": {
			"foo/bar",
			"baz",
			&Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "baz",
			},
		},
		"invalid image": {
			"foo-bar",
			"baz",
			nil,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			img, _ := ParseImage(test.repo, test.tag)
			assert.Equal(t, test.expected, img)
		})
	}
}
