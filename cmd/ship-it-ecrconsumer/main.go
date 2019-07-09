package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ship-it/internal/ecrconsumer"
	"ship-it/internal/ecrconsumer/config"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/dogstatsd"
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

	s, err := session.NewSession(cfg.AWS())
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	gitClient := ecrconsumer.NewGitHub(cfg.GithubToken, context.Background(), cfg.GithubOrg, "miranda")
	consumer, err := ecrconsumer.New(logger, dd.NewTiming("worker.time", 1.0).With("worker", "ecrconsumer", "queue", cfg.QueueName), cfg.QueueName, sqs.New(s), gitClient, cfg.ResourcePath)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	consumer.Run(ctx)
}
