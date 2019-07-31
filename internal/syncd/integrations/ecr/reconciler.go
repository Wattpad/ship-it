package ecr

import (
	"context"
	"fmt"
	"path/filepath"

	"ship-it/internal"

	"github.com/google/go-github/v26/github"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type RepositoriesService interface {
	UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error)
	GetFile(branch string, path string) (string, error)
}

type ImageReconciler struct {
	Org               string
	RegistryChartPath string
	Branch            string
	RepoService       RepositoriesService
}

func NewReconciler(org string, prefix string, branch string, r RepositoriesService) *ImageReconciler {
	return &ImageReconciler{
		Org:               org,
		RegistryChartPath: prefix,
		Branch:            branch,
		RepoService:       r,
	}
}

func (r *ImageReconciler) Reconcile(ctx context.Context, image *internal.Image) error {
	path := filepath.Join(r.RegistryChartPath, image.Repository+".yaml")
	resourceStr, err := r.RepoService.GetFile(r.Branch, path)
	if err != nil {
		return errors.Wrapf(err, "failed to download custom resource file for path: %s", path)
	}

	rls := make(map[string]interface{})

	err = yaml.Unmarshal([]byte(resourceStr), rls)
	if err != nil {
		return errors.Wrapf(err, "failed to parse the YAML file")
	}

	updatedRls := internal.WithImage(*image, rls)

	updatedBytes, err := yaml.Marshal(updatedRls)
	if err != nil {
		return err
	}

	_, err = r.RepoService.UpdateFile(fmt.Sprintf("Image Tag updated to: %s", image.Tag), r.Branch, path, updatedBytes)
	return err
}
