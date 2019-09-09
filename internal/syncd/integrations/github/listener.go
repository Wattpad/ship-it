package github

import (
	"context"

	"ship-it/internal/syncd"
	"ship-it/internal/syncd/middleware"

	"github.com/Wattpad/sqsconsumer"
	sqsmiddleware "github.com/Wattpad/sqsconsumer/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
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
	handler := NewRegistryChartEventHandler(l.downloader, r)
	stack := sqsmiddleware.ApplyDecoratorsToHandler(
		handler.HandleMessage,
		middleware.Timer(l.timer),
		middleware.Logger(l.logger),
	)
	return sqsconsumer.NewConsumer(l.service, stack).Run(ctx)
}
