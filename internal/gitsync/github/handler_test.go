package github

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type mockDownloader struct {
	mock.Mock
}

func (m *mockDownloader) BufferDirectory(ctx context.Context, repo, path, ref string) ([]*chartutil.BufferedFile, error) {
	args := m.Called(ctx, repo, path, ref)
	var ret0 []*chartutil.BufferedFile
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.([]*chartutil.BufferedFile)
	}

	return ret0, args.Error(1)
}

type mockReleaseSyncer struct {
	mock.Mock
}

func (m *mockReleaseSyncer) UpdateOrInstallFromChart(ctx context.Context, chart *chart.Chart) error {
	args := m.Called(ctx, chart)
	return args.Error(0)
}

func TestHandleMessageChartDownloadFails(t *testing.T) {
	errNotFound := errors.New("not found")
	testEvent := pushEvent{
		Ref:        "ref",
		Path:       "path",
		Repository: "repository",
	}

	var md mockDownloader
	md.On("BufferDirectory", mock.Anything, testEvent.Repository, testEvent.Path, testEvent.Ref).Return(nil, errNotFound)

	handler := newHandler(&md, nil)

	err := handler.handleMessage(context.Background(), testEvent)
	assert.Error(t, errNotFound, err)
}

func TestHandleMessage(t *testing.T) {
	testEvent := pushEvent{
		Ref:        "ref",
		Path:       "path",
		Repository: "repository",
	}

	testChartFiles := []*chartutil.BufferedFile{
		{
			Name: "Chart.yaml",
			Data: []byte(`
apiVersion: v1
name: foo
description: This is a foo.
`),
		},
	}

	var md mockDownloader
	md.On("BufferDirectory", mock.Anything, testEvent.Repository, testEvent.Path, testEvent.Ref).Return(testChartFiles, nil)

	var mrs mockReleaseSyncer
	mrs.On("UpdateOrInstallFromChart", mock.Anything, mock.Anything).Return(nil)

	handler := newHandler(&md, &mrs)
	err := handler.handleMessage(context.Background(), testEvent)
	assert.NoError(t, err)
}
