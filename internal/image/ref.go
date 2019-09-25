package image

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// Ref is a reference to a container image in a remote repository
type Ref struct {
	Registry   string
	Repository string
	Tag        string
}

// String formats the most canonical string representation of the image reference
func (r Ref) String() string {
	var tag string
	if r.Tag != "" {
		tag = ":" + r.Tag
	}

	return r.URI() + tag
}

func (r Ref) URI() string {
	return r.Registry + "/" + r.Repository
}

// Refs match when their registries and repositories match
func (r Ref) Matches(other Ref) bool {
	return r.Registry == other.Registry && r.Repository == other.Repository
}

func Parse(repo string, tag string) (*Ref, error) {
	arr := strings.Split(repo, "/")
	if len(arr) != 2 {
		return nil, fmt.Errorf("invalid image repo: %s", repo)
	}

	return &Ref{
		Registry:   arr[0],
		Repository: arr[1],
		Tag:        tag,
	}, nil
}

func FromYaml(item yaml.MapItem) (*Ref, error) {
	if key, ok := item.Key.(string); !ok || key != "image" {
		return nil, errors.New("invalid yaml: missing \"image\" key")
	}

	var (
		repo string
		tag  string
	)

	if obj, ok := item.Value.(yaml.MapSlice); ok {
		for _, item := range obj {
			if key, ok := item.Key.(string); ok {
				if key == "repository" {
					repo = item.Value.(string)
				} else if key == "tag" {
					tag = item.Value.(string)
				}
			}
		}
	}

	if repo == "" {
		return nil, errors.New("invalid yaml: missing \"repository\" key")
	}

	if tag == "" {
		return nil, errors.New("invalid yaml: missing \"tag\" key")
	}

	return Parse(repo, tag)
}
