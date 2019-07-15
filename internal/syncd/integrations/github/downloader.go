package github

import (
	"context"

	"github.com/google/go-github/v26/github"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
)

type RepositoriesService interface {
	GetContents(ctx context.Context, org, repo, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error)
}

type downloader struct {
	Organization string
	repositories RepositoriesService
}

func newDownloader(svc RepositoriesService, org string) *downloader {
	return &downloader{
		Organization: org,
		repositories: svc,
	}
}

func (d *downloader) BufferDirectory(ctx context.Context, repo, path, ref string) ([]*chartutil.BufferedFile, error) {
	file, dir, _, err := d.repositories.GetContents(ctx, d.Organization, repo, path, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	if err != nil {
		return nil, err
	}

	var files []*chartutil.BufferedFile

	if file != nil {
		content, err := file.GetContent()
		if err != nil {
			return nil, errors.Wrap(err, "unable to get file contents")
		}

		files = append(files, &chartutil.BufferedFile{
			Name: file.GetName(),
			Data: []byte(content),
		})
	} else {
		for _, subDir := range dir {
			subFiles, err := d.BufferDirectory(ctx, repo, subDir.GetPath(), ref)
			if err != nil {
				return nil, err
			}

			files = append(files, subFiles...)
		}
	}

	return files, nil
}
