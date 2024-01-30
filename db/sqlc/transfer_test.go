package db

import (
	"context"
	"testing"

	"go-bank/util"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	params := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomAmount(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), params)
	require.NoError(t, err)

	require.Equal(t, params.FromAccountID, transfer.FromAccountID)
	require.Equal(t, params.ToAccountID, transfer.ToAccountID)
	require.Equal(t, params.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))
}

func TestGetTransfer(t *testing.T) {
	transfer := createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))

	gotTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	require.Equal(t, transfer.ID, gotTransfer.ID)
	require.Equal(t, transfer.FromAccountID, gotTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, gotTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, gotTransfer.Amount)

	require.NotZero(t, transfer.CreatedAt, gotTransfer.CreatedAt)
}

func TestListTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	params := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}
