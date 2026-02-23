package messaging

import (
	"context"
	"github.com/victorotene80/bus_lib/message"
	libbus "github.com/victorotene80/bus_lib"
)

type Command = message.Command

type CommandHandler[TCommand Command, TResult any] interface{
	Handle(ctx context.Context, cmd TCommand) (TResult, error)
}

type Middleware = libbus.MiddlewareAny
type HandlerFunc = libbus.HandlerFuncAny

/*
type Command interface{}

type CommandHandler[TCommand Command, TResult any] interface {
	Handle(ctx context.Context, cmd TCommand) (TResult, error)
}

type Middleware func(next HandlerFunc) HandlerFunc

type HandlerFunc func(ctx context.Context, cmd Command) (any, error)
*/