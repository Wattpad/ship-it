package ecr

import (
	"context"
	"fmt"

	"ship-it/internal"

	"ship-it/pkg/apis/k8s.wattpad.com/v1alpha1"

	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
)

type RepositoriesService interface {
	UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error)
	GetFile(branch string, path string) ([]byte, error)
}

type ImageReconciler struct {
	Org          string
	ResourcePath string
	Branch       string
	RepoService  RepositoriesService
}

func NewReconciler(org string, valPath string, r RepositoriesService) *ImageReconciler {
	return &ImageReconciler{
		Org:          org,
		ResourcePath: valPath,
		RepoService:  r,
	}
}

func (r *ImageReconciler) Reconcile(ctx context.Context, image *internal.Image) error {
	resourceBytes, err := r.RepoService.GetFile(r.Branch, r.ResourcePath)
	if err != nil {
		return err
	}

	rls, err := v1alpha1.LoadRelease(resourceBytes)
	if err != nil {
		return err
	}

	updatedRls := internal.WithImage(*image, *rls)

	updatedBytes, err := yaml.Marshal(updatedRls)
	if err != nil {
		return nil
	}

	_, err = r.RepoService.UpdateFile(fmt.Sprintf("Image Tag updated to: %s", image.Tag), r.Branch, r.ResourcePath, updatedBytes)
	return err
}
