package ecrconsumer

import (
	"testing"

	"ship-it/internal"

	shipitv1beta1 "ship-it-operator/api/v1beta1"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

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
			APIVersion: "helmreleases.shipit.wattpad.com/v1beta1",
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
		outputRls, err := WithImage(expectedImg, rls)
		assert.NoError(t, err)
		outputValues := outputRls.HelmValues()

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
		outputRls, err := WithImage(expectedImg, rls)
		assert.NoError(t, err)
		inputValues := rls.HelmValues()
		outputValues := outputRls.HelmValues()
		assert.Exactly(t, inputValues, outputValues)
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
