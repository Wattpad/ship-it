package ecrconsumer

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
)

type GitCommands interface {
	UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error)
	GetFile(branch string, path string) ([]byte, error)
}

type Image struct {
	service    string
	Registry   string
	Repository string
	Tag        string
}

func parseImage(serviceName string, repo string, tag string) (*Image, error) {
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
		service:    serviceName,
		Registry:   registry,
		Repository: repository,
		Tag:        tag,
	}, nil
}

func getKeys(vals chartutil.Values) []string {
	var keys []string
	for k := range vals {
		keys = append(keys, k)
	}
	fmt.Println(keys)
	if len(keys) == 0 {
		return nil
	}
	return keys
}

func checkForImageKey(keys []string) bool {
	for k := range keys {
		if keys[k] == "image" {
			return true
		}
	}
	return false
}

func getVal(keys []string, deepMap map[string]interface{}) interface{} {
	if len(keys) == 1 {
		return deepMap[keys[0]]
	}

	return getVal(keys[1:], deepMap[keys[0]].(map[string]interface{}))
}

func breadthSearch(graph map[string]interface{}, start string, end string) string {
	// queue := make([]string, 0)
	// queue = append(queue, start)

	// for len(queue) != 0 {
	// 	path := queue[0]
	// 	queue = queue[1:]

	// }
	return ""
}

//func findImages() interface{} {
// m := chartutil.Values{
// 	"apples": []string{"delicious", "green", "red"},
// 	"oranges": map[string]interface{}{
// 		"foo": 123456,
// 		"image": map[string]interface{}{
// 			"repo": "bar",
// 			"tag":  "hello, world",
// 		},
// 	},
// 	// "image": map[string]interface{}{
// 	// 	"foo": 123456,
// 	// 	"bar": "hello, world",
// 	// },
// }
// keys := getKeys(m)
// imageKeys := make([]string, 0)
// values := m
// for i := range keys {
// 	if keys[i] == "image" {
// 		imageKeys = append(imageKeys, keys[i])
// 	} else {
// 		tabled, _ := values.Table(keys[i])
// 		breadthSearch([]string{"image"}, tabled)
// 	}
// }
//return breadthSearch([]string{"image"}, m)
//}

func walk(v reflect.Value) {
	fmt.Printf("Visiting %v\n", v)
	// Indirect through pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			walk(v.Index(i))
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			walk(v.MapIndex(k))
		}
	default:
		// handle other types
	}
}

func LoadImage(serviceName string, client GitCommands) (*Image, error) {
	// we assume this file path until a location in miranda for custom resources is decided upon
	resourceBytes, err := client.GetFile("master", filepath.Join("k8s/custom-resources", serviceName+".yaml"))
	if err != nil {
		return nil, err
	}

	var customResource CRYaml

	err = yaml.Unmarshal(resourceBytes, &customResource)
	if err != nil {
		return nil, err
	}
	getKeys(customResource.Spec.Values)
	walk(reflect.ValueOf(customResource.Spec.Values))
	return nil, nil
	//return parseImage(serviceName, customResource.Spec.Values.Image.Repository, customResource.Spec.Values.Image.Tag)
}

func (i *Image) Update(client GitCommands) error {
	//msg := fmt.Sprintf("Image Update --> Registry: %s, Repository: %s, Tag: %s", i.Registry, i.Repository, i.Tag)
	//path := filepath.Join("k8s/custom-resources", i.service+".yaml")

	// we assume this file path until a location in miranda for custom resources is decided upon
	//resourceBytes, err := client.GetFile("master", filepath.Join("k8s/custom-resources", i.service+".yaml"))
	//if err != nil {
	//	return err
	//}

	//customResource := CRYaml{}

	// err = yaml.Unmarshal(resourceBytes, &customResource)
	// if err != nil {
	// 	return err
	// }

	// Replace get Image helper func calls
	//customResource.Spec.Values.Image.Repository = filepath.Join(i.Registry, i.Repository)
	//customResource.Spec.Values.Image.Tag = i.Tag

	// updatedBytes, err := yaml.Marshal(customResource)
	// if err != nil {
	// 	return err
	// }

	// _, err = client.UpdateFile(msg, "master", path, updatedBytes)
	// if err != nil {
	// 	return err
	// }

	return nil
}
