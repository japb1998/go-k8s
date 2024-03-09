package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store defines all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {

	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
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

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// TODO: update account balance
		// get the account  ->  update its balance -> return the account
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)

		if err != nil {
			return err
		}
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{ID: account1.ID, Balance: account1.Balance - arg.Amount})

		if err != nil {
			return err
		}

		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)

		if err != nil {
			return err
		}
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{ID: account2.ID, Balance: account2.Balance + arg.Amount})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:      accountID1,
		Balance: amount1,
	})
	if err != nil {
		return
	}

	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:      accountID2,
		Balance: amount2,
	})
	if err != nil {
		return
	}
	return
}
