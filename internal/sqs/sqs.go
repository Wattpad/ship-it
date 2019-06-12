package ecrconsumer

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Wattpad/sqsconsumer"
	"github.com/Wattpad/sqsconsumer/middleware"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	// Metric Imports
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/dogstatsd"
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

func (s *SQSConfig) NewSQSConsumer(logger kitlog.Logger) (*sqsconsumer.Consumer, error) {
	// Create SQS service
	service, err := sqsconsumer.NewSQSService(s.Name, s.Svc)
	if err != nil {
		return nil, err
	}

	// Initialize Data Dog

	address := net.JoinHostPort("", "8125") // Add DD info to SQS config

	dd := dogstatsd.New("wattpad.", logger)
	go dd.SendLoop(time.NewTicker(time.Second).C, "udp", address)

	// Configure Middleware
	hist := dd.NewTiming("worker.time", 1.0).With("worker", "ship-it-worker", "queue", s.Name)
	track := dataDogTimeTracker(hist)
	wrappedLogger := loggerMiddleware(logger)
	handler := middleware.ApplyDecoratorsToHandler(processMessage, track, wrappedLogger)

	// Create and return SQS consumer
	consumer := sqsconsumer.NewConsumer(service, handler)

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

func handleCancel(ctx context.Context, msg string) error {
	select {
	case <-ctx.Done():
		log.Println("Context done so aborting processing message:", msg)
		return ctx.Err()
	default:
		return nil
	}
}

func processMessage(ctx context.Context, msg string) error {

	err := handleCancel(ctx, msg)
	if err != nil {
		return err
	}

	fmt.Println(msg)
	return nil
}

func dataDogTimeTracker(hist metrics.Histogram) middleware.MessageHandlerDecorator {
	return func(fn sqsconsumer.MessageHandlerFunc) sqsconsumer.MessageHandlerFunc {
		return func(ctx context.Context, msg string) error {
			start := time.Now()

			err := fn(ctx, msg)

			var status string
			if err != nil {
				status = "failure"
			} else {
				status = "success"
			}
			hist.With("status", status).Observe(float64(time.Since(start).Seconds() * 1000))

			return err
		}
	}
}

func loggerMiddleware(logger kitlog.Logger) middleware.MessageHandlerDecorator {
	return func(fn sqsconsumer.MessageHandlerFunc) sqsconsumer.MessageHandlerFunc {
		return func(ctx context.Context, msg string) error {
			err := fn(ctx, msg)
			if err != nil {
				logger.Log("error", err)
			}
			return err
		}
	}
}
