package pgclient

import "time"

type Option func(*Client)

func WithMaxConns(conns int32) Option {
	return func(c *Client) {
		c.maxConns = conns
	}
}

func WithConnAttempts(attempts uint8) Option {
	return func(c *Client) {
		c.connAttempts = attempts
	}
}

func WithConnTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.connTimeout = timeout
	}
}
