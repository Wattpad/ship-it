package github

import (
	"context"
	"encoding/json"

	"ship-it/internal/syncd"
	"ship-it/internal/syncd/middleware"

	"github.com/Wattpad/sqsconsumer"
	sqsmiddleware "github.com/Wattpad/sqsconsumer/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
)

type githubDownloader interface {
	BufferDirectory(ctx context.Context, repo, path, ref string) ([]*chartutil.BufferedFile, error)
}

type RegistryChartListener struct {
	downloader githubDownloader
	logger     log.Logger
	service    *sqsconsumer.SQSService
	timer      metrics.Histogram
}

func NewListener(l log.Logger, h metrics.Histogram, org string, r RepositoriesService, queue string, sqs sqsconsumer.SQSAPI) (*RegistryChartListener, error) {
	svc, err := sqsconsumer.NewSQSService(queue, sqs)
	if err != nil {
		return nil, err
	}

	return &RegistryChartListener{
		downloader: newDownloader(r, org),
		logger:     log.With(l, "worker", "github"),
		service:    svc,
		timer:      h.With("worker", "github"),
	}, nil
}

func (l *RegistryChartListener) Listen(ctx context.Context, r syncd.RegistryChartReconciler) error {
	stack := sqsmiddleware.ApplyDecoratorsToHandler(
		l.handler(r),
		middleware.Timer(l.timer),
		middleware.Logger(l.logger),
	)
	return sqsconsumer.NewConsumer(l.service, stack).Run(ctx)
}

type pushEvent struct {
	Ref        string `json:"ref"`
	Path       string `json:"path"`
	Repository string `json:"repository"`
}

func (l *RegistryChartListener) handler(r syncd.RegistryChartReconciler) sqsconsumer.MessageHandlerFunc {
	return func(ctx context.Context, msg string) error {
		var event pushEvent
		if err := json.Unmarshal([]byte(msg), &event); err != nil {
			return errors.Wrap(err, "failed to unmarshal github push event")
		}

		chartFiles, err := l.downloader.BufferDirectory(ctx, event.Repository, event.Path, event.Ref)
		if err != nil {
			return errors.Wrap(err, "failed to download chart directory")
		}

		chart, err := chartutil.LoadFiles(chartFiles)
		if err != nil {
			return errors.Wrap(err, "failed to load chart files")
		}

		return r.Reconcile(ctx, chart)
	}
}
