package helmrelease

import (
	"fmt"
	"io/ioutil"
	"testing"

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

	target := &HelmRelease{}
	d.Decode(expectedBytes, nil, target)

	outString := target.ToYaml()

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
