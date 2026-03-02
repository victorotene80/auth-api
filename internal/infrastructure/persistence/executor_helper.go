package persistence

import (
	"context"
	"database/sql"

)

type Executor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func ChooseExecutor(ctx context.Context, db *sql.DB) Executor {
	if tx, err := GetTx(ctx); err == nil && tx != nil {
		return tx
	}
	return db
}