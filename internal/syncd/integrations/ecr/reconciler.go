package ecr

import (
	"context"
	"ship-it/internal"

	"github.com/go-kit/kit/log"
	"github.com/google/go-github/v26/github"
	"k8s.io/apimachinery/pkg/types"
)

// As an example, if you wanted to commit a change to a file in your repository, you would:
// 1. Get the current commit object
// 2. Retrieve the tree it points to
// 3. Retrieve the content of the blob object that tree has for that particular file path
// 4. Change the content somehow and post a new blob object with that new content, getting a blob SHA back
// 5. Post a new tree object with that file path pointer replaced with your new blob SHA getting a tree SHA back
// 6. Create a new commit object with the current commit SHA as the parent and the new tree SHA, getting a commit SHA back
// 7. Update the reference of your branch to point to the new commit SHA

type GithubCommitter interface {
	GetCommit(ctx context.Context, owner string, repo string, sha string) (*github.Commit, *github.Response, error)
	CreateCommit(ctx context.Context, owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error)
}

// ReleaseEditor remotely updates the values files of the releases to use the
// provided image reference.
type ReleaseEditor interface {
	Edit(ctx context.Context, releases []types.NamespacedName, image *internal.Image) error
}

// ReleaseIndexer provides access to an index of container images and the
// releases that are deployed using them.
type ReleaseIndexer interface {
	Lookup(image *internal.Image) ([]types.NamespacedName, error)
}

type ImageReconciler struct {
	editor  ReleaseEditor
	indexer ReleaseIndexer
	logger  log.Logger
}

func NewReconciler(l log.Logger, g GithubCommitter, i ReleaseIndexer, org, repo, ref string) *ImageReconciler {
	return &ImageReconciler{
		editor:  &releaseEditor{g, org, repo, ref},
		indexer: i,
		logger:  l,
	}
}

func (r *ImageReconciler) Reconcile(ctx context.Context, image *internal.Image) error {
	releases, err := r.indexer.Lookup(image)
	if err != nil {
		return err
	}

	if len(releases) == 0 {
		return nil
	}

	return r.editor.Edit(ctx, releases, image)
}
