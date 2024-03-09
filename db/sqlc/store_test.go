package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < n; i++ {

		go func() {
			r, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- r
		}()
	}

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		r := <-results
		require.NotEmpty(t, r)

		require.NotZero(t, r.Transfer.ID)
		require.Equal(t, amount, r.Transfer.Amount)

		_, err = store.GetTransfer(context.Background(), r.Transfer.ID)

		require.NoError(t, err)

		fromEntry := r.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		// TODO: Check other entry

		fromAccount := r.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := r.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)

		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)

	}
}
