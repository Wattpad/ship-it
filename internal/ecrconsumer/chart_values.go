package ecrconsumer

// TODO
// Handle nil val path case
// load path dynamically (should not be dependent on miranda folder structure)
// // Take image type return image type in set image tag function

import (
	"fmt"
	"path/filepath"
	"strings"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type ChartEditor interface {
	Image() *Image
	SetImageTag(tag string) error
}

type HelmChart struct {
	Metadata   chart.Metadata
	Templates  []*chart.Template
	Values     chartutil.Values
	ChartPath  string
	LocalPath  string
	ValuesPath string
	AutoDeploy bool
	GitRepo    string
}

type Image struct {
	Registry   string
	Repository string
	Tag        string
}

func LoadChart(serviceName string, localPath string, client GitHub) (*HelmChart, error) {
	// Grab metadata.yml
	err := client.SaveDirectory("master", "k8s/clusters/prod-v3", localPath)
	if err != nil {
		return nil, err
	}

	autoDeploy, chartPath, valPath, gitRepo, err := readMetadata("k8s/clusters/prod-v3/metadata.yml", serviceName)
	if err != nil {
		return nil, err
	}

	err = client.SaveDirectory("master", chartPath, localPath)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(localPath, chartPath)
	fmt.Println(path)
	c, err := chartutil.Load(path)
	if err != nil {
		return nil, err
	}

	vals, err := chartutil.ReadValuesFile(valPath)
	if err != nil {
		return nil, err
	}

	return &HelmChart{
		Metadata:   *c.Metadata,
		Templates:  c.Templates,
		Values:     vals,
		ChartPath:  chartPath,
		ValuesPath: valPath,
		LocalPath:  path,
		AutoDeploy: autoDeploy,
		GitRepo:    gitRepo,
	}, nil
}

func (c *HelmChart) Image() (*Image, error) {
	image, err := c.Values.Table("image") // this path will change depending on the chart need to handle these cases
	if err != nil {
		return nil, err
	}

	str, err := getChartValue(image, "repository")
	if err != nil {
		return nil, err
	}

	arr := strings.Split(str, "/")

	repository := arr[1]
	registry := arr[0]

	tag, err := getChartValue(image, "tag")
	if err != nil {
		return nil, err
	}

	return &Image{
		Registry:   registry,
		Repository: repository,
		Tag:        tag,
	}, nil
}

func (c *HelmChart) imageValues() chartutil.Values {
	return c.Values["image"].(map[string]interface{})
}

func readMetadata(path string, serviceName string) (bool, string, string, string, error) {
	v, err := chartutil.ReadValuesFile(path)
	if err != nil {
		return false, "", "", "", err
	}

	service, err := v.Table("services." + serviceName)
	if err != nil {
		return false, "", "", "", err
	}

	return service["autoDeploy"].(bool), service["helmChartPath"].(string), service["helmValuesPath"].(string), service["gitRepo"].(string), nil
}

func (c *HelmChart) SetImageTag(tag string, client GitHub) error {
	c.imageValues()["image"].(map[string]interface{})["tag"] = tag

	// encode new values with updated into bytes
	str := chartutil.ToYaml(c.Values)

	valueData := []byte(str)
	_, err := client.UpdateFile("Update Image Tag: "+tag, "master", c.ValuesPath, valueData) // change to a path join call

	if err != nil {
		return err
	}

	return nil
}

func getChartValue(values chartutil.Values, key string) (string, error) { // From Kube Deploy v1
	val, err := values.PathValue(key)
	if err != nil {
		return "", err
	}

	valStr, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("failed to convert %s to string: %v", key, val)
	}

	return valStr, nil
}