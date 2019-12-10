package ecr

import (
	"context"
	"encoding/base64"
	"path"
	"testing"

	"ship-it/internal/image"

	"github.com/google/go-github/v26/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/types"
)

type mockGitService struct {
	mock.Mock
}

func (m *mockGitService) CreateCommit(ctx context.Context, owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error) {
	args := m.Called(ctx, owner, repo, commit)
	return args.Get(0).(*github.Commit), args.Get(1).(*github.Response), args.Error(2)
}

func (m *mockGitService) CreateTree(ctx context.Context, owner string, repo string, tree string, entries []github.TreeEntry) (*github.Tree, *github.Response, error) {
	args := m.Called(ctx, owner, repo, tree, entries)
	return args.Get(0).(*github.Tree), args.Get(1).(*github.Response), args.Error(2)
}

func (m *mockGitService) GetBlob(ctx context.Context, owner string, repo string, sha string) (*github.Blob, *github.Response, error) {
	args := m.Called(ctx, owner, repo, sha)
	return args.Get(0).(*github.Blob), args.Get(1).(*github.Response), args.Error(2)
}

func (m *mockGitService) GetCommit(ctx context.Context, owner string, repo string, sha string) (*github.Commit, *github.Response, error) {
	args := m.Called(ctx, owner, repo, sha)
	return args.Get(0).(*github.Commit), args.Get(1).(*github.Response), args.Error(2)
}

func (m *mockGitService) GetRef(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error) {
	args := m.Called(ctx, owner, repo, ref)
	return args.Get(0).(*github.Reference), args.Get(1).(*github.Response), args.Error(2)
}

func (m *mockGitService) GetTree(ctx context.Context, owner string, repo string, sha string, recursive bool) (*github.Tree, *github.Response, error) {
	args := m.Called(ctx, owner, repo, sha, recursive)
	return args.Get(0).(*github.Tree), args.Get(1).(*github.Response), args.Error(2)
}

func (m *mockGitService) UpdateRef(ctx context.Context, owner string, repo string, ref *github.Reference, force bool) (*github.Reference, *github.Response, error) {
	args := m.Called(ctx, owner, repo, ref, force)
	return args.Get(0).(*github.Reference), args.Get(1).(*github.Response), args.Error(2)
}

func TestEdit(t *testing.T) {
	ctx := context.Background()

	testOrg, testRepo, testRef, testPath := "testOrg", "testRepo", "testRef", "path/to/registry/chart"

	unmodifiedTreeEntries := []github.TreeEntry{
		{
			SHA:  github.String("sha0"),
			Path: github.String("not/the/path"),
		},
		{
			SHA:  github.String("sha1"),
			Path: github.String("still/not/the/path"),
		},
	}

	// entries contained in the registry chart path should be modified
	modifiedTreeEntries := []github.TreeEntry{
		{
			SHA:  github.String("sha2"),
			Path: github.String(path.Join(testPath, "/templates/foo.yaml")),
		},
		{
			SHA:  github.String("sha3"),
			Path: github.String(path.Join(testPath, "/templates/bar.yaml")),
		},
	}

	treeEntries := append(unmodifiedTreeEntries, modifiedTreeEntries...)

	blobContents := [][]byte{
		[]byte(`kind: HelmRelease
apiVersion: shipit.wattpad.com/v1beta1
metadata:
    name: foo-release
spec:
    values:
        image:
            repository: hub.docker.com/foo,
            tag: fooreleaseoldtag
`),
		[]byte(`kind: HelmRelease
apiVersion: shipit.wattpad.com/v1beta1
metadata:
    name: bar-release
spec:
    values:
        image:
            repository: hub.docker.com/bar,
            tag: barreleaseoldtag
`),
	}

	releases := []types.NamespacedName{
		{
			Name:      "foo",
			Namespace: "foo-namespace",
		},
		{
			Name:      "bar",
			Namespace: "bar-namespace",
		},
	}

	desired := image.Ref{
		Registry:   "test-registry",
		Repository: "test-repository",
		Tag:        "new-tag",
	}

	baseSHA := "base-sha"
	headSHA := "head-sha"

	mockGit := new(mockGitService)

	mockGit.On("GetRef", ctx, testOrg, testRepo, "refs/heads/"+testRef).Return(
		&github.Reference{
			Object: &github.GitObject{
				SHA: github.String(baseSHA),
			},
		}, &github.Response{}, nil,
	) // so much implementation detail...

	mockGit.On("GetCommit", ctx, testOrg, testRepo, baseSHA).Return(
		&github.Commit{
			SHA: github.String(baseSHA),
		}, &github.Response{}, nil,
	)

	mockGit.On("GetTree", ctx, testOrg, testRepo, baseSHA, true /* recursive */).Return(
		&github.Tree{
			SHA:     github.String(baseSHA),
			Entries: treeEntries,
		}, &github.Response{}, nil,
	)

	// expect GetBlob per changed file
	for i, entry := range modifiedTreeEntries {
		blob64 := base64.StdEncoding.EncodeToString(blobContents[i])
		mockGit.On("GetBlob", ctx, testOrg, testRepo, entry.GetSHA()).Return(
			&github.Blob{
				Content: github.String(blob64),
			}, &github.Response{}, nil,
		)
	}

	mockGit.On("CreateTree", ctx, testOrg, testRepo, baseSHA, mock.AnythingOfType("[]github.TreeEntry")).Return(
		&github.Tree{
			SHA: github.String(headSHA),
		}, &github.Response{}, nil,
	)

	mockGit.On("CreateCommit", ctx, testOrg, testRepo, mock.AnythingOfType("*github.Commit")).Return(
		&github.Commit{
			SHA: github.String(headSHA),
		}, &github.Response{}, nil,
	)

	mockGit.On("UpdateRef", ctx, testOrg, testRepo, &github.Reference{
		Object: &github.GitObject{
			SHA: github.String(headSHA),
		},
	}, false).Return(
		&github.Reference{}, &github.Response{}, nil,
	)

	editor := NewChartEditor(mockGit, testOrg, testRepo, testRef, testPath)

	assert.NoError(t, editor.Edit(ctx, releases, &desired))
}

func TestEditYaml(t *testing.T) {
	original := yaml.MapSlice{
		{
			Key: "values",
			Value: yaml.MapSlice{
				{
					Key: "image",
					Value: yaml.MapSlice{
						{
							Key:   "repository",
							Value: "foo/bar",
						},
						{

							Key:   "tag",
							Value: "oldtag",
						},
					},
				},
			},
		},
	}

	expected := yaml.MapSlice{
		{
			Key: "values",
			Value: yaml.MapSlice{
				{
					Key: "image",
					Value: yaml.MapSlice{
						{
							Key:   "repository",
							Value: "foo/bar",
						},
						{

							Key:   "tag",
							Value: "newtag",
						},
					},
				},
			},
		},
	}

	desired := image.Ref{
		Registry:   "foo",
		Repository: "bar",
		Tag:        "newtag",
	}

	assert.Equal(t, expected, editYaml(original, &desired))
}
