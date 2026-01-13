package persistence

import (
	"context"
	"database/sql"
	"errors"
)

type txKey struct{}

type SqlUnitOfWork struct {
	db *sql.DB
}

func NewSqlUnitOfWork(db *sql.DB) *SqlUnitOfWork {
	return &SqlUnitOfWork{db: db}
}

func (u *SqlUnitOfWork) WithinTransaction(ctx context.Context, fn func(exec context.Context) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func GetTx(ctx context.Context) (*sql.Tx, error) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	if !ok || tx == nil {
		return nil, errors.New("no transaction found in context")
	}
	return tx, nil
}
