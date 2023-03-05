package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()

}

type TransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferResult struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *Store) TransferTx(ctx context.Context, arg TransferParams) (TransferResult, error) {
	var result TransferResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		//Todo update account
		//This query is used to block the other transaction that read the data concurrently before do thing else.
		ta, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		///This query is used to block the other transaction that read the data concurrently before do thing else.
		fa, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		//Also we need to order update account from lower id to higher id because we not allow to transfer accId 1 to accId 2 and accId2 to accId1 concurrently due to a deadlock.
		//Always need to update lower account first! so that will be blocked by same account in diff transaction[after commit will unblock].

		if arg.ToAccountID < arg.FromAccountID {
			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.ToAccountID,
				Balance: ta.Balance + arg.Amount,
			})

			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: fa.Balance - arg.Amount,
			})

		} else {
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: ta.Balance - arg.Amount,
			})

			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.ToAccountID,
				Balance: fa.Balance + arg.Amount,
			})

		}

		return nil

	})

	return result, err
}
