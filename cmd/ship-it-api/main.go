package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ship-it/internal/api"
	"ship-it/internal/api/config"
	"ship-it/internal/api/integrations/k8s"
	"ship-it/internal/api/service"

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

	k8s, err := k8s.New()
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	svc := service.New(k8s)
	srv := http.Server{
		Addr:    ":" + cfg.ServicePort,
		Handler: api.New(svc, dd.NewTiming("api.time", 1.0)),
	}

	exit := make(chan error)
	go func() { exit <- srv.ListenAndServe() }()

	select {
	case <-ctx.Done():
		logger.Log("event", "service.exit", "error", ctx.Err())
	case err := <-exit:
		cancel()
		logger.Log("event", "service.exit", "error", err)
	}

	srv.Shutdown(context.Background())
}
