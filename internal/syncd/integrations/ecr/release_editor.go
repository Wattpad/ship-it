package ecr

import (
	"context"
	"fmt"
	"ship-it/internal"
	"strings"

	"github.com/google/go-github/v26/github"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/types"
)

type releaseEditor struct {
	github     GithubCommitter
	Org        string
	Repository string
	Ref        string
}

func (e *releaseEditor) Edit(ctx context.Context, releases []types.NamespacedName, image *internal.Image) error {
	parent, _, err := e.github.GetCommit(ctx, e.Org, e.Repository, e.Ref)
	if err != nil {
		return errors.Wrap(err, "failed to get github commit")
	}

	commit := &github.Commit{
		Message: github.String(commitMessage(releases, image)),
		Tree:    editTree(releases, image, parent.GetTree()),
		Parents: []github.Commit{*parent},
	}

	_, _, err = e.github.CreateCommit(ctx, e.Org, e.Repository, commit)
	return errors.Wrap(err, "failed to create github commit")
}

func commitMessage(releases []types.NamespacedName, image *internal.Image) string {
	names := make([]string, 0, len(releases))

	for _, r := range releases {
		names = append(names, r.Name)
	}

	return fmt.Sprintf("Updated %d helm charts using image %s\n\n%s", len(names), image.Tagged(), strings.Join(names, "\n"))
}

func isValuesYamlForRelease(entry github.TreeEntry, name types.NamespacedName) bool {
	return strings.HasSuffix(entry.GetPath(), "values.yaml") && strings.Contains(entry.GetPath(), name.Name)
}

func editTree(names []types.NamespacedName, image *internal.Image, tree *github.Tree) *github.Tree {
	entries := make([]github.TreeEntry, 0, len(tree.Entries))

	for _, e := range tree.Entries {
		for _, name := range names {
			if isValuesYamlForRelease(e, name) {
				content, err := editValues(e.GetContent(), image)
				if err == nil {
					e.Content = github.String(content)
				}
				// TODO handle error?
			}

			entries = append(entries, e)
		}

	}

	return &github.Tree{
		SHA:     tree.SHA,
		Entries: entries,
	}
}

func editValues(content string, image *internal.Image) (string, error) {
	values := make(map[string]interface{})

	if err := yaml.Unmarshal([]byte(content), &values); err != nil {
		return "", errors.Wrapf(err, "failed to parse the YAML file")
	}

	edited := internal.WithImage(*image, values)

	bytes, err := yaml.Marshal(edited)
	if err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal the edited YAML file")
	}

	return string(bytes), nil
}
