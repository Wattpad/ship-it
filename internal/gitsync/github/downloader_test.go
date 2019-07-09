package github

import (
	"context"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/helm/pkg/chartutil"
)

type mockGithubRepositoriesClient struct {
	mock.Mock
}

func (m *mockGithubRepositoriesClient) GetContents(ctx context.Context, org, repo, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error) {
	args := m.Called(ctx, org, repo, path, opts)

	var ret0 *github.RepositoryContent
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.(*github.RepositoryContent)
	}

	var ret1 []*github.RepositoryContent
	if args1 := args.Get(1); args1 != nil {
		ret1 = args1.([]*github.RepositoryContent)
	}

	var ret2 *github.Response
	if args2 := args.Get(2); args2 != nil {
		ret2 = args2.(*github.Response)
	}

	return ret0, ret1, ret2, args.Error(3)
}

func TestBufferDirectorySingleFile(t *testing.T) {
	testOrg, testRepo, testPath, testRef := "testOrg", "testRepo", "test/path", "master"

	testFilename, testContent := "filename", "content"

	var m mockGithubRepositoriesClient

	m.On("GetContents", mock.Anything, testOrg, testRepo, testPath, &github.RepositoryContentGetOptions{
		Ref: testRef,
	}).Return(&github.RepositoryContent{
		Content: github.String(testContent),
		Name:    github.String(testFilename),
	}, nil, nil, nil)

	downloader := &githubDownloader{
		Organization: testOrg,
		repositories: &m,
	}

	files, err := downloader.BufferDirectory(context.Background(), testRepo, testPath, testRef)

	if assert.NoError(t, err) {
		assert.Len(t, files, 1)
		assert.Equal(t, files[0], &chartutil.BufferedFile{
			Name: testFilename,
			Data: []byte(testContent),
		})
	}
}

func TestBufferDirectoryFlatDirectory(t *testing.T) {
	testOrg, testRepo, testPath, testRef := "testOrg", "testRepo", "test/path", "master"

	testPath1, testPath2, testPath3 := "test/path/1", "test/path/2", "test/path/3"

	testFilename1, testContent1 := "filename1", "content1"
	testFilename2, testContent2 := "filename2", "content2"
	testFilename3, testContent3 := "filename3", "content3"

	var m mockGithubRepositoriesClient

	// first call returns references to the directory contents
	m.On("GetContents", mock.Anything, testOrg, testRepo, testPath, &github.RepositoryContentGetOptions{
		Ref: testRef,
	}).Return(nil, []*github.RepositoryContent{
		{
			Path: github.String(testPath1),
		},
		{
			Path: github.String(testPath2),
		},
		{
			Path: github.String(testPath3),
		},
	}, nil, nil)

	// following calls return references to file contents
	m.On("GetContents", mock.Anything, testOrg, testRepo, testPath1, &github.RepositoryContentGetOptions{
		Ref: testRef,
	}).Return(&github.RepositoryContent{
		Content: github.String(testContent1),
		Name:    github.String(testFilename1),
	}, nil, nil, nil)

	m.On("GetContents", mock.Anything, testOrg, testRepo, testPath2, &github.RepositoryContentGetOptions{
		Ref: testRef,
	}).Return(&github.RepositoryContent{
		Content: github.String(testContent2),
		Name:    github.String(testFilename2),
	}, nil, nil, nil)

	m.On("GetContents", mock.Anything, testOrg, testRepo, testPath3, &github.RepositoryContentGetOptions{
		Ref: testRef,
	}).Return(&github.RepositoryContent{
		Content: github.String(testContent3),
		Name:    github.String(testFilename3),
	}, nil, nil, nil)

	downloader := &githubDownloader{
		Organization: testOrg,
		repositories: &m,
	}

	files, err := downloader.BufferDirectory(context.Background(), testRepo, testPath, testRef)

	if assert.NoError(t, err) {
		assert.Len(t, files, 3)
		assert.EqualValues(t, files, []*chartutil.BufferedFile{
			{
				Name: testFilename1,
				Data: []byte(testContent1),
			},
			{
				Name: testFilename2,
				Data: []byte(testContent2),
			},
			{
				Name: testFilename3,
				Data: []byte(testContent3),
			},
		})
	}
}
