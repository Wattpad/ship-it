package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ship-it/internal/api"

	"github.com/go-kit/kit/log"
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

	srv := http.Server{
		Addr:    ":80",
		Handler: api.New(),
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
