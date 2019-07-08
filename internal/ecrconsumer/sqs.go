package ecrconsumer

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/Wattpad/sqsconsumer"
	"github.com/Wattpad/sqsconsumer/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/google/go-github/github"
)

type GitCommands interface {
	GetFile(branch string, path string) ([]byte, error)
	UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error)
}

// Will delete after VELO-1453 is merged as that PR already has this type included
type Image struct {
	Registry   string
	Repository string
	Tag        string
}

type SQSMessage struct {
	Detail struct {
		EventTime        time.Time `json:"eventTime"`
		ResponseElements struct {
			Image struct {
				RepositoryName string `json:"repositoryName"`
				ImageID        struct {
					ImageDigest string `json:"imageDigest"`
				}
			} `json:"image"`
		} `json:"responseElements"`
	} `json:"detail"`
}

// New returns a SQS consumer which processes ECR PushImage events by updating
// the associated service's chart values in a remote Git repository.
func New(logger log.Logger, hist metrics.Histogram, name string, svc sqsconsumer.SQSAPI, client GitCommands) (*sqsconsumer.Consumer, error) {
	service, err := sqsconsumer.NewSQSService(name, svc)
	if err != nil {
		return nil, err
	}

	track := dataDogTimeTracker(hist)
	wrappedLogger := loggerMiddleware(logger)
	handler := middleware.ApplyDecoratorsToHandler(processMessage(client), track, wrappedLogger)
	consumer := sqsconsumer.NewConsumer(service, handler)

	return consumer, nil
}

func parseSHA(digest string) string {
	arr := strings.Split(digest, ":")
	return arr[1]
}

func parseMsg(msg string) (*SQSMessage, error) {
	js := &SQSMessage{}
	err := json.Unmarshal([]byte(msg), js)
	if err != nil {
		return nil, err
	}
	return js, nil
}

func makeImage(repoName string, tag string) Image {
	return Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: repoName,
		Tag:        tag,
	}
}

func processMessage(client GitCommands) sqsconsumer.MessageHandlerFunc {
	return func(ctx context.Context, msg string) error {
		// Handle Git Commits Here
		sqsMessage, err := parseMsg(msg)
		if err != nil {
			return err
		}

		tag := parseSHA(sqsMessage.Detail.ResponseElements.Image.ImageID.ImageDigest)
		_ = makeImage(sqsMessage.Detail.ResponseElements.Image.RepositoryName, tag)

		return nil
		// Get CR Bytes
		// Unmarshal
		// Attach updated Image
		// Marshal
		// Commit new Bytes to GitHub
	}
}

func dataDogTimeTracker(hist metrics.Histogram) middleware.MessageHandlerDecorator {
	return func(fn sqsconsumer.MessageHandlerFunc) sqsconsumer.MessageHandlerFunc {
		return func(ctx context.Context, msg string) error {
			start := time.Now()

			err := fn(ctx, msg)

			var status string
			if err != nil {
				status = "failure"
			} else {
				status = "success"
			}
			hist.With("status", status).Observe(float64(time.Since(start).Seconds() * 1000))

			return err
		}
	}
}

func loggerMiddleware(logger log.Logger) middleware.MessageHandlerDecorator {
	return func(fn sqsconsumer.MessageHandlerFunc) sqsconsumer.MessageHandlerFunc {
		return func(ctx context.Context, msg string) error {
			err := fn(ctx, msg)
			if err != nil {
				logger.Log("error", err)
			}
			return err
		}
	}
}
