package github

import (
	"context"

	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
)

// GitHub handles integrations with the GitHub API
type GitHub struct {
	Client *github.Client
}

// New creates a new GitHub integration
func New(accessToken string) (*GitHub, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHub{
		Client: client,
	}, nil
}

// GetTravisBuildURLForRef uses the Checks API to find the URL to the Travis build for the supplied ref
func (github *GitHub) GetTravisBuildURLForRef(ref string) (string, error) {
	return "https://stub", nil
}
