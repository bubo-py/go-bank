package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n transactions concurrently
	n := 5
	amount := int64(10)

	errsCh := make(chan error)
	resultsCh := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errsCh <- err
			resultsCh <- result
		}()
	}

	// check results from channels
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errsCh
		require.NoError(t, err)

		result := <-resultsCh

		// transfer
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, amount, result.Transfer.ID)
		require.NotZero(t, amount, result.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// from entry
		require.Equal(t, account1.ID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)

		_, err = store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		// to entry
		require.Equal(t, account2.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.ID)

		_, err = store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		// accounts
		require.Equal(t, result.FromAccount.ID, account1.ID)
		require.Equal(t, result.ToAccount.ID, account2.ID)

		// balances
		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	// final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n transactions concurrently
	n := 10
	amount := int64(10)

	errsCh := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errsCh <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errsCh
		require.NoError(t, err)
	}

	// final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// balances shouldn't change since we have 5 transactions per each account with the same amount
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
