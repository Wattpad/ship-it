package release

import (
	"fmt"
	"io/ioutil"
	"testing"

	"ship-it/pkg/apis/helmreleases.k8s.wattpad.com/v1alpha1"

	"github.com/stretchr/testify/assert"
)

func TestToYaml(t *testing.T) {
	expectedBytes, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	expectedStr := string(expectedBytes)

	d := NewDecoder()

	target := &v1alpha1.HelmRelease{}
	d.Decode(expectedBytes, nil, target)

	outString := ToYaml(*target)

	var tests = []struct {
		expected string
		actual   string
	}{
		{
			expectedStr,
			outString,
		},
	}

	for i := range tests {
		assert.Equal(t, tests[i].expected, tests[i].actual)
	}
}
