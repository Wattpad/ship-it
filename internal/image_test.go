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
		tc := test
		t.Run(name, func(t *testing.T) {
			img, _ := ParseImage(tc.repo, tc.tag)
			assert.Equal(t, tc.expected, img)
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
		tc := test
		t.Run(name, func(t *testing.T) {
			output := getImagePath(tc.inputMap, tc.serviceName)
			assert.Equal(t, tc.expected, output)
		})
	}
}

func TestTable(t *testing.T) {
	type testCase struct {
		inputMap map[string]interface{}
		path     []string
		expected map[string]interface{}
	}

	testCases := map[string]testCase{
		"tabling a nested map level": {
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
		},
		"tabling the top level of the map": {
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
	for name, test := range testCases {
		tc := test
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, table(tc.inputMap, tc.path))
		})
	}
}

func TestWithImage(t *testing.T) {
	type testCase struct {
		inputMap    map[string]interface{}
		inputImage  Image
		expectedMap map[string]interface{}
	}

	testCases := map[string]testCase{
		"matching image found": {
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
		},
		"no matching image found": {
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

	for name, test := range testCases {
		tc := test
		t.Run(name, func(t *testing.T) {
			WithImage(tc.inputImage, tc.inputMap)
			assert.Equal(t, tc.expectedMap, tc.inputMap)
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
	assert.Equal(t, expectedMap, cleanUpStringMap(inputMap))
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
	inputMap["a"].(map[string]interface{})["b"] = "some new value"

	assert.NotEqual(t, inputMap, copiedMap)
}
