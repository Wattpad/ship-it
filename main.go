package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"ship-it/internal/api"
	"ship-it/internal/ecrconsumer"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
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

	envConf, err := FromEnv()
	if err != nil {
		logger.Log("Error getting environment variables %s", err)
		os.Exit(1)
	}

	dd := dogstatsd.New("wattpad.ship-it.", logger)
	go dd.SendLoop(time.Tick(time.Second), "udp", envConf.DataDogAddress())
	hist := dd.NewTiming("worker.time", 1.0).With("worker", "worker:ecrconsumer", "queue", envConf.QueueName)

	s, err := session.NewSession(&aws.Config{
		Region: &envConf.AWSRegion,
	})
	if err != nil {
		logger.Log("Error opening AWS session %s", err)
		os.Exit(1)
	}

	consumer, err := ecrconsumer.NewSQSConsumer(logger, hist, envConf.QueueName, sqs.New(s))
	if err != nil {
		logger.Log("Error creating SQS consumer %s", err)
		os.Exit(1)
	}

	go consumer.Run(ctx)

	go http.ListenAndServe(":"+envConf.ServicePort, api.New())

	<-ctx.Done()
	logger.Log("event", "service.exit", "error", ctx.Err())
}
