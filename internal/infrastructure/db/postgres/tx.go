package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type txType string

const keyTx txType = "txKey"

func (r *Repository) contextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, keyTx, tx)
}

func (r *Repository) txFromContext(ctx context.Context) (pgx.Tx, error) {
	switch v := ctx.Value(keyTx).(type) {
	case pgx.Tx:
		return v, nil
	default:
		return nil, errors.New("Cannot find transaction in context")
	}
}

func (r *Repository) connection(ctx context.Context) Querier {
	if tx, err := r.txFromContext(ctx); err == nil {
		return tx.Conn()
	}
	return r.db
}

// Exec упаковывыет транзакцию в контекст, откатывает транзакцию в случае неудачи, комитит изменения если все ОК
func (r *Repository) Exec(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	ctx = r.contextWithTx(ctx, tx)
	defer func(tx pgx.Tx, ctx context.Context) {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			// TODO: логировать
			fmt.Println("failed to rollback transaction", rollbackErr)
		}
	}(tx, ctx)

	err = f(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to execute transaction")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
