package ecrconsumer

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type ImageEditor interface {
	Update() *error
}

type Image struct {
	Registry   string
	Repository string
	Tag        string
}

func LoadImage(serviceName string, client GitHub) (*Image, error) {
	imgBytes, err := client.DownloadFile("master", filepath.Join("k8s/custom-resources", serviceName+".yaml"))
	if err != nil {
		return nil, err
	}

	customResource := CRYaml{}

	err = yaml.Unmarshal(imgBytes, &customResource)
	if err != nil {
		return nil, err
	}

	arr := strings.Split(customResource.Spec.Values.Image.Repository, "/")

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
		Tag:        customResource.Spec.Values.Image.Tag,
	}, nil
}
