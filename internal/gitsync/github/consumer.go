package github

import (
	"context"
	"time"

	"ship-it/internal/gitsync/middleware"

	"github.com/Wattpad/sqsconsumer"
	sqsmiddleware "github.com/Wattpad/sqsconsumer/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"k8s.io/helm/pkg/helm"
)

type Config struct {
	GithubOrg   string
	GithubToken string
	HelmTimeout time.Duration
	Namespace   string
	Queue       string
	ReleaseName string
}

func New(l log.Logger, h metrics.Histogram, sqs sqsconsumer.SQSAPI, cfg Config) (*sqsconsumer.Consumer, error) {
	svc, err := sqsconsumer.NewSQSService(cfg.Queue, sqs)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: cfg.GithubToken,
		},
	))

	githubClient := github.NewClient(oauthClient)

	handler := newHandler(
		newGithubDownloader(githubClient, cfg.GithubOrg),
		newReleaseSyncer(helm.NewClient(), cfg.Namespace, cfg.ReleaseName, cfg.HelmTimeout),
	)

	stack := sqsmiddleware.ApplyDecoratorsToHandler(
		handler.HandleMessage,
		middleware.Timer(h.With("worker", "github", "queue", cfg.Queue)),
		middleware.Logger(l),
	)

	return sqsconsumer.NewConsumer(svc, stack), nil
}
