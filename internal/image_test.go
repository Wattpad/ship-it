package internal

import (
	"testing"

	"ship-it/pkg/apis/k8s.wattpad.com/v1alpha1"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		name string
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
	rls := v1alpha1.HelmRelease{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HelmRelease",
			APIVersion: "apiVersion: helmreleases.k8s.wattpad.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "example-microservice",
		},
		Spec: v1alpha1.HelmReleaseSpec{
			ReleaseName: "example-release",
			Chart: v1alpha1.ChartSpec{
				Repository: "wattpad.s3.amazonaws.com/helm-charts",
				Path:       "microservice",
				Revision:   "HEAD",
			},
			Values: map[string]interface{}{
				"image": map[string]interface{}{
					"repository": "bar/foo",
					"tag":        "baz",
				},
			},
		},
	}

	t.Run("Matching Image Case", func(t *testing.T) {
		expectedImg := Image{
			Registry:   "bar",
			Repository: "foo",
			Tag:        "new-tag",
		}
		outputRls := WithImage(expectedImg, rls)

		outputImg := outputRls.Spec.Values["image"].(map[string]interface{})
		assert.Equal(t, expectedImg.Tag, outputImg["tag"].(string))
	})

	// Test No Matching Image Case
	t.Run("No Matching Image Case", func(t *testing.T) {
		expectedImg := Image{
			Registry:   "bar",
			Repository: "oof",
			Tag:        "new-tag",
		}
		outputRls := WithImage(expectedImg, rls)
		assert.Exactly(t, rls, outputRls)
	})
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
