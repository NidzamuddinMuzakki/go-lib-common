package data_source

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TransactionRunner struct {
	DB *sqlx.DB
}

type TxFunc func(tx *sqlx.Tx) error
type TxOpt func(t *TransactionRunner)

func SetDB(db *sqlx.DB) TxOpt {
	return func(t *TransactionRunner) {
		t.DB = db
	}
}

func NewTransactionRunner(db *sqlx.DB) *TransactionRunner {
	return &TransactionRunner{
		DB: db,
	}
}

func (t *TransactionRunner) WithTx(ctx context.Context, txFunc TxFunc, opts *sql.TxOptions) error {
	tx, err := StartTx(ctx, t.DB, opts)
	if err != nil {
		return err
	}

	err = txFunc(tx)
	if err != nil {
		errRb := RollbackTx(tx)
		if errRb != nil {
			return errRb
		}
		return err
	}

	return CommitTx(tx)
}

func StartTx(ctx context.Context, db *sqlx.DB, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return db.BeginTxx(ctx, opts)
}

func RollbackTx(tx *sqlx.Tx) error {
	return tx.Rollback()
}

func CommitTx(tx *sqlx.Tx) error {
	return tx.Commit()
}
