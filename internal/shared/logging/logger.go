package logging

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey struct{}

type LoggerProvider struct {
	logger *zap.Logger
}

func NewLoggerProvider() (*LoggerProvider, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &LoggerProvider{logger: l}, nil
}

func (p *LoggerProvider) Logger() *zap.Logger {
	return p.logger
}

func (p *LoggerProvider) Sync() {
	_ = p.logger.Sync()
}

func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok && l != nil {
		return l
	}

	return zap.NewNop()
}
