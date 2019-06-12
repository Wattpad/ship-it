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

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/dogstatsd"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	// Setup a cancel on interrupt context
	ctx, cancel := context.WithCancel(context.Background())
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-term
		log.Println("Starting graceful shutdown")
		cancel()
	}()

	// Load Environment Variables
	envConf, err := ecrconsumer.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	sqsConf, err := ecrconsumer.NewSQSConfig("kube-deploy-events.fifo", envConf)
	if err != nil {
		log.Fatal(err)
	}

	// DataDog Setup
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	dd := dogstatsd.New("wattpad.", logger)
	go dd.SendLoop(time.NewTicker(time.Second).C, "udp", envConf.DataDogAddress())
	hist := dd.NewTiming("worker.time", 1.0).With("worker", "ship-it-worker", "queue", sqsConf.Name)

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
	consumer, err := ecrconsumer.NewSQSConsumer(logger, dd, hist, sqsConf.Name, svc)
	if err != nil {
		log.Println(err)
		return
	}

	go consumer.Run(ctx)
	http.ListenAndServe(":80", api.New())
}
