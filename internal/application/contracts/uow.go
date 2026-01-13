package contracts

import "context"

type UnitOfWork interface {
	WithinTransaction(
		ctx context.Context,
		fn func(exec context.Context) error,
	) error
}