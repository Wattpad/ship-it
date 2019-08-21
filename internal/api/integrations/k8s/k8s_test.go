package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	shipitv1beta1 "ship-it-operator/api/v1beta1"
	"ship-it/internal/api/models"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/helm/pkg/proto/hapi/release"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newTestK8sClient(objs ...runtime.Object) *K8sClient {
	scheme := runtime.NewScheme()
	shipitv1beta1.AddToScheme(scheme)

	return &K8sClient{
		client: fake.NewFakeClientWithScheme(scheme, objs...),
	}
}

func TestGetAndList(t *testing.T) {
	autodeploy := true
	chartPath := "chartPath"
	chartRepo := "chartRepo"
	chartVersion := "chartVersion"
	created := metav1.Unix(42, 0) // arbitrary, but non-zero nanoseconds causes rounding issues
	dashboard := "dashboard"
	dockerImage := "docker/foo"
	dockerImageTag := "aoeuhtns"
	github := "github"
	releaseName := "releaseName"
	releaseStatus := release.Status_DEPLOYED
	slack := "slack"
	squad := "squad"
	sumologic := "sumologic"

	expectedRelease := models.Release{
		Name:       releaseName,
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
		Status: releaseStatus.String(),
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
			Name:              releaseName,
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
			ReleaseName: releaseName,
			Chart: shipitv1beta1.ChartSpec{
				Repository: chartRepo,
				Name:       chartPath,
				Version:    chartVersion,
			},
			Values: runtime.RawExtension{
				Raw: valuesRaw,
			},
		},
		Status: shipitv1beta1.HelmReleaseStatus{
			Conditions: []shipitv1beta1.HelmReleaseCondition{
				{
					Type: releaseStatus.String(),
				},
			},
		},
	}

	client := newTestK8sClient(&k8sRelease)

	apiReleases, err := client.List(context.Background(), v1.NamespaceAll)
	if assert.NoError(t, err) {
		assert.Len(t, apiReleases, 1)
		assert.Equal(t, expectedRelease, apiReleases[0])
	}

	apiRelease, err := client.Get(context.Background(), v1.NamespaceAll, releaseName)
	if assert.NoError(t, err) && assert.NotNil(t, apiRelease) {
		assert.Equal(t, expectedRelease, *apiRelease)
	}
}
