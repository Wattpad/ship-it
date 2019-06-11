package ecrconsumer

import (
	"context"
	"fmt"

	"github.com/Wattpad/sqsconsumer"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSConfig struct {
	Name   string
	Region string
	Svc    *sqs.SQS
}

func NewSQSConfig(name string, region string) *SQSConfig {
	conf := &aws.Config{
		Region: &region,
	}
	svc := sqs.New(session.New(conf))

	return &SQSConfig{
		Name:   name,
		Region: region,
		Svc:    svc,
	}
}

func (s *SQSConfig) NewSQSConsumer() (*sqsconsumer.Consumer, error) {
	// Create SQS service
	service, err := sqsconsumer.SQSObjectForQueue(s.Name, s.Svc)
	if err != nil {
		return nil, err
	}

	// Create and return SQS consumer
	consumer := sqsconsumer.NewConsumer(service, processMessage)
	return consumer, nil
}

func processMessage(ctx context.Context, msg string) error {
	fmt.Println(msg)
	return nil
}
