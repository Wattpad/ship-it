package main

import (
	"log"
	"net/http"
	"ship-it/internal/api"
	ecrconsumer "ship-it/internal/sqs"
)

func main() {
	http.ListenAndServe(":80", api.New())

	// ECR SQS Consumer Setup
	conf, err := ecrconsumer.NewSQSConfig("kube-deploy-events.fifo", "us-east-1")
	if err != nil {
		log.Println(err)
		return
	}

	consumer, err := conf.NewSQSConsumer()
	if err != nil {
		log.Println(err)
		return
	}

	ecrconsumer.RunConsumer(consumer)
}
