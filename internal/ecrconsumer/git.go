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
	CreatePullRequest(title string, body string, branchToMerge string) (*github.PullRequest, error)
	CreateBranch(name string, base string) (*github.Branch, error)
	UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error)
	CreateLabel(name string, color string, description string) (*github.Label, error)
	DeleteLabel(name string) (*github.Response, error)
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

func (c *GitHub) CreatePullRequest(title string, body string, branchToMerge string) (*github.PullRequest, error) {
	ctx := context.Background()

	newPR := &github.NewPullRequest{
		Title:               github.String(title),
		Head:                github.String(branchToMerge),
		Base:                github.String("master"),
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := c.Client.PullRequests.Create(ctx, c.Organization, c.Respository, newPR)
	if err != nil {
		return nil, err
	}

	// Before returning the pull request add the kube deploy label to it
	_, _, err = c.Client.Issues.AddLabelsToIssue(ctx, c.Organization, c.Respository, *pr.Number, []string{"kube"})
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (c *GitHub) CreateBranch(name string, base string) (*github.Branch, error) {
	ctx := context.Background()
	// Get Reference of most recent commit to base branch
	r, _, err := c.Client.Git.GetRef(ctx, c.Organization, c.Respository, "refs/heads/"+base)
	if err != nil {
		return nil, err
	}

	newName := "refs/heads/" + name
	r.Ref = &newName // set new name
	// Use reference to create a new branch
	_, _, err = c.Client.Git.CreateRef(ctx, c.Organization, c.Respository, r)
	if err != nil {
		return nil, err
	}

	// Get Branch Object
	branch, _, err := c.Client.Repositories.GetBranch(ctx, c.Organization, c.Respository, name)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

func (c *GitHub) UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error) {
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

func (c *GitHub) CreateLabel(name string, color string, description string) (*github.Label, error) {
	newLabel := &github.Label{
		Name:        &name,
		Color:       &color,
		Description: &description,
	}
	label, _, err := c.Client.Issues.CreateLabel(context.Background(), c.Organization, c.Respository, newLabel)
	if err != nil {
		return nil, err
	}

	return label, nil
}

func (c *GitHub) DeleteLabel(name string) (*github.Response, error) {
	resp, err := c.Client.Issues.DeleteLabel(context.Background(), c.Organization, c.Respository, name)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GitHub) DownloadFile(branch string, path string) ([]byte, error) {
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

func (c *GitHub) GetFile(branch string, path string) (*github.RepositoryContent, error) {
	contents, _, _, err := c.Client.Repositories.GetContents(context.Background(), c.Organization, c.Respository, path, &github.RepositoryContentGetOptions{Ref: "refs/heads/" + branch})
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func (c *GitHub) GetDirectory(branch string, path string) ([]*github.RepositoryContent, error) {
	_, directory, _, err := c.Client.Repositories.GetContents(context.Background(), c.Organization, c.Respository, path, &github.RepositoryContentGetOptions{Ref: "refs/heads/" + branch})
	if err != nil {
		return nil, err
	}
	return directory, nil
}

func (c *GitHub) SaveDirectory(branch string, repoPath string, localPath string) error { // saves to the same relative path as the path in the github repo
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
