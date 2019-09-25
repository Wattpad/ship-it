package unstructured

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
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

func TestVisitOne(t *testing.T) {
	expected := "newvalue"

	obj := yaml.MapSlice{
		{
			Key:   "foo",
			Value: "value",
		},
		{
			Key: "foo",
			Value: yaml.MapSlice{
				{
					Key:   "bar",
					Value: "oldvalue",
				},
			},
		},
		{
			Key:   "qux",
			Value: 2,
		},
	}

	pred := func(item yaml.MapItem) bool {
		key, ok := item.Key.(string)
		return ok && key == "bar"
	}

	visit := func(item *yaml.MapItem) {
		item.Value = expected
	}

	VisitOne(obj, pred, visit)
	assert.Equal(t, expected, obj[1].Value.(yaml.MapSlice)[0].Value)
}
