package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"ship-it/internal"
	"ship-it/internal/api/models"

	shipitv1beta1 "ship-it-operator/api/v1beta1"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAnnotationFor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"word-counts", "helmreleases.shipit.wattpad.com/word-counts"},
		{"", "helmreleases.shipit.wattpad.com/"},
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

func newTestK8sClient(objs ...runtime.Object) *K8sClient {
	scheme := runtime.NewScheme()
	shipitv1beta1.AddToScheme(scheme)

	return &K8sClient{
		client: fake.NewFakeClientWithScheme(scheme, objs...),
	}
}

func TestListAll(t *testing.T) {
	autodeploy := true
	created := metav1.Unix(42, 0) // arbitrary, but non-zero nanoseconds causes rounding issues
	dashboard := "dashboard"
	github := "github"
	release := "release"
	slack := "slack"
	squad := "squad"
	sumologic := "sumologic"
	chartRepo := "chartRepo"
	chartPath := "chartPath"
	chartVersion := "chartVersion"
	dockerImage := "docker/foo"
	dockerImageTag := "aoeuhtns"

	expectedRelease := models.Release{
		Name:       release,
		Created:    created.Time,
		AutoDeploy: autodeploy,
		Code: models.SourceCode{
			Github: github,
		},
		Monitoring: models.Monitoring{
			Datadog: models.Datadog{
				Dashboard: dashboard,
			},
			Sumologic: sumologic,
		},
		Artifacts: models.Artifacts{
			Chart: models.HelmArtifact{
				Path:       chartPath,
				Repository: chartRepo,
				Version:    chartVersion,
			},
			Docker: []models.DockerArtifact{
				{
					Image: dockerImage,
					Tag:   dockerImageTag,
				},
			},
		},
		Owner: models.Owner{
			Squad: squad,
			Slack: slack,
		},
	}

	values := map[string]interface{}{
		"foo": map[string]interface{}{
			// this should appear in the final DockerArtifacts
			"image": map[string]interface{}{
				"repository": dockerImage,
				"tag":        dockerImageTag,
			},
		},
		"bar": map[string]interface{}{
			// this should not, since the image fields are invalid
			"image": map[string]interface{}{
				"badbad":  "docker/bar",
				"notgood": "htns",
			},
		},
	}

	valuesRaw, err := json.Marshal(values)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	k8sRelease := shipitv1beta1.HelmRelease{
		TypeMeta: metav1.TypeMeta{
			Kind: "HelmRelease",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              release,
			Namespace:         v1.NamespaceDefault,
			CreationTimestamp: created,
			Annotations: map[string]string{
				"helmreleases.shipit.wattpad.com/code":       github,
				"helmreleases.shipit.wattpad.com/autodeploy": fmt.Sprintf("%t", autodeploy),
				"helmreleases.shipit.wattpad.com/squad":      squad,
				"helmreleases.shipit.wattpad.com/slack":      slack,
				"helmreleases.shipit.wattpad.com/datadog":    dashboard,
				"helmreleases.shipit.wattpad.com/sumologic":  sumologic,
			},
		},
		Spec: shipitv1beta1.HelmReleaseSpec{
			ReleaseName: release,
			Chart: shipitv1beta1.ChartSpec{
				Repository: chartRepo,
				Path:       chartPath,
				Revision:   chartVersion,
			},
			Values: runtime.RawExtension{
				Raw: valuesRaw,
			},
		},
	}

	client := newTestK8sClient(&k8sRelease)

	apiReleases, err := client.ListAll(context.Background(), v1.NamespaceAll)
	if assert.NoError(t, err) {
		assert.Len(t, apiReleases, 1)
		assert.Equal(t, expectedRelease, apiReleases[0])
	}
}
