package ecrconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ship-it/pkg/apis/k8s.wattpad.com/v1alpha1"

	"github.com/Wattpad/sqsconsumer"
	"github.com/Wattpad/sqsconsumer/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
)

type GitCommands interface {
	GetFile(branch string, path string) ([]byte, error)
	UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error)
}

type ImageID struct {
	Digest string `json:"imageDigest"`
}

type ImageData struct {
	RepositoryName string `json:"repositoryName"`
	ID ImageID `json:"imageId"`
}

type ResponseElements struct {
	Image ImageData `json:"image"`
}

type Detail struct {
	EventTime time.Time `json:"eventTime"`
	Response ResponseElements `json:"responseElements"`
}

type SQSMessage struct {
	Detail Detail `json:"detail"`
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
	if len(arr) == 2 {
		return arr[1]
	}
	return ""
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

		tag := parseSHA(sqsMessage.Detail.Response.Image.ID.Digest)
		newImage := makeImage(sqsMessage.Detail.Response.Image.RepositoryName, tag)

		resourceBytes, err := client.GetFile("master", "/custom-resources")
		if err != nil {
			return err
		}

		rls, err := v1alpha1.LoadRelease(resourceBytes)
		if err != nil {
			return err
		}

		updatedRls := WithImage(newImage, *rls)

		updatedBytes, err := yaml.Marshal(updatedRls)
		if err != nil {
			return nil
		}

		_, err = client.UpdateFile(fmt.Sprintf("Image Tag updated to: %s", newImage.Tag), "master", "/custom-resources", updatedBytes)
		if err != nil {
			return err
		}
		return nil
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
