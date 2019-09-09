package github

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type mockDirectoryBufferer struct {
	Files []*chartutil.BufferedFile
	Err   error
}

func (m *mockDirectoryBufferer) BufferDirectory(ctx context.Context, repo, path, ref string) ([]*chartutil.BufferedFile, error) {
	return m.Files, m.Err
}

type mockChartReconciler struct {
	ReconcileCalled bool
}

func (m *mockChartReconciler) Reconcile(context.Context, *chart.Chart) error {
	m.ReconcileCalled = true
	return nil
}

func TestEventHandlerChartDownloadFails(t *testing.T) {
	errNotFound := errors.New("not found")
	msg := makePushEvent(t)
	bufferer := &mockDirectoryBufferer{Err: errNotFound}
	reconciler := &mockChartReconciler{}
	handler := NewRegistryChartEventHandler(bufferer, reconciler)

	err := handler.HandleMessage(context.Background(), msg)

	assert.Error(t, errNotFound, err)
}

func TestEventHandlerCallsReconciler(t *testing.T) {
	msg := makePushEvent(t)
	testChartFiles := makeChartFiles()
	downloader := &mockDirectoryBufferer{Files: testChartFiles}
	reconciler := &mockChartReconciler{}
	handler := NewRegistryChartEventHandler(downloader, reconciler)

	err := handler.HandleMessage(context.Background(), msg)

	assert.NoError(t, err)
	assert.True(t, reconciler.ReconcileCalled)
}

func makePushEvent(t *testing.T) string {
	testEvent := pushEvent{
		Ref:        "ref",
		Path:       "path",
		Repository: "repository",
	}
	eventBytes, err := json.Marshal(testEvent)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return string(eventBytes)
}

func makeChartFiles() []*chartutil.BufferedFile {
	return []*chartutil.BufferedFile{
		{
			Name: "Chart.yaml",
			Data: []byte(`
apiVersion: v1
name: foo
description: This is a foo chart.
`),
		},
	}
}
