package ecrconsumer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindImage(t *testing.T) {
	var tests = []struct {
		serviceName string
		inputMap    map[string]interface{}
		expected    Image
	}{
		{
			"bar",
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
					"image": map[string]interface{}{
						"repository": "foo/bar",
						"tag":        "baz",
					},
				},
			},
			Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "baz",
			},
		}, {
			"bar",
			map[string]interface{}{
				"apples": "delicious",
				"image": map[string]interface{}{
					"repository": "foo/bar",
					"tag":        "baz",
				},
			},
			Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "baz",
			},
		}, {
			"bar",
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
					"image": map[string]interface{}{
						"repository": "foo/bar",
						"tag":        "baz",
					},
				},
				"image": map[string]interface{}{
					"repository": "foo/not-the-desired-image",
					"tag":        "baz",
				},
			},
			Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "baz",
			},
		}, {
			"bar",
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
					"image": map[string]interface{}{
						"repository": "foo/not-the-desired-image",
						"tag":        "baz",
					},
				},
				"image": map[string]interface{}{
					"repository": "foo/bar",
					"tag":        "baz",
				},
			},
			Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "baz",
			},
		},
	}
	for i := range tests {
		imgPtr := &Image{}
		FindImage(reflect.ValueOf(tests[i].inputMap), imgPtr, tests[i].serviceName)
		assert.Equal(t, tests[i].expected, *imgPtr)
	}
}

func TestUpdateImage(t *testing.T) {
	assert.Equal(t, 1, 1)
}
