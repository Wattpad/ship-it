package ecrconsumer

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
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

func findImages(v reflect.Value, arr *[]Image) {
	fmt.Printf("Visiting %v\n", v)
	// Indirect through pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			findImages(v.Index(i), arr)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if k.Interface().(string) == "image" {
				i := v.MapIndex(k).Interface().(map[interface{}]interface{})
				img, err := parseImage(i["repository"].(string), i["tag"].(string))
				if err != nil {
					fmt.Println(err)
					return
				}
				*arr = append(*arr, Image{
					Registry:   img.Registry,
					Repository: img.Repository,
					Tag:        img.Tag,
				})
			}
			findImages(v.MapIndex(k), arr)
		}
	default:
		// handle other types
	}
}

func update(v reflect.Value, images []Image, iterator *int) map[string]interface{} {
	// Indirect through pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			update(v.Index(i), images, iterator)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if k.Interface().(string) == "image" {
				v.MapIndex(k).Interface().(map[interface{}]interface{})["repository"] = images[*iterator].Registry + "/" + images[*iterator].Repository
				v.MapIndex(k).Interface().(map[interface{}]interface{})["tag"] = images[*iterator].Tag
				*iterator++
				if *iterator == len(images) { // when all images are updated
					return v.Interface().(map[string]interface{})
				}
			}
			update(v.MapIndex(k), images, iterator)
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

	var customResource HelmRelease

	err = yaml.Unmarshal(resourceBytes, &customResource)
	if err != nil {
		return nil, err
	}

	images := make([]Image, 0)
	findImages(reflect.ValueOf(customResource.Spec.Values), &images)
	fmt.Println(images)
	images[0].Tag = "this tag is updated img 0"
	images[1].Tag = "this tag is updated img 1"
	i := 0
	fmt.Println(update(reflect.ValueOf(customResource.Spec.Values), images, &i))
	return nil, nil
}

func (r *HelmRelease) WithImages(imgs []Image) (*HelmRelease, error) {
	// reassign values field of r and return r
	return nil, nil
}
