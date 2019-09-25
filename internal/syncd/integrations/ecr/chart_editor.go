package ecr

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"ship-it/internal/image"
	"ship-it/internal/unstructured"

	"github.com/google/go-github/v26/github"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/types"
)

type GitService interface {
	CreateCommit(ctx context.Context, owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error)
	CreateTree(ctx context.Context, owner string, repo string, baseTree string, entries []github.TreeEntry) (*github.Tree, *github.Response, error)
	GetBlob(ctx context.Context, owner string, repo string, sha string) (*github.Blob, *github.Response, error)
	GetCommit(ctx context.Context, owner string, repo string, sha string) (*github.Commit, *github.Response, error)
	GetRef(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error)
	GetTree(ctx context.Context, owner string, repo string, sha string, recursive bool) (*github.Tree, *github.Response, error)
	UpdateRef(ctx context.Context, owner string, repo string, ref *github.Reference, force bool) (*github.Reference, *github.Response, error)
}

type chartEditor struct {
	github     GitService
	ChartPath  string
	Org        string
	Repository string
	Ref        string
}

func NewChartEditor(g GitService, org, repo, ref, path string) *chartEditor {
	return &chartEditor{
		github:     g,
		ChartPath:  path,
		Org:        org,
		Repository: repo,
		Ref:        ref,
	}
}

func (c *chartEditor) Edit(ctx context.Context, releases []types.NamespacedName, image *image.Ref) error {
	// Get commit SHA of the targeted branch (ref)
	ref, _, err := c.github.GetRef(ctx, c.Org, c.Repository, "refs/heads/"+c.Ref)
	if err != nil {
		return errors.Wrap(err, "failed to get reference")
	}

	// Get HEAD commit object for the branch
	parent, _, err := c.github.GetCommit(ctx, c.Org, c.Repository, ref.GetObject().GetSHA())
	if err != nil {
		return errors.Wrap(err, "failed to get commit")
	}

	// Get content tree pointed to by the commit
	tree, _, err := c.github.GetTree(ctx, c.Org, c.Repository, ref.GetObject().GetSHA(), true /* recursive */)
	if err != nil {
		return errors.Wrap(err, "failed to get tree")
	}

	// Modify the content tree that the commit points to
	entries := c.editTreeEntries(ctx, releases, image, tree.Entries)

	// Create a new content tree with the modified entries, computing
	// Merkle-esque SHAs and so forth
	tree, _, err = c.github.CreateTree(ctx, c.Org, c.Repository, ref.GetObject().GetSHA(), entries)
	if err != nil {
		return errors.Wrap(err, "failed to create tree")
	}

	// Create a new commit object with the current commit as the parent and
	// the new tree, getting a new commit back
	commit := &github.Commit{
		Message: github.String(commitMessage(releases, image)),
		Tree:    tree,
		Parents: []github.Commit{*parent},
	}

	commit, _, err = c.github.CreateCommit(ctx, c.Org, c.Repository, commit)
	if err != nil {
		return errors.Wrap(err, "failed to create commit")
	}

	// Update the reference of the branch to point to the new commit SHA
	ref.Object.SHA = commit.SHA
	_, _, err = c.github.UpdateRef(ctx, c.Org, c.Repository, ref, false /* force */)
	return errors.Wrap(err, "failed to update reference")
}

func commitMessage(releases []types.NamespacedName, image *image.Ref) string {
	names := make([]string, 0, len(releases))

	for _, r := range releases {
		names = append(names, r.Name)
	}

	return fmt.Sprintf("Updated %d helm charts using image %s\n\n%s", len(names), image, strings.Join(names, "\n"))
}

func (c *chartEditor) editTreeEntries(ctx context.Context, names []types.NamespacedName, image *image.Ref, entries []github.TreeEntry) []github.TreeEntry {
	var edited []github.TreeEntry

	for _, e := range entries {
		// // skip entries that don't belong to the chart
		if !strings.HasPrefix(e.GetPath(), c.ChartPath) {
			continue
		}

		for _, name := range names {
			// skip tree entries that don't have the expected
			// HelmRelease filename for the associated release
			valuesYaml := fmt.Sprintf("%s.yaml", name.Name)

			if !strings.HasSuffix(e.GetPath(), valuesYaml) {
				continue
			}

			content, err := c.editBlob(ctx, e.GetSHA(), image)
			if err == nil {
				edited = append(edited, github.TreeEntry{
					Content: github.String(content),
					Mode:    e.Mode,
					Path:    e.Path,
					Size:    e.Size,
					Type:    e.Type,
				})
			}

			break
		}
	}

	return edited
}

func (c *chartEditor) editBlob(ctx context.Context, sha string, image *image.Ref) (string, error) {
	blob, _, err := c.github.GetBlob(ctx, c.Org, c.Repository, sha)
	if err != nil {
		return "", errors.Wrap(err, "failed to get blob")
	}

	content, err := base64.StdEncoding.DecodeString(blob.GetContent())
	if err != nil {
		return "", errors.Wrapf(err, "failed to decode blob content")
	}

	// use MapSlice to preserve the order of fields
	var values yaml.MapSlice
	if err := yaml.Unmarshal(content, &values); err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal values YAML")
	}

	edited := editYaml(values, image)

	bytes, err := yaml.Marshal(edited)
	if err != nil {
		return "", errors.Wrapf(err, "failed to marshal values YAML")
	}

	return string(bytes), nil
}

func editYaml(spec yaml.MapSlice, desired *image.Ref) yaml.MapSlice {
	// this predicate selects the matching image block
	pred := func(item yaml.MapItem) bool {
		other, err := image.FromYaml(item)
		if err != nil {
			return false
		}

		return desired.Matches(*other)
	}

	// this visitor mutates the selected image block by setting the desired tag
	visit := func(item *yaml.MapItem) {
		if obj, ok := item.Value.(yaml.MapSlice); ok {
			for i := range obj {
				if k, ok := obj[i].Key.(string); ok && k == "tag" {
					obj[i].Value = desired.Tag
					break
				}
			}
		}
	}

	unstructured.VisitOne(spec, pred, visit)
	return spec
}
