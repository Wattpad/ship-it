package ecr

import (
	"context"
	"log"

	"ship-it/internal/syncd"
	"ship-it/internal/syncd/middleware"

	"github.com/Wattpad/sqsconsumer"
	sqsmiddleware "github.com/Wattpad/sqsconsumer/middleware"
	"k8s.io/component-base/metrics"
)

type ImageListener struct {
	logger  log.Logger
	service *sqsconsumer.SQSService
	timer   metrics.Histogram
}

func NewListener(l log.Logger, h metrics.Histogram, queue string, sqs sqsconsumer.SQSAPI) (*ImageListener, error) {
	svc, err := sqsconsumer.NewSQSService(queue, sqs)
	if err != nil {
		return nil, err
	}

	return &ImageListener{
		logger:  log.With(l, "worker", "ecr"),
		service: svc,
		timer:   h.With("worker", "ecr", "queue", queue),
	}, nil
}

func (l *ImageListener) Listen(ctx context.Context, r syncd.ImageReconciler) error {
	stack := sqsmiddleware.ApplyDecoratorsToHandler(
		l.handler(r),
		middleware.Timer(l.timer),
		middleware.Logger(l.logger),
	)
	return sqsconsumer.NewConsumer(l.service, stack).Run(ctx)
}

func (l *ImageListener) handler(r syncd.ImageReconciler) sqsconsumer.MessageHandlerFunc {
	return nil
}
