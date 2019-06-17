package github

import (
	"context"
	"errors"

	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
)

// GitHub handles integrations with the GitHub API
type GitHub struct {
	Org    string
	Client *github.Client
}

// New creates a new GitHub integration
func New(org string, accessToken string) (*GitHub, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHub{
		Client: client,
		Org:    org,
	}, nil
}

// GetTravisCIBuildURLForRef uses the Checks API to find the URL to the Travis build for the specified ref
func (github *GitHub) GetTravisCIBuildURLForRef(repo string, ref string) (string, error) {
	ctx := context.Background()

	checks, _, err := github.Client.Checks.ListCheckRunsForRef(ctx, github.Org, repo, ref, nil)

	if err != nil {
		return "", err
	}

	for _, checkRun := range checks.CheckRuns {
		if *checkRun.App.Name == "Travis CI" {
			return *checkRun.DetailsURL, nil
		}
	}

	return "", errors.New("Could not find a TravisCI build for specified ref")
}
