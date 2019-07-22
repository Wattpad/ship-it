package ecrconsumer

import (
	"encoding/json"
	"testing"

	"ship-it/internal"

	shipitv1beta1 "ship-it-operator/api/v1beta1"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

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
		inputMap map[string]interface{}
		path     []string
		expected map[string]interface{}
	}{
		{
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
		assert.Equal(t, test.expected, table(test.inputMap, test.path))
	}
}

func TestUpdateImage(t *testing.T) {
	var tests = []struct {
		newImage    internal.Image
		inputMap    map[string]interface{}
		expectedMap map[string]interface{}
	}{
		{
			internal.Image{
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
			internal.Image{
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
			internal.Image{
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
	rls := shipitv1beta1.HelmRelease{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HelmRelease",
			APIVersion: "apiVersion: helmreleases.shipit.wattpad.com/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "example-microservice",
		},
		Spec: shipitv1beta1.HelmReleaseSpec{
			ReleaseName: "example-release",
			Chart: shipitv1beta1.ChartSpec{
				Repository: "wattpad.s3.amazonaws.com/helm-charts",
				Path:       "microservice",
				Revision:   "HEAD",
			},
			Values: runtime.RawExtension{
				Raw: []byte(`{
					"image": {
						"repository": "bar/foo",
						"tag": "baz"
					}
				}`),
			},
		},
	}

	t.Run("Matching Image Case", func(t *testing.T) {
		expectedImg := internal.Image{
			Registry:   "bar",
			Repository: "foo",
			Tag:        "new-tag",
		}
		outputRls, _ := WithImage(expectedImg, rls)
		outputValues := getChartValues(outputRls)

		outputImg := outputValues["image"].(map[string]interface{})
		assert.Equal(t, expectedImg.Tag, outputImg["tag"].(string))
	})

	// Test No Matching Image Case
	t.Run("No Matching Image Case", func(t *testing.T) {
		expectedImg := internal.Image{
			Registry:   "bar",
			Repository: "oof",
			Tag:        "new-tag",
		}
		outputRls, _ := WithImage(expectedImg, rls)
		inputValues := getChartValues(rls)
		outputValues := getChartValues(outputRls)
		assert.Exactly(t, inputValues, outputValues)
	})
}

func getChartValues(r shipitv1beta1.HelmRelease) map[string]interface{} {
	var v map[string]interface{}
	json.Unmarshal(r.Spec.Values.Raw, &v)
	return v
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
