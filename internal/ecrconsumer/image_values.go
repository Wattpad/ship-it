package ecrconsumer

import (
	"fmt"

	"ship-it/internal"

	"ship-it/pkg/apis/k8s.wattpad.com/v1alpha1"
)

func getImagePath(obj map[string]interface{}, serviceName string) []string {
	for key, val := range obj {
		if key == "image" {
			if img, ok := val.(map[string]interface{}); ok {
				image, err := internal.ParseImage(img["repository"].(string), img["tag"].(string))
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

func update(vals map[string]interface{}, img internal.Image) {
	path := getImagePath(vals, img.Repository)
	if len(path) == 0 {
		return
	}

	imgVals := table(vals, path)
	if imgVals == nil {
		return
	}
	imgVals["repository"] = img.URI()
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

func WithImage(img internal.Image, r v1alpha1.HelmRelease) v1alpha1.HelmRelease {
	copy := r.DeepCopy()

	cleanMap := cleanUpStringMap(copy.Spec.Values)
	copy.Spec.Values = v1alpha1.HelmValues(cleanMap)

	update(copy.Spec.Values, img)

	return *copy
}
