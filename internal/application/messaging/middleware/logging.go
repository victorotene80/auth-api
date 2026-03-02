package middleware

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/victorotene80/authentication_api/internal/application/messaging"
)

// Logging returns a middleware that logs command execution with zap.
func Logging(logger *zap.Logger) messaging.Middleware {
	if logger == nil {
		logger = zap.NewNop()
	}

	return func(next messaging.HandlerFunc) messaging.HandlerFunc {
		return func(ctx context.Context, cmd messaging.Command) (any, error) {
			start := time.Now()
			cmdName := fmt.Sprintf("%T", cmd)

			logger.Info("command started",
				zap.String("command", cmdName),
			)

			res, err := next(ctx, cmd)

			elapsed := time.Since(start)

			if err != nil {
				logger.Error("command failed",
					zap.String("command", cmdName),
					zap.Duration("duration", elapsed),
					zap.Error(err),
				)
			} else {
				logger.Info("command finished",
					zap.String("command", cmdName),
					zap.Duration("duration", elapsed),
				)
			}

			return res, err
		}
	}
}

func AttachLogging(bus *messaging.CommandBus, logger *zap.Logger) {
	if bus == nil {
		return
	}
	bus.Use(Logging(logger))
}