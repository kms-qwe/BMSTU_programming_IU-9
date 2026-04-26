package utils

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txContextKey struct{}

type TxHandler func(ctx context.Context) error

type IExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type IExecutorProvider interface {
	GetExecutor(ctx context.Context) IExecutor
}

type TxManager interface {
	IExecutorProvider
	ReadCommitted(ctx context.Context, fn TxHandler) error
	RepeatableRead(ctx context.Context, fn TxHandler) error
	Serializable(ctx context.Context, fn TxHandler) error
}

type PGXTxManager struct {
	db *pgxpool.Pool
}

func NewTxManager(db *pgxpool.Pool) *PGXTxManager {
	return &PGXTxManager{db: db}
}

func NewTxContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txContextKey{}).(pgx.Tx)
	return tx, ok
}

func (m *PGXTxManager) GetExecutor(ctx context.Context) IExecutor {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}

	return m.db
}

func (m *PGXTxManager) ReadCommitted(ctx context.Context, fn TxHandler) error {
	return m.withinTxOpts(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, fn)
}

func (m *PGXTxManager) RepeatableRead(ctx context.Context, fn TxHandler) error {
	return m.withinTxOpts(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}, fn)
}

func (m *PGXTxManager) Serializable(ctx context.Context, fn TxHandler) error {
	return m.withinTxOpts(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}, fn)
}

func (m *PGXTxManager) withinTxOpts(ctx context.Context, opts pgx.TxOptions, fn TxHandler) (err error) {
	if _, ok := TxFromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := m.db.BeginTx(ctx, opts)
	if err != nil {
		return InternalError("failed to start transaction", fmt.Errorf("can't begin transaction: %w", err))
	}

	ctx = NewTxContext(ctx, tx)

	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("panic recovered: %v", recovered)
		}

		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != pgx.ErrTxClosed {
				err = InternalError("failed to rollback transaction", fmt.Errorf("rollback failed: %v: %w", rollbackErr, err))
			}

			return
		}

		if commitErr := tx.Commit(ctx); commitErr != nil {
			err = InternalError("failed to commit transaction", fmt.Errorf("tx commit failed: %w", commitErr))
		}
	}()

	if err = fn(ctx); err != nil {
		return err
	}

	return nil
}
