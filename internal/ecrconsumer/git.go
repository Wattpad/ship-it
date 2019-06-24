package ecrconsumer

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitCommands interface {
	UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error)
	DownloadFile(branch string, path string) ([]byte, error)
}

type GitHub struct {
	Client       *github.Client
	Organization string
	Respository  string
}

func NewGitHub(token string, ctx context.Context, org string, repo string) GitHub {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(ctx, tokenSource)

	client := GitHub{
		Client:       github.NewClient(tokenClient),
		Organization: org,
		Respository:  repo,
	}

	return client
}

func (c GitHub) UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error) {
	ctx := context.Background()
	// Get File's Blob SHA
	contents, _, _, err := c.Client.Repositories.GetContents(ctx, c.Organization, c.Respository, path, &github.RepositoryContentGetOptions{Ref: "refs/heads/" + branch})
	if err != nil {
		return nil, err
	}
	fileSha := contents.GetSHA()

	options := &github.RepositoryContentFileOptions{ // Add commit author
		Message: &msg,
		Content: fileContent,
		SHA:     &fileSha,
		Branch:  &branch,
	}

	resp, _, err := c.Client.Repositories.UpdateFile(ctx, c.Organization, c.Respository, path, options)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c GitHub) DownloadFile(branch string, path string) ([]byte, error) {
	contents, _, _, err := c.Client.Repositories.GetContents(context.Background(), c.Organization, c.Respository, path, &github.RepositoryContentGetOptions{Ref: "refs/heads/" + branch})
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(*contents.Content)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c GitHub) GetFile(branch string, path string) (*github.RepositoryContent, error) {
	contents, _, _, err := c.Client.Repositories.GetContents(context.Background(), c.Organization, c.Respository, path, &github.RepositoryContentGetOptions{Ref: "refs/heads/" + branch})
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func (c GitHub) GetDirectory(branch string, path string) ([]*github.RepositoryContent, error) {
	_, directory, _, err := c.Client.Repositories.GetContents(context.Background(), c.Organization, c.Respository, path, &github.RepositoryContentGetOptions{Ref: "refs/heads/" + branch})
	if err != nil {
		return nil, err
	}
	return directory, nil
}

func (c GitHub) SaveDirectory(branch string, repoPath string, localPath string) error { // saves to the same relative path as the path in the github repo
	// Save Directory to local folder from the git repository
	file, dir, _, err := c.Client.Repositories.GetContents(context.Background(), c.Organization, c.Respository, repoPath, &github.RepositoryContentGetOptions{Ref: "refs/heads/" + branch})
	if err != nil {
		return err
	}

	for _, d := range dir { // since we are using range the loop does not execute in the nil case
		if err := c.SaveDirectory(branch, *d.Path, localPath); err != nil {
			return err
		}
	}
	if dir == nil {
		path := filepath.Join(localPath, *file.Path)

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil { // create all folders in the relative path
			return err
		}

		fileContent, err := file.GetContent()
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(path, []byte(fileContent), 0700); err != nil { // write the files in the corrects spots in the directory tree
			return err
		}
	}
	return nil
}
