package ecr

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"ship-it/internal"
	"ship-it/internal/syncd"
	"ship-it/internal/syncd/middleware"

	"github.com/Wattpad/sqsconsumer"
	sqsmiddleware "github.com/Wattpad/sqsconsumer/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/pkg/errors"
)

type ImageListener struct {
	logger  log.Logger
	service *sqsconsumer.SQSService
	timer   metrics.Histogram
}

type ecrPushEvent struct {
	EventTime      time.Time `json:"eventTime"`
	RepositoryName string    `json:"repositoryName"`
	Tag            string    `json:"tag"`
	RegistryId     string    `json:"registryId"`
}

var validImageTagRegex = regexp.MustCompile("^[0-9a-f]{40}$")

func NewListener(l log.Logger, h metrics.Histogram, queue string, sqs sqsconsumer.SQSAPI) (*ImageListener, error) {
	svc, err := sqsconsumer.NewSQSService(queue, sqs)
	if err != nil {
		return nil, err
	}

	return &ImageListener{
		logger:  log.With(l, "worker", "ecr"),
		service: svc,
		timer:   h.With("worker", "ecr"),
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

func makeImage(repoName string, tag string, registryId string) internal.Image {
	return internal.Image{
		Registry:   registryId + ".dkr.ecr.us-east-1.amazonaws.com",
		Repository: repoName,
		Tag:        tag,
	}
}

func (l *ImageListener) handler(r syncd.ImageReconciler) sqsconsumer.MessageHandlerFunc {
	return func(ctx context.Context, msg string) error {
		var ecrEvent ecrPushEvent
		err := json.Unmarshal([]byte(msg), &ecrEvent)
		if err != nil {
			return errors.Wrap(err, "failure to parse ecr push event")
		}

		if !validImageTagRegex.MatchString(ecrEvent.Tag) {
			return fmt.Errorf("malformed image tag")
		}

		newImage := makeImage(ecrEvent.RepositoryName, ecrEvent.Tag, ecrEvent.RegistryId)

		return r.Reconcile(ctx, &newImage)
	}
}
