// Package apperror is used for storing generic application errors so other services can reuse them.
package apperror

import (
	"errors"
	"fmt"
)

const format = "%s: %s, %w"

var (
	// ErrNil args error must be returned if one of interface arguments are nil,
	// used to prevent panic.
	ErrNilArgs = errors.New("one or more arguments are nil")
)

// New creates new wrapped erros in format: "{msg}: {original_error}, {error_for_client}".
func New(msg string, err, clientErr error) error {
	return fmt.Errorf(format, msg, err.Error(), clientErr)
}
