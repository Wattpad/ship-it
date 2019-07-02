package helmrelease

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestEncode(t *testing.T) {
	expectedBytes, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	d := NewDecoder()
	gvk := schema.FromAPIVersionAndKind("helmreleases.k8s.wattpad.com/v1alpha1", "HelmRelease")
	target := &HelmRelease{}
	d.Decode(expectedBytes, &gvk, target)
	fmt.Println(target)
	outBytes := target.Encode()

	var tests = []struct {
		expected []byte
		actual   []byte
	}{
		{
			expectedBytes,
			outBytes,
		},
	}

	for i := range tests {
		assert.Equal(t, tests[i].expected, tests[i].actual)
	}
}
