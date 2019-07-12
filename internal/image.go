package internal

import (
	"fmt"
	"strings"
)

type Image struct {
	Registry   string
	Repository string
	Tag        string
}

func ParseImage(repo string, tag string) (*Image, error) {
	arr := strings.Split(repo, "/")

	var registry string
	var repository string
	if len(arr) == 2 {
		repository = arr[1]
		registry = arr[0]
	} else {
		return nil, fmt.Errorf("malformed repo: %s", repo)
	}

	return &Image{
		Registry:   registry,
		Repository: repository,
		Tag:        tag,
	}, nil
}

func (i Image) URI() string {
	return i.Registry + "/" + i.Repository
}

func GetImagePath(obj map[string]interface{}, serviceName string) []string {
	for key, val := range obj {
		if key == "image" {
			if img, ok := val.(map[string]interface{}); ok {
				image, err := ParseImage(img["repository"].(string), img["tag"].(string))
				if err != nil {
					return []string{}
				}
				if image.Repository == serviceName {
					return []string{key}
				}
			}
		}

		if nested, ok := val.(map[string]interface{}); ok {
			path := GetImagePath(nested, serviceName)
			if len(path) == 0 {
				continue
			}
			return append([]string{key}, path...)
		}
	}

	return []string{}
}

func Table(vals map[string]interface{}, path []string) map[string]interface{} {
	tabled := vals
	for _, p := range path {
		currentMap, ok := tabled[p].(map[string]interface{})
		if !ok {
			return nil
		}

		tabled = currentMap
	}
	return tabled
}
