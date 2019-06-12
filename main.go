package main

import (
	"context"
	"log"
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
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/dogstatsd"
)

func main() {
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)

	logger.Log("event", "service.start")
	defer logger.Log("event", "service.stop")

	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		cancel()
	}()

	// Load Environment Variables
	envConf, err := FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	dd := dogstatsd.New("wattpad.ship-it.", logger)
	go dd.SendLoop(time.NewTicker(time.Second).C, "udp", envConf.DataDogAddress())
	hist := dd.NewTiming("consumer.time", 1.0).With("consumer", "ship-it-consumer", "queue", envConf.QueueName)

	// AWS Setup
	conf := &aws.Config{
		Region: &envConf.AWSRegion,
	}

	s, err := session.NewSession(conf)
	if err != nil {
		log.Fatal(err)
	}

	svc := sqs.New(s)

	// ECR SQS Consumer Setup
	consumer, err := ecrconsumer.NewSQSConsumer(logger, dd, hist, envConf.QueueName, svc)
	if err != nil {
		log.Println(err)
		return
	}

	go consumer.Run(ctx)

	go http.ListenAndServe(":"+envConf.ServicePort, api.New())

	<-ctx.Done()
	logger.Log("event", "service.exit", "error", ctx.Err())
}
