package github

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type mockDownloader struct {
	mock.Mock
}

func (m *mockDownloader) BufferDirectory(ctx context.Context, repo, path, prefix, ref string) ([]*chartutil.BufferedFile, error) {
	args := m.Called(ctx, repo, path, prefix, ref)
	var ret0 []*chartutil.BufferedFile
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.([]*chartutil.BufferedFile)
	}

	return ret0, args.Error(1)
}

type mockReconciler struct {
	mock.Mock
}

func (m *mockReconciler) Reconcile(ctx context.Context, chart *chart.Chart) error {
	args := m.Called(ctx, chart)
	return args.Error(0)
}

func TestHandlerChartDownloadFails(t *testing.T) {
	errNotFound := errors.New("not found")
	testEvent := pushEvent{
		Ref:        "ref",
		Path:       "path",
		Repository: "repository",
	}

	eventBytes, err := json.Marshal(testEvent)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	var md mockDownloader
	md.On("BufferDirectory", mock.Anything, testEvent.Repository, testEvent.Path, testEvent.Path, testEvent.Ref).Return(nil, errNotFound)

	listener := &RegistryChartListener{downloader: &md}
	handler := listener.handler(nil)

	err = handler(context.Background(), string(eventBytes))
	assert.Error(t, errNotFound, err)
}

func TestHandlerCallsReconciler(t *testing.T) {
	testEvent := pushEvent{
		Ref:        "ref",
		Path:       "path",
		Repository: "repository",
	}

	eventBytes, err := json.Marshal(testEvent)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	testChartFiles := []*chartutil.BufferedFile{
		{
			Name: "Chart.yaml",
			Data: []byte(`
apiVersion: v1
name: foo
description: This is a foo chart.
`),
		},
	}

	var md mockDownloader
	md.On("BufferDirectory", mock.Anything, testEvent.Repository, testEvent.Path, testEvent.Path, testEvent.Ref).Return(testChartFiles, nil)

	var mr mockReconciler
	mr.On("Reconcile", mock.Anything, mock.Anything).Return(nil)

	listener := &RegistryChartListener{downloader: &md}
	handler := listener.handler(&mr)

	err = handler(context.Background(), string(eventBytes))
	assert.NoError(t, err)

	mr.AssertCalled(t, "Reconcile", mock.Anything, mock.Anything)
}
