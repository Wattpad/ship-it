package ecrconsumer

import (
	"fmt"
	"path/filepath"
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

func FindImage(v reflect.Value, image *Image, serviceName string) {
	// Indirect through pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			FindImage(v.Index(i), image, serviceName)
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
			FindImage(v.MapIndex(k), image, serviceName)
		}
	default:
		// handle other types
	}
}

func update(v reflect.Value, image Image, target *map[string]interface{}) {
	// Indirect through pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			update(v.Index(i), image, target)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if k.Interface().(string) == "image" {
				// parse image and check for service name match only then can the field be updated and returned
				foundImage, err := parseImage(v.MapIndex(k).Interface().(map[string]interface{})["repository"].(string), v.MapIndex(k).Interface().(map[string]interface{})["tag"].(string))
				if err != nil {
					fmt.Println(err)
					return
				}
				if foundImage.Repository == image.Repository {
					v.MapIndex(k).Interface().(map[string]interface{})["repository"] = image.Registry + "/" + image.Repository
					v.MapIndex(k).Interface().(map[string]interface{})["tag"] = image.Tag

					*target = v.Interface().(map[string]interface{})
					return
				}
			}
			update(v.MapIndex(k), image, target)
		}
	default:
		// handle other types
	}
}

func iterativeUpdate(vals map[string]interface{}, img Image) map[string]interface{} {
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
	fmt.Println(target)
	fmt.Print("\n\n")

	outBytes := target.Encode()
	_, err = client.UpdateFile("diff test", "master", "k8s/custom-resources/loki.yaml", outBytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(reflect.DeepEqual(resourceBytes, outBytes))

	image := Image{}
	FindImage(reflect.ValueOf(target.Spec.Values.Object), &image, serviceName)

	image.Tag = "This is a new tag" // change a value and print
	_ = WithImage(image, *target)
	//fmt.Println(changed.Spec.Values.Object)
	// Try encoding back to bytes to prep git commit

	// changed.TypeMeta.APIVersion
	// changed.TypeMeta.Kind

	return nil, nil
}

func WithImage(img Image, r helmrelease.HelmRelease) helmrelease.HelmRelease {
	newVals := make(map[string]interface{})
	update(reflect.ValueOf(r.Spec.Values.Object), img, &newVals)

	r.Spec.Values.Object = newVals
	return r
}
