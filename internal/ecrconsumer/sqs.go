package ecrconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"ship-it/pkg/apis/k8s.wattpad.com/v1alpha1"

	"github.com/Wattpad/sqsconsumer"
	"github.com/Wattpad/sqsconsumer/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/google/go-github/v26/github"
	"gopkg.in/yaml.v2"
)

type GitCommands interface {
	GetFile(branch string, path string) ([]byte, error)
	UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error)
}

type SQSMessage struct {
	EventTime      time.Time
	RepositoryName string
	Tag            string
	RegistryId     string
}

// New returns a SQS consumer which processes ECR PushImage events by updating
// the associated service's chart values in a remote Git repository.
func New(logger log.Logger, hist metrics.Histogram, name string, svc sqsconsumer.SQSAPI, client GitCommands, resourcePath string) (*sqsconsumer.Consumer, error) {
	service, err := sqsconsumer.NewSQSService(name, svc)
	if err != nil {
		return nil, err
	}

	track := dataDogTimeTracker(hist)
	wrappedLogger := loggerMiddleware(logger)
	handler := middleware.ApplyDecoratorsToHandler(processMessage(client, resourcePath), track, wrappedLogger)
	consumer := sqsconsumer.NewConsumer(service, handler)

	return consumer, nil
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

func validateTag(tag string) bool {
	matched, err := regexp.MatchString("^[0-9a-f]{40}$", tag)
	if err != nil {
		return false
	}
	return matched
}

func processMessage(client GitCommands, resourcePath string) sqsconsumer.MessageHandlerFunc {
	return func(ctx context.Context, msg string) error {
		sqsMessage, err := parseMsg(msg)
		if err != nil {
			return err
		}

		if !validateTag(sqsMessage.Tag) {
			return fmt.Errorf("Malformed Image Tag")
		}

		newImage := makeImage(sqsMessage.RepositoryName, sqsMessage.Tag)

		resourceBytes, err := client.GetFile("master", resourcePath)
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

		_, err = client.UpdateFile(fmt.Sprintf("Image Tag updated to: %s", newImage.Tag), "master", filepath.Join(resourcePath, newImage.Repository)+".yaml", updatedBytes)
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
