package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"ship-it/internal/api"
	"syscall"

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

	go http.ListenAndServe(":80", api.New())

	<-ctx.Done()
	logger.Log("event", "service.exit", "error", ctx.Err())
}
