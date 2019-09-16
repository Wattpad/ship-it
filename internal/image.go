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

func (i Image) Tagged() string {
	var tag string
	if i.Tag != "" {
		tag = ":" + i.Tag
	}

	return i.URI() + tag
}

func (i Image) URI() string {
	return i.Registry + "/" + i.Repository
}

func ParseImage(repo string, tag string) (*Image, error) {
	arr := strings.Split(repo, "/")
	if len(arr) != 2 {
		return nil, fmt.Errorf("malformed repo: %s", repo)
	}

	return &Image{
		Registry:   arr[0],
		Repository: arr[1],
		Tag:        tag,
	}, nil
}

func getImagePath(obj map[string]interface{}, serviceName string) []string {
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
			path := getImagePath(nested, serviceName)
			if len(path) == 0 {
				continue
			}
			return append([]string{key}, path...)
		}
	}

	return []string{}
}

func table(vals map[string]interface{}, path []string) map[string]interface{} {
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

func update(vals map[string]interface{}, img Image) {
	path := getImagePath(vals, img.Repository)
	if len(path) == 0 {
		return
	}

	imgVals := table(vals, path)
	if imgVals == nil {
		return
	}
	imgVals["tag"] = img.Tag
}

func cleanUpInterfaceArray(in []interface{}) []interface{} {
	result := make([]interface{}, len(in))
	for i, v := range in {
		result[i] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpStringMap(in map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(in))
	for k, v := range in {
		result[fmt.Sprintf("%v", k)] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(in))
	for k, v := range in {
		result[fmt.Sprintf("%v", k)] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanUpInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanUpInterfaceMap(v)
	default:
		return v
	}
}

func DeepCopyMap(in map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(in))
	for key, val := range in {
		if nested, ok := val.(map[string]interface{}); ok {
			result[key] = DeepCopyMap(nested)
		}
		result[key] = val
	}
	return result
}

func WithImage(img Image, rlsMap map[string]interface{}) map[string]interface{} {
	cleanMap := cleanUpStringMap(rlsMap)
	copiedMap := DeepCopyMap(cleanMap)

	update(copiedMap, img)

	return copiedMap
}
