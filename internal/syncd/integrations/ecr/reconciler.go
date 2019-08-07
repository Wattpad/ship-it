package ecr

import (
	"context"
	"ship-it/internal"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
)

type HelmReleaseEditor interface {
	UpdateAndReplace(ctx context.Context, releaseName string, image *internal.Image) error
}

type IndexerService interface {
	Lookup(repo string) ([]types.NamespacedName, error)
}

type ImageReconciler struct {
	Editor       HelmReleaseEditor
	IndexService IndexerService
}

func NewReconciler(r HelmReleaseEditor, i IndexerService) *ImageReconciler {
	return &ImageReconciler{
		Editor:       r,
		IndexService: i,
	}
}

func (r *ImageReconciler) Reconcile(ctx context.Context, image *internal.Image) error {
	releases, err := r.IndexService.Lookup(image.Repository)
	if err != nil {
		return errors.Wrapf(err, "failed to obtain the releases corresponding to the repository: %s", image.Repository)
	}
	for _, release := range releases {
		err := r.Editor.UpdateAndReplace(ctx, release.Name, image)
		if err != nil {
			return err
		}
	}
	return nil
}
