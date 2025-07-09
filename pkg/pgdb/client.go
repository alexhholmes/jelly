package pgdb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"jelly/pkg/api/photo"
)

// If running as a lambda, the db connection will be shared across invocations.
var db *sql.DB

type Client struct {
	db *sql.DB
}

func NewClient() (client *Client, err error) {
	if db != nil {
		// If the db connection is already open, return it (warm lambda)
		return &Client{db: db}, nil
	}

	// Open a connection to pgdb
	connStr := "host=postgres user=username dbname=jelly sslmode=disable password" +
		"=password"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return client, fmt.Errorf("failed to open a db connection: %w", err)
	}

	// Verify the connection
	if err = db.Ping(); err != nil {
		return client, fmt.Errorf("failed to ping db: %w", err)
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}
