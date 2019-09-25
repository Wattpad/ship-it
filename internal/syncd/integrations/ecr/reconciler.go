package ecr

import (
	"context"

	"ship-it/internal/image"

	"github.com/go-kit/kit/log"
	"k8s.io/apimachinery/pkg/types"
)

// ChartEditor edits a remote chart containing HelmReleases. It updates the
// named release specs using the given image reference.
type ChartEditor interface {
	Edit(ctx context.Context, releases []types.NamespacedName, image *image.Ref) error
}

// ReleaseIndexer provides access to an index of container images and the
// deployed releases that are using them.
type ReleaseIndexer interface {
	Lookup(image *image.Ref) ([]types.NamespacedName, error)
}

type ImageReconciler struct {
	editor  ChartEditor
	indexer ReleaseIndexer
	logger  log.Logger
}

func NewReconciler(l log.Logger, e ChartEditor, i ReleaseIndexer) *ImageReconciler {
	return &ImageReconciler{
		editor:  e,
		indexer: i,
		logger:  l,
	}
}

func (r *ImageReconciler) Reconcile(ctx context.Context, image *image.Ref) error {
	releases, err := r.indexer.Lookup(image)
	if err != nil {
		return err
	}

	if len(releases) == 0 {
		return nil
	}

	return r.editor.Edit(ctx, releases, image)
}
