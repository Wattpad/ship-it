package ecrconsumer

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"ship-it/internal/helmrelease"

	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
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

func findImage(v reflect.Value, image *Image, serviceName string) {
	// Indirect through pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			findImage(v.Index(i), image, serviceName)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if k.Interface().(string) == "image" {
				i := v.MapIndex(k).Interface().(map[string]interface{})
				img, err := parseImage(i["repository"].(string), i["tag"].(string))
				if err != nil {
					fmt.Println(err)
					return
				}
				if img.Repository == serviceName {
					*image = *img
					return
				}
			}
			findImage(v.MapIndex(k), image, serviceName)
		}
	default:
		// handle other types
	}
}

func update(v reflect.Value, image Image) map[string]interface{} {
	// Indirect through pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			update(v.Index(i), image)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if k.Interface().(string) == "image" {
				// parse image and check for service name match only then can the field be updated and returned
				repo := v.MapIndex(k).Interface().(map[string]interface{})["repository"].(string)
				tag := v.MapIndex(k).Interface().(map[string]interface{})["tag"].(string)
				foundImage, err := parseImage(repo, tag)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				if foundImage.Repository == image.Repository {
					v.MapIndex(k).Interface().(map[string]interface{})["repository"] = image.Registry + "/" + image.Repository
					v.MapIndex(k).Interface().(map[string]interface{})["tag"] = image.Tag

					return v.Interface().(map[string]interface{})
				}
			}
			update(v.MapIndex(k), image)
		}
	default:
		// handle other types
	}
	return nil
}

func LoadImage(serviceName string, client GitCommands) (*Image, error) {
	// we assume this file path until a location in miranda for custom resources is decided upon
	resourceBytes, err := client.GetFile("master", filepath.Join("k8s/custom-resources", serviceName+".yaml"))
	if err != nil {
		return nil, err
	}

	target := &helmrelease.HelmRelease{}
	d := helmrelease.NewDecoder()
	gvk := schema.FromAPIVersionAndKind("helmreleases.k8s.wattpad.com/v1alpha1", "HelmRelease")
	d.Decode(resourceBytes, &gvk, target)
	fmt.Println(target.Spec.Values.Object)
	fmt.Println("\n")

	newBytes, err := yaml.Marshal(target)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if reflect.DeepEqual(newBytes, resourceBytes) {
		fmt.Println("same bytes")
		return nil, nil
	}
	newVals := make(map[string]interface{})
	err = yaml.Unmarshal(newBytes, &newVals)
	if err != nil {
		return nil, err
	}
	fmt.Println("\n new vals below \n")
	fmt.Println(newVals)
	fmt.Println("\n \n")

	image := Image{}
	findImage(reflect.ValueOf(target.Spec.Values.Object), &image, serviceName)

	image.Tag = "This is a new tag" // change a value and print
	changed := WithImage(image, *target).Spec.Values.Object
	fmt.Println(changed)
	// Try encoding back to bytes to prep git commit

	return nil, nil
}

func WithImage(img Image, r helmrelease.HelmRelease) helmrelease.HelmRelease {
	newVals := update(reflect.ValueOf(r.Spec.Values.Object), img)

	r.Spec.Values.Object = newVals
	return r
}
