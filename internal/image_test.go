package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseImage(t *testing.T) {
	var tests = []struct {
		repo     string
		tag      string
		expected *Image
	}{
		{
			"foo/bar",
			"baz",
			&Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "baz",
			},
		}, {
			"foo-bar",
			"baz",
			nil,
		},
	}
	for _, test := range tests {
		img, _ := ParseImage(test.repo, test.tag)
		assert.Equal(t, test.expected, img)
	}
}
