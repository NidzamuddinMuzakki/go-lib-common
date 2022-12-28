package data_source

import (
	"context"
	"github.com/jmoiron/sqlx"
)

// Exec wrapping multiple queries or single query without transaction.
func Exec(ctx context.Context, db *sqlx.DB, statements ...*Statement) error {
	err := run(ctx, db, statements...)
	if err != nil {
		return err
	}

	return nil
}

// ExecTx wrapping multiple queries or single query in a transaction.
func ExecTx(ctx context.Context, db *sqlx.DB, statements ...*Statement) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err = runTx(ctx, tx, statements...); err != nil {
		if er := tx.Rollback(); er != nil {
			return er
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
