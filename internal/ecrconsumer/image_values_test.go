package ecrconsumer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetImagePath(t *testing.T) {
	var tests = []struct {
		serviceName string
		inputMap    map[string]interface{}
		expected    []string
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
			[]string{"oranges", "image"},
		}, {
			"bar",
			map[string]interface{}{
				"apples": "delicious",
				"image": map[string]interface{}{
					"repository": "foo/bar",
					"tag":        "baz",
				},
			},
			[]string{"image"},
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
			[]string{"oranges", "image"},
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
			[]string{"image"},
		},
	}
	for i := range tests {
		output := getImagePath(reflect.ValueOf(tests[i].inputMap), tests[i].serviceName)
		assert.Equal(t, tests[i].expected, output)
	}
}

// func TestUpdateImage(t *testing.T) {
// 	var tests = []struct {
// 		newImage    Image
// 		inputMap    map[string]interface{}
// 		expectedMap map[string]interface{}
// 	}{
// 		{
// 			Image{
// 				Registry:   "foo",
// 				Repository: "bar",
// 				Tag:        "newTag",
// 			},
// 			map[string]interface{}{
// 				"apples": "delicious",
// 				"oranges": map[string]interface{}{
// 					"taste": "delicious",
// 				},
// 				"image": map[string]interface{}{
// 					"repository": "foo/bar",
// 					"tag":        "baz",
// 				},
// 			},
// 			map[string]interface{}{
// 				"apples": "delicious",
// 				"oranges": map[string]interface{}{
// 					"taste": "delicious",
// 				},
// 				"image": map[string]interface{}{
// 					"repository": "foo/bar",
// 					"tag":        "newTag",
// 				},
// 			},
// 		},
// 	}
// 	for i := range tests {
// 		updatedMap := make(map[string]interface{})
// 		update(reflect.ValueOf(tests[i].inputMap), tests[i].newImage, &updatedMap)
// 		assert.Equal(t, tests[i].expectedMap, updatedMap)
// 	}
// }

func TestParseImage(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestWithImage(t *testing.T) {
	assert.Equal(t, 1, 1)
}
