package ecr

import (
	"context"
	"fmt"
	"path/filepath"
	"ship-it/internal"

	"github.com/google/go-github/v26/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type GitHub struct {
	client            *github.Client
	Organization      string
	Branch            string
	Repository        string
	RegistryChartPath string
}

func NewGitHub(ctx context.Context, token string, org string, repo string, branch string, registryPath string) GitHub {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	return GitHub{
		client:            github.NewClient(oauth2.NewClient(ctx, tokenSource)),
		Organization:      org,
		Branch:            branch,
		Repository:        repo,
		RegistryChartPath: registryPath,
	}
}

func transformBytes(in []byte, image *internal.Image) ([]byte, error) {
	rls := make(map[string]interface{})

	err := yaml.Unmarshal(in, rls)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse the YAML file")
	}

	updatedRls := internal.WithImage(*image, rls)

	return yaml.Marshal(updatedRls)
}

func (c GitHub) getFile(ctx context.Context, path string) (*github.RepositoryContent, error) {
	contents, _, _, err := c.client.Repositories.GetContents(ctx, c.Organization, c.Repository, path, &github.RepositoryContentGetOptions{Ref: "refs/heads/" + c.Branch})
	return contents, err
}

func (c GitHub) updateFile(ctx context.Context, path string, fileContent []byte, SHA string, msg string) error {
	options := &github.RepositoryContentFileOptions{
		Message: github.String(msg),
		Content: fileContent,
		SHA:     github.String(SHA),
		Branch:  github.String(c.Branch),
	}

	_, _, err := c.client.Repositories.UpdateFile(ctx, c.Organization, c.Repository, path, options)
	return err
}

func (c GitHub) UpdateAndReplace(ctx context.Context, releaseName string, image *internal.Image) error {
	path := filepath.Join(c.RegistryChartPath, releaseName+".yaml")
	contents, err := c.getFile(ctx, path)
	if err != nil {
		return err
	}

	resourceStr, err := contents.GetContent()
	if err != nil {
		return err
	}

	fileContent, err := transformBytes([]byte(resourceStr), image)
	if err != nil {
		return errors.Wrap(err, "failed to update the image tag in file")
	}

	return c.updateFile(ctx, path, fileContent, contents.GetSHA(), fmt.Sprintf("Image Tag updated to %s", image.Tag))
}
