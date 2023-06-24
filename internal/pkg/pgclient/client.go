package pgclient

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxConns     = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Client struct {
	maxConns     int32
	connAttempts uint8
	connTimeout  time.Duration
	tracer       pgx.QueryTracer

	Builder sq.StatementBuilderType
	Pool    *pgxpool.Pool
}

func New(connString string, opts ...Option) (*Client, error) {
	c := &Client{
		maxConns:     defaultMaxConns,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	poolConf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	if c.tracer != nil {
		poolConf.ConnConfig.Tracer = c.tracer
	}

	poolConf.MaxConns = c.maxConns

	c.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConf)
	if err != nil {
		return nil, err
	}

	for c.connAttempts > 0 {
		if err = c.Pool.Ping(context.TODO()); err == nil {
			break
		}

		log.Printf("trying connecting to postgres, attempts left: %d", c.connAttempts)
		time.Sleep(c.connTimeout)
		c.connAttempts--
	}

	if err != nil {
		return nil, err
	}

	c.Builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return c, nil
}

func (p *Client) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
