// internal/application/messaging/middleware/logging.go
package middleware

import (
	"context"
	"log"

	"github.com/victorotene80/authentication_api/internal/application/messaging"
)

//zap log

func LoggingMiddleware(next messaging.HandlerFunc) messaging.HandlerFunc {
	return func(ctx context.Context, msg any) (any, error) {
		if cmd, ok := msg.(messaging.Command); ok {
			log.Printf("[command] name=%s version=%d\n", cmd.Name(), cmd.Version())
		} else {
			log.Printf("[command] unknown type %T\n", msg)
		}

		res, err := next(ctx, msg)
		if err != nil {
			log.Printf("[command] ERROR type=%T err=%v\n", msg, err)
		}

		return res, err
	}
}