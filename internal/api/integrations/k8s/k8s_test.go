package k8s

import (
	"encoding/json"
	"testing"

	shipitv1beta1 "ship-it-operator/api/v1beta1"
	"ship-it/internal/api/models"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestDockerArtifacts(t *testing.T) {
	values := map[string]interface{}{
		"foo": map[string]interface{}{
			"image": map[string]interface{}{
				"repository": "docker/foo",
				"tag":        "aoeu",
			},
		},
		"bar": map[string]interface{}{
			"image": map[string]interface{}{
				"repository": "docker/bar",
				"tag":        "htns",
			},
		},
	}

	valuesRaw, err := json.Marshal(values)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	hr := shipitv1beta1.HelmRelease{
		Spec: shipitv1beta1.HelmReleaseSpec{
			Values: runtime.RawExtension{
				Raw: valuesRaw,
			},
		},
	}

	expected := []models.DockerArtifact{
		{
			Image: "docker/foo",
			Tag:   "aoeu",
		},
		{
			Image: "docker/bar",
			Tag:   "htns",
		},
	}

	assert.ElementsMatch(t, expected, dockerArtifacts(hr))
}
