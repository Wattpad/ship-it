package middleware

import (
	"context"

	"github.com/Wattpad/sqsconsumer"
	"github.com/Wattpad/sqsconsumer/middleware"
	"github.com/go-kit/kit/log"
)

func Logger(logger log.Logger) middleware.MessageHandlerDecorator {
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
