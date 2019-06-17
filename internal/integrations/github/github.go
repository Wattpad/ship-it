package github

import (
	"context"
	"errors"

	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
)

// Github handles integrations with the Github API
type Github struct {
	Org    string
	Client *github.Client
}

// New creates a new GitHub integration
func New(ctx context.Context, org string, accessToken string) *Github {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &Github{
		Client: client,
		Org:    org,
	}
}

// GetTravisCIBuildURLForRef uses the Checks API to find the URL to the Travis build for the specified ref
func (g *Github) GetTravisCIBuildURLForRef(ctx context.Context, repo string, ref string) (string, error) {
	checks, _, err := g.Client.Checks.ListCheckRunsForRef(ctx, g.Org, repo, ref, nil)

	if err != nil {
		return "", err
	}

	for _, checkRun := range checks.CheckRuns {
		if checkRun.GetApp().GetName() == "Travis CI" {
			return checkRun.GetDetailsURL(), nil
		}
	}

	return "", errors.New("could not find a TravisCI build for specified ref")
}
