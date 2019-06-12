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
	conf, err := ecrconsumer.NewSQSConfig("kube-deploy-events.fifo", "us-east-1")
	if err != nil {
		log.Println(err)
		return
	}

	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))

	consumer, err := conf.NewSQSConsumer(logger)
	if err != nil {
		log.Println(err)
		return
	}

	ecrconsumer.RunConsumer(consumer)
}
