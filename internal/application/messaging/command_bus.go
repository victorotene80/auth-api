package messaging

import (
	"context"
	"errors"

	libbus "github.com/victorotene80/bus_lib/bus"
	"github.com/victorotene80/bus_lib/message"

	"github.com/victorotene80/authentication_api/internal/application"
)

type Command = message.Command

type CommandHandler[TCommand Command, TResult any] interface{
	Handle(ctx context.Context, cmd TCommand) (TResult, error)
}

type Middleware = libbus.MiddlewareAny
type HandlerFunc = libbus.HandlerFuncAny

type CommandBus struct {
	bus *libbus.Bus
	builder *libbus.Builder
}


func NewCommandBus() *CommandBus {
	return &CommandBus{
		builder: libbus.NewBuilder(),
	}
}

func (cb *CommandBus) buildOnce() error{
	if cb.bus != nil {
		return nil
	}

	b, err := cb.builder.Build()
	if err != nil {
		return err
	}

	cb.bus = b
	return nil
}

func Register[TCommand Command, TResult any](
	cb *CommandBus,
	handler CommandHandler[TCommand, TResult],
) error{
	return libbus.RegisterCommand[TCommand, TResult](
		cb.builder,
		func (ctx context.Context, cmd TCommand) (TResult, error){
			return handler.Handle(ctx, cmd)
		}
	)
}

func MustRegister[TCommand Command, TResult any](
	cb *CommandBus,
	handler CommandHandler[TCommand, TResult],
){
	if err := Register(cb, handler); err != nil {
		panic(err)
	}
}

func Execute[TCommand Command, TResult any](
	cb *CommandBus,
	ctx context.Context,
	cmd TCommand,
)(TResult, error){
	var zero TResult

	if any(cmd) == nil{
		return zero, application.ErrNilCommand
	}

	if err := cb.buildOnce(); err != nil {
		return zero, err
	}

	var out TResult
	if err := libbus.DispatchCommandInto[TCommand, TResult](cb.bus, ctx, cmd, &out); err != nil {
		if errors.Is(err, libbus.ErrHandlerNotFound) {
			return zero, application.ErrHandlerNotFound
		}
		if errors.Is(err, libbus.ErrNilMessage) {
			return zero, application.ErrNilCommand
		}
		if errors.Is(err, libbus.ErrResponseTypeMismatch) {
			return zero, application.ErrInvalidResult
		}

		return zero, err
	}
	return out, nil
}

func (cb *CommandBus) Use(mw Middleware){
	if err := cb.builder.Use(mw); err != nil {
		panic(err)
	}
}
/*import (
	"context"
	"fmt"
	"sync"

	"github.com/victorotene80/authentication_api/internal/application"
)

type CommandBus struct {
	mu         sync.RWMutex
	handlers   map[string]any
	middleware []Middleware
}

func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers:   make(map[string]any),
		middleware: []Middleware{},
	}
}

func Register[TCommand Command, TResult any](
	bus *CommandBus,
	handler CommandHandler[TCommand, TResult],
) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	key := getTypeKey[TCommand]()
	if _, exists := bus.handlers[key]; exists {
		return application.ErrHandlerExists
	}

	bus.handlers[key] = handler
	return nil
}

func MustRegister[TCommand Command, TResult any](
	bus *CommandBus,
	handler CommandHandler[TCommand, TResult],
) {
	if err := Register(bus, handler); err != nil {
		panic(err)
	}
}

func Execute[TCommand Command | *any, TResult any](
	bus *CommandBus,
	ctx context.Context,
	cmd TCommand,
) (TResult, error) {
	var zero TResult
	if any(cmd) == nil {
		return zero, application.ErrNilCommand
	}

	key := getTypeKey[TCommand]()

	bus.mu.RLock()
	rawHandler, exists := bus.handlers[key]
	bus.mu.RUnlock()

	if !exists {
		return zero, application.ErrHandlerNotFound
	}

	handler, ok := rawHandler.(CommandHandler[TCommand, TResult])
	if !ok {
		return zero, application.ErrHandlerNotFound
	}

	finalHandler := func(ctx context.Context, c Command) (any, error) {
		return handler.Handle(ctx, c.(TCommand))
	}

	// apply middleware
	for i := len(bus.middleware) - 1; i >= 0; i-- {
		finalHandler = bus.middleware[i](finalHandler)
	}

	result, err := finalHandler(ctx, cmd)
	if err != nil {
		return zero, err
	}

	typedResult, ok := result.(TResult)
	if !ok {
		return zero, application.ErrInvalidResult
	}

	return typedResult, nil
}

func (bus *CommandBus) Use(mw Middleware) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.middleware = append(bus.middleware, mw)
}

func getTypeKey[TCommand Command]() string {
	var t TCommand
	return fmt.Sprintf("%T", t)
}*/
