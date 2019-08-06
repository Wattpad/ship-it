package ecr

import (
	"context"
	"fmt"
	"path/filepath"
	"ship-it/internal"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
)

type RepositoriesService interface {
	UpdateAndReplace(ctx context.Context, path string, image *internal.Image, msg string) error
}

type IndexerService interface {
	Lookup(repo string) ([]types.NamespacedName, error)
}

type ImageReconciler struct {
	RegistryChartPath string
	RepoService       RepositoriesService
	IndexService      IndexerService
}

func NewReconciler(prefix string, r RepositoriesService, i IndexerService) *ImageReconciler {
	return &ImageReconciler{
		RegistryChartPath: prefix,
		RepoService:       r,
		IndexService:      i,
	}
}

func (r *ImageReconciler) Reconcile(ctx context.Context, image *internal.Image) error {
	releases, err := r.IndexService.Lookup(image.Repository)
	if err != nil {
		return errors.Wrapf(err, "failed to obtain the releases corresponding to the repository: %s", image.Repository)
	}
	for _, release := range releases {
		err := r.RepoService.UpdateAndReplace(ctx, filepath.Join(r.RegistryChartPath, release.Name+".yaml"), image, fmt.Sprintf("Image Tag updated to: %s", image.Tag))
		if err != nil {
			return err
		}
	}
	return nil
}
