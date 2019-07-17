package ecr

import (
	"context"

	"ship-it/internal"

	"github.com/google/go-github/v26/github"
)

type RepositoriesService interface {
	GetContents(ctx context.Context, org, repo, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error)
	UpdateFile(ctx context.Context, owner, repo, path string, opt *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
}

type ImageReconciler struct {
	// TODO
}

func NewReconciler(org string, r RepositoriesService) *ImageReconciler {
	return &ImageReconciler{
		// TODO
	}
}

func (r *ImageReconciler) Reconcile(ctx context.Context, image *internal.Image) error {
	return nil
}
