package main

import (
	"context"
	"fmt"
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

// also create DD object here + histogram
// Make NewSQSConsumer independent function with extra args
// Create context stuff in main and run sqsconsumer.Run in go routine in main.go
// Create SVC object in main and take interface for SQSAPI as args in NewSQSConfig instead of constructing the SVC in the function
// Data Dog event for successful message processing
// raname package directory to ecrconsumer and merge in the git chart editor code

func main() {
	http.ListenAndServe(":80", api.New())

	// Setup a cancel on interrupt context
	fetchCtx, cancelFetch := context.WithCancel(context.Background())
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-term
		log.Println("Starting graceful shutdown")
		cancelFetch()
	}()

	// Load Environment Variables
	envConf, err := ecrconsumer.ConfigFromEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	sqsConf, err := ecrconsumer.NewSQSConfig("kube-deploy-events.fifo", envConf)
	if err != nil {
		log.Println(err)
		return
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
		fmt.Println(err)
		return
	}

	svc := sqs.New(s)

	// ECR SQS Consumer Setup
	consumer, err := ecrconsumer.NewSQSConsumer(logger, dd, hist, sqsConf.Name, svc)
	if err != nil {
		log.Println(err)
		return
	}

	go consumer.Run(fetchCtx)
}
