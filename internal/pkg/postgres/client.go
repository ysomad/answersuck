// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Client struct {
	maxPoolSize          int
	connAttempts         int
	connTimeout          time.Duration
	preferSimpleProtocol bool

	Pool *pgxpool.Pool
}

func NewClient(uri string, opts ...Option) (*Client, error) {
	c := &Client{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	poolConfig, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(c.maxPoolSize)
	poolConfig.ConnConfig.PreferSimpleProtocol = c.preferSimpleProtocol

	for c.connAttempts > 0 {
		c.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}
		defer c.Close()

		log.Printf("Postgres is trying to connect, attempts left: %d", c.connAttempts)

		time.Sleep(c.connTimeout)

		c.connAttempts--
	}
	if err != nil {
		return nil, fmt.Errorf("postgres - NewClient - connAttempts == 0: %w", err)
	}

	return c, nil
}

func (p *Client) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
