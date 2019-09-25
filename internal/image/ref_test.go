package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestParse(t *testing.T) {
	expected := Ref{
		Registry:   "foo",
		Repository: "bar",
		Tag:        "baz",
	}

	image, err := Parse("foo/bar", "baz")
	assert.NoError(t, err)
	assert.Equal(t, expected, *image)
}

func TestParseInvalid(t *testing.T) {
	_, err := Parse("foo-bar", "baz")
	assert.Error(t, err)
}

func TestRefString(t *testing.T) {
	type testCase struct {
		in  Ref
		out string
	}

	testCases := map[string]testCase{
		"without tag": {
			in:  Ref{"registry", "repository", ""},
			out: "registry/repository",
		},

		"with tag": {
			in:  Ref{"registry", "repository", "tag"},
			out: "registry/repository:tag",
		},
	}

	for name, tc := range testCases {
		tc := tc // scopelint
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.out, tc.in.String())
		})
	}
}

func TestRefMatches(t *testing.T) {
	type testCase struct {
		this    Ref
		that    Ref
		matches bool
	}

	testCases := map[string]testCase{
		"matches": {
			this:    Ref{"foo", "bar", "baz"},
			that:    Ref{"foo", "bar", "qux"},
			matches: true,
		},
		"doesn't match": {
			this:    Ref{"foo", "bar", "baz"},
			that:    Ref{"foo", "baz", "qux"},
			matches: false,
		},
	}

	for name, tc := range testCases {
		tc := tc // scopelint
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.matches, tc.this.Matches(tc.that))
		})
	}
}

func TestFromYaml(t *testing.T) {
	yamlMap := yaml.MapItem{
		Key: "image",
		Value: yaml.MapSlice{
			{
				Key:   "repository",
				Value: "registry/repository",
			},
			{
				Key:   "tag",
				Value: "tag",
			},
		},
	}

	expected := Ref{
		Registry:   "registry",
		Repository: "repository",
		Tag:        "tag",
	}

	ref, err := FromYaml(yamlMap)
	assert.NoError(t, err)
	assert.Equal(t, expected, *ref)
}

func TestFromYamlInvalid(t *testing.T) {
	testCases := map[string]yaml.MapItem{
		"not an image block": yaml.MapItem{
			Key: "something-else",
			Value: yaml.MapSlice{
				{
					Key:   "repository",
					Value: "registry/repository",
				},
				{
					Key:   "tag",
					Value: "tag",
				},
			},
		},
		"invalid repository field": yaml.MapItem{
			Key: "image",
			Value: yaml.MapSlice{
				{
					Key:   "repository",
					Value: "repository",
				},
				{
					Key:   "tag",
					Value: "tag",
				},
			},
		},
		"missing repository field": yaml.MapItem{
			Key: "image",
			Value: yaml.MapSlice{
				{
					Key:   "tag",
					Value: "tag",
				},
			},
		},
		"missing tag field": yaml.MapItem{
			Key: "image",
			Value: yaml.MapSlice{
				{
					Key:   "repository",
					Value: "repository",
				},
			},
		},
	}

	for name, tc := range testCases {
		tc := tc // scopelint
		t.Run(name, func(t *testing.T) {
			_, err := FromYaml(tc)
			assert.Error(t, err)
		})
	}
}
