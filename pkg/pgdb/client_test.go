package pgdb

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

func InitEmptyTables(t *testing.T, connStr string) {
	client, err := NewClient(connStr)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	// Read tables.sql file
	query, err := os.ReadFile("../../migrations/tables.sql")
	require.NoError(t, err)

	_, err = client.db.Exec(string(query))
	require.NoError(t, err)
}

func InitTables(t *testing.T, connStr string) {
	client, err := NewClient(connStr)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	// Read all migration files
	err = filepath.WalkDir("../../migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".sql") {
			return nil
		}

		sqlContent, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		_, execErr := client.db.Exec(string(sqlContent))
		return execErr
	})
	require.NoError(t, err)
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

	err = client.Close()
	require.NoError(t, err)
}

func TestInitEmptyTables(t *testing.T) {
	endpoint := WithPostgres(t)

	InitEmptyTables(t, endpoint)

	client, err := NewClient(endpoint)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	// Verify that the photos table exists
	var count int
	err = client.db.Get(&count, "SELECT COUNT(*) FROM photos")
	require.NoError(t, err)
	require.Equal(t, 0, count, "Expected newsletters table to be empty after initialization")
}
