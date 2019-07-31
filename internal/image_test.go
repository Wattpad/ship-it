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

func TestGetImagePath(t *testing.T) {
	type testCase struct {
		serviceName string
		inputMap    map[string]interface{}
		expected    []string
	}

	testCases := map[string]testCase{
		"nested image found": {
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
		},
		"un-nested image found": {
			"bar",
			map[string]interface{}{
				"apples": "delicious",
				"image": map[string]interface{}{
					"repository": "foo/bar",
					"tag":        "baz",
				},
			},
			[]string{"image"},
		},
		"matches desired nested image": {
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
		},
		"matches desired un-nested image": {
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
		"desired image repo not found": {
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
			},
			[]string{},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			output := getImagePath(test.inputMap, test.serviceName)
			assert.Equal(t, test.expected, output)
		})
	}
}

func TestTable(t *testing.T) {
	var tests = []struct {
		name     string
		inputMap map[string]interface{}
		path     []string
		expected map[string]interface{}
	}{
		{
			"tabling a nested field",
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
					"image": map[string]interface{}{
						"repository": "foo/not-the-desired-image",
						"tag":        "baz",
					},
				},
			},
			[]string{"oranges", "image"},
			map[string]interface{}{
				"repository": "foo/not-the-desired-image",
				"tag":        "baz",
			},
		}, {
			"tabling a top level field",
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
					"image": map[string]interface{}{
						"repository": "foo/not-the-desired-image",
						"tag":        "baz",
					},
				},
			},
			[]string{"oranges"},
			map[string]interface{}{
				"taste": "delicious",
				"image": map[string]interface{}{
					"repository": "foo/not-the-desired-image",
					"tag":        "baz",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, table(test.inputMap, test.path))
		})
	}
}

func TestUpdateImage(t *testing.T) {
	var tests = []struct {
		newImage    Image
		inputMap    map[string]interface{}
		expectedMap map[string]interface{}
	}{
		{
			Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "newTag",
			},
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
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
				},
				"image": map[string]interface{}{
					"repository": "foo/bar",
					"tag":        "newTag",
				},
			},
		}, {
			Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "newTag",
			},
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
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
					"image": map[string]interface{}{
						"repository": "foo/bar",
						"tag":        "newTag",
					},
				},
			},
		}, {
			Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "newTag",
			},
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
				},
			},
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
				},
			},
		},
	}
	for _, test := range tests {
		update(test.inputMap, test.newImage)
		assert.Equal(t, test.expectedMap, test.inputMap)
	}
}

func TestWithImage(t *testing.T) {
	type testCase struct {
		name        string
		inputMap    map[string]interface{}
		inputImage  Image
		expectedMap map[string]interface{}
	}

	testCases := []testCase{
		{
			name: "Matching Image Case",
			inputMap: map[string]interface{}{
				"image": map[string]interface{}{
					"repository": "foo/bar",
					"tag":        "baz",
				},
			},
			inputImage: Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "a-new-tag",
			},
			expectedMap: map[string]interface{}{
				"image": map[string]interface{}{
					"repository": "foo/bar",
					"tag":        "a-new-tag",
				},
			},
		}, {
			name: "No Matching Image Case",
			inputMap: map[string]interface{}{
				"image": map[string]interface{}{
					"repository": "foo/bar",
					"tag":        "baz",
				},
			},
			inputImage: Image{
				Registry:   "foo",
				Repository: "oof",
				Tag:        "a-new-tag",
			},
			expectedMap: map[string]interface{}{
				"image": map[string]interface{}{
					"repository": "foo/bar",
					"tag":        "baz",
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			WithImage(test.inputImage, test.inputMap)
			assert.Equal(t, test.expectedMap, test.inputMap)
		})
	}
}

func TestStringMapCleanup(t *testing.T) {
	inputMap := map[string]interface{}{
		"foo": map[interface{}]interface{}{
			"bar": "baz",
		},
	}
	expectedMap := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
		},
	}
	assert.Equal(t, expectedMap, CleanUpStringMap(inputMap))
}

func TestDeepCopyMap(t *testing.T) {
	inputMap := map[string]interface{}{
		"a": map[string]interface{}{
			"b": "a-fake-value",
		},
		"c": "fake",
	}

	copiedMap := DeepCopyMap(inputMap)
	assert.Equal(t, inputMap, copiedMap)

	// Modify the input map and check that the copy does not change
	inputMap["c"] = "some new value"

	assert.NotEqual(t, inputMap, copiedMap)
}
