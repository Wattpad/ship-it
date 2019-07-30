package ecrconsumer

import (
	"encoding/json"
	"fmt"

	"ship-it/internal"

	shipitv1beta1 "ship-it-operator/api/v1beta1"
)

func update(vals map[string]interface{}, img internal.Image) {
	path := internal.GetImagePath(vals, img.Repository)
	if len(path) == 0 {
		return
	}

	imgVals := internal.Table(vals, path)
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

func WithImage(img internal.Image, r shipitv1beta1.HelmRelease) (shipitv1beta1.HelmRelease, error) {
	copy := r.DeepCopy()

	cleanMap := cleanUpStringMap(r.HelmValues())
	update(cleanMap, img)

	newJSON, _ := json.Marshal(cleanMap)
	copy.Spec.Values.Raw = newJSON

	return *copy, nil
}
