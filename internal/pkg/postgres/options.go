package postgres

import "time"

type Option func(*Client)

func MaxPoolSize(size int) Option {
	return func(c *Client) {
		c.maxPoolSize = size
	}
}

func ConnAttempts(attempts int) Option {
	return func(c *Client) {
		c.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.connTimeout = timeout
	}
}

func PreferSimpleProtocol(p bool) Option {
	return func(c *Client) {
		c.preferSimpleProtocol = p
	}
}
