package messaging

import "context"

type Command interface{}

type CommandHandler[TCommand Command, TResult any] interface {
	Handle(ctx context.Context, cmd TCommand) (TResult, error)
}

type Middleware func(next HandlerFunc) HandlerFunc

type HandlerFunc func(ctx context.Context, cmd Command) (any, error)
