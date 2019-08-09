package ecr

import (
	"context"
	"ship-it/internal"

	"github.com/go-kit/kit/log"
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
	l            log.Logger
}

func NewReconciler(r HelmReleaseEditor, i IndexerService, logger log.Logger) *ImageReconciler {
	return &ImageReconciler{
		Editor:       r,
		IndexService: i,
		l:            logger,
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
			r.l.Log("error", err)
		}
	}
	return nil
}
