package github

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
)

type repositoriesService interface {
	GetContents(ctx context.Context, org, repo, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error)
}

type githubDownloader struct {
	Organization string
	repositories repositoriesService
}

func newGithubDownloader(c *github.Client, org string) *githubDownloader {
	return &githubDownloader{
		Organization: org,
		repositories: c.Repositories,
	}
}

func (g *githubDownloader) BufferDirectory(ctx context.Context, repo, path, ref string) ([]*chartutil.BufferedFile, error) {
	file, dir, _, err := g.repositories.GetContents(ctx, g.Organization, repo, path, &github.RepositoryContentGetOptions{
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
			subFiles, err := g.BufferDirectory(ctx, repo, subDir.GetPath(), ref)
			if err != nil {
				return nil, err
			}

			files = append(files, subFiles...)
		}
	}

	return files, nil
}
