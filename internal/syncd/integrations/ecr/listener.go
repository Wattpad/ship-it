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
)

type ImageListener struct {
	logger  log.Logger
	service *sqsconsumer.SQSService
	timer   metrics.Histogram
}

type SQSMessage struct {
	EventTime      time.Time
	RepositoryName string
	Tag            string
	RegistryId     string
}

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

func parseMsg(msg string) (*SQSMessage, error) {
	js := &SQSMessage{}
	err := json.Unmarshal([]byte(msg), js)
	if err != nil {
		return nil, err
	}
	return js, nil
}

func makeImage(repoName string, tag string) internal.Image {
	return internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: repoName,
		Tag:        tag,
	}
}

func validateTag(tag string) bool {
	matched, err := regexp.MatchString("^[0-9a-f]{40}$", tag)
	if err != nil {
		return false
	}
	return matched
}

func (l *ImageListener) handler(r syncd.ImageReconciler) sqsconsumer.MessageHandlerFunc {
	return func(ctx context.Context, msg string) error {
		sqsMessage, err := parseMsg(msg)
		if err != nil {
			return err
		}

		if !validateTag(sqsMessage.Tag) {
			return fmt.Errorf("Malformed Image Tag")
		}

		newImage := makeImage(sqsMessage.RepositoryName, sqsMessage.Tag)

		return r.Reconcile(ctx, &newImage)
	}
}
