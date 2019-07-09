package github

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type pushEvent struct {
	Ref        string
	Path       string
	Repository string
}

type handler struct {
	downloader GithubDownloader
	release    ReleaseSyncer
}

type GithubDownloader interface {
	BufferDirectory(ctx context.Context, repo, path, ref string) ([]*chartutil.BufferedFile, error)
}

type ReleaseSyncer interface {
	UpdateOrInstallFromChart(ctx context.Context, chart *chart.Chart) error
}

func newHandler(d GithubDownloader, r ReleaseSyncer) *handler {
	return &handler{
		downloader: d,
		release:    r,
	}
}

func (h *handler) HandleMessage(ctx context.Context, msg string) error {
	var event pushEvent
	if err := json.Unmarshal([]byte(msg), &event); err != nil {
		return errors.Wrap(err, "failed to unmarshal github push event")
	}

	return h.handleMessage(ctx, event)
}

func (h *handler) handleMessage(ctx context.Context, event pushEvent) error {
	chartFiles, err := h.downloader.BufferDirectory(ctx, event.Repository, event.Path, event.Ref)
	if err != nil {
		return errors.Wrap(err, "failed to download chart contents")
	}

	chart, err := chartutil.LoadFiles(chartFiles)
	if err != nil {
		return errors.Wrap(err, "failed to load chart files")
	}

	return h.release.UpdateOrInstallFromChart(ctx, chart)
}
