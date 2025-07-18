package pgdb

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	db *sqlx.DB
}

func NewClient(endpoint string) (client *Client, err error) {
	// Open a connection to pgdb
	db, err := sqlx.Connect("postgres", endpoint)
	if err != nil {
		return client, fmt.Errorf("failed to open a db connection: %w", err)
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}
