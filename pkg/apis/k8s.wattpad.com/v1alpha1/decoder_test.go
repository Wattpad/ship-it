package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
)

const crYaml = `apiVersion: helmreleases.k8s.wattpad.com/v1alpha1
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
    autoscaler:
      maxPods: 50
      minPods: 30
      targetCPUUtilizationPercent: 60
    containerPort: 80
    cronjob:
      closeoutAppLabel: loki-closeout
      image:
        repository: 723255503624.dkr.ecr.us-east-1.amazonaws.com/kube-tools
        tag: deda27d
      schedule: 0 0 * * *
    image:
      repository: 723255503624.dkr.ecr.us-east-1.amazonaws.com/loki
      tag: cc064f8a3d3fa0fe938e95d961ad0278770fa5d2
    microservice:
      nameOverride: loki
    nodePort: 31828
    resources:
      limits:
        cpu: 500m
        memory: 256Mi
      requests:
        cpu: 500m
        memory: 128Mi
    securityContext:
      privileged: true
    serviceAccountName: loki
    servicePort: 80
status: {}
`

func LoadRelease(fileData []byte) (*HelmRelease, error) {
	rls := &HelmRelease{}
	d := Decoder

	err := runtime.DecodeInto(d, fileData, rls)
	if err != nil {
		return nil, err
	}

	return rls, nil
}

func TestDeepCopy(t *testing.T) {
	tests := []struct {
		original         HelmValues
		transformer      func(v HelmValues) HelmValues
		expectedCopy     HelmValues
		expectedOriginal HelmValues
	}{
		// reassignment
		{
			original: nil,
			transformer: func(v HelmValues) HelmValues {
				return HelmValues{}
			},
			expectedCopy:     HelmValues{},
			expectedOriginal: nil,
		},
		// mutation
		{
			original: HelmValues{},
			transformer: func(v HelmValues) HelmValues {
				v["foo"] = "bar"
				return v
			},
			expectedCopy:     HelmValues{"foo": "bar"},
			expectedOriginal: HelmValues{},
		},
		{
			original: HelmValues{"foo": HelmValues{"bar": "baz"}},
			transformer: func(v HelmValues) HelmValues {
				v["foo"] = HelmValues{"bar": "oof"}
				return v
			},
			expectedCopy:     HelmValues{"foo": HelmValues{"bar": "oof"}},
			expectedOriginal: HelmValues{"foo": HelmValues{"bar": "baz"}},
		},
	}
	for i, tc := range tests {
		output := HelmValues(make(map[string]interface{}))
		tc.original.DeepCopyInto(&output)
		assert.Exactly(t, tc.expectedCopy, tc.transformer(output), "copy was not mutated. test case: %d", i)
		assert.Exactly(t, tc.expectedOriginal, tc.original, "original was mutated. test case: %d", i)
	}
}

func TestSerializeRoundTrip(t *testing.T) {
	rls, err := LoadRelease([]byte(crYaml))
	assert.NoError(t, err)

	outBytes, err := yaml.Marshal(rls)
	assert.NoError(t, err)

	assert.Equal(t, crYaml, string(outBytes))
}
