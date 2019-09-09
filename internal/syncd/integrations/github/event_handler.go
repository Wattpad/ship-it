package github

import (
	"context"
	"encoding/json"

	"ship-it/internal/syncd"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
)

type DirectoryBufferer interface {
	BufferDirectory(ctx context.Context, repo, path, ref string) ([]*chartutil.BufferedFile, error)
}

type RegistryChartEventHandler struct {
	downloader DirectoryBufferer
	reconciler syncd.RegistryChartReconciler
}

func NewRegistryChartEventHandler(downloader DirectoryBufferer, reconciler syncd.RegistryChartReconciler) *RegistryChartEventHandler {
	return &RegistryChartEventHandler{
		downloader: downloader,
		reconciler: reconciler,
	}
}

type pushEvent struct {
	Ref        string `json:"ref"`
	Path       string `json:"path"`
	Repository string `json:"repository"`
}

func (h *RegistryChartEventHandler) HandleMessage(ctx context.Context, msg string) error {
	var event pushEvent
	if err := json.Unmarshal([]byte(msg), &event); err != nil {
		return errors.Wrap(err, "failed to unmarshal github push event")
	}

	chartFiles, err := h.downloader.BufferDirectory(ctx, event.Repository, event.Path, event.Ref)
	if err != nil {
		return errors.Wrap(err, "failed to download chart directory")
	}

	chart, err := chartutil.LoadFiles(chartFiles)
	if err != nil {
		return errors.Wrap(err, "failed to load chart files")
	}

	return h.reconciler.Reconcile(ctx, chart)
}
