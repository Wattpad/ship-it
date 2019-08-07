package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ship-it/internal/syncd"
	"ship-it/internal/syncd/config"
	"ship-it/internal/syncd/integrations/ecr"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/dogstatsd"
	"k8s.io/apimachinery/pkg/types"
)

type NOPIndexer struct{}

func (n NOPIndexer) Lookup(repo string) ([]types.NamespacedName, error) {
	return nil, fmt.Errorf("error not implemented")
}

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

	imageListener, err := ecr.NewListener(logger, dd.NewTiming("syncd.time", 1.0), cfg.EcrQueue, sqs.New(session.Must(session.NewSession())))
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	gitClient := ecr.NewGitHub(ctx, cfg.GithubToken, cfg.GithubOrg, cfg.OperationsRepoName, cfg.ReleaseBranch, cfg.RegistryChartPath)
	imageReconciler := ecr.NewReconciler(gitClient, NOPIndexer{})

	// TODO: Allow configurable image/chart sync implementations. For now
	// we'll just use our specific ecr+sqs/github+sqs implmentations.
	syncd := syncd.New(nil, nil, imageListener, imageReconciler)
	if err := syncd.Run(ctx); err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}
}
