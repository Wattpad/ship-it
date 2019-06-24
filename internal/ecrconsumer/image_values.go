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
	service    string
	Registry   string
	Repository string
	Tag        string
}

func LoadImage(serviceName string, client GitCommands) (*Image, error) {
	resourceBytes, err := client.DownloadFile("master", filepath.Join("k8s/custom-resources", serviceName+".yaml")) // we assume this file path until a location in miranda for custom resources is decided upon
	if err != nil {
		return nil, err
	}

	customResource := CRYaml{}

	err = yaml.Unmarshal(resourceBytes, &customResource)
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
		service:    serviceName,
		Registry:   registry,
		Repository: repository,
		Tag:        customResource.Spec.Values.Image.Tag,
	}, nil
}

func (i *Image) Update(client GitCommands) error {
	msg := fmt.Sprintf("Image Update --> Registry: %s, Repository: %s, Tag: %s", i.Registry, i.Repository, i.Tag)
	path := filepath.Join("k8s/custom-resources", i.service+".yaml")

	resourceBytes, err := client.DownloadFile("master", filepath.Join("k8s/custom-resources", i.service+".yaml")) // we assume this file path until a location in miranda for custom resources is decided upon
	if err != nil {
		return err
	}

	customResource := CRYaml{}

	err = yaml.Unmarshal(resourceBytes, &customResource)
	if err != nil {
		return err
	}

	customResource.Spec.Values.Image.Repository = filepath.Join(i.Registry, i.Repository)
	customResource.Spec.Values.Image.Tag = i.Tag

	updatedBytes, err := yaml.Marshal(customResource)
	if err != nil {
		return err
	}

	_, err = client.UpdateFile(msg, "master", path, updatedBytes)
	if err != nil {
		return err
	}

	return nil
}
