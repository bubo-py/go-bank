package db

import (
	"context"
	"database/sql"
	"testing"

	"go-bank/util"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	params := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)
	require.NoError(t, err)

	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, params.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	gotAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)

	require.Equal(t, account.ID, gotAccount.ID)
	require.Equal(t, account.Owner, gotAccount.Owner)
	require.Equal(t, account.Balance, gotAccount.Balance)
	require.Equal(t, account.Currency, gotAccount.Currency)

	require.NotZero(t, gotAccount.CreatedAt)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	params := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	params := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomAmount(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), params)
	require.NoError(t, err)

	require.Equal(t, account.ID, updatedAccount.ID)
	require.Equal(t, account.Owner, updatedAccount.Owner)
	require.Equal(t, params.Balance, updatedAccount.Balance)
	require.Equal(t, account.Currency, updatedAccount.Currency)

	require.NotZero(t, updatedAccount.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	gotAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Empty(t, gotAccount)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
