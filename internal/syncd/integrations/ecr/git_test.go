package ecr

import (
	"testing"

	"ship-it/internal"

	"github.com/stretchr/testify/assert"
)

const inputYaml = `apiVersion: helmreleases.k8s.wattpad.com/v1alpha1
kind: HelmRelease
metadata:
  creationTimestamp: null
  name: example-microservice
spec:
  chart:
    path: microservice
    repository: wattpad.s3.amazonaws.com/helm-charts
    revision: HEAD
  releaseName: example-release
  values:
    image:
      repository: 723255503624.dkr.ecr.us-east-1.amazonaws.com/bar
      tag: baz
status: {}
`

const expectedYaml = `apiVersion: helmreleases.k8s.wattpad.com/v1alpha1
kind: HelmRelease
metadata:
  creationTimestamp: null
  name: example-microservice
spec:
  chart:
    path: microservice
    repository: wattpad.s3.amazonaws.com/helm-charts
    revision: HEAD
  releaseName: example-release
  values:
    image:
      repository: 723255503624.dkr.ecr.us-east-1.amazonaws.com/bar
      tag: 78bc9ccf64eb838c6a0e0492ded722274925e2bd
status: {}
`

func TestTransformBytes(t *testing.T) {
	inputImage := &internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}

	t.Run("Valid yaml bytes", func(t *testing.T) {
		outBytes, err := transformBytes([]byte(inputYaml), inputImage)
		assert.NoError(t, err)
		assert.Equal(t, []byte(expectedYaml), outBytes)
	})

	t.Run("Invalid yaml bytes", func(t *testing.T) {
		outBytes, err := transformBytes([]byte("some invalid yaml"), inputImage)
		assert.Error(t, err)
		assert.Nil(t, outBytes)
	})
}
