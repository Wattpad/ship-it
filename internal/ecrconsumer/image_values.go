package ecrconsumer

import (
	"encoding/json"
	"fmt"

	"ship-it/internal"

	shipitv1beta1 "ship-it-operator/api/v1beta1"
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

func WithImage(img internal.Image, r shipitv1beta1.HelmRelease) shipitv1beta1.HelmRelease {
	copy := r.DeepCopy()

	var values map[string]interface{}

	// FIXME handle error
	json.Unmarshal(copy.Spec.Values.Raw, &values)

	cleanMap := cleanUpStringMap(values)

	update(cleanMap, img)

	newJSON, _ := json.Marshal(cleanMap)

	copy.Spec.Values.Raw = newJSON

	return *copy
}
