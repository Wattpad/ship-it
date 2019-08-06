package unstructured

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	obj := map[string]interface{}{
		"bar": "happy",
		"foo": map[string]interface{}{
			"bar": "path",
		},
		"qux": 2,
	}

	var values []string
	FindAll(obj, "bar", func(x interface{}) {
		values = append(values, x.(string))
	})

	assert.ElementsMatch(t, []string{"happy", "path"}, values)
}
