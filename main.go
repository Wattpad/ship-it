package main

import (
	"log"
	"net/http"
	"os"
	"ship-it/internal/api"
	ecrconsumer "ship-it/internal/sqs"

	kitlog "github.com/go-kit/kit/log"
)

func main() {
	http.ListenAndServe(":80", api.New())
	// ECR SQS Consumer Setup
	conf, err := ecrconsumer.NewSQSConfig("kube-deploy-events.fifo")
	if err != nil {
		log.Println(err)
		return
	}

	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	// also create DD object here + histogram
	// Make NewSQSConsumer independent function with extra args
	// Create context stuff in main and run sqsconsumer.Run in go routine in main.go
	// Create SVC object in main and take interface for SQSAPI as args in NewSQSConfig instead of constructing the SVC in the function
	// Data Dog event for successful message processing

	consumer, err := conf.NewSQSConsumer(logger)
	if err != nil {
		log.Println(err)
		return
	}

	ecrconsumer.RunConsumer(consumer)
}
