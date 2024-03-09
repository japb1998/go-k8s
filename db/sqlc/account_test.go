package db

import (
	"context"
	"testing"
	"time"

	"simplebank/util"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.TODO(), arg)

	if err != nil {
		t.Error(err)
	}
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	a := createRandomAccount(t)
	a2, err := testQueries.GetAccount(context.Background(), a.ID)

	require.NoError(t, err)
	require.NotEmpty(t, a2)

	require.Equal(t, a.ID, a2.ID)
	require.Equal(t, a.Owner, a2.Owner)
	require.Equal(t, a.Balance, a2.Balance)
	require.Equal(t, a.Currency, a2.Currency)
	require.WithinDuration(t, a.CreatedAt, a2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	a, _ := testQueries.UpdateAccount(context.Background(), arg)

	require.Equal(t, a.ID, account1.ID)
	require.Equal(t, a.Owner, account1.Owner)
	require.Equal(t, a.Balance, arg.Balance)
	require.Equal(t, a.Currency, account1.Currency)
	require.WithinDuration(t, a.CreatedAt, account1.CreatedAt, time.Second)
	require.Equal(t, a.Balance, arg.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	n := 10
	for range n {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
