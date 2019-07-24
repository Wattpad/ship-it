package ecr

import (
	"context"

	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
)

type GitHub struct {
	client       *github.Client
	Organization string
	Repository   string
}

func NewGitHub(ctx context.Context, token string, org string, repo string) GitHub {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	client := GitHub{
		client:       github.NewClient(oauth2.NewClient(ctx, tokenSource)),
		Organization: org,
		Repository:   repo,
	}

	return client
}

func (c GitHub) UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error) {
	ctx := context.Background()
	// Get File's Blob SHA
	contents, _, _, err := c.client.Repositories.GetContents(ctx, c.Organization, c.Repository, path, &github.RepositoryContentGetOptions{Ref: "refs/heads/" + branch})
	if err != nil {
		return nil, err
	}

	options := &github.RepositoryContentFileOptions{ // Add commit author
		Message: github.String(msg),
		Content: fileContent,
		SHA:     github.String(contents.GetSHA()),
		Branch:  github.String(branch),
	}

	resp, _, err := c.client.Repositories.UpdateFile(ctx, c.Organization, c.Repository, path, options)

	return resp, err
}

func (c GitHub) GetFile(branch string, path string) (string, error) {
	contents, _, _, err := c.client.Repositories.GetContents(context.Background(), c.Organization, c.Repository, path, &github.RepositoryContentGetOptions{Ref: "refs/heads/" + branch})
	if err != nil {
		return "", err
	}

	return contents.GetContent()
}
