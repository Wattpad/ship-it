package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ship-it/internal/syncd"
	"ship-it/internal/syncd/config"

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

	// TODO: Allow configurable image/chart sync implementations. For now
	// we'll just use our specific ecr+sqs/github+sqs implmentations.
	syncd := syncd.New(nil, nil, nil, nil)
	if err := syncd.Run(ctx); err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}
}
