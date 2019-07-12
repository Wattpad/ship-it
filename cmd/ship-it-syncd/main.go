package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ship-it/internal/syncd"
	"ship-it/internal/syncd/config"
	"ship-it/internal/syncd/integrations/github"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/dogstatsd"
	gogithub "github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
	"k8s.io/helm/pkg/helm"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	logger.Log("event", "service.start")
	defer logger.Log("event", "service.stop")

	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		cancel()
	}()

	cfg, err := config.FromEnv()
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	dd := dogstatsd.New("wattpad.ship-it.", logger)
	go dd.SendLoop(time.Tick(time.Second), "udp", cfg.DataDogAddress())

	chartListener, err := initChartListener(logger, dd, cfg)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	chartReconciler := github.NewReconciler(
		helm.NewClient(),
		cfg.Namespace,
		cfg.ReleaseName,
		cfg.HelmTimeout(),
	)

	// TODO: Allow configurable image/chart sync implementations. For now
	// we'll just use our specific ecr+sqs/github+sqs implmentations.
	syncd := syncd.New(chartListener, chartReconciler, nil, nil)
	if err := syncd.Run(ctx); err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}
}

func initChartListener(l log.Logger, dd *dogstatsd.Dogstatsd, cfg *config.Config) (syncd.ChartListener, error) {
	workerTime := dd.NewTiming("worker.time", 1)

	awsSession := session.New(cfg.AWS())
	sqsClient := sqs.New(awsSession)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: cfg.GithubToken,
		},
	)

	oauthClient := oauth2.NewClient(context.Background(), ts)
	githubClient := gogithub.NewClient(oauthClient)

	return github.NewListener(
		l,
		workerTime,
		cfg.GithubOrg,
		githubClient.Repositories,
		cfg.GithubQueue,
		sqsClient,
	)
}
