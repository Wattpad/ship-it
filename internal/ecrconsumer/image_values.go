package ecrconsumer

import (
	"fmt"
	"reflect"
	"strings"

	"ship-it/internal/helmrelease"

	"github.com/google/go-github/github"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type GitCommands interface {
	UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error)
	GetFile(branch string, path string) ([]byte, error)
}

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
		return nil, fmt.Errorf("repository field is invalid! (lacking either the full registry name or repository name)")
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
			img := v.MapIndex(iter.Key()).Interface().(map[string]interface{})
			image, err := parseImage(img["repository"].(string), img["tag"].(string))
			if err != nil {
				fmt.Println("Invalid image")
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
	for i := range path {
		_, ok := tabled[path[i]].(map[string]interface{})
		if !ok {
			fmt.Println("invalid path")
			return nil
		}

		tabled = tabled[path[i]].(map[string]interface{})
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

func LoadRelease(serviceName string, fileData []byte) (*helmrelease.HelmRelease, error) {

	rls := &helmrelease.HelmRelease{}
	d := helmrelease.NewDecoder()
	gvk := schema.FromAPIVersionAndKind("helmreleases.k8s.wattpad.com/v1alpha1", "HelmRelease")
	d.Decode(fileData, &gvk, rls)

	return rls, nil
}

func WithImage(img Image, r helmrelease.HelmRelease, path []string) helmrelease.HelmRelease {
	newVals := update(r.Spec.Values.Object, img, path)
	r.Spec.Values.Object = newVals

	return r
}
