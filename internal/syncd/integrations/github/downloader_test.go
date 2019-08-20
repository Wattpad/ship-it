package github

import (
	"context"
	"path"
	"testing"

	"github.com/google/go-github/v26/github"
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
	testOrg, testRepo, testRef := "testOrg", "testRepo", "master"

	testPath, testContent := "test/file", "content"

	var m mockGithubRepositoriesClient

	m.On("GetContents", mock.Anything, testOrg, testRepo, testPath, &github.RepositoryContentGetOptions{
		Ref: testRef,
	}).Return(&github.RepositoryContent{
		Content: github.String(testContent),
		Name:    github.String(testPath),
	}, nil, nil, nil)

	downloader := newDownloader(&m, testOrg)
	files, err := downloader.BufferDirectory(context.Background(), testRepo, testPath, testRef)

	if assert.NoError(t, err) {
		assert.Len(t, files, 1)
		assert.Equal(t, &chartutil.BufferedFile{
			Name: testPath,
			Data: []byte(testContent),
		}, files[0])
	}
}

func TestBufferChartDirectory(t *testing.T) {
	testOrg, testRepo, testPath, testRef := "testOrg", "testRepo", "test/path", "master"

	templates := "templates"
	templatesPath := path.Join(testPath, templates)

	k8sYaml := "templates/k8s.yaml"
	k8sPath := path.Join(testPath, k8sYaml)

	chartYaml := "Chart.yaml"
	chartPath := path.Join(testPath, chartYaml)

	valuesYaml := "values.yaml"
	valuesPath := path.Join(testPath, valuesYaml)

	var m mockGithubRepositoriesClient

	// first call returns references to the directory contents
	m.On("GetContents", mock.Anything, testOrg, testRepo, testPath, &github.RepositoryContentGetOptions{
		Ref: testRef,
	}).Return(nil, []*github.RepositoryContent{
		{
			Path: github.String(templatesPath),
		},
		{
			Path: github.String(chartPath),
		},
		{
			Path: github.String(valuesPath),
		},
	}, nil, nil)

	// templates directory
	m.On("GetContents", mock.Anything, testOrg, testRepo, templatesPath, &github.RepositoryContentGetOptions{
		Ref: testRef,
	}).Return(nil, []*github.RepositoryContent{
		{
			Path: github.String(k8sPath),
		},
	}, nil, nil)

	m.On("GetContents", mock.Anything, testOrg, testRepo, k8sPath, &github.RepositoryContentGetOptions{
		Ref: testRef,
	}).Return(&github.RepositoryContent{
		Content: github.String(k8sYaml),
		Name:    github.String(k8sYaml),
	}, nil, nil, nil)

	m.On("GetContents", mock.Anything, testOrg, testRepo, valuesPath, &github.RepositoryContentGetOptions{
		Ref: testRef,
	}).Return(&github.RepositoryContent{
		Content: github.String(valuesYaml),
		Name:    github.String(valuesYaml),
	}, nil, nil, nil)

	m.On("GetContents", mock.Anything, testOrg, testRepo, chartPath, &github.RepositoryContentGetOptions{
		Ref: testRef,
	}).Return(&github.RepositoryContent{
		Content: github.String(chartYaml),
		Name:    github.String(chartYaml),
	}, nil, nil, nil)

	downloader := newDownloader(&m, testOrg)
	files, err := downloader.BufferDirectory(context.Background(), testRepo, testPath, testRef)

	if assert.NoError(t, err) {
		assert.Len(t, files, 3)
		assert.ElementsMatch(t, []*chartutil.BufferedFile{
			{
				Name: k8sYaml,
				Data: []byte(k8sYaml),
			},
			{
				Name: valuesYaml,
				Data: []byte(valuesYaml),
			},
			{
				Name: chartYaml,
				Data: []byte(chartYaml),
			},
		}, files)
	}
}
