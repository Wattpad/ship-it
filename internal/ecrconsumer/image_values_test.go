package ecrconsumer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
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

func TestGetImagePath(t *testing.T) {
	var tests = []struct {
		serviceName string
		inputMap    map[string]interface{}
		expected    []string
	}{
		{
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
		}, {
			"bar",
			map[string]interface{}{
				"apples": "delicious",
				"image": map[string]interface{}{
					"repository": "foo/bar",
					"tag":        "baz",
				},
			},
			[]string{"image"},
		}, {
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
		}, {
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
		}, {
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
		}, {
			"bar",
			map[string]interface{}{
				"apples": "delicious",
				"oranges": map[string]interface{}{
					"taste": "delicious",
					"image": map[string]interface{}{
						"repository": "foo",
						"tag":        "baz",
					},
				},
			},
			[]string{},
		},
	}
	for _, test := range tests {
		output := getImagePath(reflect.ValueOf(test.inputMap), test.serviceName)
		assert.Equal(t, test.expected, output)
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
			nil,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.expectedMap, update(test.inputMap, test.newImage))
	}
}

func TestParseImage(t *testing.T) {
	var tests = []struct {
		repo     string
		tag      string
		expected *Image
	}{
		{
			"foo/bar",
			"baz",
			&Image{
				Registry:   "foo",
				Repository: "bar",
				Tag:        "baz",
			},
		}, {
			"foo-bar",
			"baz",
			nil,
		},
	}
	for _, test := range tests {
		img, _ := parseImage(test.repo, test.tag)
		assert.Equal(t, test.expected, img)
	}
}

func TestWithImage(t *testing.T) {
	expectedImg := Image{
		Registry: "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "loki",
		Tag: "cc064f8a3d3fa0fe938e95d961ad0278770fa5d2",
	}

	rls, err := LoadRelease([]byte(crYaml))
	if err != nil {
		t.Fatal(err)
	}

	outputRls := WithImage(expectedImg, *rls)

	path := getImagePath(reflect.ValueOf(outputRls.Spec.Values), expectedImg.Repository)
	if path == nil {
		t.Fatal("no mathcing image found")
	}

	imgVals := table(outputRls.Spec.Values, path)

	outputImage, err := parseImage(imgVals["repository"].(string), imgVals["tag"].(string))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedImg, *outputImage)
}

func TestFullCycle(t *testing.T) {
	rls, err := LoadRelease([]byte(crYaml))
	assert.Nil(t, err)
	outBytes, err := yaml.Marshal(rls)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, crYaml, string(outBytes))
}
