package postgres

import "time"

// Option -.
type Option func(*Client)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Client) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Client) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.connTimeout = timeout
	}
}
