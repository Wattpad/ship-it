package middleware

import (
	"context"
	"time"

	"github.com/Wattpad/sqsconsumer"
	"github.com/Wattpad/sqsconsumer/middleware"
	"github.com/go-kit/kit/metrics"
)

func Timer(hist metrics.Histogram) middleware.MessageHandlerDecorator {
	return func(fn sqsconsumer.MessageHandlerFunc) sqsconsumer.MessageHandlerFunc {
		return func(ctx context.Context, msg string) error {
			start := time.Now()

			err := fn(ctx, msg)

			status := "failure"
			if err == nil {
				status = "success"
			}
			hist.With("status", status).Observe(millisecondsSince(start))

			return err
		}
	}
}

func millisecondsSince(t time.Time) float64 {
	return time.Since(t).Seconds() * 1000
}
