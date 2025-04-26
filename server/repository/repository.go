package repository

import (
	"context"

	"github.com/gocraft/dbr/v2"
)

type TxKey struct{}

type Repository struct {
	sess *dbr.Session
}

func NewRepository(sess *dbr.Session) *Repository {
	return &Repository{
		sess: sess,
	}
}

func (r *Repository) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	ctx = context.WithValue(ctx, TxKey{}, tx)

	err = fn(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func getRunner(ctx context.Context, sess *dbr.Session) dbr.SessionRunner {
	value := ctx.Value(TxKey{})
	if value == nil {
		// RunInTxの中でない場合は、sessを返す
		return sess
	}

	tx, ok := value.(*dbr.Tx)
	if !ok {
		panic("context value is not a transaction")
	}

	return tx
}
