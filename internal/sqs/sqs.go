package ecrconsumer

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

func NewSQSConfig(name string, region string) (*SQSConfig, error) {
	conf := &aws.Config{
		Region: &region,
	}

	s, err := session.NewSession(conf)
	if err != nil {
		return nil, err
	}

	svc := sqs.New(s)

	return &SQSConfig{
		Name:   name,
		Region: region,
		Svc:    svc,
	}, nil
}

func SQSObjectForQueue(name string, svc *sqs.SQS) (*sqsconsumer.SQSService, error) {
	return nil, nil
}

func (s *SQSConfig) NewSQSConsumer() (*sqsconsumer.Consumer, error) {
	// Create SQS service
	service, err := SQSObjectForQueue(s.Name, s.Svc)
	if err != nil {
		return nil, err
	}

	// Create and return SQS consumer
	consumer := sqsconsumer.NewConsumer(service, processMessage)
	consumer.SetLogger(log.Printf)
	return consumer, nil
}

// Runs a preconfigured Consumer
func RunConsumer(c *sqsconsumer.Consumer) {

	numFetchers := 1

	// set up a context which will gracefully cancel the worker on interrupt
	fetchCtx, cancelFetch := context.WithCancel(context.Background())
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-term
		log.Println("Starting graceful shutdown")
		cancelFetch()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(numFetchers)
	for i := 0; i < numFetchers; i++ {
		go func() {
			// start running the consumer with a context that will be cancelled when a graceful shutdown is requested
			c.Run(fetchCtx)

			wg.Done()
		}()
	}

	// wait for all the consumers to exit cleanly
	wg.Wait()
	log.Println("Shutdown complete")
}

func processMessage(ctx context.Context, msg string) error {
	fmt.Println(msg)
	return nil
}
