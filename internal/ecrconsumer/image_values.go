package ecrconsumer

import (
	"fmt"
	"strings"

	"ship-it/pkg/apis/k8s.wattpad.com/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
)

type Image struct {
	Registry   string
	Repository string
	Tag        string
}

func parseImage(repo string, tag string) (*Image, error) {
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

func getImagePath(obj map[string]interface{}, serviceName string) []string {
	for key, val := range obj {
		if key == "image" {
			if img, ok := val.(map[string]interface{}); ok {
				image, err := parseImage(img["repository"].(string), img["tag"].(string))
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
		_, ok := tabled[p].(map[string]interface{})
		if !ok {
			return nil
		}

		tabled = tabled[p].(map[string]interface{})
	}
	return tabled
}

func update(vals map[string]interface{}, img Image) map[string]interface{} {
	path := getImagePath(vals, img.Repository)
	if len(path) == 0 {
		return nil
	}

	imgVals := table(vals, path)
	if imgVals == nil {
		return nil
	}

	imgVals["repository"] = img.Registry + "/" + img.Repository
	imgVals["tag"] = img.Tag

	return vals
}

func LoadRelease(fileData []byte) (*v1alpha1.HelmRelease, error) {
	rls := &v1alpha1.HelmRelease{}
	d := v1alpha1.Decoder

	err := runtime.DecodeInto(d, fileData, rls)
	if err != nil {
		return nil, err
	}

	return rls, nil
}

func cleanUpInterfaceArray(in []interface{}) []interface{} {
	result := make([]interface{}, len(in))
	for i, v := range in {
		result[i] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpStringMap(in map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range in {
		result[fmt.Sprintf("%v", k)] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
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
	case string:
		return v
	default:
		return v
	}
}

func WithImage(img Image, r v1alpha1.HelmRelease) v1alpha1.HelmRelease {
	copy := r.DeepCopy()

	cleanMap := cleanUpStringMap(copy.Spec.Values)
	copy.Spec.Values = v1alpha1.HelmValues(cleanMap)

	update(copy.Spec.Values, img)

	return *copy
}
