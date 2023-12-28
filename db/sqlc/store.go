package db

import (
	"context"
	"database/sql"
	"fmt"
)

// # begin Tx
// transfer monay
// entry entry 1 in
// entry entry 2 out
// update balance
// # commit transaction.

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) executeTransaction(ctx context.Context, callBack func(q *Queries) (error)) error {
	// initialize transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(tx)
	err = callBack(queries)

	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("error rolling back transaction: %v", txErr)
		}
		return err
	}

	return tx.Commit()
}
