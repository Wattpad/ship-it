package helmrelease

import (
	"testing"
	"io/ioutil"
	"fmt"
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

	outBytes := target.Encode()
	// decode the expected bytes and re-encode check that there is still a match
	assert.Equal(t, expectedBytes, outBytes)
}


