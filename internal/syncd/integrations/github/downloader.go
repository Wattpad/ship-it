package github

import (
	"context"
	"strings"

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

// BufferDirectory recursively buffers files in a github directory. The
// filenames of buffered files are relative to the directory root, which
// is required by 'chartutils.LoadFiles'
func (d *downloader) BufferDirectory(ctx context.Context, repo, path, ref string) ([]*chartutil.BufferedFile, error) {
	trimPrefix := func(p string) string {
		if p == path {
			return p
		}
		return strings.TrimPrefix(strings.TrimPrefix(p, path), "/")
	}

	return d.bufferDirectory(ctx, repo, path, trimPrefix, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
}

func (d *downloader) bufferDirectory(ctx context.Context, repo, path string, trim func(string) string, ref *github.RepositoryContentGetOptions) ([]*chartutil.BufferedFile, error) {
	file, dir, _, err := d.repositories.GetContents(ctx, d.Organization, repo, path, ref)
	if err != nil {
		return nil, err
	}

	if file != nil {
		content, err := file.GetContent()
		if err != nil {
			return nil, errors.Wrap(err, "unable to get file contents")
		}

		return []*chartutil.BufferedFile{
			{
				Name: trim(path),
				Data: []byte(content),
			},
		}, nil
	}

	var files []*chartutil.BufferedFile

	for _, subDir := range dir {
		subFiles, err := d.bufferDirectory(ctx, repo, subDir.GetPath(), trim, ref)
		if err != nil {
			return nil, err
		}

		files = append(files, subFiles...)
	}

	return files, nil
}
