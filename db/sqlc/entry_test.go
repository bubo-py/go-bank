package db

import (
	"context"
	"testing"

	"go-bank/util"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	params := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomAmount(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), params)
	require.NoError(t, err)

	require.Equal(t, params.AccountID, entry.AccountID)
	require.Equal(t, params.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t, createRandomAccount(t))
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t, createRandomAccount(t))

	gotEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	require.Equal(t, entry.ID, gotEntry.ID)
	require.Equal(t, entry.AccountID, gotEntry.AccountID)
	require.Equal(t, entry.Amount, gotEntry.Amount)
	require.NotZero(t, entry.CreatedAt, gotEntry.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	params := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, params.AccountID, entry.AccountID)
	}
}
