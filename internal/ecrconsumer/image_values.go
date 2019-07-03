package ecrconsumer

import (
	"fmt"
	"reflect"
	"strings"

	"ship-it/pkg/apis/helmreleases.k8s.wattpad.com/v1alpha1"

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

func getImagePath(v reflect.Value, serviceName string) []string {
	if v.Kind() != reflect.Map {
		return []string{}
	}

	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		if key == "image" {
			img := iter.Value().Interface().(map[string]interface{})
			image, err := parseImage(img["repository"].(string), img["tag"].(string))
			if err != nil {
				return []string{}
			}
			if image.Repository == serviceName {
				return []string{key}
			}
		}

		val := reflect.ValueOf(iter.Value().Interface())
		if val.Kind() == reflect.Map {
			path := getImagePath(val, serviceName)
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

func update(vals map[string]interface{}, img Image, path []string) map[string]interface{} {
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
	d := NewDecoder()

	err := runtime.DecodeInto(d, fileData, rls)
	if err != nil {
		return nil, err
	}

	return rls, nil
}

func WithImage(img Image, r v1alpha1.HelmRelease, path []string) v1alpha1.HelmRelease {
	newVals := update(r.Spec.Values.Object, img, path)
	r.Spec.Values.Object = newVals

	return r
}