package pgdb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func WithPostgres(t *testing.T) string {
	ctx := context.Background()
	postgresC, err := postgres.Run(ctx,
		"postgres:16",
		postgres.WithDatabase("jelly"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
		postgres.WithInitScripts("../../migrations/tables.sql"),
	)
	require.NoError(t, err)

	connStr, err := postgresC.ConnectionString(ctx, "sslmode=disable")
	t.Cleanup(func() {
		if err = postgresC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate postgres container: %v", err)
		}
	})
	require.NoError(t, err)

	return connStr
}

func TestNewClient(t *testing.T) {
	endpoint := WithPostgres(t)

	client, err := NewClient(endpoint)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Test closing the client
	err = client.Close()
	require.NoError(t, err)

	// Test that the client can be reopened
	client, err = NewClient(endpoint)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Verify that the photos table exists
	var count int
	err = client.db.Get(&count, "SELECT COUNT(*) FROM photos")
	require.NoError(t, err)
	require.Equal(t, 0, count, "Expected newsletters table to be empty after initialization")

	// Check non-existent table
	err = client.db.Get(&count, "SELECT COUNT(*) FROM test_non_existent_table")
	require.Error(t, err)

	err = client.Close()
	require.NoError(t, err)
}
