package ecr

import (
	"testing"

	"ship-it/internal"

	"github.com/stretchr/testify/assert"
)

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
