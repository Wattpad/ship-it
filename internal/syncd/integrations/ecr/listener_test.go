package ecr

import (
	"testing"

	"ship-it/internal"

	"github.com/stretchr/testify/assert"
)

func TestMakeImage(t *testing.T) {
	assert.Exactly(t, internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "ship-it",
		Tag:        "shipped",
	}, makeImage("ship-it", "shipped", "723255503624"))
}

func TestValidateTag(t *testing.T) {
	tests := []struct {
		inputString string
		expected    bool
	}{
		{"78bc9ccf64eb838c6a0e0492ded722274925e2bd", true},
		{"latest", false},
		{"78bc9ccf64eb838c6a0e0492ded722274925E2ND", false},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, validImageTagRegex.MatchString(test.inputString))
	}
}
